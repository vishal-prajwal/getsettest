package imagecompressorV2

import (
	"context"
	"fmt"
	"time"

	"github.com/h2non/bimg"
	nrf "github.com/newrelic/go-agent/v3/newrelic"
)

// It will compress the image size by reducing the pixels
func ImageCompressor(ctx context.Context, imageData []byte) ([]byte, error) {

	start := time.Now()

	defer nrf.FromContext(ctx).StartSegment("ImageCompressor").End()
	var err error
	newImage, err := bimg.NewImage(imageData).Resize(800, 600)
	if err != nil {
		return nil, err
	}

	fmt.Println(time.Since(start))

	return newImage, nil
}
