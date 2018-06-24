package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type handler struct {
}

func (h *handler) HandleHealth(w http.ResponseWriter, req *http.Request) {
	message := "OK\n"
	w.Write([]byte(message))
}

func (h *handler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	message := "hello\n"
	w.Write([]byte(message))
}

func consulRegistration(name string, addr string, port int, tags []string) {
	s, err := service.New(name, addr, port, tags)
	if err != nil {
		timer1 := time.NewTimer(2 * time.Second)
		<-timer1.C
		log.Println("couldn't connect to consul, retrying...")
		consulRegistration(name, addr, port, tags)
	} else {
		log.Println("registered service in consul")
		go s.Check()
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

func main() {
	port := flag.Int("port", 8080, "service port to run on")
	consul := flag.Bool("consul", false, "enable consul support (default: false)")
	flag.Parse()

	addr := GetLocalIP()

	if *consul {
		go consulRegistration("pingpong", addr, *port, []string{"tbn-cluster"})
	}

	httpHandler := handler{}

	log.Printf("listening on %d\n", *port)
	http.Handle("/pong", prometheus.InstrumentHandler("pong", http.HandlerFunc(httpHandler.HandleRequest)))
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/health", prometheus.InstrumentHandler("health", http.HandlerFunc(httpHandler.HandleHealth)))

	l := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(l, nil))
}
