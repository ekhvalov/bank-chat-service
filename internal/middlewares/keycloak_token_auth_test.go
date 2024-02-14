package middlewares_test

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	internaljwt "github.com/ekhvalov/bank-chat-service/internal/jwt"
	jwtmocks "github.com/ekhvalov/bank-chat-service/internal/jwt/mocks"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const (
	issuer           = "http://localhost:3010/realms/Bank"
	requiredResource = "chat-ui-client"
	requiredRole     = "support-chat-client"
	secWsProtocol    = "chat-service-protocol"
)

func TestNewKeycloakTokenAuth(t *testing.T) {
	suite.Run(t, new(KeycloakTokenAuthSuite))
}

type KeycloakTokenAuthSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	keyFunc         jwt.Keyfunc
	keyFuncProvider *jwtmocks.MockKeyFuncProvider
	authMdlwr       echo.MiddlewareFunc
	req             *http.Request
	resp            *httptest.ResponseRecorder
	ctx             echo.Context
}

func (s *KeycloakTokenAuthSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	var err error
	s.publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	s.Require().NoError(err)
	s.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	s.Require().NoError(err)

	s.keyFunc = func(token *jwt.Token) (interface{}, error) { return s.publicKey, nil }
	s.keyFuncProvider = jwtmocks.NewMockKeyFuncProvider(s.ctrl)
	jwtParser, err := internaljwt.NewJWTParser(internaljwt.NewJWTParserOptions(s.keyFuncProvider, issuer))
	s.Require().NoError(err)
	s.authMdlwr = middlewares.NewKeycloakTokenAuth(jwtParser, requiredResource, requiredRole, secWsProtocol)
	s.Require().NoError(err)

	s.req = httptest.NewRequest(http.MethodPost, "/getHistory", bytes.NewBufferString(`{"pageSize": 100, "cursor": ""}`))
	s.resp = httptest.NewRecorder()
	s.ctx = echo.New().NewContext(s.req, s.resp)
}

func (s *KeycloakTokenAuthSuite) TearDownTest() {
	s.ctrl.Finish()
}

// Positive.

func (s *KeycloakTokenAuthSuite) TestValidToken_AudString() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtY2xpZW50Il19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	var uid types.UserID

	err := s.authMdlwr(func(c echo.Context) error {
		uid = middlewares.MustUserID(c)
		return nil
	})(s.ctx)
	s.Require().NoError(err)
	s.Equal("5cb40dc0-a249-4783-a301-9e1f3cf3ea41", uid.String())
}

func (s *KeycloakTokenAuthSuite) TestValidToken_AudList() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	var uid types.UserID

	err := s.authMdlwr(func(c echo.Context) error {
		uid = middlewares.MustUserID(c)
		return nil
	})(s.ctx)
	s.Require().NoError(err)
	s.Equal("5cb40dc0-a249-4783-a301-9e1f3cf3ea41", uid.String())
}

// Negative.

func (s *KeycloakTokenAuthSuite) TestNoAuthorizationHeader() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add("Authentication", "Bearer "+token)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusBadRequest)
}

func (s *KeycloakTokenAuthSuite) TestNotBearerAuth() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Basic "+token)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusBadRequest)
}

func (s *KeycloakTokenAuthSuite) TestKeyFuncProviderError() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtY2xpZW50Il19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(nil, context.Canceled)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, context.Canceled)
}

func (s *KeycloakTokenAuthSuite) TestInvalidExpiresAt() {
	const claims = "eyJleHAiOjE2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)

	// var jwtErr *jwt.ValidationError
	s.Require().ErrorIs(err, jwt.ErrTokenExpired)
	// s.Empty(jwtErr.Errors ^ jwt.ValidationErrorExpired)
}

func (s *KeycloakTokenAuthSuite) TestInvalidIssuedAt() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MjY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)

	// var jwtErr *jwt.ValidationError
	s.Require().ErrorIs(err, jwt.ErrTokenUsedBeforeIssued)
	// s.Empty(jwtErr.Errors ^ jwt.ValidationErrorIssuedAt)
}

func (s *KeycloakTokenAuthSuite) TestInvalidIssuer() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MjY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vd3JvbmcuY29tOjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSIsImFsZyI6IkhTMjU2In0" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)

	// var jwtErr *jwt.ValidationError
	s.Require().ErrorIs(err, jwt.ErrTokenInvalidIssuer)
	// s.Empty(jwtErr.Errors ^ jwt.ValidationErrorIssuedAt)
}

func (s *KeycloakTokenAuthSuite) TestNoSubject() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrSubjectNotDefined)
}

func (s *KeycloakTokenAuthSuite) TestSubjectIsZeroUUID() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrSubjectNotDefined)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceAccess_NoKey() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9uZDAwNyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiYm9uZDAwN0B1ay5jb20ifQ" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoAllowedResources)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceAccess_EmptyMap() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6e30sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9uZDAwNyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiYm9uZDAwN0B1ay5jb20ifQ" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoAllowedResources)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceRole_NoNeededResource() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoRequiredResourceRole)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceRole_NoNeededRole() {
	const claims = "eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJhYnJhY2FkYWJyYSJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIiwic3VwcG9ydC1jaGF0LWNsaWVudCJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9uZDAwNyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiYm9uZDAwN0B1ay5jb20ifQ" //nolint:lll
	token := s.signClaims(claims)
	s.req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)

	s.keyFuncProvider.EXPECT().Keyfunc(gomock.Any()).Return(s.publicKey, nil)

	err := s.authMdlwr(func(c echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoRequiredResourceRole)
}

