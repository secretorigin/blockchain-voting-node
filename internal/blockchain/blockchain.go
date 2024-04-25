package blockchain

import (
	"voting-blockchain/internal/models"

	"github.com/google/uuid"
)

type Blockchain struct {
	Validator struct{}
	Storage   StorageInterface
}

type ValidatorInterface interface {
	CheckData(userUuid uuid.UUID, bytes []byte) (bool, error)
}

type StorageInterface interface {
	GetVoting() (models.Voting, error)
	IsUserAlreadyIn() (bool, error)
	IsNodeAlreadyIn() (bool, error)
	NodesCount()
	GetNodeBySerial(serial uint64) (models.NodeMeta, error)
}
