// Code generated by options-gen. DO NOT EDIT.
package sendmanagermessagejob

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	producer messageProducer,
	messagesRepo messagesRepository,
	chatsRepository chatsRepository,
	evStream eventStream,
	log *zap.Logger,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.producer = producer
	o.messagesRepo = messagesRepo
	o.chatsRepository = chatsRepository
	o.evStream = evStream
	o.log = log

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("producer", _validate_Options_producer(o)))
	errs.Add(errors461e464ebed9.NewValidationError("messagesRepo", _validate_Options_messagesRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("chatsRepository", _validate_Options_chatsRepository(o)))
	errs.Add(errors461e464ebed9.NewValidationError("evStream", _validate_Options_evStream(o)))
	errs.Add(errors461e464ebed9.NewValidationError("log", _validate_Options_log(o)))
	return errs.AsError()
}

func _validate_Options_producer(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.producer, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `producer` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_messagesRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.messagesRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `messagesRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_chatsRepository(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.chatsRepository, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `chatsRepository` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_evStream(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.evStream, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `evStream` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_log(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.log, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `log` did not pass the test: %w", err)
	}
	return nil
}