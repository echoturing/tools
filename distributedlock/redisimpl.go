package distributedlock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisLocker struct {
	client *redis.Client
	expire time.Duration
}

const (
	unlock string = `
if redis.call("get",KEYS[1]) == ARGV[1]
then
    return redis.call("del",KEYS[1])
else
    return 0
end
`
)

func (r *redisLocker) Lock(ctx context.Context, key string, value string) (bool, error) {
	success, err := r.client.SetNX(ctx, key, value, r.expire).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *redisLocker) Unlock(ctx context.Context, key string, value string) (success bool, err error) {
	res, err := r.client.Eval(ctx, unlock, []string{key}, value).Int()
	if err != nil {
		return false, nil
	}
	if res == 1 {
		success = true
	}
	return success, nil
}
