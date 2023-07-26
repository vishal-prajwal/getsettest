package main

import (
	"context"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strings"

	imagecompressor "bitbucket.org/junglee_games/getsetgo/utils/imageCompressor"
)

func main() {
	file, err := os.Open("input2.jpeg")
	fileNameSplit := strings.Split(file.Name(), ".")
	ext := fileNameSplit[len(fileNameSplit)-1]
	if err != nil {
		log.Fatal("err1: ", err)
	}
	defer file.Close()

	imageData := make([]byte, 0)
	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			imageData = append(imageData, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	resizedImg, err := imagecompressor.ImageCompressor(context.Background(), imageData, 1080, ext)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputFile, err := os.Create("output.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Encode the resized image as JPEG and write it to the output file with 80% quality (you can adjust the quality as needed)
	err = jpeg.Encode(outputFile, resizedImg, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image compressed and resized to 360p, and saved to output.jpeg")
}
