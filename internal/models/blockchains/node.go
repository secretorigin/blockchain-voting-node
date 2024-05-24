package blockchains

import (
	"voting-blockchain/internal/models/types"
)

const NODE_SIZE uint64 = NODE_META_SIZE + types.SIGNATURE_SIZE

type Node struct {
	NodeMeta
	Signature types.Signature `json:"signature"`
}

func (nd Node) Size() uint64 {
	return NODE_META_SIZE + types.SIGNATURE_SIZE
}

func (nd Node) Marshal() []byte {
	bytes := make([]byte, nd.Size())
	last_index := uint64(0)

	// node meta
	meta_bytes := nd.NodeMeta.Marshal()
	copy(bytes[last_index:last_index+NODE_META_SIZE], meta_bytes[:])
	last_index += NODE_META_SIZE
	// signature
	copy(bytes[last_index:last_index+types.SIGNATURE_SIZE], nd.Signature[:])
	last_index += types.SIGNATURE_SIZE

	return bytes
}

func (nd *Node) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	_ = nd.NodeMeta.Unmarshal(bytes)
	last_index += NODE_META_SIZE
	// signature
	copy(nd.Signature[:], bytes[last_index:last_index+types.SIGNATURE_SIZE])
	last_index += types.SIGNATURE_SIZE

	return nil
}
