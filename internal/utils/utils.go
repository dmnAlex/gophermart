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

// IsValidLuhn проверяет является ли строка корректным числом соответствующим алгоритму Луна
func IsValidLuhn(number string) bool {
	if len(number) < 2 {
		return false
	}

	var sum int64
	var alt bool

	for i := len(number) - 1; i >= 0; i-- {
		digit := int64(number[i] - '0')

		if digit < 0 || digit > 9 {
			return false
		}

		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alt = !alt
	}

	return sum%10 == 0
}
