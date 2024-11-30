package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
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
	handler, err := NewHandler(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/", handler.AppDefault)
	router.HandleFunc("/v1/get", handler.AppGet).Methods("GET")
	router.HandleFunc("/v1/set", handler.AppSet).Methods("POST")
	router.Use(handler.CorsMiddle)
	router.Use(handler.RateLimitMiddle)
	return &Server{
		httpServer: http.Server{
			Addr:           addr,
			Handler:        router,
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
