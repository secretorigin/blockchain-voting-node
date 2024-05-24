package blockchains

import (
	"encoding/binary"
	"voting-blockchain/internal/models/types"
)

const NODE_META_SIZE uint64 = types.UUID_SIZE + types.IP_SIZE + types.PORT_SIZE

type NodeMeta struct {
	Uuid types.Uuid `json:"uuid"` // = user uuid
	Ip   types.Ip   `json:"host"`
	Port types.Port `json:"port"`
}

func (nd NodeMeta) Size() uint64 {
	return NODE_META_SIZE
}

func (nd NodeMeta) Marshal() []byte {
	bytes := make([]byte, nd.Size())
	last_index := uint64(0)

	// user uuid
	copy(bytes[last_index:last_index+types.UUID_SIZE], nd.Uuid[:])
	last_index += types.UUID_SIZE
	// host
	ip_bytes := nd.Ip.Marshal()
	copy(bytes[last_index:last_index+types.IP_SIZE], ip_bytes[:])
	last_index += types.IP_SIZE
	// port
	binary.LittleEndian.PutUint16(bytes[last_index:last_index+types.PORT_SIZE], uint16(nd.Port))
	last_index += types.PORT_SIZE

	return bytes
}

func (nd *NodeMeta) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// user uuid
	copy(nd.Uuid[:], bytes[last_index:last_index+types.UUID_SIZE])
	last_index += types.UUID_SIZE
	// host
	_ = nd.Ip.Unmarshal(bytes[last_index : last_index+types.IP_SIZE])
	last_index += types.IP_SIZE
	// port
	nd.Port = types.Port(binary.LittleEndian.Uint16(bytes[last_index : last_index+types.PORT_SIZE]))
	last_index += types.PORT_SIZE

	return nil
}
