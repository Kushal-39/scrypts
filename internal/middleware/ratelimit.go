package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	requests map[string]*bucket
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type bucket struct {
	count      int
	resetTime  time.Time
	lastAccess time.Time
}

// NewRateLimiter creates a rate limiter with specified requests per window
func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string]*bucket),
		limit:    limit,
		window:   window,
	}

	// Cleanup goroutine to remove old entries
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, b := range rl.requests {
		if now.Sub(b.lastAccess) > rl.window*2 {
			delete(rl.requests, key)
		}
	}
}

func (rl *rateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.requests[key]

	if !exists || now.After(b.resetTime) {
		rl.requests[key] = &bucket{
			count:      1,
			resetTime:  now.Add(rl.window),
			lastAccess: now,
		}
		return true
	}

	b.lastAccess = now

	if b.count >= rl.limit {
		return false
	}

	b.count++
	return true
}

// RateLimit middleware limits requests per IP
func (rl *rateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP (strip port)
		ip := r.RemoteAddr

		// Check X-Forwarded-For header first (for proxied requests)
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			// Take the first IP if multiple are present
			ip = strings.Split(forwarded, ",")[0]
			ip = strings.TrimSpace(ip)
		} else {
			// Strip port from RemoteAddr
			host, _, err := net.SplitHostPort(ip)
			if err == nil {
				ip = host
			}
		}

		if !rl.Allow(ip) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
