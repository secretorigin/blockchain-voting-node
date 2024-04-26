package models

const PORT_SIZE = 2

const NODE_SIZE = NODEMETA_SIZE + SIGNATURE_SIZE

type Node struct {
	NodeMeta
	Signature Signature `json:"signature"`
}

func (nd Node) Size() uint64 {
	return NODE_SIZE
}

func (nd Node) Marshal() []byte {
	bytes := make([]byte, nd.Size())
	last_index := uint64(0)

	// node meta
	meta_bytes := nd.NodeMeta.Marshal()
	copy(bytes[last_index:last_index+NODEMETA_SIZE], meta_bytes[:])
	last_index += NODEMETA_SIZE
	// signature
	copy(bytes[last_index:last_index+SIGNATURE_SIZE], nd.Signature[:])
	last_index += SIGNATURE_SIZE

	return bytes
}

func (nd *Node) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	_ = nd.NodeMeta.Unmarshal(bytes)
	last_index += NODEMETA_SIZE
	// signature
	copy(nd.Signature[:], bytes[last_index:last_index+SIGNATURE_SIZE])
	last_index += SIGNATURE_SIZE

	return nil
}
