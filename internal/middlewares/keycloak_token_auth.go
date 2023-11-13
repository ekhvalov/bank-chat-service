package middlewares

//go:generate mockgen -source=$GOFILE -destination=mocks/introspector_mock.gen.go -package=middlewaresmocks Introspector

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const tokenCtxKey = "user-token"

var ErrNoRequiredResourceRole = errors.New("no required resource role")

type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

// NewKeycloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeycloakTokenAuth(introspector Introspector, resource, role, secWsProtocol string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		// KeyLookup:  "header:" + echo.HeaderAuthorization,
		KeyLookup: "header:Authorization,header:Sec-WebSocket-Protocol",
		// AuthScheme: "Bearer",
		Validator: func(auth string, eCtx echo.Context) (bool, error) {
			tokenStr := parseToken(auth, secWsProtocol)
			if result, err := introspector.IntrospectToken(eCtx.Request().Context(), tokenStr); err != nil {
				return false, err
			} else if !result.Active {
				return false, nil
			}
			token, c, err := parseTokenAndClaims(tokenStr)
			if err != nil {
				return false, err
			}
			if err := c.Valid(); err != nil {
				return false, err
			}
			if !c.HasResourceWithRole(resource, role) {
				return false, ErrNoRequiredResourceRole
			}
			eCtx.Set(tokenCtxKey, token)
			return true, nil
		},
	})
}

func MustUserID(eCtx echo.Context) types.UserID {
	uid, ok := userID(eCtx)
	if !ok {
		panic("no user token in request context")
	}
	return uid
}

func userID(eCtx echo.Context) (types.UserID, bool) {
	t := eCtx.Get(tokenCtxKey)
	if t == nil {
		return types.UserIDNil, false
	}

	tt, ok := t.(*jwt.Token)
	if !ok {
		return types.UserIDNil, false
	}

	userIDProvider, ok := tt.Claims.(interface{ UserID() types.UserID })
	if !ok {
		return types.UserIDNil, false
	}
	return userIDProvider.UserID(), true
}

func parseToken(auth, secWsProtocol string) string {
	if strings.HasPrefix(auth, secWsProtocol) {
		return strings.TrimLeft(auth[len(secWsProtocol):], ", ")
	}
	return auth
}
