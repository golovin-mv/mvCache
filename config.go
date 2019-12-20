// TODO: to config package
package main

import (
	"errors"
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type ConsulConfig struct {
	Enable      bool
	Host        string
	ServiceName string `yaml:"service-name"`
}

type cache struct {
	Type    string
	Ttl     int64
	Address string
}
type Config struct {
	Port  int16
	Proxy struct {
		To string
	}

	Cache cache

	Statistic struct {
		Storetime uint64
	}

	CacheErrors bool `yaml:"cache-errors"`
	Consul      ConsulConfig
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

	// если не задано куда кэшировать - паникуем
	if conf.Proxy.To == "" {
		panic(errors.New("Proxy to not set"))
	}

	// если у нас не задан кэш, по умолчанию будет в памяти
	if (cache{}) == conf.Cache {
		conf.Cache = cache{Ttl: 10, Type: "memory"}
	}

	return &conf
}
