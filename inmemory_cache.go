package main

import (
	"errors"
	"time"
)

type InMemoryCache struct {
	data map[string]CachedResponse
	ttl  int64
}

func (i *InMemoryCache) Get(key string) (error, *CachedResponse) {
	// проверим существеут ли элемент
	if val, ok := i.data[key]; ok {
		d := time.Duration(0)

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
	return nil
}
