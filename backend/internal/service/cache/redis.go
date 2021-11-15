package cache

import (
	"backend/internal/config"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type RedisDB struct {
	CacheConn *redis.Client
	cacheCtx  context.Context
}

const TTL = time.Minute * 2

func InitCache(ctx context.Context, cfg *config.Config) (*RedisDB, error) {
	cacheConn := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheAdress,
		Password: "", // no password set
		DB:       0,  // use default DB,
	})

	_, err := cacheConn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisDB{cacheCtx: ctx, CacheConn: cacheConn}, nil
}

func (r RedisDB) Close() error {
	return r.CacheConn.Close()
}

func (r RedisDB) Set(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithCancel(r.cacheCtx)
	defer cancel()
	err := r.CacheConn.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisDB) SetData(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithCancel(r.cacheCtx)
	defer cancel()

	zvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.CacheConn.Set(ctx, key, zvalue, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisDB) Get(key string) (string, error) {
	ctx, cancel := context.WithCancel(r.cacheCtx)
	defer cancel()
	resp, err := r.CacheConn.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	if resp == "" {
		return "", errors.New("empty redis get")
	}
	return resp, nil
}

func (r RedisDB) GetData(key string, data interface{}) error {
	ctx, cancel := context.WithCancel(r.cacheCtx)
	defer cancel()
	resp, err := r.CacheConn.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return err
	}
	return nil
}

func (r RedisDB) Del(key ...string) (int64, error) {
	ctx, cancel := context.WithCancel(r.cacheCtx)
	defer cancel()
	deleted, err := r.CacheConn.Del(ctx, key...).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (c RedisDB) CheckCacheStatus(uri string) (int, error) {
	ctx, cancel := context.WithCancel(c.cacheCtx)
	defer cancel()
	result, err := c.CacheConn.Get(ctx, uri).Result()
	if err != nil {
		return 0, err
	}
	status, _ := strconv.Atoi(result)
	return status, nil
}

func (c RedisDB) AddCacheStatus(uri string, status int) error {
	ctx, cancel := context.WithCancel(c.cacheCtx)
	defer cancel()
	err := c.CacheConn.Set(ctx, uri, status, TTL).Err()
	if err != nil {
		return err
	}
	return nil
}
