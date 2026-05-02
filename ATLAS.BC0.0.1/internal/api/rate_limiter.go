package api

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a per-IP token bucket rate limiter.
type RateLimiter struct {
	mu       sync.RWMutex
	visitors map[string]*visitor
	rate     int           // tokens added per interval
	burst    int           // max tokens (bucket size)
	interval time.Duration // how often tokens refill
}

type visitor struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter.
// rate: requests allowed per interval. burst: max burst size.
func NewRateLimiter(rate int, burst int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		interval: interval,
	}
	// Clean up stale entries every 3 minutes
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{tokens: rl.burst, lastSeen: time.Now()}
		rl.visitors[ip] = v
		return v
	}

	// Refill tokens based on elapsed time
	elapsed := time.Since(v.lastSeen)
	refill := int(elapsed / rl.interval) * rl.rate
	if refill > 0 {
		v.tokens += refill
		if v.tokens > rl.burst {
			v.tokens = rl.burst
		}
		v.lastSeen = time.Now()
	}

	return v
}

// Allow checks if a request from the given IP should be allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	v := rl.getVisitor(ip)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if v.tokens > 0 {
		v.tokens--
		return true
	}
	return false
}

// RateLimitMiddleware wraps an http.HandlerFunc with rate limiting.
func RateLimitMiddleware(limiter *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract client IP (supports reverse proxy X-Forwarded-For)
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		if !limiter.Allow(ip) {
			http.Error(w, `{"error":"rate limit exceeded","message":"Too many requests. Please try again later."}`, http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
