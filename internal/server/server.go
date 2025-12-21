package server

import (
	"context"
	"net/http"
	"time"

	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/handler"
	"github.com/dmnAlex/gophermart/internal/logger"
	"go.uber.org/zap"
)

type Server struct {
	cfg *config.Config
	srv *http.Server
}

func NewServer(handler *handler.Handler, cfg *config.Config) *Server {
	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: newRouter(handler, cfg),
	}

	return &Server{
		cfg: cfg,
		srv: srv,
	}
}

func (s *Server) Run() {
	logger.Log.Info("starting server", zap.String("address", s.cfg.RunAddress))
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("failed to start server", zap.Error(err))
	}
}

func (s *Server) Shutdown() {
	logger.Log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("server shutdown", zap.Error(err))
	}

	logger.Log.Info("server shut down gracefully")
}
