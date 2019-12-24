package consul

import (
	"crypto/rand"
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

type ConsulConfig struct {
	Enable      bool
	Host        string
	ServiceName string `yaml:"service-name"`
}

type ConsulClient struct {
	ServiceName string
	Host        string
	AppPort     int16
	Client      *api.Client
}

func (c *ConsulClient) Connect() {
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

	c.register(c.ServiceName)

	log.Println("Connected to Consul: " + c.Host)

}

func NewConsulClient(conf *ConsulConfig, appPort int16) *ConsulClient {
	client := ConsulClient{Host: conf.Host, ServiceName: conf.ServiceName}
	client.AppPort = appPort
	return &client
}

func (c *ConsulClient) register(name string) *api.AgentServiceRegistration {
	registration := new(api.AgentServiceRegistration)
	registration.ID = name + "-" + uuid()
	registration.Name = name
	address := hostname()
	registration.Address = address

	port := c.AppPort

	registration.Port = int(port)
	registration.Check = new(api.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, port)
	registration.Check.Interval = "10s"
	registration.Check.Timeout = "8s"

	c.Client.Agent().ServiceRegister(registration)

	return registration
}

func healtCheck(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(HealthCheckDto{"OK"})
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(res)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}

func uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
