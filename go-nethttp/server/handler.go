package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store"
)

type AppResponse struct {
	Message string `json:"message"`
}

type RateLimiter struct {
	mu          sync.Mutex
	maxRequests uint64
	window      time.Duration
	requests    uint64
	lastRequest time.Time
}

type Handler struct {
	store       *store.Store
	rateLimiter *RateLimiter
}

func NewRateLimiter(maxRequests uint64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    0,
		lastRequest: time.Now(),
	}
}

// Fixed window rate limiting
// Time is divided into fixed intervals (windows)
// A counter keeps track of the number of requests within each window
// If the counter exceeds a limit, subsequenct requests within that window are rejected until the window resets
func (r *RateLimiter) Check() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	ok := false
	if time.Since(r.lastRequest) > r.window {
		r.requests = 0
		r.lastRequest = time.Now()
	}
	if r.requests < r.maxRequests {
		r.requests += 1
		ok = true
	}
	return ok
}

func NewHandler(store *store.Store) (Handler, error) {
	handler := Handler{}
	reqPerSec, err := strconv.ParseUint(config.Get(config.ReqPerSecKey), 10, 64)
	if err != nil {
		return handler, err
	}
	rateLimiter := NewRateLimiter(reqPerSec, time.Second)
	handler.store = store
	handler.rateLimiter = rateLimiter
	return handler, nil
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
	http.Error(w, string(data), http.StatusBadRequest)
}

func (handler Handler) tooManyRequests(w http.ResponseWriter) {
	resp := &AppResponse{
		Message: http.StatusText(http.StatusTooManyRequests),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	http.Error(w, string(data), http.StatusTooManyRequests)
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
	http.Error(w, string(data), http.StatusNotFound)
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
	http.Error(w, string(data), http.StatusInternalServerError)
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
	io.WriteString(w, string(data))
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
	io.WriteString(w, string(data))
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
		if !handler.rateLimiter.Check() {
			log.Printf("rate limit exceeded (%v requests in %v)",
				handler.rateLimiter.requests, time.Since(handler.rateLimiter.lastRequest))
			handler.tooManyRequests(w)
			return
		}
		log.Print("rate limit middleware calling next handler")
		next.ServeHTTP(w, r)
	})
}
