package app

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vmkteam/appkit"
)

// registerMetrics is a function that initializes a.stat* variables and adds /metrics endpoint to echo.
func (a *App) registerMetrics() {
	// add db conn metrics
	a.echo.Use(appkit.HTTPMetrics(appkit.DefaultServerName))
	a.echo.Any("/metrics", echo.WrapHandler(promhttp.Handler()))
}
