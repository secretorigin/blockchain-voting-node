package controller

import "net/http"

type Router struct {
}

func (rt *Router) Set(mux *http.ServeMux) {
	mux.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/node", func(w http.ResponseWriter, r *http.Request) {})
}
