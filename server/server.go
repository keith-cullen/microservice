package server

import (
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

type Server struct {
	tlsEn      bool
	store      *store.Store
	apiServer  *api.Server
	echoServer *echo.Echo
	httpServer *http.Server
}

func Open(tlsEn bool) (*Server, error) {
	var addr string
	if tlsEn {
		addr = config.Get(config.HttpsAddrKey)
	} else {
		addr = config.Get(config.HttpAddrKey)
	}
	timeout, err := strconv.ParseUint(config.Get(config.HttpTimeoutKey), 10, 64)
	if err != nil {
		return nil, err
	}
	maxHeaderBytes, err := strconv.ParseUint(config.Get(config.HttpMaxHeaderBytesKey), 10, 64)
	if err != nil {
		return nil, err
	}
	reqsPerSec, err := strconv.ParseUint(config.Get(config.HttpReqsPerSecKey), 10, 64)
	if err != nil {
		return nil, err
	}
	store, err := store.Open()
	if err != nil {
		return nil, err
	}
	apiServer := api.NewServer(store)
	echoServer := echo.New()
	echoServer.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(reqsPerSec))))
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.Get((config.HttpCorsOriginKey))},
	}))
	api.RegisterHandlers(echoServer, apiServer)
	httpServer := &http.Server{
		Handler:        echoServer,
		Addr:           addr,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: int(maxHeaderBytes),
	}
	return &Server{
		tlsEn:      tlsEn,
		store:      store,
		apiServer:  apiServer,
		echoServer: echoServer,
		httpServer: httpServer,
	}, nil
}

func (s *Server) Close() {
	s.store.Close()
}

func (s *Server) Run() error {
	if s.tlsEn {
		return s.httpServer.ListenAndServeTLS(config.Get(config.CertKey), config.Get(config.PrivkeyKey))
	} else {
		return s.httpServer.ListenAndServe()
	}
}

func (s *Server) Address() string {
	return s.httpServer.Addr
}
