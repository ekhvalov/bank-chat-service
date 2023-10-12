package keycloakclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type IntrospectTokenResult struct {
	Exp    int  `json:"exp"`
	Iat    int  `json:"iat"`
	Aud    aud  `json:"aud"`
	Active bool `json:"active"`
}

type aud []string

func (a *aud) UnmarshalJSON(data []byte) error {
	var s string
	if nil == json.Unmarshal(data, &s) {
		*a = []string{s}
		return nil
	}
	var ss []string
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	*a = ss
	return nil
}

// IntrospectToken implements
// https://www.keycloak.org/docs/latest/authorization_services/index.html#obtaining-information-about-an-rpt
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectTokenResult, error) {
	url := fmt.Sprintf("realms/%s/protocol/openid-connect/token/introspect", c.realm)

	var result IntrospectTokenResult

	resp, err := c.auth(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"token_type_hint": "requesting_party_token",
			"token":           token,
		}).
		SetResult(&result).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("send request to keycloak: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("errored keycloak response: %v", resp.Status())
	}

	return &result, nil
}

func (c *Client) auth(ctx context.Context) *resty.Request {
	return c.cli.R().
		SetBasicAuth(c.username, c.password).
		SetContext(ctx)
}
