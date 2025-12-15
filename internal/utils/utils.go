package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// sha256Hex вычисляет сумму массива и кодирует в hex.
func Sha256Hex(buf []byte) string {
	hash := sha256.Sum256(buf)
	return hex.EncodeToString(hash[:])
}
