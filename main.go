package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	c := GetConfig()
	GetKey(r)
	ServeReverseProxy(c.Proxy.To, w, r)
}

func main() {
	config := GetConfig()

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
