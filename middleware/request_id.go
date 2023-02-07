package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type (
	RequestIDConfig struct {
		Skipper   Skipper
		Generator func() string
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:   DefaultSkipper,
		Generator: generator,
	}
)

func generator() string {
	return random.String(32) //nolint:gomnd
}

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			req, res := c.Request(), c.Response()
			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = config.Generator()
			}
			res.Header().Set(echo.HeaderXRequestID, requestID)
			c.Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	}
}
