package middleware

import (
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func LogRequest(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		h(ctx)
		addr := ctx.RemoteIP().String()
		zap.L().Debug("http request", zap.String("method", string(ctx.Method())), zap.String("path", string(ctx.RequestURI())), zap.String("RemoteIP", addr), zap.Duration("duration", time.Since(start)))
	}
}
