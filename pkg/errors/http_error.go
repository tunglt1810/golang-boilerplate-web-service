package errors

import "fmt"

type HTTPError struct {
	HTTPStatus int `json:"-"`
	Code       int `json:"error,omitempty"`

	Status  string      `json:"error_status,omitempty"`
	Message string      `json:"error_message,omitempty"`
	Details interface{} `json:"error_details,omitempty"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(
		"httpStatus: %d; code: %d; status:%s; message: %s", e.HTTPStatus, e.Code, e.Status, e.Message)
}

func NewWithCode(httpStatus, code int, msg ...string) *HTTPError {
	fe := &HTTPError{
		HTTPStatus: httpStatus,
		Code:       code,
	}
	if len(msg) > 0 {
		fe.Message = msg[0]
	}

	return fe
}
