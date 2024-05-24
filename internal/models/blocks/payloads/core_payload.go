package payloads

import (
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/models/votings"
	"voting-blockchain/internal/utils"
)

type CorePayload struct {
	UserUuid types.Uuid     `json:"user_uuid"` // user uuid
	Voting   votings.Voting `json:"voting"`
}

func (pl CorePayload) Size() uint64 {
	return types.UUID_SIZE + pl.Voting.Size() + types.SIGNATURE_SIZE
}

func (pl CorePayload) Marshal() []byte {
	bytes := make([]byte, pl.Size())
	last_index := uint64(0)

	// user uuid
	copy(bytes[last_index:last_index+types.UUID_SIZE], pl.UserUuid[:])
	last_index += types.UUID_SIZE
	// voting
	voting_bytes := pl.Voting.Marshal()
	copy(bytes[last_index:last_index+pl.Voting.Size()], voting_bytes)
	last_index += pl.Voting.Size()
	// signature
	copy(bytes[last_index:last_index+types.SIGNATURE_SIZE], pl.Signature[:])
	last_index += types.SIGNATURE_SIZE

	return bytes
}

func (pl *CorePayload) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	copy(pl.UserUuid[:], bytes[last_index:last_index+types.UUID_SIZE])
	last_index += types.UUID_SIZE
	// voting
	_ = pl.Voting.Unmarshal(bytes[last_index:])
	last_index += pl.Voting.Size()
	// signature
	copy(pl.Signature[:], bytes[last_index:last_index+types.SIGNATURE_SIZE])
	last_index += types.SIGNATURE_SIZE

	return nil
}

func (pl *CorePayload) Hash() []byte {
	bytes := pl.Marshal()
	return utils.GetHash(bytes)
}
