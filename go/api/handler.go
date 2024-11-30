package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/keith-cullen/microservice/store"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Store *store.Store
}

type AppDefaultResponse struct {
	Message string `json:"message"`
}

type AppGetResponse struct {
	Message string `json:"message"`
}

type AppSetResponse struct {
	Message string `json:"message"`
}

func NewServer(store *store.Store) *Server {
	return &Server{Store: store}
}

func (s *Server) AppDefault(ctx echo.Context) error {
	log.Print("Default()")
	resp := &AppGetResponse{
		Message: "404 Not Found",
	}
	return ctx.JSON(http.StatusNotFound, resp)
}

func (s *Server) AppGet(ctx echo.Context, params AppGetParams) error {
	var name string
	if params.Name == nil {
		name = ""
	} else {
		name = *params.Name
	}
	log.Printf("AppGet(name: %q)", name)
	if name == "" {
		resp := &AppGetResponse{
			Message: "400 Bad Request",
		}
		return ctx.JSON(http.StatusBadRequest, resp)
	}
	if _, err := s.Store.GetThing(ctx.Request().Context(), name); err != nil {
		resp := &AppGetResponse{
			Message: "404 Not Found",
		}
		return ctx.JSON(http.StatusNotFound, resp)
	}
	resp := &AppGetResponse{
		Message: fmt.Sprintf("Hello, %s", name),
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) AppSet(ctx echo.Context, params AppSetParams) error {
	var name string
	if params.Name == nil {
		name = ""
	} else {
		name = *params.Name
	}
	log.Printf("AppSet(name: %q)", name)
	if name == "" {
		resp := &AppSetResponse{
			Message: "400 Bad Request",
		}
		return ctx.JSON(http.StatusBadRequest, resp)
	}
	if err := s.Store.SetThing(ctx.Request().Context(), name); err != nil {
		resp := &AppSetResponse{
			Message: "500 Internal Server Error",
		}
		return ctx.JSON(http.StatusInternalServerError, resp)
	}
	resp := &AppSetResponse{
		Message: fmt.Sprintf("Hello, %s", name),
	}
	return ctx.JSON(http.StatusOK, resp)
}
