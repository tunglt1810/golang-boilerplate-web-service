package routes

import (
	"goland-boilerplate-web-service/config"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, cfg *config.Schema) {
	registerPublic(e)
	registerPrivate(e)
	registerInternal(e)
}
