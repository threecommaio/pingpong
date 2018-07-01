package main

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

// Service struct handling consul
type Service struct {
	Name        string
	Address     string
	Port        int
	Tags        []string
	ConsulAgent *consul.Agent
}

// Register a new service in consul
func (s *Service) Register() error {
	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return err
	}
	s.ConsulAgent = c.Agent()
	serviceID := s.Name + "-" + s.Address

	serviceDef := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    s.Name,
		Port:    s.Port,
		Address: s.Address,
		Tags:    s.Tags,
		Check: &consul.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", s.Address, s.Port),
			Interval: "5s",
			Timeout:  "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return err
	}

	return nil
}

// New func handling setup of consul
func (s *Service) New(name string, address string, port int, tags []string) *Service {
	s.Name = name
	s.Port = port
	s.Address = address

	return s
}
