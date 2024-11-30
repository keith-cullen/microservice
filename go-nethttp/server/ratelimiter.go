package server

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	cleanupInterval = 30 * time.Second // cleanupInterval is the time between iterations of the cleanupLoop
	maxAge          = 10 * time.Second // maxAge is the maximum time a rate.Limiter can continue to exist since it was last used
)

// rateLimiterItem combines a rate.Limiter with a last used value
type rateLimiterItem struct {
	limiter  *rate.Limiter
	lastUsed time.Time
}

// RateLimiter maintains a map of IP addresses to token bucket rate limiters
// Each token bucket has a size of burstSize, is initially full, and is refilled at a rate of reqsPerSec tokens per second
type RateLimiter struct {
	mu            sync.Mutex                  // mu synchronises access to the items map
	items         map[string]*rateLimiterItem // items is a map of IP addresses to rateLimiterItems
	reqsPerSec    int                         // reqsPerSec is the rate at which tokens are added to each token bucket
	burstSize     int                         // burstSize is the size of each token bucket
	cleanupTicker *time.Ticker                // cleanupTicker delivers events that trigger the cleanupLoop
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(reqsPerSec, burstSize int) *RateLimiter {
	r := &RateLimiter{
		items:         make(map[string]*rateLimiterItem),
		reqsPerSec:    reqsPerSec,
		burstSize:     burstSize,
		cleanupTicker: time.NewTicker(cleanupInterval),
	}
	go r.cleanupLoop()
	return r
}

// getItem finds an existing rateLimterItem for a given IP address or creates a new one
// It also updates the last used time of the rateLimterItem
// This must be called with the mu mutex already held
func (r *RateLimiter) getItem(ip string) *rateLimiterItem {
	item, exists := r.items[ip]
	if !exists {
		item = &rateLimiterItem{
			limiter: rate.NewLimiter(rate.Limit(r.reqsPerSec), r.burstSize),
		}
		r.items[ip] = item
	}
	item.lastUsed = time.Now()
	return item
}

// Allow determines if a request from the given IP address is allowed or not
func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter := r.getItem(ip).limiter
	before := int(limiter.Tokens())
	allow := limiter.Allow()
	after := int(limiter.Tokens())
	if allow {
		log.Printf("rate limiter tokens for %s reduced from %d to %d, request allowed", ip, before, after)
	} else {
		log.Printf("ratelimer tokens for %s remain at %d, request denied", ip, after)
	}
	return allow
}

// cleanupLoop periodically deletes expired rateLimiterItems
func (r *RateLimiter) cleanupLoop() {
	for range r.cleanupTicker.C { // Wait for events on the ticker channel
		r.mu.Lock()
		now := time.Now()
		for ip, item := range r.items {
			age := now.Sub(item.lastUsed)
			log.Printf("rate limiter for %s was last used %v ago", ip, age)
			if age > maxAge {
				log.Printf("rate limiter for %s has expired", ip)
				delete(r.items, ip)
			}
		}
		r.mu.Unlock()
	}
}
