package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	internaljwt "github.com/ekhvalov/bank-chat-service/internal/jwt"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const tokenCtxKey = "user-token"

var ErrNoRequiredResourceRole = errors.New("no required resource role")

// NewKeycloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeycloakTokenAuth(parser *internaljwt.JWTParser, resource, role, secWsProtocol string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Authorization,header:Sec-WebSocket-Protocol",
		Validator: func(auth string, eCtx echo.Context) (bool, error) {
			tokenStr := extractToken(auth, secWsProtocol)

			c := claims{}
			token, err := parser.ParseWithClaims(tokenStr, &c)
			if err != nil {
				// var ve *jwt.ValidationError
				// if ok := errors.As(err, &ve); ok {
				// 	switch {
				// 	case errors.Is(ve.Inner, ErrSubjectNotDefined):
				// 		return false, ErrSubjectNotDefined
				// 	case errors.Is(ve.Inner, ErrNoAllowedResources):
				// 		return false, ErrNoAllowedResources
				// 	case ve.Errors&jwt.ValidationErrorUnverifiable != 0:
				// 		return false, ve.Inner
				// 	}
				// }
				return false, fmt.Errorf("parse token: %w", err)
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

func extractToken(auth, secWsProtocol string) string {
	if strings.HasPrefix(auth, secWsProtocol) {
		return strings.TrimLeft(auth[len(secWsProtocol):], ", ")
	}
	return auth
}
