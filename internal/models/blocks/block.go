package blocks

import (
	"voting-blockchain/internal/models"
	"voting-blockchain/internal/utils"
)

type Block struct {
	Header  BlockHeader
	Payload BlockPayloadInterface
}

type BlockPayloadInterface interface {
	models.ByteForm

	Hash() []byte
}

func (bl *Block) Hash() []byte {
	bytes := bl.Header.Marshal()
	return utils.GetHash(bytes)
}
