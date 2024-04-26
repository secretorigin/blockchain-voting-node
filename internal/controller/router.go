package controller

import (
	"net/http"
	"voting-blockchain/internal/validators"
)

type Router struct {
	VoteValidator *validators.VotePayloadValidator
}

func (rt *Router) Set(mux *http.ServeMux) {
	mux.HandleFunc("/voting/{voting_uuid}/vote", rt.VotePost)
	mux.HandleFunc("/node", rt.NodePost)
}
