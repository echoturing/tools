package tools

import (
	"github.com/go-redis/redis"
	"github.com/spaolacci/murmur3"
)

// default string hash object
// can be use to hash a string
type StringHash string

func (s StringHash) Hash() (uint32, error) {
	hashObj := murmur3.New32()
	_, err := hashObj.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	position := hashObj.Sum32()
	return position, nil
}

type CanHash interface {
	Hash() (uint32, error)
}

type BloomFilterInterface interface {
	AddItem(item CanHash) error
	RemoveItem(item CanHash) error
	TestItem(item CanHash) (bool, error)
}

type redisBloomFilter struct {
	client    *redis.Client
	keySuffix string
}

func NewRedisBloomFilter(addr, password string, db int, keySuffix string) *redisBloomFilter {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &redisBloomFilter{
		client:    client,
		keySuffix: keySuffix,
	}
}

func (r *redisBloomFilter) getKey() string {
	return "bloomFilter|" + r.keySuffix
}

func (r *redisBloomFilter) AddItem(item CanHash) error {
	position, err := item.Hash()
	if err != nil {
		return err
	}
	return r.client.SetBit(r.getKey(), int64(position), 1).Err()
}

func (r *redisBloomFilter) RemoveItem(item CanHash) error {
	position, err := item.Hash()
	if err != nil {
		return err
	}
	return r.client.SetBit(r.getKey(), int64(position), 0).Err()

}

func (r *redisBloomFilter) TestItem(item CanHash) (bool, error) {
	position, err := item.Hash()
	if err != nil {
		return false, err
	}
	isTrue, err := r.client.GetBit(r.getKey(), int64(position)).Result()
	if err != nil {
		return false, err
	}
	if isTrue == 1 {
		return true, nil
	}
	return false, nil
}

var _ BloomFilterInterface = &redisBloomFilter{}
