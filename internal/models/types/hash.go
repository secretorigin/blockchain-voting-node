package types

const HASH_SIZE uint64 = 32

type BlockHash [HASH_SIZE]byte

func (hash BlockHash) Size() uint64 {
	return HASH_SIZE
}
