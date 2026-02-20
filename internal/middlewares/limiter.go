package middlewares

import (
	"net/http"

	"golang.org/x/time/rate"
)

func LimitMiddleware(next http.HandlerFunc, requestsPerSecond int) http.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond)

	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}
