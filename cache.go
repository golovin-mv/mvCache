package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CacheObject struct {
	CreatedAt time.Time
	Data      interface{}
}

type Cacher interface {
	Get(key string) (error, CacheObject)
	Add(key string, data interface{}) error
	Remove(key string)
}

// memory cache
type InMemoryCache struct {
	data map[string]CacheObject
	ttl  int16
}

func (i *InMemoryCache) Get(key string) (error, CacheObject) {
	// проверим существеут ли элемент
	if val, ok := i.data[key]; ok {
		d := val.CreatedAt.Sub(time.Now())

		if d.Seconds() > float64(i.ttl) {
			i.Remove(key)
		} else {
			return nil, val
		}
	}

	return errors.New(key + ": not exist"), CacheObject{}
}

func (i *InMemoryCache) Remove(key string) {
	delete(i.data, key)
}

func (i *InMemoryCache) Add(key string, data interface{}) error {
	i.data[key] = CacheObject{time.Now(), data}

	return nil
}

func GetKey(r *http.Request) string {
	url, err := url.Parse(r.RequestURI)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return ""
	}

	key := r.Method + url.Path + url.RawQuery + string(body)

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return getMD5Hash(strings.ToLower(key))
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateCacher(ctype string, ttl int16) *Cacher {
	var cacher Cacher

	switch ctype {
	case "memory":
		cacher = &InMemoryCache{ttl: ttl}
	default:
		log.Fatalln("Unknown cacher type")
	}

	return &cacher
}
