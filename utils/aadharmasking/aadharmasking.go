package aadharmasking

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
)

func HashAadhar(aadhar string) string {
	if !verifyAadharNumberToMask(aadhar) {
		return aadhar
	}
	return hashFunc(aadhar)
}

// HashFunc returns the SHA-256 hash of the input string.
func hashFunc(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	sha256Hash := hash.Sum(nil)
	sha256HashString := hex.EncodeToString(sha256Hash)
	return sha256HashString
}

func verifyAadharNumberToMask(aadharNumber string) bool {
	// Check if the length is exactly 12
	if len(aadharNumber) != 12 {
		return false
	}

	// Check if all characters are digits
	for _, char := range aadharNumber {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	return !strings.EqualFold(aadharNumber[0:8], "XXXXXXXX")
}
