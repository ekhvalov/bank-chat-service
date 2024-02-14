// Code generated by options-gen. DO NOT EDIT.
package server

import (
	fmt461e464ebed9 "fmt"

	internaljwt "github.com/ekhvalov/bank-chat-service/internal/jwt"
	websocketstream "github.com/ekhvalov/bank-chat-service/internal/websocket-stream"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	addr string,
	allowOrigins []string,
	accessResource string,
	accessRole string,
	secWsProtocol string,
	logger *zap.Logger,
	jwtParser *internaljwt.JWTParser,
	errorHandler echo.HTTPErrorHandler,
	wsHandler *websocketstream.HTTPHandler,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.addr = addr
	o.allowOrigins = allowOrigins
	o.accessResource = accessResource
	o.accessRole = accessRole
	o.secWsProtocol = secWsProtocol
	o.logger = logger
	o.jwtParser = jwtParser
	o.errorHandler = errorHandler
	o.wsHandler = wsHandler

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("addr", _validate_Options_addr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("allowOrigins", _validate_Options_allowOrigins(o)))
	errs.Add(errors461e464ebed9.NewValidationError("accessResource", _validate_Options_accessResource(o)))
	errs.Add(errors461e464ebed9.NewValidationError("accessRole", _validate_Options_accessRole(o)))
	errs.Add(errors461e464ebed9.NewValidationError("secWsProtocol", _validate_Options_secWsProtocol(o)))
	errs.Add(errors461e464ebed9.NewValidationError("logger", _validate_Options_logger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("jwtParser", _validate_Options_jwtParser(o)))
	errs.Add(errors461e464ebed9.NewValidationError("errorHandler", _validate_Options_errorHandler(o)))
	errs.Add(errors461e464ebed9.NewValidationError("wsHandler", _validate_Options_wsHandler(o)))
	return errs.AsError()
}

func _validate_Options_addr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.addr, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `addr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_allowOrigins(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.allowOrigins, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `allowOrigins` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_accessResource(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.accessResource, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `accessResource` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_accessRole(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.accessRole, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `accessRole` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_secWsProtocol(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.secWsProtocol, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `secWsProtocol` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_logger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.logger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `logger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_jwtParser(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.jwtParser, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `jwtParser` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_errorHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.errorHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `errorHandler` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_wsHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.wsHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `wsHandler` did not pass the test: %w", err)
	}
	return nil
}