func (s *KeycloakTokenAuthSuite) assertHTTPCode(err error, code int) {
	var httpErr *echo.HTTPError
	s.Require().ErrorAs(err, &httpErr)
	s.Equal(httpErr.Code, code)
}

func TestMustUserID_NoUID(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	assert.Panics(t, func() {
		middlewares.MustUserID(echo.New().NewContext(req, httptest.NewRecorder()))
	})
}

func (s *KeycloakTokenAuthSuite) signClaims(claims string) string {
	claimsJSON, err := base64.RawURLEncoding.DecodeString(claims)
	s.Require().NoError(err)

	var claimsMap map[string]interface{}
	err = json.Unmarshal(claimsJSON, &claimsMap)
	s.Require().NoError(err)

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claimsMap)).SignedString(s.privateKey)
	s.Require().NoError(err)

	return token
}

const (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAzf10MG/4YiDJ7M94FaVIL7sZ1z/fJKyTEm3fbJ4PgownCTv3
o3adWZYhNRdGwu/YOhKak2uSOxQUj15QwaCFjmlVCwKuaJeXbI5BNHct46Kzo0pj
aX5SiY1RhCPxiZtfGk/OaRXbiyU+yHNffY7TTvpAyLoFNTgn7OiiYPWPSCOmZ2zQ
L+1judRIyjP1Z1aIwenmD+LoyPZ+RQ9TrdZXKHi5DxgdV/f660smWHICiMBEAJ5a
kcu/uemvJbmBCJkJPoeQWz39x3t1OrMWE0G/Ocs09tUDUzdxXNes+RDLyx+b0J0O
zUIq/+m3rWJRpe+6ErWhGvj7mBHlm8aQirBerwIDAQABAoIBAH+zkjV5JP4In8ZM
tICOz9qvXozADyFYT3EMZoea0bi4FHc4EwTmwxPH69xTCs5NDLqrz+J2vNgdUcWz
zdLMJiAskslZpzA2Umy9IBVbkTpfIoin1EuRQa/+yTtnYRVTGjlgonEpWMrBk1OH
mvpm8f8zS7hlAleE8dOAQTJk6afpPTyNvj1baN9okdpNZ7+5pK9Ij+YcS+aOWLix
A+vsIm5b0W6eXXnJLZzXNr2N2O9P/iEIdOs0+cvP58rkNQ/d4flZ6AYnUCgHHei4
gZxCWZgMHXzdY/t/pFM+l3G1QJzlGM6L8sIcXToTYmE1xJEf4PCV0ILt/cHUFkqU
HGeExnECgYEA6SpWbNgClEkl2NG9qmNCsXNOVWKj7mRylPJsKVLMZq208SekmTMj
qMqeN0wnwhyhmM//nYYu6dxTxhJBvXxYJyMz0G27p/5HaZC7WAl5XYM6QPPmsky5
T4h11J2X8TLBIjrE3wQm5d/EL9i3UqZtffTfFRwv6r4r0dGJdYZ2WokCgYEA4inO
iKAVd1ERIIRL9pBb0fVCVJSlX86NR3VRvB16fyrCrFjx1IE8CaeEKcYmnE7Abe0J
/jSM5OHKSULGbl2DeofhT2FhgV+hM/wKd3G3dVaHiuMWO9lCwnwelbXq2Rt0hhN0
b1YVHkI8rWMC1RDvK8Z9cExLz9VH+VJq+41TwXcCgYEAhf4cmIQyR0EaDNXLp0VP
qGZZF9yN1IvJBSujWMQKTt94YjWj855d2bxG3ARZvMVzYDv363CXOTGyutr3CIuS
pTsnpZnKA6qvI01XPCqFomWtbnI7my9YNwp2nG7MSIIgVylqxba/G89SEST7hPW7
amz0Xk9Kgh4zVGqUEgPps/ECgYEAnhR6uCSs3Gldf0z5i637gBXd9yCvNvg45+mo
58PzC0/oIm9JGS/7twPP7SMDed3RwwQcKAKzOIhZzDtQV3Qlok+3vLRkYvlkw+E3
r6VchjelJf70W4DQmQAIoLw3GumF2PFgQTH6MNw7bTX3lNXxVre2lfe+RdbeJ/bj
sFBoaqECgYEAzK91/ea5p5Hlt5yCQLeLDKSf2ohmYspkqk0HTi8iGfji2Zo99Iir
1rFR0Oe3otPG40HXhKDi2YdhNy/D4ypaVDkr94awTBYY8zlmgAPhf/oZu48tkxCh
qIanZhvea4LFXIctQKhXDCH0qwTkR9adILLKgLBS/dTrzWG2JHBE1B8=
-----END RSA PRIVATE KEY-----`

	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzf10MG/4YiDJ7M94FaVI
L7sZ1z/fJKyTEm3fbJ4PgownCTv3o3adWZYhNRdGwu/YOhKak2uSOxQUj15QwaCF
jmlVCwKuaJeXbI5BNHct46Kzo0pjaX5SiY1RhCPxiZtfGk/OaRXbiyU+yHNffY7T
TvpAyLoFNTgn7OiiYPWPSCOmZ2zQL+1judRIyjP1Z1aIwenmD+LoyPZ+RQ9TrdZX
KHi5DxgdV/f660smWHICiMBEAJ5akcu/uemvJbmBCJkJPoeQWz39x3t1OrMWE0G/
Ocs09tUDUzdxXNes+RDLyx+b0J0OzUIq/+m3rWJRpe+6ErWhGvj7mBHlm8aQirBe
rwIDAQAB
-----END PUBLIC KEY-----`
)
