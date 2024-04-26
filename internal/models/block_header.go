package models

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/google/uuid"
)

const HASH_SIZE uint64 = 32
const MERKLE_ROOT_SIZE uint64 = 32
const UUID_SIZE uint64 = 16

/*
	Block Header structure
	1) Id - 8 bytes
	2) Epoch - 8 bytes
	3) Type - 1 byte
	4) PrevHash - 32 bytes
	5) MerkleRoot - 32 bytes
	6) Votes - len * 16 bytes
	7) Nodes - len * 22 bytes
	8) OptionsMeta - len * (8+8)
	9) UserUuid - 16 bytes
	9) Signature - 72 bytes
*/

type BlockHeader struct {
	Id          uint64                 `json:"id"`
	Epoch       uint64                 `json:"epoch"`
	Type        BlockType              `json:"type"`
	PrevHash    [HASH_SIZE]byte        `json:"prev_hash"`
	MerkleRoot  [MERKLE_ROOT_SIZE]byte `json:"merkle_root"`
	OptionsMeta []OptionMeta           `json:"options_meta"`
	Votes       []uuid.UUID            `json:"votes"`
	Nodes       []NodeMeta             `json:"nodes"`
	UserUuid    uuid.UUID              `json:"user_uuid"`
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
			uint64(len(bl.OptionsMeta))*OPTIONS_META_SIZE +
			8 + uint64(len(bl.Votes))*UUID_SIZE +
			8 + uint64(len(bl.Nodes))*NODEMETA_SIZE +
			UUID_SIZE +
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
	fmt.Println(last_index)
	// options meta
	for i := 0; i < len(bl.OptionsMeta); i++ {
		meta_bytes := bl.OptionsMeta[i].Marshal()
		copy(bytes[last_index:last_index+OPTIONS_META_SIZE], meta_bytes[:])
		last_index += OPTIONS_META_SIZE
	}
	// new votes count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(bl.Votes)))
	last_index += 8
	// new votes
	for i := 0; i < len(bl.Votes); i++ {
		copy(bytes[last_index:last_index+UUID_SIZE], bl.Votes[i][:])
		last_index += UUID_SIZE
	}
	// new nodes count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(bl.Nodes)))
	last_index += 8
	// new nodes
	for i := 0; i < len(bl.Nodes); i++ {
		meta_bytes := bl.Nodes[i].Marshal()
		copy(bytes[last_index:last_index+NODEMETA_SIZE], meta_bytes[:])
		last_index += NODEMETA_SIZE
	}
	// user uuid
	copy(bytes[last_index:last_index+UUID_SIZE], bl.UserUuid[:])
	last_index += UUID_SIZE
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
	// new votes count
	votesCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// new votes
	for i := uint64(0); i < votesCount; i++ {
		var uuid uuid.UUID
		copy(uuid[:], bytes[last_index:last_index+UUID_SIZE])
		last_index += UUID_SIZE
		bl.Votes = append(bl.Votes, uuid)
	}
	// new nodes count
	nodesCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// new nodes
	for i := uint64(0); i < nodesCount; i++ {
		var meta NodeMeta
		_ = meta.Unmarshal([]byte(bytes[last_index : last_index+NODEMETA_SIZE]))
		bl.Nodes = append(bl.Nodes, meta)
		last_index += NODEMETA_SIZE
	}
	// user uuid
	copy(bl.UserUuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// signature
	copy(bl.Signature[:], bytes[last_index:last_index+SIGNATURE_SIZE])
	last_index += SIGNATURE_SIZE

	return nil
}

func (bl *BlockHeader) Hash() []byte {
	bytes := bl.Marshal()
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}
