package models

import (
	"encoding/binary"
)

const PORT_SIZE = 2

const NODE_SIZE = UUID_SIZE + IP_SIZE + PORT_SIZE + SIGNATURE_SIZE

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

	// user uuid
	copy(bytes[last_index:last_index+UUID_SIZE], nd.Uuid[:])
	last_index += UUID_SIZE
	// host
	ip_bytes := nd.Host.Marshal()
	copy(bytes[last_index:last_index+IP_SIZE], ip_bytes[:])
	last_index += IP_SIZE
	// port
	binary.LittleEndian.PutUint16(bytes[last_index:last_index+PORT_SIZE], uint16(nd.Port))
	last_index += PORT_SIZE
	// signature
	copy(bytes[last_index:last_index+SIGNATURE_SIZE], nd.Signature[:])
	last_index += SIGNATURE_SIZE

	return bytes
}

func (nd *Node) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	copy(nd.Uuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// host
	_ = nd.Host.Unmarshal(bytes[last_index : last_index+IP_SIZE])
	last_index += IP_SIZE
	// port
	nd.Port = binary.LittleEndian.Uint16(bytes[last_index : last_index+PORT_SIZE])
	last_index += PORT_SIZE
	// signature
	copy(nd.Signature[:], bytes[last_index:last_index+SIGNATURE_SIZE])
	last_index += SIGNATURE_SIZE

	return nil
}
