package models

import (
	"encoding/binary"

	"github.com/google/uuid"
)

const NODEMETA_SIZE = UUID_SIZE + IP_SIZE + PORT_SIZE

type NodeMeta struct {
	Uuid uuid.UUID `json:"uuid"` // = user uuid
	Host IP        `json:"host"`
	Port uint16    `json:"port"`
}

func (nd NodeMeta) Size() uint64 {
	return NODEMETA_SIZE
}

func (nd NodeMeta) Marshal() []byte {
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

	return bytes
}

func (nd *NodeMeta) Unmarshal(bytes []byte) error {
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

	return nil
}
