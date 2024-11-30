package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/keith-cullen/microservice/api"
	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

const (
	maxHeaderBytes = 4096
)

func Run(store *store.Store, insecure bool) error {
	addr := config.Get(config.AddrKey)
	corsOrigin := config.Get(config.CorsOriginKey)
	reqPerSec, err := strconv.ParseUint(config.Get(config.ReqPerSecKey), 10, 64)
	if err != nil {
		return err
	}
	apiServer := api.NewServer(store)
	echoServer := echo.New()
	echoServer.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(reqPerSec))))
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{corsOrigin},
	}))
	api.RegisterHandlers(echoServer, apiServer)
	httpServer := &http.Server{
		Handler:        echoServer,
		Addr:           addr,
		MaxHeaderBytes: maxHeaderBytes,
	}
	log.Printf("listening on %s", httpServer.Addr)
	if !insecure {
		return httpServer.ListenAndServeTLS(config.Get(config.CertKey), config.Get(config.PrivkeyKey))
	} else {
		return httpServer.ListenAndServe()
	}
}
