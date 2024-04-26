package caches

import (
	"sync"
	"voting-blockchain/internal/models"

	"github.com/google/uuid"
)

type VotingCache struct {
	usersAlreadyIn map[uuid.UUID]bool
	userById       map[uint64]uuid.UUID
	lastUserId     uint64
	usersMut       *sync.RWMutex

	nodesAlreadyIn map[uuid.UUID]bool
	nodesById      map[uint64]models.NodeMeta
	lastNodeId     uint64
	nodesMut       *sync.RWMutex

	headers     map[uint64]models.BlockHeader
	lastBlockId uint64
	headersMut  *sync.RWMutex
}

func NewVotePayloadValidatorStorage() *VotingCache {
	return &VotingCache{
		usersAlreadyIn: make(map[uuid.UUID]bool),
		userById:       make(map[uint64]uuid.UUID),
		lastUserId:     0,
		usersMut:       &sync.RWMutex{},

		nodesAlreadyIn: make(map[uuid.UUID]bool),
		nodesById:      make(map[uint64]models.NodeMeta),
		lastNodeId:     0,
		nodesMut:       &sync.RWMutex{},

		headers:    make(map[uint64]models.BlockHeader),
		headersMut: &sync.RWMutex{},
	}
}

func (cache *VotingCache) AddHeader(header models.BlockHeader) {
	cache.headersMut.Lock()
	defer cache.headersMut.Unlock()
	cache.lastBlockId += 1
	cache.headers[cache.lastBlockId] = header

	// here we need to add nodes and users !!!
}

func (cache *VotingCache) GetLastBlockHeader() models.BlockHeader {
	cache.headersMut.RLock()
	defer cache.headersMut.RUnlock()
	return cache.headers[cache.lastBlockId]
}

func (cache *VotingCache) GetNodesCount() uint64 {
	cache.nodesMut.RLock()
	defer cache.nodesMut.RUnlock()
	return cache.lastNodeId + 1
}

func (cache *VotingCache) GetNodeBySerial(serial uint64) models.NodeMeta {
	return cache.nodesById[serial]
}

func (cache VotingCache) IsUserAlreadyIn(userUuid uuid.UUID) (bool, error) {
	cache.usersMut.RLock()
	defer cache.usersMut.RUnlock()
	return cache.usersAlreadyIn[userUuid], nil
}

func (cache VotingCache) IsNodeAlreadyIn(nodeUuid uuid.UUID) (bool, error) {
	cache.nodesMut.RLock()
	defer cache.nodesMut.RUnlock()
	return cache.nodesAlreadyIn[nodeUuid], nil
}

type VotingsCache struct {
	Votings map[uuid.UUID]VotingCache
	Mut     *sync.RWMutex
}

func NewVotingsCache() VotingsCache {
	return VotingsCache{
		Mut: &sync.RWMutex{},
	}
}
