package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/param108/grpc-chat/errors"
	"os"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func (R *RedisCache) Connect() error {
	R.client = redis.NewClient(&redis.Options{
		Addr: getRedisDSN(),
	})

	_, err := R.client.Ping().Result()
	if err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func getRedisDSN() string {
	return fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
}

func (R *RedisCache) SetKey(key string, val string) error {
	err := R.client.Set(key, val, 0).Err()
	if err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (R *RedisCache) SetKeyWithExpiry(key string, val string, ttlSeconds int) error {
	err := R.client.SetNX(key, val, time.Duration(ttlSeconds)*time.Second).Err()
	if err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (R *RedisCache) GetKey(key string) (string, error) {
	val, err := R.client.Get(key).Result()
	if err != nil {
		return "", errors.NewInternalError(err)
	}
	return val, err
}
