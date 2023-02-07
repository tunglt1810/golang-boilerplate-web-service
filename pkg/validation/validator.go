package validation

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(i interface{}) error
	RegisterCustomValidation(tag string, fn validator.Func, message string, override bool, callValidationEvenIfNull ...bool) error
	RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{})
}

type customValidator struct {
	validator *validator.Validate
}

func New() (Validator, error) {
	v, err := InitValidate()
	if err != nil {
		return nil, err
	}
	return &customValidator{
		validator: v,
	}, nil
}

func NewWithValidate(v *validator.Validate) Validator {
	return &customValidator{
		validator: v,
	}
}

func (v *customValidator) Validate(i interface{}) error {
	return Validate(v.validator, i)
}

func (v *customValidator) RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	v.validator.RegisterStructValidation(fn, types...)
}

func (v *customValidator) RegisterCustomValidation(
	tag string, fn validator.Func, message string, override bool, callValidationEvenIfNull ...bool) error {
	if err := v.validator.RegisterValidation(tag, fn, callValidationEvenIfNull...); err != nil {
		return err
	}

	var innerError error

	if err := v.validator.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, message, override)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(tag, fmt.Sprintf("%s", fe.Value()))
		if err != nil {
			innerError = err
		}
		return t
	}); err != nil {
		return err
	}

	return innerError
}
