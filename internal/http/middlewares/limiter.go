package middlewares

import (
	"go-api/internal"

	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// Singleton limiter object reference
var limiter *IPRateLimiter
var once sync.Once

// IPRateLimiter represents an IP rate limiter.
type IPRateLimiter struct {
	ips     map[string]*rate.Limiter
	mu      *sync.RWMutex
	limiter *rate.Limiter
}

// GetIPRateLimiter gives singleton instance of IPRateLimiter with the given rate limit.
func GetIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	once.Do(func() {
		limiter = &IPRateLimiter{
			ips:     make(map[string]*rate.Limiter),
			mu:      &sync.RWMutex{},
			limiter: rate.NewLimiter(r, burst),
		}
	})
	return limiter
}

// Allow checks if the request from the given IP is allowed.
func (lim *IPRateLimiter) Allow(ip string) bool {
	lim.mu.RLock()
	rl, exists := lim.ips[ip]
	lim.mu.RUnlock()

	if !exists {
		lim.mu.Lock()
		rl, exists = lim.ips[ip]
		if !exists {
			rl = rate.NewLimiter(lim.limiter.Limit(), lim.limiter.Burst())
			lim.ips[ip] = rl
		}
		lim.mu.Unlock()
	}

	return rl.Allow()
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Rate Limit: Max 5 times in 1 second
		limiter := GetIPRateLimiter(1, 5)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if !limiter.Allow(ip) {
			internal.APIError(w, "Middleware::RateLimiter", http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
