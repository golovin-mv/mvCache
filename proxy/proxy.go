package proxy

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/golovin-mv/mvCache/cache"
)

type ProxyConfig struct {
	Type        string
	Target      string
	Reserve     string
	CacheErrors bool
}

type Proxy interface {
	Serve(res http.ResponseWriter, req *http.Request)
}

func NewProxy(c *ProxyConfig) Proxy {
	var p Proxy
	switch c.Type {
	case "reverse":
		p = &ReverseProxy{c}
	case "retry":
		p = &RetryProxy{c}
	default:
		panic(errors.New("Unknown Proxy Type"))
	}

	return p
}

func headerToArray(header http.Header) map[string]string {
	res := make(map[string]string)
	for name, values := range header {
		for _, value := range values {
			res[name] = value
		}
	}
	return res
}

func makeHandler(key string, cacheError bool) func(r *http.Response) error {
	return func(r *http.Response) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return err
		}

		// если статус 200 или не 200 но мы кэшируем ошибки
		if isOkStatus(r.StatusCode) || (!isOkStatus(r.StatusCode) && cacheError) {
			cache.CurrentCacher.Add(key, &cache.CachedResponse{headerToArray(r.Header), body})
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		return nil
	}
}

func isOkStatus(status int) bool {
	return status >= 200 && status <= 299
}
