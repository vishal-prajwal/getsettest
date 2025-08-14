package aadharlite

import "errors"

var (
	MarshalErr                  = errors.New("MarshalError")
	CreateRequestErr            = errors.New("CreateRequestError")
	SendRequestErr              = errors.New("SendRequestError")
	ReadResponseBodyErr         = errors.New("ReadResponseBodyError")
	HealthCheckFailedErr        = errors.New("HealthCheckFailedError")
	EmptyHealthCheckResponseErr = errors.New("EmptyHealthCheckResponseError")
	UnknownHealthStatusErr      = errors.New("UnknownHealthStatusError")
)
