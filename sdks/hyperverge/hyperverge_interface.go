package hyperverge

import "bytes"

type Hyperverge interface {
	readDocument(documentType string, hypervergeRequest HypervergeRequest) (*bytes.Buffer, error)
	ReadPan(hypervergeRequest HypervergeRequest) (*PanResponse, error)
	ReadAadhar(hypervergeRequest HypervergeRequest) (*AadharResponse, error)
	ReadPassport(hypervergeRequest HypervergeRequest) (*PassportResponse, error)
	ReadVotedID(hypervergeRequest HypervergeRequest) (*VoterIdResponse, error)
}

type HypervergeConfig interface {
	GetHypervergeAppID() string
	GetHypervergeAppKey() string
	GetHypervergeEndpoint() string
	GetHypervergeTransactionId() string
}
