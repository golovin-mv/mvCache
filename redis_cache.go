package main

import (
	"github.com/go-redis/redis"
)

type RedisCache struct {
	ttl    int16
	client *redis.Client
}

func (r *RedisCache) Get(key string) (error, *CacheObject) {

}

func (r *RedisCache) Add(key string, data interface{}) error {

}

func (r *RedisCache) Remove(key string) {

}

func NewRedisCache(ttl int16) *RedisCache {
	conf := GetConfig()
	c := new(RedisCache)
	client := redis.NewClient()
}
