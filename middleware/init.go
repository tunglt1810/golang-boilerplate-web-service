package middleware

import "github.com/labstack/echo/v4"

type (
	Skipper    func(echo.Context) bool
	BeforeFunc func(echo.Context)
)

func DefaultSkipper(echo.Context) bool {
	return false
}
