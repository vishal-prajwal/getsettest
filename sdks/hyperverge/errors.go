package hyperverge

import "errors"

var ErrBlurredImage = errors.New("blurred image")
var ErrTemperedImage = errors.New("tempered image")
var ErrInvalidDoc = errors.New("Invalid doc")
var ErrUploadError = errors.New("Upload error")
var ErrHttpError = errors.New("Http error")