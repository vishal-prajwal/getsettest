package imagecompressor

import (
	"bytes"
	"context"
	"image"
	"image/png"

	nrf "github.com/newrelic/go-agent/v3/newrelic"

	"github.com/nfnt/resize"
)

// It will compress the image size by reducing the pixels
func ImageCompressor(ctx context.Context, imageData []byte, pixel int, extension string) (image.Image, error) {

	defer nrf.FromContext(ctx).StartSegment("ImageCompressor").End()
	var img image.Image
	var err error
	switch extension {
	case "jpg":
		img, _, err = image.Decode(bytes.NewReader(imageData))
	case "jpeg":
		img, _, err = image.Decode(bytes.NewReader(imageData))
	case "png":
		img, err = png.Decode(bytes.NewReader(imageData))
	default:
		return nil, image.ErrFormat
	}

	if err != nil {
		return nil, err
	}

	newWidth := pixel
	newHeight := newWidth * img.Bounds().Dy() / img.Bounds().Dx()
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	return resizedImg, nil
}
