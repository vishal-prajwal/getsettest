package imagecompressor

import (
	"bytes"
	"image"
	"image/png"

	"github.com/adrium/goheif"

	"github.com/nfnt/resize"
)

// It will compress the image size by reducing the pixels
func ImageCompressor(imageData []byte, pixel int, extension string) (image.Image, error) {
	var img image.Image
	var err error
	switch extension {
	case "jpg" :
		img, _, err = image.Decode(bytes.NewReader(imageData))
	case "jpeg":
		img, _, err = image.Decode(bytes.NewReader(imageData))
	case "png":
		img, err = png.Decode(bytes.NewReader(imageData))
	case "HEIC":
		img, err = goheif.Decode(bytes.NewReader(imageData))
	default:
		return nil, nil

	}

	if err != nil {
		return nil, err
	}

	newWidth := pixel
	newHeight := newWidth * img.Bounds().Dy() / img.Bounds().Dx()
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	return resizedImg, nil
}
