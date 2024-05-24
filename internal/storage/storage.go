package storage

import (
	"errors"
	"voting-blockchain/internal/caches"
	"voting-blockchain/internal/models"
	"voting-blockchain/internal/models/blocks"
	"voting-blockchain/internal/models/blocks/payloads"
	"voting-blockchain/internal/models/types"
	"voting-blockchain/internal/utils"

	"github.com/google/uuid"
)

type StorageConfig struct {
	BatchSize uint64 `json:"batch_size"`
}

type Storage struct {
	pg  *Postgres
	cfg *StorageConfig
}

func NewStorage(pg *Postgres, cfg *StorageConfig) (*Storage, error) {
	err := pg.PrepareMultipleFromPath(
		"./internal/storage/sql/",
		[]string{
			"select_blockchain_blocks", "select_blockchains", "select_blocks",
			"insert_block", "insert_blockchain",
		},
	)
	if err != nil {
		return nil, err
	}

	return &Storage{
		pg:  pg,
		cfg: cfg,
	}, nil
}

func (storage Storage) InitBlockchainsCache(cache *caches.BlockchainsCache) error {
	db, err := storage.pg.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(storage.pg.Get("select_blockchains"))
	if err != nil {
		return errors.New("cannot select_blockchains: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		blockchain_cache := caches.BlockchainCache{}

		var blockchain_uuid_str string
		var header_bytes []byte
		err := rows.Scan(&blockchain_uuid_str, &header_bytes)
		if err != nil {
			return errors.New("cannot scan block: " + err.Error())
		}

		err = blockchain_cache.LastBlockHeader.Unmarshal(header_bytes)
		if err != nil {
			return errors.New("cannot unmarshal last block header: " + err.Error())
		}
		blockchain_cache.LastBlockHash = utils.GetHash(header_bytes)

		blockchain_uuid, err := uuid.Parse(blockchain_uuid_str)
		if err != nil {
			return errors.New("cannot parse uuid from postgres: " + err.Error())
		}

		cache.Blockchains[types.Uuid(blockchain_uuid)] = blockchain_cache
	}
	rows.Close()

	// get blocks for each blockchain
	for key, value := range cache.Blockchains {
		var last_processed uint64 = 0
		// get blocks by batches
		for last_processed < value.LastBlockHeader.Id {
			first_index := last_processed
			last_index := last_processed + storage.cfg.BatchSize - 1
			if last_index > value.LastBlockHeader.Id {
				last_index = value.LastBlockHeader.Id
			}
			rows, err := db.Query(storage.pg.Get("select_blocks"), key, first_index, last_index)
			if err != nil {
				return errors.New("cannot select_blocks: " + err.Error())
			}
			defer rows.Close()

			// save data from each block
			for rows.Next() {
				// TODO check what we processed all blocks from 0 to last_block_id
				var id uint64
				var header_bytes []byte
				var payload_bytes NullByteSlice
				err := rows.Scan(&id, &header_bytes, &payload_bytes)
				if err != nil {
					return errors.New("cannot scan block: " + err.Error())
				}

				var header blocks.BlockHeader
				header.Unmarshal(header_bytes)
				// save data
				for _, vote := range header.Votes {
					value.Users[vote] = true
				}
				for _, node := range header.Nodes {
					value.Nodes[node.Uuid] = node
				}
			}
		}
	}

	return nil
}

func (storage Storage) InsertBlock(blockchain_uuid types.Uuid, id uint64, header blocks.BlockHeader, payload models.ByteForm) error {
	db, err := storage.pg.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	header_bytes := header.Marshal()
	payload_bytes := payload.Marshal()

	_, err = db.Query(
		storage.pg.Get("insert_block"),
		uuid.UUID(blockchain_uuid).String(),
		id,
		header_bytes,
		payload_bytes,
	)
	if err != nil {
		return errors.New("cannot insert_block: " + err.Error())
	}

	return nil
}

func (storage Storage) InsertBlockchain(uuid types.Uuid, header blocks.BlockHeader, payload payloads.CorePayload) error {
	db, err := storage.pg.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	header_bytes := header.Marshal()
	payload_bytes := payload.Marshal()

	_, err = db.Query(
		storage.pg.Get("insert_blockchain"),
		uuid.ToString,
		header_bytes,
		payload_bytes,
	)
	if err != nil {
		return errors.New("cannot insert_blockchain: " + err.Error())
	}

	return nil
}

func (storage Storage) SelectBlocks(blockchain_uuid string, first_index uint64, last_index uint64) ([]blocks.Block, error) {
	result := []blocks.Block{}

	db, err := storage.pg.Open()
	if err != nil {
		return result, err
	}
	defer db.Close()

	var last_processed uint64 = 0
	// get blocks by batches
	for last_processed < last_index {
		first_index := last_processed
		last_index := last_processed + storage.cfg.BatchSize - 1
		rows, err := db.Query(storage.pg.Get("select_blocks"), blockchain_uuid, first_index, last_index)
		if err != nil {
			return result, errors.New("cannot select_blocks: " + err.Error())
		}
		defer rows.Close()

		// save data from each block
		for rows.Next() {
			var block blocks.Block

			// TODO check what we processed all blocks from 0 to last_block_id
			var id uint64
			var header_bytes []byte
			var payload_bytes NullByteSlice
			err := rows.Scan(&id, &header_bytes, &payload_bytes)
			if err != nil {
				return result, errors.New("cannot scan block: " + err.Error())
			}

			// unmarshal header
			err = block.Header.Unmarshal(header_bytes)
			if err != nil {
				return result, errors.New("cannot unmarshal block header: " + err.Error())
			}
			// unmarshal payload
			if payload_bytes.Valid {
				var payload blocks.BlockPayloadInterface
				if block.Header.Type == blocks.BLOCK_TYPE_CORE {
					payload = &payloads.CorePayload{}
					err := payload.Unmarshal(payload_bytes.ByteSlice)
					if err != nil {
						return result, errors.New("cannot unmarshal core block payload: " + err.Error())
					}
				} else if block.Header.Type == blocks.BLOCK_TYPE_VOTE {
					payload = &payloads.VotePayload{}
					err := payload.Unmarshal(payload_bytes.ByteSlice)
					if err != nil {
						return result, errors.New("cannot unmarshal vote block payload: " + err.Error())
					}
				}
				block.Payload = payload
			}

			result = append(result, block)
		}
	}

	return result, nil
}
