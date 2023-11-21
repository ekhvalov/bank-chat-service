// Code generated by options-gen. DO NOT EDIT.
package outbox

import (
	fmt461e464ebed9 "fmt"
	"time"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	workers int,
	idleTime time.Duration,
	reserveFor time.Duration,
	jobsRepo jobsRepository,
	txtor transactor,
	lg *zap.Logger,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.workers = workers
	o.idleTime = idleTime
	o.reserveFor = reserveFor
	o.jobsRepo = jobsRepo
	o.txtor = txtor
	o.lg = lg

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("workers", _validate_Options_workers(o)))
	errs.Add(errors461e464ebed9.NewValidationError("idleTime", _validate_Options_idleTime(o)))
	errs.Add(errors461e464ebed9.NewValidationError("reserveFor", _validate_Options_reserveFor(o)))
	errs.Add(errors461e464ebed9.NewValidationError("jobsRepo", _validate_Options_jobsRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("txtor", _validate_Options_txtor(o)))
	errs.Add(errors461e464ebed9.NewValidationError("lg", _validate_Options_lg(o)))
	return errs.AsError()
}

func _validate_Options_workers(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.workers, "min=1,max=32"); err != nil {
		return fmt461e464ebed9.Errorf("field `workers` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_idleTime(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.idleTime, "min=100ms,max=10s"); err != nil {
		return fmt461e464ebed9.Errorf("field `idleTime` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_reserveFor(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.reserveFor, "min=1s,max=10m"); err != nil {
		return fmt461e464ebed9.Errorf("field `reserveFor` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_jobsRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.jobsRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `jobsRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_txtor(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.txtor, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `txtor` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_lg(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.lg, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `lg` did not pass the test: %w", err)
	}
	return nil
}
