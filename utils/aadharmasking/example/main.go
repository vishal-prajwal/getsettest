package main

import (
	"fmt"

	"bitbucket.org/junglee_games/getsetgo/utils/aadharmasking"
)

func main() {
	str := "678693525967"
	sha256HashString := aadharmasking.HashAadhar(str)
	fmt.Printf("SHA-256 hash of '%s': %s\n", str, sha256HashString)
}
