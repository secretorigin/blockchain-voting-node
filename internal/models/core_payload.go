package models

import (
	"crypto/sha256"
	"encoding/binary"
)

type CorePayload struct {
	Voting          Voting    `json:"voting"`
	CycleDuration   uint64    `json:"cycle_duration"`   // in seconds
	SendingDuration uint64    `json:"sending_duration"` // in seconds
	Signature       Signature `json:"signature"`
}

func (pl CorePayload) Size() uint64 {
	return pl.Voting.Size() + uint64(8) + uint64(8) + SIGNATURE_SIZE
}

func (pl CorePayload) Marshal() []byte {
	bytes := make([]byte, pl.Size())
	last_index := uint64(0)

	// voting
	voting_bytes := pl.Voting.Marshal()
	copy(bytes[last_index:last_index+pl.Voting.Size()], voting_bytes)
	last_index += pl.Voting.Size()
	// cycle duration
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], pl.CycleDuration)
	last_index += 8
	// sending duration
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], pl.SendingDuration)
	last_index += 8
	// signature
	copy(bytes[last_index:last_index+SIGNATURE_SIZE], pl.Signature[:])
	last_index += SIGNATURE_SIZE

	return bytes
}

func (pl *CorePayload) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// voting
	_ = pl.Voting.Unmarshal(bytes)
	last_index += pl.Voting.Size()
	// cycle duration
	pl.CycleDuration = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// sending duration
	pl.SendingDuration = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
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
