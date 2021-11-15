package cache

import "time"

type CacheService interface {
	Set(string, interface{}, time.Duration) error
	SetData(string, interface{}, time.Duration) error
	Get(string) (string, error)
	GetData(string, interface{}) error
	Del(...string) (int64, error)
	Close() error
}
