package limiter

import (
	"time"

	"github.com/ericoalmeida/go_rate-limiter/internal/configs"
)

type Limiter struct {
	store            Store
	defaultLimit     int
	defaultBlockTime time.Duration
	tokenLimits      map[string]int
	tokenBlockTime   map[string]time.Duration
}

func NewLimiter(store Store) *Limiter {
	defaultRateLimit := configs.GetEnvInt("DEFAULT_RATE_LIMIT", 5)
	defaultBlockDuration := configs.GetEnvInt("DEFAULT_BLOCK_DURATION", 300)

	return &Limiter{
		store:            store,
		defaultLimit:     defaultRateLimit,
		defaultBlockTime: time.Duration(defaultBlockDuration) * time.Second,
		tokenLimits:      make(map[string]int),
		tokenBlockTime:   make(map[string]time.Duration),
	}
}

func (l *Limiter) SetTokenLimit(token string, limit int, blockTime time.Duration) {
	l.tokenLimits[token] = limit
	l.tokenBlockTime[token] = blockTime
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
		if tLimit, ok := l.tokenLimits[key]; ok {
			limit = tLimit
		}
		if tBlock, ok := l.tokenBlockTime[key]; ok {
			blockTime = tBlock
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
