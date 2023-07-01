package hyperverge

import "errors"

var ErrBlurredImage = errors.New("blurred image")
var ErrTemperedImage = errors.New("tempered image")
var ErrInvalidDoc = errors.New("Invalid doc")
var ErrUploadError = errors.New("Upload error")
var ErrHttpError = errors.New("Http error")
var ErrFruadCheckUnauthrised = errors.New("401 Unauthorized")
var ErrSomethingWentWrong = errors.New("something went wrong")
var ErrBadRequest = errors.New("bad request")
var ErrInvalidDocID = errors.New("invalid doc id")
