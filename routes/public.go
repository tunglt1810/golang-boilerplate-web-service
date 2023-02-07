package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"goland-boilerplate-web-service/config"
)

type APIVersion struct {
	APIVersion    string `json:"api_version"`
	LastUpdatedAt string `json:"last_updated_at"`
	Status        string `json:"status"`
}

func HealthCheck(c echo.Context) error {
	version := config.Config.APIInfo.Version
	lastUpdated := config.Config.APIInfo.LastUpdatedAt
	return c.JSON(http.StatusOK, &APIVersion{
		APIVersion:    version,
		LastUpdatedAt: lastUpdated,
		Status:        "OK",
	})
}

func registerPublic(e *echo.Echo) {
	e.GET("/health", HealthCheck)
}
