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
	"time"
)

type CacheObject struct {
	CreatedAt time.Time
	Data      interface{}
}

type Cacher interface {
	Get(key string) (error, *CacheObject)
	Add(key string, data interface{}) error
	Remove(key string)
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateCacher(ctype string, ttl int16) Cacher {
	var cacher Cacher

	switch ctype {
	case "memory":
		cacher = &InMemoryCache{make(map[string]CacheObject), ttl}
	default:
		log.Fatalln("Unknown cacher type")
	}

	return cacher
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
