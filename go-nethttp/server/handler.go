package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
)

type AppResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	store       *store.Store
	rateLimiter *RateLimiter
}

func NewHandler(store *store.Store) (Handler, error) {
	handler := Handler{}
	reqPerSec, err := strconv.ParseUint(config.Get(config.ReqPerSecKey), 10, 32)
	if err != nil {
		return handler, err
	}
	burstSize, err := strconv.ParseUint(config.Get(config.BurstSizeKey), 10, 32)
	if err != nil {
		return handler, err
	}
	rateLimiter := NewRateLimiter(int(reqPerSec), int(burstSize))
	handler.store = store
	handler.rateLimiter = rateLimiter
	return handler, nil
}

// Send an error response with a JSON-encoded body
// If the JSON-encoding operation fails, then send a plain text body
func respondError(w http.ResponseWriter, status int) {
	msg := http.StatusText(status)
	resp := &AppResponse{
		Message: msg,
	}
	if data, err := json.Marshal(&resp); err == nil {
		msg = string(data)
	}
	http.Error(w, msg, status)
}

// Send an OK response with a JSON-encoded body
// If the JSON-encoding operation fails, then send a plain text body
func respondOk(w http.ResponseWriter, msg string) {
	resp := &AppResponse{
		Message: msg,
	}
	if data, err := json.Marshal(&resp); err == nil {
		msg = string(data)
	}
	io.WriteString(w, msg)
}

func (handler Handler) AppDefault(w http.ResponseWriter, r *http.Request) {
	log.Print("Default()")
	respondError(w, http.StatusNotFound)
}

func (handler Handler) AppGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	log.Printf("AppGet(%s)", name)
	if _, err := handler.store.GetThing(context.Background(), name); err != nil {
		respondError(w, http.StatusNotFound)
		return
	}
	msg := fmt.Sprintf("Hello, %s", name)
	respondOk(w, msg)
}

func (handler Handler) AppSet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	log.Printf("AppSet(%s)", name)
	if err := handler.store.SetThing(context.Background(), name); err != nil {
		respondError(w, http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("Hello, %s", name)
	respondOk(w, msg)
}

func (handler Handler) CorsMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsOrigin := config.Get(config.CorsOriginKey)
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		// if this is a preflight options request then write an empty ok response and return
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(nil)
			return
		}
		log.Print("cors middleware calling next handler")
		next.ServeHTTP(w, r)
	})
}

func (handler Handler) RateLimitMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("unable to determine IP")
			return
		}
		if !handler.rateLimiter.Allow(ip) {
			respondError(w, http.StatusTooManyRequests)
			return
		}
		log.Print("rate limit middleware calling next handler")
		next.ServeHTTP(w, r)
	})
}
