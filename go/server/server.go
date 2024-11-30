package server

import (
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

func Run(store *store.Store, secure bool) error {
	var addr, cors string
	if secure {
		addr = config.Get(config.HttpsAddrKey)
		cors = config.Get(config.HttpsCorsOriginKey)
	} else {
		addr = config.Get(config.HttpAddrKey)
		cors = config.Get(config.HttpCorsOriginKey)
	}
	timeout, err := strconv.ParseUint(config.Get(config.HttpTimeoutKey), 10, 64)
	if err != nil {
		return err
	}
	maxHeaderBytes, err := strconv.ParseUint(config.Get(config.HttpMaxHeaderBytesKey), 10, 64)
	if err != nil {
		return err
	}
	reqsPerSec, err := strconv.ParseUint(config.Get(config.HttpReqsPerSecKey), 10, 64)
	if err != nil {
		return err
	}
	apiServer := api.NewServer(store)
	echoServer := echo.New()
	echoServer.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(reqsPerSec))))
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cors},
	}))
	api.RegisterHandlers(echoServer, apiServer)
	httpServer := &http.Server{
		Handler:        echoServer,
		Addr:           addr,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: int(maxHeaderBytes),
	}
	log.Printf("listening on %s", httpServer.Addr)
	if secure {
		return httpServer.ListenAndServeTLS(config.Get(config.CertKey), config.Get(config.PrivkeyKey))
	} else {
		return httpServer.ListenAndServe()
	}
}
