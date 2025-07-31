package cache

import "time"

type Cache interface {
	Close() error
	Set(key string, value any, expiration time.Duration) error
}
