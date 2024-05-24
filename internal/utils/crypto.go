package utils

import "crypto/sha256"

func GetHash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
