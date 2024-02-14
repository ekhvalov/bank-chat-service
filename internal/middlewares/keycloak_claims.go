package middlewares

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

type claims struct {
	jwt.RegisteredClaims

	Subject        types.UserID      `json:"sub,omitempty"`
	ResourceAccess map[string]access `json:"resource_access,omitempty"`
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

// Validate checks validity of the claims (will be called by jwt-go library parser).
// possible errors:
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map, or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c claims) Validate() error {
	if nil == c.ResourceAccess || len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}

	if c.Subject.IsZero() {
		return ErrSubjectNotDefined
	}

	if nil == c.Audience || len(c.Audience) == 0 {
		return ErrNoAllowedResources
	}

	return nil
}

func (c claims) UserID() types.UserID {
	return c.Subject
}
