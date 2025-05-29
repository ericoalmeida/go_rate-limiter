package middleware

import (
	"net/http"
	"strings"

	"github.com/ericoalmeida/go_rate-limiter/internal/limiter"
)

func RateLimiterMiddleware(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("API_KEY")
			var key string
			var isToken bool

			if apiKey != "" {
				key = apiKey
				isToken = true
			} else {
				ip := strings.Split(r.RemoteAddr, ":")[0]
				key = ip
			}

			allowed, _, err := l.Allow(key, isToken)
			if err != nil || !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
