package httpcontroller

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	"voting-blockchain/internal/models/blocks"
	"voting-blockchain/internal/models/blocks/payloads"
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/models/votings"
)

func (rt *Router) VotingPost(w http.ResponseWriter, r *http.Request) {
	body_bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot read body: " + err.Error()))
		return
	}

	var voting votings.Voting

	err = json.Unmarshal(body_bytes, &voting)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot unmarshal body: " + err.Error()))
		return
	}

	payload := payloads.CorePayload{
		UserUuid: rt.userConfig.Uuid,
		Voting:   voting,
	}

	prevHash := types.BlockHash{}
	for i := 0; i < len(prevHash); i++ {
		prevHash[i] = 0
	}

	header := blocks.BlockHeader{
		Id:         0,
		Epoch:      time.Now().Unix(),
		Type:       blocks.BLOCK_TYPE_CORE,
		PrevHash:   prevHash,
		MerkleRoot: types.BlockMerkleRoot(payload.Hash()),
	}

	rt.blockchainCache.Update(voting.Uuid, header)
}
