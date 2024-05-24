package validators

import (
	"fmt"
	"voting-blockchain/internal/models/blocks/payloads"
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/requesters"
)

type CorePayloadValidator struct {
	validator *requesters.Validator
}

func NewCorePayloadValidator(validator *requesters.Validator) *CorePayloadValidator {
	return &CorePayloadValidator{
		validator: validator,
	}
}

func (vl CorePayloadValidator) ValidateCorePayload(payload payloads.CorePayload) (bool, error) {
	if len(payload.Voting.Options) == 0 {
		return false, fmt.Errorf("options count = 0")
	}

	if payload.Voting.CycleDuration < payload.Voting.SendingDuration {
		return false, fmt.Errorf("CycleDuration < SendingDuration")
	}

	bytes := payload.Marshal()
	payload_size := payload.Size()

	// signature
	ok, err := vl.validator.Validate(
		payload.UserUuid,
		payload.Voting.Uuid,
		bytes[:payload_size-types.SIGNATURE_SIZE-1],
		types.Signature(bytes[payload_size-types.SIGNATURE_SIZE:]),
	)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("core payload signature")
	}

	return ok, nil
}
