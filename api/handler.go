package api

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
	"github.com/labstack/echo/v4"
)

var (
	staticHtmlStr   = ""
	staticStyleStr  = ""
	staticScriptStr = ""
)

type Server struct {
	Store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{Store: store}
}

func (s *Server) Static(ctx echo.Context) error {
	urlPath := ctx.Request().URL.Path
	if urlPath == "/" {
		urlPath = config.Get(config.StaticHtmlFileKey)
	}
	filePath := filepath.Join(config.Get(config.StaticDirKey), urlPath)
	log.Printf("Static(url path: %q, file path: %q)", urlPath, filePath)
	return ctx.File(filePath)
}

func (s *Server) AppGet(ctx echo.Context, params AppGetParams) error {
	if params.Name == nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	name := *params.Name
	log.Printf("AppGet(name: %q), headers: %v", name, ctx.Request().Header)
	_, err := s.Store.GetThing(ctx.Request().Context(), name)
	if err != nil {
		return err
	}
	str := fmt.Sprintf("hello %s", name)
	return ctx.JSON(http.StatusOK, str)
}

func (s *Server) AppSet(ctx echo.Context, params AppSetParams) error {
	if params.Name == nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	name := *params.Name
	log.Printf("AppSet(name: %q), headers: %v", name, ctx.Request().Header)
	err := s.Store.SetThing(ctx.Request().Context(), name)
	if err != nil {
		return err
	}
	str := fmt.Sprintf("welcome %s", name)
	return ctx.JSON(http.StatusOK, str)
}
