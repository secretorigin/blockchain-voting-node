package models

import "github.com/google/uuid"

const VOTE_SIZE = UUID_SIZE + UUID_SIZE + SIGNATURE_SIZE

type Vote struct {
	UserUuid   uuid.UUID `json:"user_uuid"`
	OptionUuid uuid.UUID `json:"option_uuid"`
	Signature  Signature `json:"signature"`
}

func (vt Vote) Size() uint64 {
	return VOTE_SIZE
}

func (vt Vote) Marshal() []byte {
	bytes := make([]byte, vt.Size())
	last_index := uint64(0)

	// user uuid
	copy(bytes[last_index:last_index+UUID_SIZE], vt.UserUuid[:])
	last_index += UUID_SIZE
	// option uuid
	copy(bytes[last_index:last_index+UUID_SIZE], vt.OptionUuid[:])
	last_index += UUID_SIZE
	// signature
	copy(bytes[last_index:last_index+SIGNATURE_SIZE], vt.Signature[:])
	last_index += SIGNATURE_SIZE

	return bytes
}

func (vt *Vote) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	copy(vt.UserUuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// option uuid
	copy(vt.OptionUuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// signature
	copy(vt.Signature[:], bytes[last_index:last_index+SIGNATURE_SIZE])
	last_index += SIGNATURE_SIZE

	return nil
}
