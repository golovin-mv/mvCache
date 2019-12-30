package config

import (
	"io/ioutil"
	"sync"

	"github.com/golovin-mv/mvCache/mutation"

	"github.com/golovin-mv/mvCache/guard"

	"github.com/golovin-mv/mvCache/proxy"

	"github.com/golovin-mv/mvCache/consul"

	"github.com/golovin-mv/mvCache/cache"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Port int16

	Cache *cache.CacheConfig

	Statistic struct {
		Storetime uint64
	}

	CacheErrors bool `yaml:"cache-errors"`
	Consul      *consul.ConsulConfig
	Proxy       *proxy.ProxyConfig
	Guard       *guard.GuardCongif
	Mutation    *mutation.MutationConfig
}

var (
	conf *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		conf = createConfig()
	})

	return conf
}

func createConfig() *Config {
	data, err := ioutil.ReadFile("./config/config.yml")

	if err != nil {
		panic(err)
	}

	// парсим конфигурацию
	conf := Config{}
	err = yaml.Unmarshal(data, &conf)

	if err != nil {
		panic(err)
	}

	return &conf
}
