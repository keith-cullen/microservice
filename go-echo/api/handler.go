package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/keith-cullen/microservice/store"
	"github.com/labstack/echo/v4"
)

type AppResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (handler *Handler) AppDefault(ctx echo.Context) error {
	log.Print("Default()")
	resp := &AppResponse{
		Message: "404 Not Found",
	}
	return ctx.JSON(http.StatusNotFound, resp)
}

func (handler *Handler) AppGet(ctx echo.Context, params AppGetParams) error {
	var name string
	if params.Name == nil {
		name = ""
	} else {
		name = *params.Name
	}
	log.Printf("AppGet(name: %q)", name)
	if name == "" {
		resp := &AppResponse{
			Message: "400 Bad Request",
		}
		return ctx.JSON(http.StatusBadRequest, resp)
	}
	if _, err := handler.store.GetThing(ctx.Request().Context(), name); err != nil {
		resp := &AppResponse{
			Message: "404 Not Found",
		}
		return ctx.JSON(http.StatusNotFound, resp)
	}
	resp := &AppResponse{
		Message: fmt.Sprintf("Hello, %s", name),
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (handler *Handler) AppSet(ctx echo.Context, params AppSetParams) error {
	var name string
	if params.Name == nil {
		name = ""
	} else {
		name = *params.Name
	}
	log.Printf("AppSet(name: %q)", name)
	if name == "" {
		resp := &AppResponse{
			Message: "400 Bad Request",
		}
		return ctx.JSON(http.StatusBadRequest, resp)
	}
	if err := handler.store.SetThing(ctx.Request().Context(), name); err != nil {
		resp := &AppResponse{
			Message: "500 Internal Server Error",
		}
		return ctx.JSON(http.StatusInternalServerError, resp)
	}
	resp := &AppResponse{
		Message: fmt.Sprintf("Hello, %s", name),
	}
	return ctx.JSON(http.StatusOK, resp)
}
