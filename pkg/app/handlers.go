package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"mailtplsrv"
	"mailtplsrv/pkg/renderer"
	"mailtplsrv/pkg/rpc"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vmkteam/appkit"
	"github.com/vmkteam/rpcgen/v2"
	"github.com/vmkteam/rpcgen/v2/golang"
	"github.com/vmkteam/zenrpc/v2"
)

// runHTTPServer is a function that starts http listener using labstack/echo.
func (a *App) runHTTPServer(ctx context.Context, host string, port int) error {
	listenAddress := fmt.Sprintf("%s:%d", host, port)
	addr := "http://" + listenAddress
	a.Print(ctx, "starting http listener", "url", addr, "templates", addr+"/template/", "smdbox", addr+"/int/rpc/doc/")

	return a.echo.Start(listenAddress)
}

// registerHandlers register echo handlers.
func (a *App) registerHandlers() {
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{"Authorization", "Authorization2", "Origin", "X-Requested-With", "Content-Type", "Accept", "Platform", "Version"},
	}))

	group := a.echo.Group("")
	group.Use(a.requestLogger())

	// renderer templates
	if a.cfg.Renderer.IsEmbed {
		group.StaticFS(a.cfg.Renderer.WebPath, echo.MustSubFS(mailtplsrv.Static, renderer.StaticDir))
	} else {
		group.Static(a.cfg.Renderer.WebPath, renderer.StaticDir)
	}

	// test endpoints
	group.GET("/template/test/:path", a.renderer.EchoTestTemplate)
	group.GET("/template/", a.renderer.EchoList)

	// REST render endpoint. Use /int/rpc/ unless a plain HTTP client is required
	group.POST("/template/:name", a.renderer.EchoRenderTemplate)
}

// registerDebugHandlers adds /debug/pprof handlers into a.echo instance.
func (a *App) registerDebugHandlers() {
	dbg := a.echo.Group("/debug")

	// add pprof integration
	dbg.Any("/pprof/*", appkit.PprofHandler)

	// add healthcheck
	a.echo.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// show all routes in devel mode
	if a.cfg.Server.IsDevel {
		a.echo.GET("/", appkit.RenderRoutes(a.appName, a.echo))
	}
}

// registerAPIHandlers registers main rpc server.
func (a *App) registerAPIHandlers() {
	srv := rpc.New(a.Logger, a.cfg.Server.IsDevel, a.renderer)
	gen := rpcgen.FromSMD(srv.SMD())

	group := a.echo.Group("")
	group.Use(a.requestLogger())

	group.Any("/int/rpc/", appkit.EchoHandler(appkit.XRequestID(srv)))
	group.Any("/int/rpc/doc/", appkit.EchoHandlerFunc(zenrpc.SMDBoxHandler))
	goClient := gen.GoClient(golang.Settings{Package: a.appName, CallerNamespace: a.cfg.RPC.CallerNamespace})
	group.Any("/int/rpc/api.go", appkit.EchoHandlerFunc(rpcgen.Handler(goClient)))
}

func (a *App) requestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		HandleError:  true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("ip", v.RemoteIP),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("userAgent", v.UserAgent),
				slog.String("duration", v.Latency.String()),
				slog.String("xRequestId", v.RequestID),
			}

			if v.Error == nil {
				a.Log().LogAttrs(context.Background(), slog.LevelInfo, "http request", attrs...)
			} else {
				a.Log().LogAttrs(context.Background(), slog.LevelError, "http request error", append(attrs, slog.String("err", v.Error.Error()))...)
			}
			return nil
		},
	})
}
