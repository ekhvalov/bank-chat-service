// Code generated by options-gen. DO NOT EDIT.
package internaljwt

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptJWTParserOptionsSetter func(o *JWTParserOptions)

func NewJWTParserOptions(
	keyFuncProvider KeyFuncProvider,
	issuer string,
	options ...OptJWTParserOptionsSetter,
) JWTParserOptions {
	o := JWTParserOptions{}

	// Setting defaults from field tag (if present)

	o.keyFuncProvider = keyFuncProvider
	o.issuer = issuer

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *JWTParserOptions) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("keyFuncProvider", _validate_JWTParserOptions_keyFuncProvider(o)))
	errs.Add(errors461e464ebed9.NewValidationError("issuer", _validate_JWTParserOptions_issuer(o)))
	return errs.AsError()
}

func _validate_JWTParserOptions_keyFuncProvider(o *JWTParserOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keyFuncProvider, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keyFuncProvider` did not pass the test: %w", err)
	}
	return nil
}

func _validate_JWTParserOptions_issuer(o *JWTParserOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.issuer, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `issuer` did not pass the test: %w", err)
	}
	return nil
}
