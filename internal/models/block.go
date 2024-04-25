package models

type Block struct {
	BlockHeader
	Payload BlockPayloadInterface
}

type BlockPayloadInterface interface {
	ByteForm

	Hash() []byte
}
