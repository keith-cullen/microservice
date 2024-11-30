package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/keith-cullen/microservice/api"
	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

const (
	timeout        = 10 * time.Second
	maxHeaderBytes = 4096
)

type Server struct {
	httpServer http.Server
}

func New(store *store.Store) (*Server, error) {
	addr := config.Get(config.AddrKey)
	corsOrigin := config.Get(config.CorsOriginKey)
	reqPerSec, err := strconv.ParseUint(config.Get(config.ReqPerSecKey), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	handler := api.NewHandler(store)
	echoServer := echo.New()
	echoServer.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(reqPerSec))))
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{corsOrigin},
	}))
	api.RegisterHandlers(echoServer, handler)
	return &Server{
		httpServer: http.Server{
			Addr:           addr,
			Handler:        echoServer,
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: maxHeaderBytes,
		},
	}, nil
}

func (server *Server) Start(insecure bool) error {
	log.Printf("http server listening on %s", server.httpServer.Addr)
	var err error
	if !insecure {
		err = server.httpServer.ListenAndServeTLS(config.Get(config.CertKey), config.Get(config.PrivkeyKey))
	} else {
		err = server.httpServer.ListenAndServe()
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (server *Server) Stop() error {
	if err := server.httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	log.Print("http server stopped")
	return nil
}
