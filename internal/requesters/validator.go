package requesters

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"voting-blockchain/internal/models/types"

	"github.com/google/uuid"
)

type Validator struct {
	host string `json:"host"`
	port uint16 `json:"port"`
}

func NewValidator(host string, port uint16) *Validator {
	return &Validator{
		host: host,
		port: port,
	}
}

func (vl Validator) getPath() string {
	return "http://" + vl.host + ":" + fmt.Sprint(vl.port) + "/validate"
}

func (vl Validator) Validate(userUuid types.Uuid, votingUuid types.Uuid, data []byte, signature types.Signature) (bool, error) {
	type RequestBody struct {
		UserUuid        string `json:"user_uuid"`
		VotingUuid      string `json:"voting_uuid"`
		DataBase64      string `json:"data_base64"`
		SignatureBase64 string `json:"signature_base64"`
	}

	body_bytes, err := json.Marshal(RequestBody{
		UserUuid:        uuid.UUID(userUuid).String(),
		VotingUuid:      uuid.UUID(votingUuid).String(),
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

	fmt.Printf("userUuid = %s, StatusCode = %d", uuid.UUID(userUuid).String(), resp.StatusCode)
	return (resp.StatusCode == 200), nil
}
