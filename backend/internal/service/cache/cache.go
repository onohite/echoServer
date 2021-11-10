package cache

import "time"

type CacheService interface {
	Set(string, interface{}, time.Duration) error
	Get(string) (string, error)
	Del(...string) (int64, error)
	Close() error
}
