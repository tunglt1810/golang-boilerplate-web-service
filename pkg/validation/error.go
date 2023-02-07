package validation

import (
	"errors"
	"strings"
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

type Errors []Error

func (errs Errors) Error() string {
	errArr := make([]string, len(errs))
	for i, e := range errs {
		errArr[i] = e.Error()
	}
	return strings.Join(errArr, string('\n'))
}

func (errs Errors) Is(err error) bool {
	return errors.As(err, &errs)
}
