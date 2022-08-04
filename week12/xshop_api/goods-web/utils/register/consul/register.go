package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

type Registry struct {
	Host string
	Port int
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func  (r Registry)  Register(address string, port int, name string, tags []string, id string) error {
	//新增一个consul
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	//生成对于的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registeration := new(api.AgentServiceRegistration)
	registeration.Name = name
	registeration.ID = id
	registeration.Port = port
	registeration.Tags = tags
	registeration.Address = address
	registeration.Check = check

	err = client.Agent().ServiceRegister(registeration)
	if err != nil {
		panic(err)
	}
	return nil

}

func (r Registry) DeRegister(serviceId string)  error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	err = client.Agent().ServiceDeregister(serviceId)
	return err
}