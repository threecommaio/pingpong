package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	port := flag.Int("port", 8080, "service port to run on")
	consul := flag.Bool("consul", false, "enable consul support (default: false)")
	flag.Parse()
	addr := GetLocalIP()

	if *consul {
		go registerService("pingpong", addr, *port, []string{"tbn-cluster"})
	}

	httpHandler := handler{}

	log.Printf("listening on %d\n", *port)
	http.Handle("/pong", prometheus.InstrumentHandler("pong", http.HandlerFunc(httpHandler.HandleRequest)))
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/health", prometheus.InstrumentHandler("health", http.HandlerFunc(httpHandler.HandleHealth)))

	l := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(l, nil))
}
