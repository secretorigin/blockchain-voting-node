package blocks

import (
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"

	"voting-blockchain/internal/models/blockchains"
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/models/votings"
	"voting-blockchain/internal/utils"
)

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
	Epoch       int64                  `json:"epoch"`
	Type        BlockType              `json:"type"`
	PrevHash    types.BlockHash        `json:"prev_hash"`
	MerkleRoot  types.BlockMerkleRoot  `json:"merkle_root"`
	OptionsMeta []votings.OptionMeta   `json:"options_meta"`
	Votes       []types.Uuid           `json:"votes"`
	Nodes       []blockchains.NodeMeta `json:"nodes"`
	UserUuid    types.Uuid             `json:"user_uuid"`
	Signature   types.Signature        `json:"signature"`

	OptionsCount uint8 `json:"-"`
}

func (bl BlockHeader) Size() uint64 {
	return uint64(
		uint64(unsafe.Sizeof(bl.Id)) +
			uint64(unsafe.Sizeof(bl.Epoch)) +
			uint64(unsafe.Sizeof(bl.Type)) +
			types.HASH_SIZE +
			types.MERKLE_ROOT_SIZE +
			uint64(len(bl.OptionsMeta))*votings.OPTIONS_META_SIZE +
			8 + uint64(len(bl.Votes))*types.UUID_SIZE +
			8 + uint64(len(bl.Nodes))*blockchains.NODE_META_SIZE +
			types.UUID_SIZE +
			types.SIGNATURE_SIZE)
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
	copy(bytes[last_index:last_index+types.HASH_SIZE], bl.PrevHash[:])
	last_index += types.HASH_SIZE
	// merkle root
	copy(bytes[last_index:last_index+types.MERKLE_ROOT_SIZE], bl.MerkleRoot[:])
	last_index += types.MERKLE_ROOT_SIZE
	fmt.Println(last_index)
	// options meta
	for i := 0; i < len(bl.OptionsMeta); i++ {
		meta_bytes := bl.OptionsMeta[i].Marshal()
		copy(bytes[last_index:last_index+votings.OPTIONS_META_SIZE], meta_bytes[:])
		last_index += votings.OPTIONS_META_SIZE
	}
	// new votes count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(bl.Votes)))
	last_index += 8
	// new votes
	for i := 0; i < len(bl.Votes); i++ {
		copy(bytes[last_index:last_index+types.UUID_SIZE], bl.Votes[i][:])
		last_index += types.UUID_SIZE
	}
	// new nodes count
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(bl.Nodes)))
	last_index += 8
	// new nodes
	for i := 0; i < len(bl.Nodes); i++ {
		meta_bytes := bl.Nodes[i].Marshal()
		copy(bytes[last_index:last_index+blockchains.NODE_META_SIZE], meta_bytes[:])
		last_index += blockchains.NODE_META_SIZE
	}
	// user uuid
	copy(bytes[last_index:last_index+types.UUID_SIZE], bl.UserUuid[:])
	last_index += types.UUID_SIZE
	// signature
	copy(bytes[last_index:last_index+types.SIGNATURE_SIZE], bl.Signature[:])
	last_index += types.SIGNATURE_SIZE

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
	copy(bl.PrevHash[:], bytes[last_index:last_index+types.HASH_SIZE])
	last_index += types.HASH_SIZE
	// merkle root
	copy(bl.MerkleRoot[:], bytes[last_index:last_index+types.MERKLE_ROOT_SIZE])
	last_index += types.MERKLE_ROOT_SIZE
	// options meta
	for i := uint8(0); i < bl.OptionsCount; i++ {
		var meta votings.OptionMeta
		_ = meta.Unmarshal([]byte(bytes[last_index : last_index+votings.OPTIONS_META_SIZE]))
		bl.OptionsMeta = append(bl.OptionsMeta, meta)
		last_index += votings.OPTIONS_META_SIZE
	}
	// new votes count
	votesCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// new votes
	for i := uint64(0); i < votesCount; i++ {
		var uuid types.Uuid
		copy(uuid[:], bytes[last_index:last_index+types.UUID_SIZE])
		last_index += types.UUID_SIZE
		bl.Votes = append(bl.Votes, uuid)
	}
	// new nodes count
	nodesCount := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// new nodes
	for i := uint64(0); i < nodesCount; i++ {
		var meta blockchains.NodeMeta
		_ = meta.Unmarshal([]byte(bytes[last_index : last_index+blockchains.NODE_META_SIZE]))
		bl.Nodes = append(bl.Nodes, meta)
		last_index += blockchains.NODE_META_SIZE
	}
	// user uuid
	copy(bl.UserUuid[:], bytes[last_index:last_index+types.UUID_SIZE])
	last_index += types.UUID_SIZE
	// signature
	copy(bl.Signature[:], bytes[last_index:last_index+types.SIGNATURE_SIZE])
	last_index += types.SIGNATURE_SIZE

	return nil
}

func (bl *BlockHeader) Hash() []byte {
	bytes := bl.Marshal()
	return utils.GetHash(bytes)
}

func CreataNewBlockHeader(prev BlockHeader, merkle_root []byte) BlockHeader {

	options_meta := prev.OptionsMeta[]

	header := blocks.BlockHeader{
		Id: prev.Id + 1,
	Epoch: time.Now().Unix(), 
	Type: BLOCK_TYPE_VOTE,
	PrevHash: prev.Hash(),
	MerkleRoot: merkle_root,
	OptionsMeta: 
	Votes:       
	Nodes:      
	UserUuid:   
	Signature:  
	}
}
