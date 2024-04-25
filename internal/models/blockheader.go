package models

import (
	"encoding/binary"
	"unsafe"

	"github.com/google/uuid"
)

const HASH_SIZE = 32
const MERKLE_ROOT_SIZE = 32
const UUID_SIZE = 16

/*
	Block Header structure
	1) Id - 8 bytes
	2) Epoch - 8 bytes
	3) Type - 1 byte
	4) PrevHash - 32 bytes
	5) MerkleRoot - 32 bytes
	6) Voters - len * 16 bytes
	7) OptionsMeta - len * (8+8)
	8) Signature - 72 bytes
*/

type BlockHeader struct {
	Id          uint64                 `json:"id"`
	Epoch       uint64                 `json:"epoch"`
	Type        BlockType              `json:"type"`
	PrevHash    [HASH_SIZE]byte        `json:"prev_hash"`
	MerkleRoot  [MERKLE_ROOT_SIZE]byte `json:"merkle_root"`
	OptionsMeta []OptionMeta           `json:"options_meta"`
	Voters      []uuid.UUID            `json:"voters"`
	Signature   Signature              `json:"signature"`

	OptionsCount uint8 `json:"-"`
}

func (bl BlockHeader) Size() uint64 {
	return uint64(
		uint64(unsafe.Sizeof(bl.Id)) +
			uint64(unsafe.Sizeof(bl.Epoch)) +
			uint64(unsafe.Sizeof(bl.Type)) +
			HASH_SIZE +
			MERKLE_ROOT_SIZE +
			uint64(len(bl.OptionsMeta)*OPTIONS_META_SIZE) +
			8 + uint64(len(bl.Voters)*UUID_SIZE) +
			SIGNATURE_SIZE)
}

func (bl BlockHeader) Marshal() []byte {
	bytes := make([]byte, bl.Size())
	last_index := uint64(0)

	// id
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(bl.Id))
	last_index += 8
	// epoch
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(bl.Epoch))
	last_index += 8
	// type
	bytes[last_index] = uint8(bl.Type)
	last_index += 1
	// prev hash
	copy(bytes[last_index:last_index+HASH_SIZE], bl.PrevHash[:])
	last_index += HASH_SIZE
	// merkle root
	copy(bytes[last_index:last_index+MERKLE_ROOT_SIZE], bl.MerkleRoot[:])
	last_index += MERKLE_ROOT_SIZE
	// options meta
	for i := 0; i < len(bl.OptionsMeta); i++ {
		meta_bytes := bl.OptionsMeta[i].Marshal()
		copy(bytes[last_index:last_index+OPTIONS_META_SIZE], meta_bytes[:])
		last_index += OPTIONS_META_SIZE
	}
	// new voters count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(bl.Voters)))
	last_index += 8
	// new voters
	for i := 0; i < len(bl.Voters); i++ {
		copy(bytes[last_index:last_index+UUID_SIZE], bl.Voters[i][:])
		last_index += UUID_SIZE
	}
	// signature
	copy(bytes[last_index:last_index+SIGNATURE_SIZE], bl.Signature[:])
	last_index += SIGNATURE_SIZE

	return bytes
}

func (bl *BlockHeader) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// id
	bl.Id = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// epoch
	bl.Epoch = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// type
	bl.Type = BlockType(bytes[last_index])
	last_index += 1
	// prev hash
	copy(bl.PrevHash[:], bytes[last_index:last_index+HASH_SIZE])
	last_index += HASH_SIZE
	// merkle root
	copy(bl.MerkleRoot[:], bytes[last_index:last_index+MERKLE_ROOT_SIZE])
	last_index += MERKLE_ROOT_SIZE
	// options meta
	for i := uint8(0); i < bl.OptionsCount; i++ {
		var meta OptionMeta
		_ = meta.Unmarshal([24]byte(bytes[last_index : last_index+OPTIONS_META_SIZE]))
		bl.OptionsMeta = append(bl.OptionsMeta, meta)
		last_index += OPTIONS_META_SIZE
	}
	// new voters count
	votersCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	for i := uint64(0); i < votersCount; i++ {
		var uuid uuid.UUID
		copy(uuid[:], bytes[last_index:last_index+UUID_SIZE])
		last_index += UUID_SIZE
		bl.Voters = append(bl.Voters, uuid)
	}
	// signature
	copy(bl.Signature[:], bytes[last_index:last_index+SIGNATURE_SIZE])
	last_index += SIGNATURE_SIZE

	return nil
}
