package middlewares

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/keycloak_mock.gen.go -typed -package=middlewaresmocks
type KeyFuncProvider interface {
	Keyfunc(token *jwt.Token) (any, error)
}

//go:generate options-gen -out-filename=keycloak_jwt_parser_options.gen.go -from-struct=JWTParserOptions
type JWTParserOptions struct {
	keyFuncProvider KeyFuncProvider `option:"mandatory" validate:"required"`
	issuer          string          `option:"mandatory" validate:"required"`
}

func NewJWTParser(opts JWTParserOptions) (*JWTParser, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	return &JWTParser{
		JWTParserOptions: opts,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
			jwt.WithIssuedAt(),
			jwt.WithIssuer(opts.issuer),
		),
	}, nil
}

// q: write documentation

type JWTParser struct {
	JWTParserOptions
	parser *jwt.Parser
}

func (p *JWTParser) ParseWithClaims(token string, claims jwt.Claims) (*jwt.Token, error) {
	return p.parser.ParseWithClaims(token, claims, p.keyFuncProvider.Keyfunc)
}
