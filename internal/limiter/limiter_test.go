package limiter

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
)

func TestTokenLimitStore_SetAndGet(t *testing.T) {
	store := NewTokenLimitStore()

	token := "abc123"
	cfg := TokenLimitConfig{Limit: 10, BlockDuration: 5 * time.Second}
	store.Set(token, cfg)

	got, ok := store.Get(token)
	if !ok {
		t.Fatal("expected token config to exist")
	}
	if got.Limit != cfg.Limit || got.BlockDuration != cfg.BlockDuration {
		t.Errorf("expected %v, got %v", cfg, got)
	}
}

func TestRedisLimiter_TokenLimitOverride(t *testing.T) {
	srv, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})
	redisStore := NewRedisStore(rdb)
	limiter := NewLimiter(redisStore, 1, 1)

	token := "abc123"
	limiter.SetTokenLimit(token, 10, 5*time.Second)

	isToken := true

	for i := 0; i < 10; i++ {
		allowed, _, err := limiter.Allow(token, isToken)
		if err != nil {
			t.Fatal(err)
		}
		if !allowed {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	allowed, _, _ := limiter.Allow(token, isToken)
	if allowed {
		t.Error("expected request to be blocked after limit exceeded")
	}
}

func TestRedisLimiter_Block(t *testing.T) {
	srv, _ := miniredis.Run()
	defer srv.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})
	redisStore := NewRedisStore(rdb)
	limiter := NewLimiter(redisStore, 1, 2)

	token := ""
	isToken := false

	allowed, _, _ := limiter.Allow(token, isToken)
	if !allowed {
		t.Fatal("first request should be allowed")
	}

	allowed, _, _ = limiter.Allow(token, isToken)
	if allowed {
		t.Fatal("second request should be blocked")
	}
}
