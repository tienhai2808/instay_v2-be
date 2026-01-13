package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/router"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg  *config.Config
	http *http.Server
}

func NewServer(cfg *config.Config, ctn *container.Container) *Server {
	r := gin.New()
	_ = r.SetTrustedProxies([]string{"0.0.0.0/0"})

	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.AllowOrigins,
		AllowMethods:     cfg.Server.AllowMethods,
		AllowHeaders:     cfg.Server.AllowHeaders,
		ExposeHeaders:    cfg.Server.ExposeHeaders,
		AllowCredentials: cfg.Server.AllowCredentials,
		MaxAge:           cfg.Server.MaxAge,
	}

	r.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.New(corsConfig),
		ctn.CtxHTTPMid.ErrorHandler(),
		ctn.CtxHTTPMid.Recovery(),
	)

	api := router.NewRouter(r)
	api.Setup(cfg.Server, ctn)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	http := &http.Server{
		Addr:           addr,
		Handler:        r,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes * 1024 * 1024,
		IdleTimeout:    cfg.Server.IdleTimeout,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
	}

	return &Server{
		cfg,
		http,
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown failed: %v", err)
			return
		}
	}

	log.Println("Server stopped successfully")
}

func (s *Server) GracefulShutdown(ch <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		log.Printf("Server run failed: %v", err)
	case <-ctx.Done():
		log.Println("Server stop signal")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(shutdownCtx)
}
