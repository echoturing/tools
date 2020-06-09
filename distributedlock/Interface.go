package distributedlock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
)

type Locker interface {
	Lock(ctx context.Context, key string, value string) (bool, error)
	Unlock(ctx context.Context, key string, value string) (success bool, err error)
	option
}

type option interface {
	setExpire(expire time.Duration)
	setKeyPrefix(prefix string)
}

var _ Locker = (*redisLocker)(nil)

func NewRedisLocker(client *redis.Client, optionFunc ...optionFunc) Locker {
	locker := &redisLocker{
		client: client,
		expire: time.Minute, // TODO(xiangxu)temporary set 1min,give a set func later
	}
	for _, f := range optionFunc {
		f(locker)
	}
	return locker
}

type optionFunc func(locker Locker)

var WithExpire = func(expire time.Duration) optionFunc {
	return func(locker Locker) {
		locker.setExpire(expire)
	}
}

var WithKeyPrefix = func(prefix string) optionFunc {
	return func(locker Locker) {
		locker.setKeyPrefix(prefix)
	}
}
