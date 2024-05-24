package user

import "voting-blockchain/internal/models/types"

type Config struct {
	Uuid       types.Uuid `json:"uuid"`
	PublicKey  string     `json:"public_key"`
	PrivateKey string     `json:"private_key"`
}
