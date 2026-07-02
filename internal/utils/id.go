package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func NewID() string {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(buffer)
}
