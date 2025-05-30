package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) Increment(key string, window time.Duration) (int, error) {
	ctx := context.Background()
	pipe := r.client.TxPipeline()
	count := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return int(count.Val()), nil
}

func (r *RedisStore) Get(key string) (int, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *RedisStore) Reset(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStore) Block(key string, duration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, "block:"+key, true, duration).Err()
}

func (r *RedisStore) IsBlocked(key string) (bool, error) {
	ctx := context.Background()
	res, err := r.client.Get(ctx, "block:"+key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return res == "1" || res == "true", nil
}
