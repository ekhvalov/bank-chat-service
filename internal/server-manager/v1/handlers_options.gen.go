// Code generated by options-gen. DO NOT EDIT.
package managerv1

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	lg *zap.Logger,
	canReceiveProblemsUC canReceiveProblemsUsecase,
	freeHandsUC freeHandsUsecase,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.lg = lg
	o.canReceiveProblemsUC = canReceiveProblemsUC
	o.freeHandsUC = freeHandsUC

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("lg", _validate_Options_lg(o)))
	errs.Add(errors461e464ebed9.NewValidationError("canReceiveProblemsUC", _validate_Options_canReceiveProblemsUC(o)))
	errs.Add(errors461e464ebed9.NewValidationError("freeHandsUC", _validate_Options_freeHandsUC(o)))
	return errs.AsError()
}

func _validate_Options_lg(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.lg, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `lg` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_canReceiveProblemsUC(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.canReceiveProblemsUC, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `canReceiveProblemsUC` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_freeHandsUC(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.freeHandsUC, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `freeHandsUC` did not pass the test: %w", err)
	}
	return nil
}
