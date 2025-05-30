package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ericoalmeida/go_rate-limiter/internal/configs"
	"github.com/ericoalmeida/go_rate-limiter/internal/limiter"
)

func RateLimiterMiddleware(l *limiter.Limiter) func(http.Handler) http.Handler {
	defaultRateLimit := configs.GetEnvInt("DEFAULT_RATE_LIMIT", 5)
	defaultBlockDuration := configs.GetEnvInt("DEFAULT_BLOCK_DURATION", 300)
	rateLimitToken := configs.GetEnvInt("RATE_LIMIT_TOKEN", defaultRateLimit)
	blockDurationToken := configs.GetEnvInt("BLOCK_DURATION_TOKEN", defaultBlockDuration)

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

			l.SetTokenLimit(apiKey, rateLimitToken, time.Duration(blockDurationToken)*time.Second)

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
