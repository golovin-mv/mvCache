package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

type HealthCheckDto struct {
	Status string
}

type ConsulClient struct {
	ServiceName string
	Host        string
	Client      *api.Client
}

func (c *ConsulClient) connect() {
	config := api.DefaultConfig()
	config.Address = c.Host

	client, err := api.NewClient(config)

	if err != nil {
		panic(err)
	}

	svc, err := connect.NewService(c.ServiceName, client)

	defer svc.Close()

	if err != nil {
		panic(err)
	}
	c.Client = client

	http.HandleFunc("/healthcheck", healtCheck)

	c.register()

	log.Println("Connected to Consul: " + c.Host)

}

func NewConsulClient(conf *ConsulConfig) *ConsulClient {
	client := ConsulClient{Host: conf.Host, ServiceName: conf.ServiceName}

	return &client
}

func (c *ConsulClient) register() *api.AgentServiceRegistration {
	registration := new(api.AgentServiceRegistration)
	registration.ID = "product-service"   //replace with service id
	registration.Name = "product-service" //replace with service name
	address := hostname()
	registration.Address = address

	port := port()

	registration.Port = int(port)
	registration.Check = new(api.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"

	c.Client.Agent().ServiceRegister(registration)

	return registration
}

func healtCheck(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(HealthCheckDto{"OK"})
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(res)
}

func port() int16 {
	return GetConfig().Port
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}
