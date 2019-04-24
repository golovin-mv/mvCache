package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/jasonlvhit/gocron"
)

type Counter struct {
	Requests uint64
	Proxy    uint64
}

var count *Counter

func handler(w http.ResponseWriter, r *http.Request) {
	c := GetConfig()
	// получаем ключ
	key := GetKey(r)
	// проверим есть ли значение в кэш
	err, data := CurrentCacher.Get(key)

	if err != nil {
		ServeReverseProxy(c.Proxy.To, w, r)
		atomic.AddUint64(&count.Requests, 1)
		return
	}
	for key, val := range data.Headers {
		w.Header().Set(key, val)
	}
	w.Header().Set("Mv-Proxy", "proxy")
	w.Write(data.Body)
	atomic.AddUint64(&count.Proxy, 1)
}

func getCounter(w http.ResponseWriter, r *http.Request) {
	s, _ := json.Marshal(count)
	w.Write(s)
}

func dropCounterHandler(w http.ResponseWriter, r *http.Request) {
	dropCounter()
	getCounter(w, r)
}

func dropCounter() {
	atomic.SwapUint64(&count.Requests, 0)
	atomic.SwapUint64(&count.Proxy, 0)
}

func main() {
	c := GetConfig()
	count = &Counter{}
	go gocron.Start()
	CreateCacher(c.Cache.Type, c.Cache.Ttl)
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/counter", getCounter)
	http.HandleFunc("/api/counter/drop", dropCounterHandler)

	gocron.Every(c.Statistic.Storetime).Seconds().Do(dropCounter)
	log.Println(fmt.Sprintf("mvCache run port: %d", c.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
