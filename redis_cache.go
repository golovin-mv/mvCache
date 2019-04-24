package main

import (
	"encoding/json"
	"log"
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

func NewRedisCache(ttl int64) *RedisCache {
	conf := GetConfig()
	c := new(RedisCache)
	client := redis.NewClient(&redis.Options{Addr: conf.Cache.Address})
	c.client = client
	c.ttl = 0 - ttl
	log.Println("Connected to Redis: " + conf.Cache.Address)
	return c
}
