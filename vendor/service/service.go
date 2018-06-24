package service

import (
	"log"
	"net"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// Service struct handling consul
type Service struct {
	Name        string
	TTL         time.Duration
	Port        int
	ConsulAgent *consul.Agent
}

// UpdateTTL updates consul checks
func (s *Service) UpdateTTL() {
	ticker := time.NewTicker(s.TTL / 2)
	for range ticker.C {
		if agentErr := s.ConsulAgent.PassTTL("service:"+s.Name, ""); agentErr != nil {
			log.Print(agentErr)
		}
	}
}

// GetLocalIP returns the local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// New func handling setup of consul
func New(name string, port int, ttl time.Duration) (*Service, error) {
	s := new(Service)
	s.Name = name
	s.TTL = ttl
	s.Port = port

	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	s.ConsulAgent = c.Agent()
	serviceDef := &consul.AgentServiceRegistration{
		Name:    s.Name,
		Port:    s.Port,
		Address: GetLocalIP(),
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}
	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}

	return s, nil
}
