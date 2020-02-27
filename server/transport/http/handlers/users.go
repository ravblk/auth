package handlers

import (
	"auth/model"
	"context"
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (h *Handlers) Registration(ctx *fasthttp.RequestCtx) {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(h.s.Cfg.API.TTL)*time.Second)
	defer cancel()
	req := &model.User{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		h.s.Log.Warn("", zap.Error(err))
		h.responseError(ctx, err)
		return
	}
	s := &model.Session{}
	s.IP = IPGet(ctx)

	ua, err := UAGet(ctx)
	if err != nil {
		h.s.Log.Warn("", zap.Error(err))
		h.responseError(ctx, err)
		return
	}
	s.UserAgent = ua
	if err := h.s.UsrSvc.UserRegistration(c, req, s); err != nil {
		h.s.Log.Warn("", zap.Error(err))
		h.responseError(ctx, err)
		return
	}
	h.responseJSON(ctx, &s.SessionClient)
}

func (h *Handlers) Session(ctx *fasthttp.RequestCtx) {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(h.s.Cfg.API.TTL)*time.Second)
	defer cancel()
	req := &model.Login{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		h.s.Log.Warn("", zap.Error(err))
		h.responseError(ctx, err)
		return
	}
	s := &model.Session{}
	s.IP = IPGet(ctx)
	ua, err := UAGet(ctx)
	if err != nil {
		h.s.Log.Warn("", zap.Error(err))
		h.responseError(ctx, err)
		return
	}
	s.UserAgent = ua
	if err := h.s.UsrSvc.UserLogin(c, req, s); err != nil {
		h.s.Log.Warn("", zap.Error(err))
		h.responseError(ctx, err)
		return
	}
	h.responseJSON(ctx, &s.SessionClient)
}
