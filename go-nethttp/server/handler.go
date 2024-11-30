package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/keith-cullen/microservice/store"
)

type AppResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) Handler {
	return Handler{store: store}
}

func (handler Handler) badRequest(w http.ResponseWriter) {
	resp := &AppResponse{
		Message: http.StatusText(http.StatusBadRequest),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	text := string(data)
	http.Error(w, text, http.StatusBadRequest)
}

func (handler Handler) notFound(w http.ResponseWriter) {
	resp := &AppResponse{
		Message: http.StatusText(http.StatusNotFound),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	text := string(data)
	http.Error(w, text, http.StatusNotFound)
}

func (handler Handler) internalServerError(w http.ResponseWriter) {
	resp := &AppResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	text := string(data)
	http.Error(w, text, http.StatusInternalServerError)
}

func (handler Handler) AppDefault(w http.ResponseWriter, r *http.Request) {
	log.Print("Default()")
	handler.notFound(w)
}

func (handler Handler) AppGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handler.badRequest(w)
		return
	}
	name := r.URL.Query().Get("name")
	log.Printf("AppGet(%s)", name)
	if _, err := handler.store.GetThing(context.Background(), name); err != nil {
		handler.notFound(w)
		return
	}
	resp := &AppResponse{
		Message: fmt.Sprintf("Hello, %s", name),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error: %v", err)
		handler.internalServerError(w)
		return
	}
	text := string(data)
	io.WriteString(w, text)
}

func (handler Handler) AppSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.badRequest(w)
		return
	}
	name := r.URL.Query().Get("name")
	log.Printf("AppSet(%s)", name)
	if err := handler.store.SetThing(context.Background(), name); err != nil {
		handler.internalServerError(w)
		return
	}
	resp := &AppResponse{
		Message: fmt.Sprintf("Hello, %s", name),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error: %v", err)
		handler.internalServerError(w)
		return
	}
	text := string(data)
	io.WriteString(w, text)
}
