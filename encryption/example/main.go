package main

import (
	"os"

	"bitbucket.org/junglee_games/getsetgo/encryption"
	"bitbucket.org/junglee_games/getsetgo/utils/files"
)

func main() {
	bytes, err := os.ReadFile("test.pdf")
	if err != nil {
		panic(err)
	}
	bytes, err = encryption.ProtectPDFBytes(bytes, "test1")
	if err != nil {
		panic(err)
	}
	err = files.SaveFile(".", "test.pdf", string(bytes))
	if err != nil {
		panic(err)
	}

}
