package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ericoalmeida/go_rate-limiter/internal/configs"
	"github.com/ericoalmeida/go_rate-limiter/internal/limiter"
	"github.com/ericoalmeida/go_rate-limiter/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	configs.LoadConfig()

	port := configs.GetEnv("PORT")
	redisHost := configs.GetEnv("REDIS_HOST")
	redisPassword := configs.GetEnv("REDIS_PASSWORD")
	redisDbInt := configs.GetEnvInt("REDIS_DB", 0)
	defaultRateLimit := configs.GetEnvInt("DEFAULT_RATE_LIMIT", 5)
	defaultBlockDuration := configs.GetEnvInt("DEFAULT_BLOCK_DURATION", 300)

	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		Password:     redisPassword,
		DB:           redisDbInt,
		PoolSize:     20,
		MinIdleConns: 10,
	})

	store := limiter.NewRedisStore(redisClient)

	rateLimiter := limiter.NewLimiter(store, defaultRateLimit, defaultBlockDuration)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})

	handler := middleware.RateLimiterMiddleware(rateLimiter)(mux)
	log.Println("Server running at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
