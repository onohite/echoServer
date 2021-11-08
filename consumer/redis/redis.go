package redis

import (
	"consumer/config"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type CacheDB struct {
	CacheConn *redis.Client
}

const TTL = time.Minute * 2

func (c *CacheDB) InitCache(cfg *config.Config) {
	c.CacheConn = redis.NewClient(&redis.Options{
		Addr:     cfg.CacheAdress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (c CacheDB) CheckCacheStatus(cfg *config.Config, ctx context.Context, uri string) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result, err := c.CacheConn.Get(ctx, uri).Result()
	if err != nil {
		return 0, err
	}
	status, _ := strconv.Atoi(result)
	return status, nil
}

func (c CacheDB) AddCacheStatus(cfg *config.Config, ctx context.Context, uri string, status int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err := c.CacheConn.Set(ctx, uri, status, TTL).Err()
	if err != nil {
		return err
	}
	return nil
}
