package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/golovin-mv/mvCache/guard"

	"github.com/golovin-mv/mvCache/consul"

	"github.com/golovin-mv/mvCache/proxy"

	"github.com/golovin-mv/mvCache/config"

	"github.com/golovin-mv/mvCache/cache"
	"github.com/jasonlvhit/gocron"
)

type Counter struct {
	Requests uint64
	Proxy    uint64
}

type ProxyCount struct {
	Count uint64
}

var count *Counter
var p proxy.Proxy
var cacher cache.Cacher
var consulClient *consul.ConsulClient
var gu *guard.Guard

func handler(w http.ResponseWriter, r *http.Request) {
	if gu != nil {
		if !gu.Guard(w, r) {
			return
		}
	}
	// получаем ключ
	key := cache.GetKey(r)
	// проверим есть ли значение в кэш
	err, data := cache.CurrentCacher.Get(key)

	if err != nil {
		p.Serve(w, r)
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
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
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

func getCachedCount(w http.ResponseWriter, r *http.Request) {
	s, _ := json.Marshal(ProxyCount{uint64(cacher.Count())})
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(s)
}

func dropCache(w http.ResponseWriter, r *http.Request) {
	cacher.Clear()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}

func main() {
	c := config.GetConfig()

	count = &Counter{}
	p = proxy.NewProxy(c.Proxy)
	cacher = cache.CreateCacher(c.Cache)

	http.HandleFunc("/", handler)
	http.HandleFunc("/api/stat", getCounter)
	http.HandleFunc("/api/stat/drop", dropCounterHandler)
	http.HandleFunc("/api/cache/stat", getCachedCount)
	http.HandleFunc("/api/cache/drop", dropCache)

	if c.Statistic.Storetime > 0 {
		go gocron.Start()
		gocron.Every(c.Statistic.Storetime).Seconds().Do(dropCounter)
	}

	if c.Consul.Enable {
		consulClient := consul.NewConsulClient(c.Consul, c.Port)
		consulClient.Connect()
	}

	if c.Guard.Enable {
		gu = guard.NewGuard(c.Guard)
	}

	log.Println(fmt.Sprintf("mvCache run port: %d", c.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
