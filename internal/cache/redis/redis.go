package redis

import (
	"context"
	"fmt"
	"github.com/berkinyildiran/insider-case/internal/cache"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Redis struct {
	config *cache.Config

	client  *redis.Client
	context context.Context
	wg      sync.WaitGroup
}

func NewRedis(config *cache.Config, context context.Context) *Redis {
	options := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}

	return &Redis{
		config: config,

		client:  redis.NewClient(options),
		context: context,
	}
}

func (r *Redis) Close() error {
	r.wg.Wait()
	return r.client.Close()
}

func (r *Redis) Set(key string, value any, expiration time.Duration) error {
	r.wg.Add(1)
	defer r.wg.Done()

	return r.client.Set(r.context, key, value, expiration).Err()
}
