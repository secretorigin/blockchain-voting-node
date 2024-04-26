package blockchain

import (
	"voting-blockchain/internal/models"
	"voting-blockchain/internal/validators"

	"github.com/google/uuid"
)

type Blockchain struct {
	HeaderValidator validators.BlockHeaderValidator
}

type ValidatorInterface interface {
	CheckData(userUuid uuid.UUID, bytes []byte) (bool, error)
}

type StorageInterface interface {
	GetVoting() (models.Voting, error)
	NodesCount()
	GetNodeBySerial(serial uint64) (models.NodeMeta, error)
}
