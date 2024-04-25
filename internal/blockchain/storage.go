package blockchain

type ValidatorStorage struct {
	Users map[uuid.UUID]bool
	Nodes map[uuid.UUID]bool
}

func NewStorage() {
	return 
} 

IsUserAlreadyIn(userUuid uuid.UUID) (bool, error)
IsNodeAlreadyIn(nodeUuid uuid.UUID) (bool, error)