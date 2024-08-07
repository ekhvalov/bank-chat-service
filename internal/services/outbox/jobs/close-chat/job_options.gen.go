// Code generated by options-gen. DO NOT EDIT.
package closechatjob

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	problemRepo problemsRepository,
	messageRepo messagesRepository,
	eventStream eventStream,
	log *zap.Logger,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.problemRepo = problemRepo
	o.messageRepo = messageRepo
	o.eventStream = eventStream
	o.log = log

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("problemRepo", _validate_Options_problemRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("messageRepo", _validate_Options_messageRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("eventStream", _validate_Options_eventStream(o)))
	errs.Add(errors461e464ebed9.NewValidationError("log", _validate_Options_log(o)))
	return errs.AsError()
}

func _validate_Options_problemRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.problemRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `problemRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_messageRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.messageRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `messageRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_eventStream(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.eventStream, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `eventStream` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_log(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.log, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `log` did not pass the test: %w", err)
	}
	return nil
}
