package http

import (
	"auth/server/transport/http/handlers"
	"auth/server/transport/http/middleware"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	Root = "/api/v1"
)

type Server struct {
	HTTPServer *fasthttp.Server
	router     *fasthttprouter.Router
	Handlers   *handlers.Handlers
}

func NewServer(hs *handlers.Handlers, rs int) (*Server, error) {

	s := &Server{
		router:   &fasthttprouter.Router{},
		Handlers: hs,
	}
	h := s.router.Handler
	s.HTTPServer = newHTTPServer(h, rs)
	return s, nil
}

func newHTTPServer(h fasthttp.RequestHandler, rs int) *fasthttp.Server {
	return &fasthttp.Server{
		Handler:            h,
		MaxRequestBodySize: rs,
	}

}

func (s *Server) Run(port string) error {
	return s.HTTPServer.ListenAndServe(":" + port)
}

func (s *Server) RoutesInit() {
	s.router.POST(Root+"/users", middleware.LogRequest(s.Handlers.Registration))
	s.router.POST(Root+"/sessions", middleware.LogRequest(s.Handlers.Session))
}
