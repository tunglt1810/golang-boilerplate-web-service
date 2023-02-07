package middleware

import (
	_errors "errors"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"goland-boilerplate-web-service/pkg/errors"
	"goland-boilerplate-web-service/pkg/validation"
)

func GlobalErrorHandler(err error, c echo.Context) {
	code := errors.HTTPErrInternalServerError.HTTPStatus
	var body = errors.NewWithCode(code, code, err.Error())
	body.Status = err.Error()
	body.HTTPStatus = errors.HTTPErrBadRequest.HTTPStatus

	switch he := err.(type) {
	case *echo.HTTPError:
		body.Code = he.Code
		body.Details = he.Message
		body.Status = he.Error()
		body.HTTPStatus = he.Code
		if _errors.Is(he.Internal, middleware.ErrJWTMissing) || _errors.Is(he.Internal, middleware.ErrJWTInvalid) {
			body.HTTPStatus = http.StatusUnauthorized
		}
	case *echo.BindingError:
		body.Code = errors.HTTPErrInvalidArgument.Code
		body.Message = "Invalid data payload format"
		body.Details = he.Message
		body.Status = he.Error()
		body.HTTPStatus = he.Code
	case errors.Error:
		if herr := errors.HTTPErrFromMapping(he); herr != nil {
			body = herr
		} else {
			body.Code = errors.HTTPErrInternalServerError.Code
			body.Details = he.Details()
			body.Message = he.Msg()
			body.Status = he.Error()
		}
	case *jwt.ValidationError:
		body.HTTPStatus = http.StatusUnauthorized
		body.Code = http.StatusUnauthorized
		body.Message = he.Error()
		body.Status = "ERR_INVALID_AUTHORIZATION_CODE"
	case *errors.HTTPError:
		body = he
	case *pgconn.PgError:
		if he.Code == "23505" {
			body = errors.HTTPErrAlreadyExists
		}
	case validation.Errors:
		body = &errors.HTTPError{
			HTTPStatus: http.StatusBadRequest,
			Code:       errors.HTTPErrInvalidArgument.Code,
			Status:     errors.ErrInvalidArgument.Error(),
			Message:    he.Error(),
			Details:    he,
		}
	default:
		body.Message = err.Error()
		c.Logger().Error(errors.Wrap(err).StackTrace())
	}

	if err := c.JSON(body.HTTPStatus, body); err != nil {
		c.Logger().Error(err)
	}
}
