// Code generated by options-gen. DO NOT EDIT.
package keycloakclient

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	basePath string,
	realm string,
	username string,
	password string,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.basePath = basePath
	o.realm = realm
	o.username = username
	o.password = password

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithDebugMode(opt bool) OptOptionsSetter {
	return func(o *Options) {
		o.debugMode = opt
	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("basePath", _validate_Options_basePath(o)))
	errs.Add(errors461e464ebed9.NewValidationError("realm", _validate_Options_realm(o)))
	errs.Add(errors461e464ebed9.NewValidationError("username", _validate_Options_username(o)))
	errs.Add(errors461e464ebed9.NewValidationError("password", _validate_Options_password(o)))
	return errs.AsError()
}

func _validate_Options_basePath(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.basePath, "required,http_url"); err != nil {
		return fmt461e464ebed9.Errorf("field `basePath` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_realm(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.realm, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `realm` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_username(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.username, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `username` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_password(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.password, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `password` did not pass the test: %w", err)
	}
	return nil
}
