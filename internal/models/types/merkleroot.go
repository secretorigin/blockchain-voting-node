package types

const MERKLE_ROOT_SIZE uint64 = 32

type BlockMerkleRoot [MERKLE_ROOT_SIZE]byte

func (hash BlockMerkleRoot) Size() uint64 {
	return MERKLE_ROOT_SIZE
}
