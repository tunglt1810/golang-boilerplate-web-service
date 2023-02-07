package errors

import "net/http"

// List common errors.
var (
	ErrUnknown                   = New("ERR_UNKNOWN")
	ErrBadRequest                = New("ERR_BAD_REQUEST")
	ErrRequestBindingFailed      = New("ERR_REQUEST_BINDING_FAILED")
	ErrRequestTimeout            = New("ERR_REQUEST_TIMEOUT")
	ErrHTTPUnsupportedMediaType  = New("ERR_HTTP_UNSUPPORTED_MEDIA_TYPE")
	ErrUnauthorized              = New("ERR_UNAUTHORIZED")
	ErrInsufficientPermission    = New("ERR_INSUFFICIENT_PERMISSION")
	ErrNotFound                  = New("ERR_NOT_FOUND")
	ErrRouteNotFound             = New("ERR_ROUTE_NOT_FOUND")
	ErrNotImplemented            = New("ERR_NOT_IMPLEMENTED")
	ErrSMSignatureCannotBeVerify = New("ERR_SM_SIGNATURE_CANNOT_BE_VERIFY")
	ErrBearerTokenInvalid        = New("ERR_BEARER_INVALID_TOKEN")
	ErrTooManyRequests           = New("ERR_TOO_MANY_REQUESTS")
	ErrServiceUnavailable        = New("ERR_SERVICE_UNAVAILABLE")
	ErrGatewayTimeout            = New("ERR_GATEWAY_TIMEOUT")

	ErrAlreadyExists    = New("ERR_ALREADY_EXISTS")
	ErrResourceNotFound = New("ERR_RESOURCE_NOT_FOUND")

	// Validation error
	ErrInvalidArgument = New("ERR_INVALID_ARGUMENT")
)

// Common HTTP Errors used in services.
//
//nolint:gomnd
var (
	// ========================================================================
	// Client errors
	// ========================================================================

	// HTTPErrBadRequest is a common bad request error.
	HTTPErrBadRequest = NewWithCode(http.StatusBadRequest, 400000, "Bad Request")

	HTTPErrInvalidArgument = NewWithCode(http.StatusBadRequest, 400002, "Invalid Argument")

	// HTTPErrRequestBindingFailed is an error when failed to bind http error
	// request.
	HTTPErrRequestBindingFailed = NewWithCode(http.StatusBadRequest, 400003, "Request Binding Failed")

	// HTTPErrRequestTimeout is an error when a http error timeout.
	HTTPErrRequestTimeout = NewWithCode(http.StatusRequestTimeout, 408000, "Request Timeout")

	// HTTPErrUnsupportedMediaType is an error when a media type of http error is
	// unsupported.
	HTTPErrUnsupportedMediaType = NewWithCode(http.StatusUnsupportedMediaType, 415000, "Unsupported Media Type")

	// HTTPErrUnauthorized is a common authorization error.
	HTTPErrUnauthorized = NewWithCode(http.StatusUnauthorized, 401000, "Unauthorized")

	HTTPErrTooManyRequests = NewWithCode(http.StatusTooManyRequests, 429000, "Too many requests")

	// HTTPErrInsufficientPermission is an error which occurred when a request
	// does not have permission to access a API.
	HTTPErrInsufficientPermission = NewWithCode(http.StatusForbidden, 403001, "Insufficient Permission")

	// HTTPErrRouteNotFound is a common error not found.
	HTTPErrRouteNotFound = NewWithCode(http.StatusNotFound, 404000, "Route Not Found")

	// HTTPErrResourceNotFound is error for querying the database for a document
	// which does not exist.
	HTTPErrResourceNotFound = NewWithCode(http.StatusNotFound, 404001, "Resource Not Found")

	// HTTPErrAlreadyExists is error for inserting a document to database with
	// a duplicate id or values which are marked as unique index.
	HTTPErrAlreadyExists = NewWithCode(http.StatusConflict, 409001, "Already Exists")

	// HTTPErrSMSignatureCannotBeVerify is error for validating signature of
	// the request with SM Signature.
	HTTPErrSMSignatureCannotBeVerify = NewWithCode(http.StatusUnauthorized, 401001, "SM Signature Cannot Be Verify")

	// HTTPErrTokenInvalid is error for invalid parsing bearer token.
	HTTPErrTokenInvalid = NewWithCode(http.StatusUnauthorized, 401002, "Token Invalid")

	// ========================================================================
	// Server errors
	// ========================================================================

	// HTTPErrInternalServerError is common internal error in server.
	HTTPErrInternalServerError = NewWithCode(http.StatusInternalServerError, 500000, "Internal Server Error")

	// HTTPErrsNotImplemented is common internal error in server.
	HTTPErrsNotImplemented = NewWithCode(http.StatusNotImplemented, 501000, "Not Implemented")
)
