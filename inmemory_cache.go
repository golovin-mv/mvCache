package main

import (
	"errors"
	"time"
)

type InMemoryCache struct {
	data map[string]CacheObject
	ttl  int16
}

func (i *InMemoryCache) Get(key string) (error, *CacheObject) {
	// проверим существеут ли элемент
	if val, ok := i.data[key]; ok {
		d := time.Now().Sub(val.CreatedAt)

		if d.Seconds() > float64(i.ttl) {
			i.Remove(key)
		} else {
			return nil, &val
		}
	}

	return errors.New(key + ": not exist"), nil
}

func (i *InMemoryCache) Remove(key string) {
	delete(i.data, key)
}

func (i *InMemoryCache) Add(key string, data interface{}) error {
	i.data[key] = CacheObject{time.Now(), data}

	return nil
}
