package redis

import (
	"consmer/config"
	"github.com/go-redis/redis/v8"
)

type CacheDB struct {
	CacheConn *redis.Client
}

func (c *CacheDB) InitCache(cfg *config.Config) {
	c.CacheConn = redis.NewClient(&redis.Options{
		Addr:     cfg.CacheAdress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
