package models

type ByteForm interface {
	Size() uint64
	Marshal() []byte
	Unmarshal(bytes []byte) error
}
