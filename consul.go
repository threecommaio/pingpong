package main

import (
	"log"
	"time"
)

func registerService(name string, addr string, port int, tags []string) {
	s := new(Service)
	s.New(name, addr, port, tags)
	err := s.Register()

	if err != nil {
		timer1 := time.NewTimer(2 * time.Second)
		<-timer1.C
		log.Println("couldn't connect to consul, retrying...")
		registerService(name, addr, port, tags)
	} else {
		log.Println("registered service in consul")
	}
}
