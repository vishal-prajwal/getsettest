package aadharmasking

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashAadhar(aadhar string) string {
	return hashFunc(aadhar) + "_" + aadhar[len(aadhar)-4:]
}

// HashFunc returns the SHA-256 hash of the input string.
func hashFunc(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	sha256Hash := hash.Sum(nil)
	sha256HashString := hex.EncodeToString(sha256Hash)
	return sha256HashString
}
