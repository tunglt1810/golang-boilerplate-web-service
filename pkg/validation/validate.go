package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	enTran "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance of Validate, it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func InitValidate() (*validator.Validate, error) {
	var _ *validator.Validate
	enLocale := en.New()
	uni = ut.New(enLocale, enLocale)
	_ = validator.New()

	trans, _ = uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return name
	})
	err := enTran.RegisterDefaultTranslations(v, trans)

	return v, err
}

func RegisterCustomValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func ValidateStruct(req interface{}) (bool, error) {
	err := Validate(validate, req)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Validate(v *validator.Validate, i interface{}) error {
	// Validate input data
	var errs Errors
	if err := v.Struct(i); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		for _, fe := range validateErrs {
			errs = append(errs, Error{
				Field:   fe.Field(),
				Message: fe.Translate(trans),
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
