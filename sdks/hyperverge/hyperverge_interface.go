package hyperverge

import (
	"bytes"
	"context"
)

type Hyperverge interface {
	readDocument(documentType string, hypervergeRequest HypervergeRequest) (*bytes.Buffer, error)
	ReadPan(ctx context.Context, hypervergeRequest HypervergeRequest, traceID string) (*PanResponse, error)
	ReadAadhar(ctx context.Context, hypervergeRequest HypervergeRequest, traceID string) (*AadharResponse, error)
	ReadPassport(ctx context.Context, hypervergeRequest HypervergeRequest, traceID string) (*PassportResponse, error)
	ReadVotedID(ctx context.Context, hypervergeRequest HypervergeRequest, traceID string) (*VoterIdResponse, error)
	FraudCheckPan(ctx context.Context, fraudCheckPanRequest FraudCheckPanRequest, txnID string, traceID string) (*FraudCheckPanResponse, error)
	FraudCheckPanV2(ctx context.Context, NSDLPanRequest NSDLPanRequest, txnID string, traceID string) (*NSDLPanResponse, error)
	FraudCheckDl(ctx context.Context, fraudCheckDlRequest FraudCheckDlRequest, txnID string, traceID string) (*FraudCheckDlResponse, error)
	FraudCheckVoter(ctx context.Context, fraudCheckVoterRequest FraudCheckVoterRequest, txnID string, traceID string) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(ctx context.Context, fraudCheckPassportRequest FraudCheckPassportRequest, txnID string, traceID string) (*FraudCheckPassportResponse, error)
	FraudCheckAadhar(ctx context.Context, fraudCheckAadharRequest FraudCheckAadharRequest, txnID string, traceID string) (*FraudCheckAadharResponse, string, error)
}

type HypervergeConfig interface {
	GetHypervergeAppID() string
	GetHypervergeAppKey() string
	GetHypervergeEndpoint() string
	GetHypervergeFraudCheckEndpoint() string
	GetHypervergeNSDLUrl() string
}
