package rummykyc

import "errors"

var (
	ErrCreatingRequest     = errors.New("ErrCreatingRequest")
	ErrCallingKYC          = errors.New("ErrCallingKYC")
	ErrReadingResponseBody = errors.New("ErrReadingResponseBody")
	ErrUnmarshlingResponse = errors.New("ErrUnmarshlingResponse")
	ErrReqValidate         = errors.New("ERROR : kyc request validate")
	ErrHVServer            = errors.New("ERROR : kyc Server")
	ErrNotFound            = errors.New("ERROR : Not found")
)
