package votings

import (
	"encoding/binary"
	"errors"
	"voting-blockchain/internal/models/types"
)

const OPTIONS_META_SIZE uint64 = types.UUID_SIZE + 8

type OptionMeta struct {
	Uuid    types.Uuid
	Counter uint64
}

func (meta OptionMeta) Size() uint64 {
	return types.UUID_SIZE + 8
}

func (meta OptionMeta) Marshal() []byte {
	var bytes []byte = make([]byte, meta.Size())
	last_index := uint64(0)

	copy(bytes[last_index:last_index+types.UUID_SIZE], meta.Uuid[:])
	last_index += types.UUID_SIZE
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(meta.Counter))
	last_index += 8

	return bytes
}

func (meta *OptionMeta) Unmarshal(bytes []byte) error {
	if uint64(len(bytes)) == meta.Size() {
		return errors.New("invalid option meta length")
	}
	last_index := uint64(0)

	// uuid
	copy(meta.Uuid[:], bytes[last_index:last_index+types.UUID_SIZE])
	last_index += types.UUID_SIZE
	// counter
	meta.Counter = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8

	return nil
}
