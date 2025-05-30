package limiter

import (
	"time"

	"github.com/ericoalmeida/go_rate-limiter/internal/configs"
)

type Limiter struct {
	store            Store
	defaultLimit     int
	defaultBlockTime time.Duration
	tokenLimits      *TokenLimitStore
}

func NewLimiter(store Store) *Limiter {
	defaultRateLimit := configs.GetEnvInt("DEFAULT_RATE_LIMIT", 5)
	defaultBlockDuration := configs.GetEnvInt("DEFAULT_BLOCK_DURATION", 300)

	return &Limiter{
		store:            store,
		defaultLimit:     defaultRateLimit,
		defaultBlockTime: time.Duration(defaultBlockDuration) * time.Second,
		tokenLimits:      NewTokenLimitStore(),
	}
}

func (l *Limiter) SetTokenLimit(token string, limit int, blockTime time.Duration) {
	l.tokenLimits.Set(token, TokenLimitConfig{
		Limit:         limit,
		BlockDuration: blockTime,
	})
}

func (l *Limiter) Allow(key string, isToken bool) (bool, string, error) {
	blocked, err := l.store.IsBlocked(key)
	if err != nil {
		return false, "", err
	}
	if blocked {
		return false, key, nil
	}

	limit := l.defaultLimit
	blockTime := l.defaultBlockTime
	if isToken {
		if tLimit, ok := l.tokenLimits.Get(key); ok {
			limit = tLimit.Limit
		}
		if tBlock, ok := l.tokenLimits.Get(key); ok {
			blockTime = tBlock.BlockDuration
		}
	}

	count, err := l.store.Increment(key, time.Second)
	if err != nil {
		return false, "", err
	}

	if count > limit {
		_ = l.store.Block(key, blockTime)
		return false, key, nil
	}

	return true, key, nil
}
