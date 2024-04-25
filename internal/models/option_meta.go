package models

import (
	"encoding/binary"

	"github.com/google/uuid"
)

const OPTIONS_META_SIZE = UUID_SIZE + 8

type OptionMeta struct {
	Uuid    uuid.UUID
	Counter uint64
}

func (meta OptionMeta) Marshal() [OPTIONS_META_SIZE]byte {
	var bytes [OPTIONS_META_SIZE]byte
	last_index := 0

	copy(bytes[last_index:last_index+UUID_SIZE], meta.Uuid[:])
	last_index += UUID_SIZE
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(meta.Counter))
	last_index += 8

	return bytes
}

func (meta *OptionMeta) Unmarshal(bytes [OPTIONS_META_SIZE]byte) error {
	last_index := 0

	// uuid
	copy(meta.Uuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// counter
	meta.Counter = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8

	return nil
}
