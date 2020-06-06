package distributedlock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Locker interface {
	Lock(ctx context.Context, key string, value string) (bool, error)
	Unlock(ctx context.Context, key string, value string) (success bool, err error)
}

var _ Locker = (*redisLocker)(nil)

func NewRedisLocker(client *redis.Client) Locker {
	return &redisLocker{
		client: client,
		expire: time.Minute, // TODO(xiangxu)temporary set 1min,give a set func later
	}
}
