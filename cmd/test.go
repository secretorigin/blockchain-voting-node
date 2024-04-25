package main

import (
	"fmt"
	"voting-blockchain/internal/models"

	"github.com/google/uuid"
)

func gen(size int) []byte {
	var array []byte
	for i := 0; i < size; i++ {
		array = append(array, byte(i))
	}

	return array
}

func testByteForm(object models.ByteForm) {
	bytes := object.Marshal()
	fmt.Println(object.Size())

	fmt.Println(bytes)

	fmt.Println(object)

	_ = object.Unmarshal(bytes)

	fmt.Println(object)
}

func testBlockHeader() {
	option1 := uuid.New()
	option2 := uuid.New()
	voter1 := uuid.New()
	voter2 := uuid.New()
	voter3 := uuid.New()

	header := models.BlockHeader{}

	header.Id = 0
	header.Epoch = 123
	header.Type = 1
	header.PrevHash = [models.HASH_SIZE]byte(gen(models.HASH_SIZE))
	header.MerkleRoot = [models.MERKLE_ROOT_SIZE]byte(gen(models.MERKLE_ROOT_SIZE))
	header.OptionsMeta = append(header.OptionsMeta, models.OptionMeta{
		Uuid:    option1,
		Counter: 255,
	})
	header.OptionsMeta = append(header.OptionsMeta, models.OptionMeta{
		Uuid:    option2,
		Counter: 254,
	})

	header.OptionsCount = 2

	header.Voters = append(header.Voters, voter1)
	header.Voters = append(header.Voters, voter2)
	header.Voters = append(header.Voters, voter3)
	header.Signature = [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE))

	testByteForm(&header)
}

func MakeVoting() models.Voting {
	option1 := models.Option{
		Uuid:  uuid.New(),
		Title: "option1",
	}
	option2 := models.Option{
		Uuid:  uuid.New(),
		Title: "option2",
	}
	option3 := models.Option{
		Uuid:  uuid.New(),
		Title: "option3",
	}
	option4 := models.Option{
		Uuid:  uuid.New(),
		Title: "option4",
	}

	voting := models.Voting{
		Uuid:    uuid.New(),
		Title:   "voting",
		Options: []models.Option{option1, option2, option3, option4},
	}

	return voting
}

func testVoting() {
	voting := MakeVoting()

	testByteForm(&voting)
}

func MakeCorePayload() models.CorePayload {
	payload := models.CorePayload{
		Voting:          MakeVoting(),
		CycleDuration:   1000,
		SendingDuration: 100,
		Signature:       [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE)),
	}

	return payload
}

func testCorePaylaod() {
	payload := MakeCorePayload()

	testByteForm(&payload)
}

func MakeVotePayload() models.VotePaylaod {
	var payload models.VotePaylaod

	vote1 := models.Vote{
		UserUuid:   uuid.New(),
		OptionUuid: uuid.New(),
		Signature:  [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE)),
	}

	vote2 := models.Vote{
		UserUuid:   uuid.New(),
		OptionUuid: uuid.New(),
		Signature:  [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE)),
	}

	vote3 := models.Vote{
		UserUuid:   uuid.New(),
		OptionUuid: uuid.New(),
		Signature:  [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE)),
	}

	payload.Votes = append(payload.Votes, vote1)
	payload.Votes = append(payload.Votes, vote2)
	payload.Votes = append(payload.Votes, vote3)

	node1 := models.Node{
		NodeMeta: models.NodeMeta{
			Uuid: uuid.New(),
			Host: "1.2.3.4",
			Port: 80,
		},
		Signature: [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE)),
	}

	node2 := models.Node{
		NodeMeta: models.NodeMeta{
			Uuid: uuid.New(),
			Host: "255.255.0.255",
			Port: 80,
		},
		Signature: [models.SIGNATURE_SIZE]byte(gen(models.SIGNATURE_SIZE)),
	}

	payload.Nodes = append(payload.Nodes, node1)
	payload.Nodes = append(payload.Nodes, node2)

	return payload
}

func testVotePayload() {
	payload := MakeVotePayload()

	testByteForm(&payload)
}

func main() {
	testVotePayload()
}
