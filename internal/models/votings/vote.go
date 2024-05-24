package votings

import (
	"voting-blockchain/internal/models/types"
)

const VOTE_SIZE uint64 = types.UUID_SIZE + types.UUID_SIZE + types.SIGNATURE_SIZE

type Vote struct {
	UserUuid   types.Uuid `json:"user_uuid"`
	OptionUuid types.Uuid `json:"option_uuid"`
	// VotingUuid types.Uuid `json:"voting_uuid"` // TODO
	Signature types.Signature `json:"signature"`
}

func (vt Vote) Size() uint64 {
	return VOTE_SIZE
}

func (vt Vote) Marshal() []byte {
	bytes := make([]byte, vt.Size())
	last_index := uint64(0)

	// user uuid
	copy(bytes[last_index:last_index+types.UUID_SIZE], vt.UserUuid[:])
	last_index += types.UUID_SIZE
	// option uuid
	copy(bytes[last_index:last_index+types.UUID_SIZE], vt.OptionUuid[:])
	last_index += types.UUID_SIZE
	// signature
	copy(bytes[last_index:last_index+types.SIGNATURE_SIZE], vt.Signature[:])
	last_index += types.SIGNATURE_SIZE

	return bytes
}

func (vt *Vote) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	copy(vt.UserUuid[:], bytes[last_index:last_index+types.UUID_SIZE])
	last_index += types.UUID_SIZE
	// option uuid
	copy(vt.OptionUuid[:], bytes[last_index:last_index+types.UUID_SIZE])
	last_index += types.UUID_SIZE
	// signature
	copy(vt.Signature[:], bytes[last_index:last_index+types.SIGNATURE_SIZE])
	last_index += types.SIGNATURE_SIZE

	return nil
}
