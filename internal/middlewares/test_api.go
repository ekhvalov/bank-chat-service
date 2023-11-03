package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func SetToken(c echo.Context, uid types.UserID) {
	c.Set(tokenCtxKey, jwt.NewWithClaims(jwt.SigningMethodNone, claims{
		StandardClaims: jwt.StandardClaims{},
		Aud:            []string{"something"},
		AuthTime:       0,
		Subject:        uid,
		ResourceAccess: map[string]access{"something": {Roles: []string{"role"}}},
	}))
}
