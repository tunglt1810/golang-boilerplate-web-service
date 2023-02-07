package routes

import (
	"github.com/labstack/echo/v4"
	"goland-boilerplate-web-service/config"
)

func SetupRoutes(e *echo.Echo, cfg *config.Schema) {
	registerPublic(e)
	registerPrivate(e)
	registerInternal(e)
}
