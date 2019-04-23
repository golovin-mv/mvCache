package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var cacher Cacher
var cRes *CachedResponse

type CachedResponse struct {
	headers http.Header
	body    string
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := GetConfig()
	// получаем ключ
	key := GetKey(r)
	// проверим есть ли значение в кэш
	err, data := cacher.Get(key)

	if err != nil {
		ServeReverseProxy(c.Proxy.To, w, r, saveBodyToCache)
		cacher.Add(key, cRes)
	}

	log.Print(data)
}

func saveBodyToCache(r *http.Response) error {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	cRes = &CachedResponse{r.Header, string(body)}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return nil
}

func main() {
	c := GetConfig()
	cacher = CreateCacher(c.Cache.Type, c.Cache.Ttl)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
