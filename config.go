// TODO: to config package
package main

import (
	"errors"
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Port  int16
	Proxy struct {
		To string
	}

	Cache *CacheConfig

	Statistic struct {
		Storetime uint64
	}

	CacheErrors bool `yaml:"cache-errors"`
	Consul      *ConsulConfig
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

	return &conf
}
