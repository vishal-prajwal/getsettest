package main

import (
	"context"
	"fmt"
	"os"

	imagecompressorV2 "bitbucket.org/junglee_games/getsetgo/utils/imageCompressorV2"
	"github.com/h2non/bimg"
)

func main() {
	file, err := os.Open("utils/imageCompressor/example/megha_aadharcard_back.jpg")
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

	resizedImg, err := imagecompressorV2.ImageCompressor(context.Background(), imageData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Encode the resized image as JPEG and write it to the output file with 80% quality (you can adjust the quality as needed)
	bimg.Write("new.jpg", resizedImg)

	fmt.Println("Image compressed and resized to 360p, and saved to output.jpeg")
}
