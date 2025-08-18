package idfy

import "errors"

var (
	ErrInternalServerError  = errors.New("Internal Server Error")
	ErrTimeout              = errors.New("Server Request Timeout after")
	ErrImageNotAccessible   = errors.New("Server Unable to Access Image Either bcz of pdf or bluerred image")
	ErrBadRequest           = errors.New("Bad Request")
	ErrAddharLiteFetchError = errors.New("Error Fetching Aadhar Lite Data")
	ErrDobMismatch          = errors.New("DOB mismatch")
)
