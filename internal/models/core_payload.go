package models

import (
	"crypto/sha256"

	"github.com/google/uuid"
)

type CorePayload struct {
	UserUuid  uuid.UUID `json:"user_uuid"` // user uuid
	Voting    Voting    `json:"voting"`
	Signature Signature `json:"signature"`
}

func (pl CorePayload) Size() uint64 {
	return UUID_SIZE + pl.Voting.Size() + SIGNATURE_SIZE
}

func (pl CorePayload) Marshal() []byte {
	bytes := make([]byte, pl.Size())
	last_index := uint64(0)

	// user uuid
	copy(bytes[last_index:last_index+UUID_SIZE], pl.UserUuid[:])
	last_index += UUID_SIZE
	// voting
	voting_bytes := pl.Voting.Marshal()
	copy(bytes[last_index:last_index+pl.Voting.Size()], voting_bytes)
	last_index += pl.Voting.Size()
	// signature
	copy(bytes[last_index:last_index+SIGNATURE_SIZE], pl.Signature[:])
	last_index += SIGNATURE_SIZE

	return bytes
}

func (pl *CorePayload) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	copy(pl.UserUuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// voting
	_ = pl.Voting.Unmarshal(bytes[last_index:])
	last_index += pl.Voting.Size()
	// signature
	copy(pl.Signature[:], bytes[last_index:last_index+SIGNATURE_SIZE])
	last_index += SIGNATURE_SIZE

	return nil
}

func (pl *CorePayload) Hash() []byte {
	bytes := pl.Marshal()
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}
