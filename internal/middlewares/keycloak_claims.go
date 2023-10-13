package middlewares

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"

	keycloakclient "github.com/ekhvalov/bank-chat-service/internal/clients/keycloak"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
	parser                = &jwt.Parser{}
)

type claims struct {
	jwt.StandardClaims
	Aud            keycloakclient.Audition `json:"aud"`
	AuthTime       int64                   `json:"auth_time"`
	ResourceAccess map[string]access       `json:"resource_access,omitempty"`
	userID         types.UserID
}

func (c claims) HasResourceWithRole(resource, role string) bool {
	if nil == c.ResourceAccess || len(c.ResourceAccess) == 0 {
		return false
	}
	if _, ok := c.ResourceAccess[resource]; ok {
		return c.ResourceAccess[resource].HasRole(role)
	}
	return false
}

type access struct {
	Roles []string `json:"roles"`
}

func (a access) HasRole(name string) bool {
	for _, role := range a.Roles {
		if role == name {
			return true
		}
	}
	return false
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if nil == c.ResourceAccess || len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}
	if c.userID.IsZero() {
		return ErrSubjectNotDefined
	}
	if nil == c.Aud {
		return ErrNoAllowedResources
	}
	return nil
}

func (c claims) UserID() types.UserID {
	return c.userID
}

func parseTokenAndClaims(tokenStr string) (*jwt.Token, *claims, error) {
	c := claims{}
	t, _, err := parser.ParseUnverified(tokenStr, &c)
	if err != nil {
		return nil, nil, fmt.Errorf("parse token: %v", err)
	}
	uid, err := types.Parse[types.UserID](c.Subject)
	if err != nil {
		return nil, nil, ErrSubjectNotDefined
	}
	c.userID = uid
	return t, &c, nil
}
