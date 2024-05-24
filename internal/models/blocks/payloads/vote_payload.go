package payloads

import (
	"encoding/binary"
	"voting-blockchain/internal/models/blockchains"
	"voting-blockchain/internal/models/votings"
	"voting-blockchain/internal/utils"
)

type VotePayload struct {
	Votes []votings.Vote     `json:"votes"`
	Nodes []blockchains.Node `json:"nodes"`
}

func (pl VotePayload) Size() uint64 {
	return uint64(8 + uint64(len(pl.Votes))*votings.VOTE_SIZE + 8 + uint64(len(pl.Nodes))*blockchains.NODE_SIZE)
}

func (pl VotePayload) Marshal() []byte {
	bytes := make([]byte, pl.Size())
	last_index := uint64(0)

	// votes count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(pl.Votes)))
	last_index += 8
	// votes
	for i := 0; i < len(pl.Votes); i++ {
		meta_bytes := pl.Votes[i].Marshal()
		copy(bytes[last_index:last_index+votings.VOTE_SIZE], meta_bytes[:])
		last_index += votings.VOTE_SIZE
	}
	// nodes count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(pl.Nodes)))
	last_index += 8
	// nodes
	for i := 0; i < len(pl.Nodes); i++ {
		meta_bytes := pl.Nodes[i].Marshal()
		copy(bytes[last_index:last_index+blockchains.NODE_SIZE], meta_bytes[:])
		last_index += blockchains.NODE_SIZE
	}

	return bytes
}

func (pl *VotePayload) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// votes count
	votesCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// votes
	for i := uint64(0); i < votesCount; i++ {
		var vote votings.Vote
		_ = vote.Unmarshal(bytes[last_index : last_index+votings.VOTE_SIZE])
		pl.Votes = append(pl.Votes, vote)
		last_index += votings.VOTE_SIZE
	}
	// nodes count
	nodesCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// nodes
	for i := uint64(0); i < nodesCount; i++ {
		var node blockchains.Node
		_ = node.Unmarshal(bytes[last_index : last_index+blockchains.NODE_SIZE])
		pl.Nodes = append(pl.Nodes, node)
		last_index += blockchains.NODE_SIZE
	}

	return nil
}

func (pl *VotePayload) Hash() []byte {
	bytes := pl.Marshal()
	return utils.GetHash(bytes)
}
