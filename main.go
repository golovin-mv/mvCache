package main

import (
	"fmt"
	"log"
	"net/http"
)

var cacher *Cacher

func handler(w http.ResponseWriter, r *http.Request) {
	c := GetConfig()
	// проверим есть ли значение
	ServeReverseProxy(c.Proxy.To, w, r)
}

func main() {
	c := GetConfig()
	cacher = CreateCacher(c.Cache.Type, c.Cache.Ttl)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
