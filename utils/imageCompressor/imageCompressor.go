package imagecompressor

import (
	"bytes"
	"image"

	"github.com/nfnt/resize"
)

// It will compress the image size by reducing the pixels
func ImageCompressor(imageData []byte, pixel int) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	newWidth := pixel
	newHeight := newWidth * img.Bounds().Dy() / img.Bounds().Dx()
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	return resizedImg, nil
}
