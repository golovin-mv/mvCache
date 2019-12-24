package cache

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	ttl    int64
	client *redis.Client
}

func (r *RedisCache) Get(key string) (error, *CachedResponse) {
	co := CachedResponse{}
	s, err := r.client.Get(key).Result()

	if err != nil {
		return err, nil
	}

	err = json.Unmarshal([]byte(s), &co)

	if err != nil {
		return err, nil
	}

	return nil, &co
}

func (r *RedisCache) Add(key string, data interface{}) error {
	s, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = r.client.Set(key, s, time.Duration(r.ttl)*time.Second).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Remove(key string) {
	r.client.Del(key)
}

func (r *RedisCache) Count() int {
	c := r.makeCommand("DBSIZE")
	res, err := c.Result()

	if err != nil {
		return -1
	}

	co, err := strconv.ParseInt(res, 10, 32)

	if err != nil {
		return -1
	}

	return int(co)
}

func (r *RedisCache) Clear() {
	r.makeCommand("FLUSHDB")
}

func (r *RedisCache) makeCommand(com string) *redis.StringCmd {
	cmd := redis.NewStringCmd(com)
	r.client.Process(cmd)

	return cmd
}

func NewRedisCache(conf *CacheConfig) *RedisCache {
	c := new(RedisCache)
	client := redis.NewClient(&redis.Options{Addr: conf.Address})
	c.client = client
	c.ttl = conf.Ttl
	_, err := client.Ping().Result()

	if err != nil {
		panic(errors.New("Failed connect to Redis : " + err.Error()))
	}
	log.Println("Connected to Redis: " + conf.Address)
	return c
}
