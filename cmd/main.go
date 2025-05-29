package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ericoalmeida/go_rate-limiter/internal/configs"
	"github.com/ericoalmeida/go_rate-limiter/internal/limiter"
	"github.com/ericoalmeida/go_rate-limiter/internal/middleware"
)

func main() {
	configs.LoadConfig()

	redisAddr := configs.GetEnv("REDIS_ADDR")
	redisPassword := configs.GetEnv("REDIS_PASSWORD")
	redisDbInt := configs.GetEnvInt("REDIS_DB", 0)

	store := limiter.NewRedisStore(
		redisAddr,
		redisPassword,
		redisDbInt,
	)

	rateLimiter := limiter.NewLimiter(store)

	// exemplo de token personalizado com 100 req/s e bloqueio de 1 min
	rateLimiter.SetTokenLimit("abc123", 100, 60*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})

	handler := middleware.RateLimiterMiddleware(rateLimiter)(mux)
	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
