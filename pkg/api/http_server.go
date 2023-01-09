package api

import (
	"context"
	"net/http"
	"time"

	"github.com/magicLian/gostarter/pkg/middleware"
	"github.com/magicLian/gostarter/pkg/services/health"
	"github.com/magicLian/gostarter/pkg/services/s3"
	"github.com/magicLian/gostarter/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/magicLian/gostarter/pkg/log"
)

type HTTPServer struct {
	log       log.Logger
	context   context.Context
	ginEngine *gin.Engine

	Cfg        *setting.Cfg
	middleware *middleware.MiddleWare
	health     health.Health
	s3Svc      s3.S3Service
}

func ProvideHttpServer(Cfg *setting.Cfg, middleware *middleware.MiddleWare, health health.Health, s3Svc s3.S3Service) *HTTPServer {
	httpServer := &HTTPServer{
		log:        log.New("http.server"),
		Cfg:        Cfg,
		middleware: middleware,
		ginEngine:  gin.Default(),
		health:     health,
		s3Svc:      s3Svc,
	}
	return httpServer
}

func (hs *HTTPServer) Run(ctx context.Context) error {
	hs.context = ctx
	hs.ginEngine.Use(hs.healthHandler)
	hs.ginEngine.Use(hs.middleware.InitContextHandler())
	hs.ginEngine.Use(middleware.CrossRequestCheck())
	hs.apiRegister()
	httpSrv := &http.Server{
		Addr:    ":" + setting.HttpPort,
		Handler: hs.ginEngine,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		hs.log.Debugf("strating http server", "port", setting.HttpPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hs.log.Errorf("Could not start listener", err.Error())
		}
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			time.Sleep(1000)
		}
	}
	if err := httpSrv.Shutdown(ctx); err != nil {
		hs.log.Errorf("Server forced to shutdown:", err)
	}
	return nil
}
