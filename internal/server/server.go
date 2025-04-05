package server

import (
	"context"
	"docker-compose-training/config"
	"github.com/valyala/fasthttp"
)

type Server struct {
	httpServer *fasthttp.Server
}

// Run the server.
func (s *Server) Run(handler fasthttp.RequestHandler, cfg config.Server) error {
	s.httpServer = &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s.httpServer.ListenAndServe(cfg.Address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.ShutdownWithContext(ctx)
}
