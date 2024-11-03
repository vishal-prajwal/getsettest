package hyperverge

import "bytes"

type Hyperverge interface {
	readDocument(documentType string, hypervergeRequest HypervergeRequest) (*bytes.Buffer, error)
	ReadPan(hypervergeRequest HypervergeRequest) (*PanResponse, error)
	ReadAadhar(hypervergeRequest HypervergeRequest) (*AadharResponse, error)
	ReadPassport(hypervergeRequest HypervergeRequest) (*PassportResponse, error)
	ReadVotedID(hypervergeRequest HypervergeRequest) (*VoterIdResponse, error)
	FraudCheckPan(fraudCheckPanRequest FraudCheckPanRequest, txnID string) (*FraudCheckPanResponse, error)
	FraudCheckDl(fraudCheckDlRequest FraudCheckDlRequest, txnID string) (*FraudCheckDlResponse, error)
	FraudCheckVoter(fraudCheckVoterRequest FraudCheckVoterRequest, txnID string) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(fraudCheckPassportRequest FraudCheckPassportRequest, txnID string) (*FraudCheckPassportResponse, error)
	FraudCheckAadhar(fraudCheckAadharRequest FraudCheckAadharRequest, txnID string) (*FraudCheckAadharResponse, error)
}

type HypervergeConfig interface {
	GetHypervergeAppID() string
	GetHypervergeAppKey() string
	GetHypervergeEndpoint() string
	GetHypervergeFraudCheckEndpoint() string
	GetHypervergeNSDLUrl() string
}
