package models

import (
	"encoding/binary"

	"github.com/google/uuid"
)

type Option struct {
	Uuid  uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
}

func (op Option) Size() uint64 {
	return UUID_SIZE + uint64(8) + uint64(len(op.Title))
}

func (op Option) Marshal() []byte {
	bytes := make([]byte, op.Size())
	last_index := uint64(0)

	// uuid
	copy(bytes[last_index:last_index+UUID_SIZE], op.Uuid[:])
	last_index += UUID_SIZE
	// title length
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(op.Title)))
	last_index += 8
	// title
	copy(bytes[last_index:last_index+uint64(len(op.Title))], op.Title[:])
	last_index += uint64(len(op.Title))

	return bytes
}

func (op *Option) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// uuid3]
	copy(op.Uuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// title legth
	title_length := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// title
	op.Title = string(bytes[last_index : last_index+title_length])
	last_index += title_length

	return nil
}
