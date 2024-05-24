package types

const PORT_SIZE uint64 = 2

type Port uint16

func (port Port) Size() uint64 {
	return PORT_SIZE
}
