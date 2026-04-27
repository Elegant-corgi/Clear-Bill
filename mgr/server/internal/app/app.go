package app

import (
	"context"
	"net/http"
	"time"

	"clearbill/mgr/server/internal/app/config"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func NewServer() *Server {
	cfg := config.Load()
	injector, err := BuildInjector(cfg)
	if err != nil {
		panic(err)
	}

	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           injector.Engine,
			ReadHeaderTimeout: time.Duration(httpServeTimeout(cfg)) * time.Second,
		},
		shutdownTimeout: time.Duration(httpShutdownTimeout(cfg)) * time.Second,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		_ = s.httpServer.Shutdown(shutdownCtx)
	}()

	return s.httpServer.ListenAndServe()
}

func httpServeTimeout(cfg config.Config) int {
	if cfg.HTTP.ServeTimeout > 0 {
		return cfg.HTTP.ServeTimeout
	}

	return 5
}

func httpShutdownTimeout(cfg config.Config) int {
	if cfg.HTTP.ShutdownTimeout > 0 {
		return cfg.HTTP.ShutdownTimeout
	}

	return 5
}
