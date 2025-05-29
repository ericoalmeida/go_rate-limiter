package limiter

import "time"

type Store interface {
	Increment(key string, window time.Duration) (int, error)
	Get(key string) (int, error)
	Reset(key string) error
	Block(key string, duration time.Duration) error
	IsBlocked(key string) (bool, error)
}
