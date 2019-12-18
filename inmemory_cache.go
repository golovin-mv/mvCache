package main

import (
	"errors"
	"time"
)

type InMemoryCache struct {
	data map[string]storedData
	ttl  int64
}

type storedData struct {
	Data interface{}
	Time time.Time
}

func (i *InMemoryCache) Get(key string) (error, *CachedResponse) {
	// проверим существеут ли элемент
	if val, ok := i.data[key]; ok {

		if i.isExpired(val.Time) {
			i.Remove(key)
		} else {
			return nil, val.Data.(*CachedResponse)
		}
	}

	return errors.New(key + ": not exist"), nil
}

func (i *InMemoryCache) Remove(key string) {
	delete(i.data, key)
}

func (i *InMemoryCache) Add(key string, data interface{}) error {
	i.data[key] = storedData{data, time.Now()}
	return nil
}

func (i *InMemoryCache) Count() int {
	return len(i.data)
}

func (i *InMemoryCache) Clear() {
	i.data = make(map[string]storedData)
}

func (i *InMemoryCache) isExpired(t time.Time) bool {
	if i.ttl < 0 {
		return false
	}

	d := time.Now().Sub(t).Seconds()
	return d > float64(i.ttl)
}
