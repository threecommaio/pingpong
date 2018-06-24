package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type handler struct {
}

func (h *handler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	message := "hello\n"
	w.Write([]byte(message))
}

func consulRegistration(name string, port int, ttl time.Duration) {
	s, err := service.New(name, port, ttl)
	if err != nil {
		timer1 := time.NewTimer(2 * time.Second)
		<-timer1.C
		log.Println("couldn't connect to consul, retrying...")
		consulRegistration(name, port, ttl)
	} else {
		log.Println("registered service in consul")
		go s.UpdateTTL()
	}
}

func main() {

	port := flag.Int("port", 8080, "service port to run on")
	consul := flag.Bool("consul", false, "enable consult support (default: false)")
	ttl := flag.Duration("ttl", time.Second*15, "service ttl check duration")
	flag.Parse()

	if *consul {
		go consulRegistration("pingpong", *port, *ttl)
	}

	httpHandler := handler{}

	log.Printf("listening on %d\n", *port)
	http.Handle("/pong", prometheus.InstrumentHandler("pong", http.HandlerFunc(httpHandler.HandleRequest)))
	http.Handle("/metrics", promhttp.Handler())

	l := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(l, nil))
}
