package keycloakclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type IntrospectTokenResult struct {
	Exp    int      `json:"exp"`
	Iat    int      `json:"iat"`
	Aud    []string `json:"aud"`
	Active bool     `json:"active"`
}

func (i *IntrospectTokenResult) UnmarshalJSON(p []byte) error {
	var tmp map[string]interface{}
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if val, ok := tmp["exp"]; ok {
		switch v := val.(type) {
		case float64:
			i.Exp = int(v)
		case float32:
			i.Exp = int(v)
		case int:
			i.Exp = v
		case int64:
			i.Exp = int(v)
		case int32:
			i.Exp = int(v)
		}
	}
	if val, ok := tmp["iat"]; ok {
		switch v := val.(type) {
		case float64:
			i.Iat = int(v)
		case float32:
			i.Iat = int(v)
		case int:
			i.Iat = v
		case int64:
			i.Iat = int(v)
		case int32:
			i.Iat = int(v)
		}
	}
	if val, ok := tmp["active"]; ok {
		if v, ok := val.(bool); ok {
			i.Active = v
		}
	}
	if val, ok := tmp["aud"]; ok {
		switch v := val.(type) {
		case string:
			i.Aud = []string{v}
		case []interface{}:
			for _, vv := range v {
				if vvv, ok := vv.(string); ok {
					i.Aud = append(i.Aud, vvv)
				}
			}
		}
	}

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
