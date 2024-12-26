package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"goland-boilerplate-web-service/config"
	inMiddleware "goland-boilerplate-web-service/middleware"
	"goland-boilerplate-web-service/pkg/validation"
	"goland-boilerplate-web-service/routes"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Init()

	e := echo.New()
	e.HideBanner = true
	// Integrate common middleware
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	// Integrate custom middleware
	e.Use(inMiddleware.RequestID())
	// Integrate common middleware

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Config.Sentry.DSN,
		Environment:      config.Config.Sentry.ENV,
		TracesSampleRate: config.Config.Sentry.TracesSampleRate,
		Debug:            config.Config.Sentry.Debug,
	}); err != nil {
		log.Info().Caller().Msg("Sentry initialization failed")
	}

	e.Use(sentryecho.New(sentryecho.Options{}))

	v, err := validation.New()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("error when init validator")
	}
	e.Validator = v
	e.HTTPErrorHandler = inMiddleware.GlobalErrorHandler

	// Register routes
	routes.SetupRoutes(e, &config.Config)

	eMetrics, prom := setupMetrics(&config.Config.Server)
	e.Use(prom.HandlerFunc)

	// Start webserver with graceful shutdown
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", config.Config.Server.Port)); err != nil {
			log.Info().Caller().Err(err).Msg("shutting down")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	if err := eMetrics.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func setupMetrics(cfg *config.Server) (*echo.Echo, *prometheus.Prometheus) {
	e := echo.New()
	e.HideBanner = true
	p := prometheus.NewPrometheus("echo", nil)
	p.SetMetricsPath(e)
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Metrics)); err != nil {
			log.Info().Caller().Err(err).Send()
			log.Info().Caller().Msg("shutting down the prometheus server")
		}
	}()

	return e, p
}
