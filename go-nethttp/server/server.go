package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
)

const (
	maxHeaderBytes = 4096
)

type Server struct {
	httpServer http.Server
}

func New(store *store.Store) *Server {
	addr := config.Get(config.AddrKey)
	handler := NewHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.AppDefault)
	mux.HandleFunc("/v1/get", handler.AppGet)
	mux.HandleFunc("/v1/set", handler.AppSet)
	return &Server{
		httpServer: http.Server{
			Addr:           addr,
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: maxHeaderBytes,
		},
	}
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

func (server *Server) Stop() {
	if err := server.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown HTTP server: %v", err)
	} else {
		log.Print("http server stopped")
	}
}
