package validator

import "github.com/google/uuid"

type Storage struct {
	Users map[uuid.UUID]bool
	Nodes map[uuid.UUID]bool
}

func NewStorage() *Storage {
	return &Storage{
		Users: make(map[uuid.UUID]bool),
		Nodes: make(map[uuid.UUID]bool),
	}
}

func (st Storage) IsUserAlreadyIn(userUuid uuid.UUID) (bool, error) {
	return st.Users[userUuid], nil
}

func (st Storage) IsNodeAlreadyIn(nodeUuid uuid.UUID) (bool, error) {
	return st.Nodes[nodeUuid], nil
}
