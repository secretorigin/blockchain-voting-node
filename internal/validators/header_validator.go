package validators

import (
	"fmt"
	"voting-blockchain/internal/models"
	"voting-blockchain/internal/requesters"
	"voting-blockchain/internal/utils"

	"github.com/google/uuid"
)

type BlockHeaderValidatorStorageInterface interface {
	GetLastBlockHeader() models.BlockHeader
	GetNodesCount() uint64
	GetNodeBySerial(uint64) models.NodeMeta
}

type BlockHeaderValidator struct {
	validator            *requesters.Validator
	voting               *models.Voting
	corePayloadValidator *CorePayloadValidator
	votePayloadValidator *VotePayloadValidator
	storage              BlockHeaderValidatorStorageInterface
}

func NewBlockValidator(
	validator *requesters.Validator,
	voting *models.Voting,
	corePayloadValidator *CorePayloadValidator,
	votePayloadValidator *VotePayloadValidator,
	storage BlockHeaderValidatorStorageInterface,
) *BlockHeaderValidator {
	return &BlockHeaderValidator{
		validator:            validator,
		voting:               voting,
		corePayloadValidator: corePayloadValidator,
		votePayloadValidator: votePayloadValidator,
		storage:              storage,
	}
}

func (vl BlockHeaderValidator) Validate(header models.BlockHeader, users []uuid.UUID, votingUuid uuid.UUID) (bool, error) {
	last_block_header := vl.storage.GetLastBlockHeader()
	last_block_hash := last_block_header.Hash()

	// check id
	if header.Id != last_block_header.Id+1 {
		return false, fmt.Errorf("Block id is not valid")
	}

	// check prev hash
	if utils.SliceEq(last_block_hash, header.PrevHash[:]) {
		return false, fmt.Errorf("PrevHash is not valid")
	}

	// check id
	if header.Epoch != last_block_header.Epoch+vl.voting.CycleDuration {
		return false, fmt.Errorf("Block epoch is not valid")
	}

	serial := utils.CalculateNextBlockCreator(last_block_hash, vl.storage.GetNodesCount())
	creator := vl.storage.GetNodeBySerial(serial)

	// check creator
	if creator.Uuid != header.UserUuid {
		return false, fmt.Errorf("new block creator must be %s, but it is %s", creator.Uuid, header.UserUuid)
	}

	// check signature
	header_bytes := header.Marshal()
	ok, err := vl.validator.Validate(header.UserUuid, votingUuid, header_bytes[:header.Size()-models.SIGNATURE_SIZE-1], models.Signature(header_bytes[header.Size()-models.SIGNATURE_SIZE:]))
	if !ok || err != nil {
		return ok, err
	}

	return true, nil
}
