package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address

	check := &api.AgentServiceCheck{
		// CheckID:                        id,
		GRPC: fmt.Sprintf("%s:%d", address, port), //
		// HTTP:                           "http://" + address + ":" + strconv.Itoa(port) + "/health", // 检查服务的地址
		Timeout:                        "5s",   // 健康检查的超时时间
		Interval:                       "5s",   // 健康检查的间隔时间
		DeregisterCriticalServiceAfter: "400s", // 如果健康检查失败次数超过这个时间，则注销服务
	}

	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
