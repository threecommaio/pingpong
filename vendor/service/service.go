package service

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

// Service struct handling consul
type Service struct {
	Name        string
	Address     string
	Port        int
	ConsulAgent *consul.Agent
}

func (s *Service) Check() {

}

// New func handling setup of consul
func New(name string, address string, port int, tags []string) (*Service, error) {
	s := new(Service)
	s.Name = name
	s.Port = port
	s.Address = address

	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	s.ConsulAgent = c.Agent()
	serviceID := s.Name + "-" + s.Address

	serviceDef := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    s.Name,
		Port:    s.Port,
		Address: s.Address,
		Tags:    tags,
		Check: &consul.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", s.Address, s.Port),
			Interval: "5s",
			Timeout:  "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}

	return s, nil
}
