package httpcontroller

import (
	"net/http"
	"voting-blockchain/internal/models/blocks"
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/user"
)

type Router struct {
	blockchainCache BlockchainsCacheInterface
	userConfig      *user.Config
}

func New(blockchainCache BlockchainsCacheInterface, userConfig *user.Config) *Router {
	return &Router{
		blockchainCache: blockchainCache,
		userConfig:      userConfig,
	}
}

func (rt *Router) Set(mux *http.ServeMux) {
	mux.HandleFunc("/voting/{voting_uuid}/vote", rt.VotePost)
	mux.HandleFunc("/node", rt.NodePost)
}

type BlockchainsCacheInterface interface {
	Update(blockchain types.Uuid, new_block_header blocks.BlockHeader)
}
