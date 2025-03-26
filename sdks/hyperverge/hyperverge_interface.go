package hyperverge

import (
	"bytes"
	"context"
)

type Hyperverge interface {
	readDocument(documentType string, hypervergeRequest HypervergeRequest) (*bytes.Buffer, error)
	ReadPan(ctx context.Context, hypervergeRequest HypervergeRequest) (*PanResponse, error)
	ReadAadhar(ctx context.Context, hypervergeRequest HypervergeRequest) (*AadharResponse, error)
	ReadPassport(ctx context.Context, hypervergeRequest HypervergeRequest) (*PassportResponse, error)
	ReadVotedID(ctx context.Context, hypervergeRequest HypervergeRequest) (*VoterIdResponse, error)
	FraudCheckPan(ctx context.Context, fraudCheckPanRequest FraudCheckPanRequest, txnID string) (*FraudCheckPanResponse, error)
	FraudCheckPanV2(ctx context.Context, NSDLPanRequest NSDLPanRequest, txnID string) (*NSDLPanResponse, error)
	FraudCheckDl(ctx context.Context, fraudCheckDlRequest FraudCheckDlRequest, txnID string) (*FraudCheckDlResponse, error)
	FraudCheckVoter(ctx context.Context, fraudCheckVoterRequest FraudCheckVoterRequest, txnID string) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(ctx context.Context, fraudCheckPassportRequest FraudCheckPassportRequest, txnID string) (*FraudCheckPassportResponse, error)
	FraudCheckAadhar(ctx context.Context, fraudCheckAadharRequest FraudCheckAadharRequest, txnID string) (*FraudCheckAadharResponse, string, error)
}

type HypervergeConfig interface {
	GetHypervergeAppID() string
	GetHypervergeAppKey() string
	GetHypervergeEndpoint() string
	GetHypervergeFraudCheckEndpoint() string
	GetHypervergeNSDLUrl() string
}
