package rpc

import (
	"net/http"

	"mailtplsrv/pkg/renderer"

	"github.com/vmkteam/embedlog"
	zm "github.com/vmkteam/zenrpc-middleware"
	"github.com/vmkteam/zenrpc/v2"
)

var (
	ErrNotImplemented = zenrpc.NewStringError(http.StatusInternalServerError, "not implemented")
	ErrInternal       = zenrpc.NewStringError(http.StatusInternalServerError, "internal error")
)

var allowDebugFn = func() zm.AllowDebugFunc {
	return func(req *http.Request) bool {
		return req != nil && req.FormValue("__level") == "5"
	}
}

//go:generate go tool zenrpc

// New returns new zenrpc Server.
func New(logger embedlog.Logger, isDevel bool, r *renderer.Renderer) *zenrpc.Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
	})

	rpc.Use(
		zm.WithDevel(isDevel),
		zm.WithHeaders(),
		zm.WithSentry(zm.DefaultServerName),
		zm.WithNoCancelContext(),
		zm.WithMetrics(zm.DefaultServerName),
		zm.WithTiming(isDevel, allowDebugFn()),
		zm.WithSLog(logger.Print, zm.DefaultServerName, nil),
		zm.WithErrorSLog(logger.Print, zm.DefaultServerName, nil),
	)

	// services
	rpc.RegisterAll(map[string]zenrpc.Invoker{
		"render": NewRenderService(r),
	})

	return rpc
}

//nolint:unused
func newInternalError(err error) *zenrpc.Error {
	return zenrpc.NewError(http.StatusInternalServerError, err)
}
