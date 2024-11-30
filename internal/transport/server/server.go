package server

import (
	"context"
	"fmt"
	"musicAPI/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg config.Config, handler http.Handler) error {
	op := "server.Run"
	s.httpServer = &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 16,
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
	}
	return fmt.Errorf("%w:%s", s.httpServer.ListenAndServe(), op)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
