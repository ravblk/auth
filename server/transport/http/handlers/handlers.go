package handlers

import (
	"auth/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	ErrUserAgent        = errors.New("User Agent missing")
	MIMEApplicationJSON = "application/json"
)

type Error struct {
	Message string
}

type Handlers struct {
	s services.Auth
}

func New(svc services.Auth) *Handlers {
	return &Handlers{
		s: svc,
	}
}

func IPGet(ctx *fasthttp.RequestCtx) string {
	return ctx.RemoteIP().String()
}
func UAGet(ctx *fasthttp.RequestCtx) (string, error) {
	if ctx.Request.Header.Peek("User Agent") == nil {
		return "", ErrUserAgent
	}
	return string(ctx.Request.Header.Peek("User Agent")), nil
}

func (h *Handlers) responseError(ctx *fasthttp.RequestCtx, err error) {
	e := &Error{err.Error()}
	buf, err := json.Marshal(e)
	if err != nil {
		h.s.Log.Error("", zap.Error(err))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusBadRequest)
	ctx.SetContentType(MIMEApplicationJSON)
	ctx.SetBody(buf)
}

func (h *Handlers) responseJSON(ctx *fasthttp.RequestCtx, res interface{}) {
	buf, err := json.Marshal(res)
	if err != nil {
		h.s.Log.Error("", zap.Error(err))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType(MIMEApplicationJSON)
	ctx.SetBody(buf)
}
