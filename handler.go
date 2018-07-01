package main

import "net/http"

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
