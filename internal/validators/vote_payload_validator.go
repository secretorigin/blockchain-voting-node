package validators

import (
	"fmt"
	"voting-blockchain/internal/models"
	"voting-blockchain/internal/requesters"

	"github.com/google/uuid"
)

type VotePaylaodValidatorStorageInterface interface {
	IsUserAlreadyIn(userUuid uuid.UUID) (bool, error)
	IsNodeAlreadyIn(nodeUuid uuid.UUID) (bool, error)
}

type VotePayloadValidator struct {
	validator *requesters.Validator
	storage   VotePaylaodValidatorStorageInterface
}

func NewVotePayloadValidator(validator *requesters.Validator, storage VotePaylaodValidatorStorageInterface) *VotePayloadValidator {
	return &VotePayloadValidator{
		validator: validator,
		storage:   storage,
	}
}

func (vl VotePayloadValidator) ValidateVote(vote models.Vote, votingUuid uuid.UUID) (bool, error) {
	bytes := vote.Marshal()

	ok, err := vl.validator.Validate(
		vote.UserUuid,
		votingUuid,
		bytes[:models.VOTE_SIZE-models.SIGNATURE_SIZE-1],
		models.Signature(bytes[models.VOTE_SIZE-models.SIGNATURE_SIZE:]),
	)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("vote %s signature", vote.UserUuid)
	}

	return ok, nil
}

func (vl VotePayloadValidator) ValidateNode(node models.Node, votingUuid uuid.UUID) (bool, error) {
	bytes := node.Marshal()

	ok, err := vl.validator.Validate(
		node.Uuid,
		votingUuid,
		bytes[:models.VOTE_SIZE-models.SIGNATURE_SIZE-1],
		models.Signature(bytes[models.VOTE_SIZE-models.SIGNATURE_SIZE:]),
	)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("node %s signature", node.Uuid)
	}

	return ok, nil
}

func (vl *VotePayloadValidator) ValidateVotePayload(payload models.VotePayload, votingUuid uuid.UUID) (bool, error) {
	// votes
	votes := make(map[uuid.UUID]bool)
	for i := 0; i < len(payload.Votes); i++ {
		ok, err := vl.ValidateVote(payload.Votes[i], votingUuid)
		if !ok || err != nil {
			return ok, err
		}
		// check existence
		userUuid := payload.Votes[i].UserUuid
		if votes[userUuid] {
			return false, fmt.Errorf("duplicated user %s", userUuid)
		}
		ok, err = vl.storage.IsUserAlreadyIn(userUuid)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, fmt.Errorf("duplicated user (storage) %s", userUuid)
		}
		// save
		votes[userUuid] = true
	}

	// nodes
	nodes := make(map[uuid.UUID]bool)
	for i := 0; i < len(payload.Nodes); i++ {
		ok, err := vl.ValidateNode(payload.Nodes[i], votingUuid)
		if !ok || err != nil {
			return ok, err
		}
		// check existence
		nodeUuid := payload.Nodes[i].Uuid
		if nodes[nodeUuid] {
			return false, fmt.Errorf("duplicated node %s", nodeUuid)
		}
		ok, err = vl.storage.IsNodeAlreadyIn(nodeUuid)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, fmt.Errorf("duplicated node (storage) %s", nodeUuid)
		}
		// save
		nodes[nodeUuid] = true
	}

	return true, nil
}
