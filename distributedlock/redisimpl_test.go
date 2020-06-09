package distributedlock

import (
	"testing"
	"time"

	"github.com/echoturing/log"
	"github.com/go-redis/redis/v7"
)

func NewConnection(addr, password string, db, poolSize int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 5,
		PoolSize:     poolSize,
	})
	return client
}

func TestRedisLocker_Locker(t *testing.T) {
	locker := NewRedisLocker(NewConnection("127.0.0.1:6379", "", 0, 2))

	key1 := "key1"
	value1 := "value1"
	value2 := "value2"
	ctx := log.NewDefaultContext()
	success, err := locker.Lock(ctx, key1, value1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !success {
		t.Error("acquire lock failed")
		return
	}

	success, err = locker.Unlock(ctx, key1, value2)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if success {
		t.Error("unlock error")
		return
	}
	success, err = locker.Unlock(ctx, key1, value1)
	if err != nil {
		t.Error("unlock error")
		return
	}
	if !success {
		t.Error("unlock error")
		return
	}
}

func TestNewRedisLocker(t *testing.T) {
	prefix := "prefix_xxx"
	expire := time.Hour
	locker := NewRedisLocker(NewConnection("127.0.0.1:6379", "", 0, 2), WithExpire(expire), WithKeyPrefix(prefix))
	rLocker := locker.(*redisLocker)
	if rLocker.keyPrefix != prefix {
		t.Error("prefix not set!")
		return
	}
	if rLocker.expire != expire {
		t.Error("expire not set")
		return
	}
}
