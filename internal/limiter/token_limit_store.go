package limiter

import (
	"sync"
	"time"
)

type TokenLimitConfig struct {
	Limit         int
	BlockDuration time.Duration
}

type TokenLimitStore struct {
	mu     sync.RWMutex
	config map[string]TokenLimitConfig
}

func NewTokenLimitStore() *TokenLimitStore {
	return &TokenLimitStore{
		config: make(map[string]TokenLimitConfig),
	}
}

func (s *TokenLimitStore) Get(token string) (TokenLimitConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg, ok := s.config[token]
	return cfg, ok
}

func (s *TokenLimitStore) Set(token string, cfg TokenLimitConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config[token] = cfg
}
