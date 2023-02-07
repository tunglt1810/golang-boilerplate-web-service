package errors

import "errors"

var (
	mapError       map[error]HTTPError
	mapErrorStatus map[string]error
)

func Init() {
	mapError = make(map[error]HTTPError)
	mapErrorStatus = make(map[string]error)
	goErrErrMappings := []struct {
		goErr error
		err   *HTTPError
	}{
		{ErrUnknown, HTTPErrInternalServerError},
		{ErrNotImplemented, HTTPErrsNotImplemented},
		{ErrBadRequest, HTTPErrBadRequest},
		{ErrRequestBindingFailed, HTTPErrRequestBindingFailed},
		{ErrRequestTimeout, HTTPErrRequestTimeout},
		{ErrHTTPUnsupportedMediaType, HTTPErrUnsupportedMediaType},
		{ErrUnauthorized, HTTPErrUnauthorized},
		{ErrInsufficientPermission, HTTPErrInsufficientPermission},
		{ErrRouteNotFound, HTTPErrRouteNotFound},
		{ErrAlreadyExists, HTTPErrAlreadyExists},
		{ErrResourceNotFound, HTTPErrResourceNotFound},
		{ErrSMSignatureCannotBeVerify, HTTPErrSMSignatureCannotBeVerify},
		{ErrBearerTokenInvalid, HTTPErrTokenInvalid},
		{ErrTooManyRequests, HTTPErrTooManyRequests},
	}
	for _, m := range goErrErrMappings {
		AddGoErrorErrorMapping(m.goErr, m.err)
	}
}

func AddGoErrorErrorMapping(err error, he *HTTPError) {
	heCpy := *he
	heCpy.Status = err.Error()
	if heCpy.Message == "" {
		heCpy.Message = err.Error()
	}

	if existedErr, exists := mapErrorStatus[heCpy.Status]; exists && !errors.Is(err, existedErr) {
		panic("different errors with same error status")
	}

	mapErrorStatus[heCpy.Status] = err
	mapError[err] = heCpy
}

func HTTPErrFromMapping(goErr error) *HTTPError {
	if goErr == nil {
		return nil
	}

	if herr, exists := mapError[goErr]; exists {
		var err Error
		if errors.As(goErr, &err) {
			herr.Details = goErr.(Error).Details()
		}
		return &herr
	}

	return nil
}
