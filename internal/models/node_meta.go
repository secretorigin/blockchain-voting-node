package models

import "github.com/google/uuid"

type NodeMeta struct {
	Uuid uuid.UUID `json:"uuid"` // = user uuid
	Host IP        `json:"host"`
	Port uint16    `json:"port"`
}
