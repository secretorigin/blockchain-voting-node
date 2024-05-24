package main

import (
	"fmt"
	"voting-blockchain/internal/models"
	"voting-blockchain/internal/models/blockchains"
	"voting-blockchain/internal/models/blocks"
	"voting-blockchain/internal/models/blocks/payloads"
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/models/votings"

	"github.com/google/uuid"
)

func gen(size uint64) []byte {
	var array []byte
	for i := uint64(0); i < size; i++ {
		array = append(array, byte(i))
	}

	return array
}

func testByteForm(object models.ByteForm, emptyobject models.ByteForm) {
	bytes := object.Marshal()
	fmt.Println(object.Size())

	fmt.Println(bytes)

	fmt.Println(object)

	_ = emptyobject.Unmarshal(bytes)

	fmt.Println(emptyobject)
}

func testBlockHeader() {
	option1 := types.Uuid(types.Uuid(uuid.New()))
	option2 := types.Uuid(types.Uuid(uuid.New()))
	voter1 := types.Uuid(types.Uuid(uuid.New()))
	voter2 := types.Uuid(types.Uuid(uuid.New()))
	voter3 := types.Uuid(types.Uuid(uuid.New()))

	header := blocks.BlockHeader{}

	header.Id = 0
	header.Epoch = 123
	header.Type = 1
	header.PrevHash = [types.HASH_SIZE]byte(gen(types.HASH_SIZE))
	header.MerkleRoot = [types.MERKLE_ROOT_SIZE]byte(gen(types.MERKLE_ROOT_SIZE))
	header.OptionsMeta = append(header.OptionsMeta, votings.OptionMeta{
		Uuid:    option1,
		Counter: 255,
	})
	header.OptionsMeta = append(header.OptionsMeta, votings.OptionMeta{
		Uuid:    option2,
		Counter: 254,
	})

	header.OptionsCount = 2

	header.Votes = append(header.Votes, voter1)
	header.Votes = append(header.Votes, voter2)
	header.Votes = append(header.Votes, voter3)

	header.UserUuid = types.Uuid(uuid.New())
	header.Signature = [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE))

	testByteForm(&header, &blocks.BlockHeader{OptionsCount: header.OptionsCount})
}

func MakeVoting() votings.Voting {
	option1 := votings.Option{
		Uuid:  types.Uuid(uuid.New()),
		Title: "option1",
	}
	option2 := votings.Option{
		Uuid:  types.Uuid(uuid.New()),
		Title: "option2",
	}
	option3 := votings.Option{
		Uuid:  types.Uuid(uuid.New()),
		Title: "option3",
	}
	option4 := votings.Option{
		Uuid:  types.Uuid(uuid.New()),
		Title: "option4",
	}

	voting := votings.Voting{
		Uuid:            types.Uuid(uuid.New()),
		Title:           "voting",
		Options:         []votings.Option{option1, option2, option3, option4},
		CycleDuration:   1000,
		SendingDuration: 100,
	}

	return voting
}

func testVoting() {
	voting := MakeVoting()

	testByteForm(&voting, &votings.Voting{})
}

func MakeCorePayload() payloads.CorePayload {
	payload := payloads.CorePayload{
		UserUuid:  types.Uuid(uuid.New()),
		Voting:    MakeVoting(),
		Signature: [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE)),
	}

	return payload
}

func testCorePaylaod() {
	payload := MakeCorePayload()

	testByteForm(&payload, &payloads.CorePayload{})
}

func MakeVotePayload() payloads.VotePayload {
	var payload payloads.VotePayload

	vote1 := votings.Vote{
		UserUuid:   types.Uuid(uuid.New()),
		OptionUuid: types.Uuid(uuid.New()),
		Signature:  [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE)),
	}

	vote2 := votings.Vote{
		UserUuid:   types.Uuid(uuid.New()),
		OptionUuid: types.Uuid(uuid.New()),
		Signature:  [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE)),
	}

	vote3 := votings.Vote{
		UserUuid:   types.Uuid(uuid.New()),
		OptionUuid: types.Uuid(uuid.New()),
		Signature:  [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE)),
	}

	payload.Votes = append(payload.Votes, vote1)
	payload.Votes = append(payload.Votes, vote2)
	payload.Votes = append(payload.Votes, vote3)

	node1 := blockchains.Node{
		NodeMeta: blockchains.NodeMeta{
			Uuid: types.Uuid(uuid.New()),
			Ip:   "1.2.3.4",
			Port: 80,
		},
		Signature: [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE)),
	}

	node2 := blockchains.Node{
		NodeMeta: blockchains.NodeMeta{
			Uuid: types.Uuid(uuid.New()),
			Ip:   "255.255.0.255",
			Port: 80,
		},
		Signature: [types.SIGNATURE_SIZE]byte(gen(types.SIGNATURE_SIZE)),
	}

	payload.Nodes = append(payload.Nodes, node1)
	payload.Nodes = append(payload.Nodes, node2)

	return payload
}

func testVotePayload() {
	payload := MakeVotePayload()

	testByteForm(&payload, &payloads.VotePayload{})
}

func main() {
	fmt.Println("TEST BLOCK HEADER")
	testBlockHeader()
	fmt.Println("TEST CORE PAYLOAD")
	testCorePaylaod()
	fmt.Println("TEST VOTE PAYLOAD")
	testVotePayload()
}
