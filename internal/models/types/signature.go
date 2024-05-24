package types

const SIGNATURE_SIZE uint64 = 72

type Signature [SIGNATURE_SIZE]byte

func (sign Signature) Size() uint64 {
	return SIGNATURE_SIZE
}
