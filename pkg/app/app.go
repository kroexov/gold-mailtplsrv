package app

import (
	"context"
	"time"

	"mailtplsrv/pkg/renderer"

	"github.com/labstack/echo/v4"
	"github.com/vmkteam/appkit"
	"github.com/vmkteam/embedlog"
)

type Config struct {
	Server struct {
		Host    string
		Port    int
		IsDevel bool
	}
	Renderer renderer.Config
	RPC      struct {
		// CallerNamespace is the project name used in the generated Go client caller name:
		// "myproject" gives "myproject-mailtplsrv", empty value gives "mailtplsrv".
		CallerNamespace string
	}
	Sentry struct {
		Environment string
		DSN         string
	}
}

type App struct {
	embedlog.Logger
	appName  string
	cfg      Config
	echo     *echo.Echo
	renderer *renderer.Renderer
}

func New(appName string, sl embedlog.Logger, cfg Config) *App {
	a := &App{
		appName:  appName,
		cfg:      cfg,
		echo:     appkit.NewEcho(),
		Logger:   sl,
		renderer: renderer.New(cfg.Renderer, sl),
	}

	return a
}

// Run is a function that runs application.
func (a *App) Run(ctx context.Context) error {
	a.registerMetrics()
	a.registerHandlers()
	a.registerDebugHandlers()
	a.registerAPIHandlers()
	a.registerMetadata()

	return a.runHTTPServer(ctx, a.cfg.Server.Host, a.cfg.Server.Port)
}

// Shutdown is a function that gracefully stops HTTP server.
func (a *App) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return a.echo.Shutdown(ctx)
}

// registerMetadata is a function that registers meta info from service. Must be updated.
func (a *App) registerMetadata() {
	opts := appkit.MetadataOpts{
		HasPublicAPI:  true,
		HasPrivateAPI: true,
	}

	md := appkit.NewMetadataManager(opts)
	md.RegisterMetrics()

	a.echo.GET("/debug/metadata", md.Handler)
}
