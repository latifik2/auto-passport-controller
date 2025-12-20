package utils

import (
	"crypto/sha3"
	"encoding/hex"
)

func GetSHA128String(rawBytes []byte) string {
	outputLength := 16

	hashBytes := sha3.SumSHAKE128(rawBytes, outputLength)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
