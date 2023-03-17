package main

import (
	"log"
	"os"

	"bitbucket.org/junglee_games/getsetgo/utils/files"
	"bitbucket.org/junglee_games/getsetgo/utils/zipper"
)

func zip() {

	zipper := zipper.NewZipper()

	err := zipper.AddFile("test/pavan.txt", "Hello there !!!")
	if err != nil {
		panic(err)
	}

	err = zipper.AddFile("test/pavan1.txt", "Hello there !!!")
	if err != nil {
		panic(err)
	}

	zipper.Close()

	// get a zip file content as string
	content := zipper.GetString()
	log.Printf("Zip file content:%s", content)

	// you can save it to file as well
	zipper.SaveFile("myzip.zip")
}

func gunzip() {
	filename := "files.zip.gz"
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	name, rcontent, err := zipper.GunZip(filename, string(content))
	if err != nil {
		panic(err)
	}
	files.SaveFile(".", name, rcontent)
}

func gzip() {
	filename := "files.zip"
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	name, rcontent, err := zipper.GZip(filename, string(content))
	if err != nil {
		panic(err)
	}
	files.SaveFile(".", name, rcontent)
}

func main() {
	// zip()
	// gunzip()
	gzip()
}
