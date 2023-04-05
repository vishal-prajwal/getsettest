package hyperverge

import "bytes"

type Hyperverge interface {
	readDocument(documentType string, hypervergeRequest HypervergeRequest) (*bytes.Buffer, error)
	ReadPan(hypervergeRequest HypervergeRequest) (*PanResponse, error)
	ReadAadhar(hypervergeRequest HypervergeRequest) (*AadharResponse, error)
	ReadPassport(hypervergeRequest HypervergeRequest) (*PassportResponse, error)
	ReadVotedID(hypervergeRequest HypervergeRequest) (*VoterIdResponse, error)
	FraudCheckPan(fraudCheckPanRequest FraudCheckPanRequest) (*FraudCheckPanResponse, error)
	FraudCheckDl(fraudCheckDlRequest FraudCheckDlRequest) (*FraudCheckDlResponse, error)
	FraudCheckVoter(fraudCheckVoterRequest FraudCheckVoterRequest) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(fraudCheckPassportRequest FraudCheckPassportRequest) (*FraudCheckPassportResponse, error)
	FraudCheckAadhar(fraudCheckAadharRequest FraudCheckAadharRequest) (*FraudCheckAadharResponse, error)
}

type HypervergeConfig interface {
	GetHypervergeAppID() string
	GetHypervergeAppKey() string
	GetHypervergeEndpoint() string
	GetHypervergeFraudCheckEndpoint() string
	GetHypervergeFraudCheckAadharEndpoint() string
}
