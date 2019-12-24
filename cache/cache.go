package cache

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

type CacheConfig struct {
	Type    string
	Ttl     int64
	Address string
}

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

func CreateCacher(c *CacheConfig) Cacher {
	if c == nil {
		c = defaultConfig()
	}

	switch c.Type {
	case "memory":
		CurrentCacher = &InMemoryCache{make(map[string]storedData), c.Ttl}
	case "redis":
		CurrentCacher = NewRedisCache(c)
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

func defaultConfig() *CacheConfig {
	conf := CacheConfig{Ttl: 10, Type: "memory"}

	return &conf
}
