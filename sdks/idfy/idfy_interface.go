package idfy

import (
	"context"
)

type Idfy interface {
	ExtractPan(ctx context.Context, idfyrequest IdfyRequest) (*IdfyPanResponse, error)
	ExtractAadhar(ctx context.Context, idfyrequest IdfyRequest) (*IdfyAadharResponse, error)
	ExtractDl(ctx context.Context, idfyrequest IdfyRequest) (*IdfyDlResponse, error)
	ExtractVoter(ctx context.Context, idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error)
	ExtractPassport(ctx context.Context, idfyrequest IdfyRequest) (*IdfyPassportResponse, error)
	FraudCheckPan(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckPanResponse, error)
	FraudCheckAadhar(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckAadharResponse, string, error)
	FraudCheckDl(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckDlResponse, error)
	FraudCheckVoter(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckPassportResponse, error)
	CheckTemperedImage(ctx context.Context, req CheckTemperedReq) (bool, error)
	PostFruadValidationReq(ctx context.Context, documentType string, fraudCheckRequest FraudCheckRequest) (*string, string, error)
	FetchPostedReq(ctx context.Context, requestID string) (*FraudCheckAadharResponse, string, error)
	Healthcheck() (*HealthCheckRes, error)
	MaskAadharDoc(ctx context.Context, id string, maskAadharDocRequest MaskAadharDocRequest) (*MaskAadharDocResponse, error)
}

type IdfyConfig interface {
	GetIdfyAccountId() string
	GetIdfyApiKey() string
	GetIdfyEndpoint() string
	GetIdfyHealthCheckEndpoint() string
	GetIdfyRetryAttemps() int
}
