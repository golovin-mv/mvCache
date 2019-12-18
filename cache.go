package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type CachedResponse struct {
	Headers map[string]string
	Body    []byte
}

var CurrentCacher Cacher

type Cacher interface {
	Get(key string) (error, *CachedResponse)
	Add(key string, data interface{}) error
	Remove(key string)
	Count() int
	Clear()
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateCacher(ctype string, ttl int64) Cacher {
	switch ctype {
	case "memory":
		CurrentCacher = &InMemoryCache{make(map[string]storedData), ttl}
	case "redis":
		CurrentCacher = NewRedisCache(ttl)
	default:
		log.Fatalln("Unknown cacher type")
	}

	return CurrentCacher
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
