package main

import (
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Port  int16
	Proxy struct {
		To string
	}
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
