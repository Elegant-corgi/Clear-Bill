package app

import (
	"context"
	"net/http"
	"time"

	"clearbill/mgr/server/api/router"
	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/config"
	"clearbill/mgr/server/internal/app/dal"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func NewServer() *Server {
	cfg := config.Load()
	systemDAL := dal.NewSystemDAL()
	billingDAL := dal.NewBillingDAL()

	systemService := bll.NewSystemService(systemDAL)
	billingService := bll.NewBillingService(billingDAL)

	systemAction := action.NewSystemAction(systemService)
	billingAction := action.NewBillingAction(billingService)

	switch cfg.RunMode {
	case gin.ReleaseMode, gin.TestMode, gin.DebugMode:
		gin.SetMode(cfg.RunMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	router.New(cfg, systemAction, billingAction).Register(engine)

	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           engine,
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
