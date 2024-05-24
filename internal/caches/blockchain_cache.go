package caches

import (
	"errors"
	"voting-blockchain/internal/models/blockchains"
	"voting-blockchain/internal/models/blocks"
	"voting-blockchain/internal/models/types"

	"github.com/google/uuid"
)

type BlockchainCache struct {
	LastBlockHeader blocks.BlockHeader
	LastBlockHash   []byte
	Users           map[types.Uuid]bool
	Nodes           map[types.Uuid]blockchains.NodeMeta
}

type BlockchainsCache struct {
	Blockchains map[types.Uuid]BlockchainCache
}

func NewBlockchainsCache() *BlockchainsCache {
	return &BlockchainsCache{
		Blockchains: make(map[types.Uuid]BlockchainCache),
	}
}

func (cache BlockchainsCache) ExistsBlockchain(blockchain types.Uuid) bool {
	_, exists := cache.Blockchains[blockchain]
	return exists
}

func (cache BlockchainsCache) ExistsUser(blockchain types.Uuid, user types.Uuid) (bool, error) {
	value, exists := cache.Blockchains[blockchain]
	if !exists {
		return false, errors.New("blockchain does not exists: " + uuid.UUID(blockchain).String())
	}

	_, exists = value.Users[user]

	return exists, nil
}

func (cache BlockchainsCache) GetLastBlockHeader(blockchain types.Uuid) (blocks.BlockHeader, error) {
	var result blocks.BlockHeader
	value, exists := cache.Blockchains[blockchain]
	if !exists {
		return result, errors.New("blockchain does not exists: " + uuid.UUID(blockchain).String())
	}

	result = value.LastBlockHeader

	return result, nil
}

func (cache BlockchainsCache) GetLastBlockHash(blockchain types.Uuid) ([]byte, error) {
	var result []byte
	value, exists := cache.Blockchains[blockchain]
	if !exists {
		return result, errors.New("blockchain does not exists: " + uuid.UUID(blockchain).String())
	}

	result = value.LastBlockHash

	return result, nil
}

func (cache *BlockchainsCache) Update(blockchain types.Uuid, new_block_header blocks.BlockHeader) error {
	var users map[types.Uuid]bool
	var nodes map[types.Uuid]blockchains.NodeMeta

	_, exists := cache.Blockchains[blockchain]
	if exists {
		users = cache.Blockchains[blockchain].Users
		nodes = cache.Blockchains[blockchain].Nodes
	} else {
		users = make(map[types.Uuid]bool)
		nodes = make(map[types.Uuid]blockchains.NodeMeta)
	}

	for _, val := range new_block_header.Votes {
		users[val] = true
	}
	for _, val := range new_block_header.Nodes {
		nodes[val.Uuid] = val
	}

	cache.Blockchains[blockchain] = BlockchainCache{
		LastBlockHeader: new_block_header,
		LastBlockHash:   new_block_header.Hash(),
		Users:           users,
		Nodes:           nodes,
	}

	return nil
}
