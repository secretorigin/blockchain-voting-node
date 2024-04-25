package validator

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"voting-blockchain/internal/models"

	"github.com/google/uuid"
)

type Validator struct {
	host    string `json:"host"`
	port    uint16 `json:"port"`
	storage StorageInterface
}

func NewValidator(host string, port uint16, storage StorageInterface) *Validator {
	return &Validator{
		host:    host,
		port:    port,
		storage: storage,
	}
}

type StorageInterface interface {
	IsUserAlreadyIn(userUuid uuid.UUID) (bool, error)
	IsNodeAlreadyIn(nodeUuid uuid.UUID) (bool, error)
}

func (vl Validator) getPath() string {
	return "http://" + vl.host + ":" + string(vl.port) + "/validate"
}

func (vl Validator) validateBytesInValidator(userUuid uuid.UUID, votingUuid uuid.UUID, data []byte, signature models.Signature) (bool, error) {
	type RequestBody struct {
		UserUuid        string `json:"user_uuid"`
		VotingUuid      string `json:"voting_uuid"`
		DataBase64      string `json:"data_base64"`
		SignatureBase64 string `json:"signature_base64"`
	}

	body_bytes, err := json.Marshal(RequestBody{
		UserUuid:        userUuid.String(),
		VotingUuid:      votingUuid.String(),
		DataBase64:      base64.StdEncoding.EncodeToString(data),
		SignatureBase64: base64.StdEncoding.EncodeToString(signature[:]),
	})
	if err != nil {
		return false, err
	}

	resp, err := http.Post(vl.getPath(), "application/json", bytes.NewBuffer(body_bytes))
	if err != nil {
		return false, err
	}

	fmt.Printf("userUuid = %s, StatusCode = %d", userUuid.String(), resp.StatusCode)
	return (resp.StatusCode == 200), nil
}

func (vl Validator) CorePayload(payload models.CorePayload) (bool, error) {
	if len(payload.Voting.Options) == 0 {
		return false, fmt.Errorf("options count = 0")
	}

	if payload.CycleDuration < payload.SendingDuration {
		return false, fmt.Errorf("CycleDuration < SendingDuration")
	}

	bytes := payload.Marshal()
	payload_size := payload.Size()

	// signature
	ok, err := vl.validateBytesInValidator(
		payload.UserUuid,
		payload.Voting.Uuid,
		bytes[:payload_size-models.SIGNATURE_SIZE-1],
		models.Signature(bytes[payload_size-models.SIGNATURE_SIZE:]),
	)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("core payload signature")
	}

	return ok, nil
}

func (vl Validator) Vote(vote models.Vote, votingUuid uuid.UUID) (bool, error) {
	bytes := vote.Marshal()

	ok, err := vl.validateBytesInValidator(
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

func (vl Validator) Node(node models.Node, votingUuid uuid.UUID) (bool, error) {
	bytes := node.Marshal()

	ok, err := vl.validateBytesInValidator(
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

func (vl Validator) VotePayload(payload models.VotePayload, votingUuid uuid.UUID) (bool, error) {
	// votes
	votes := make(map[uuid.UUID]bool)
	for i := 0; i < len(payload.Votes); i++ {
		ok, err := vl.Vote(payload.Votes[i], votingUuid)
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
		ok, err := vl.Node(payload.Nodes[i], votingUuid)
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
