package redis

import (
	"context"
	"fmt"

	"github.com/exceptioon/tiktok-fav-publisher/internal"
	"github.com/go-redis/redis/v9"
)

const keySet = "pubVideos"

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) (internal.Cache, error) {
	redis := redisCache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}

	err := redis.client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &redis, nil
}

func (r *redisCache) Add(value string) error {
	err := r.client.SAdd(context.Background(), keySet, value).Err()
	if err != nil {
		return err
	}
	if !r.IsExist(value) {
		return fmt.Errorf("value %s not added", value)
	}
	return nil
}

func (r *redisCache) IsExist(value string) bool {
	return r.client.SIsMember(context.Background(), keySet, value).Val()
}
