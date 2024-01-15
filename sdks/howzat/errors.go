package howzat

import "errors"

var (
	ErrCreatingRequest     = errors.New("ErrCreatingRequest")
	ErrCallingHowzat       = errors.New("ErrCallingHowzat")
	ErrReadingResponseBody = errors.New("ErrReadingResponseBody")
	ErrUnmarshlingResponse = errors.New("ErrUnmarshlingResponse")
	ErrReqValidate         = errors.New("ERROR : howzat request validate")
	ErrHVServer            = errors.New("ERROR : howzat Server")
	ErrNotFound            = errors.New("ERROR : not found")
)
