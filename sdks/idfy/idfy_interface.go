package idfy

import "context"

type Idfy interface {
	ExtractPan(ctx context.Context, idfyrequest IdfyRequest, traceID string) (*IdfyPanResponse, error)
	ExtractAadhar(ctx context.Context, idfyrequest IdfyRequest, traceID string) (*IdfyAadharResponse, error)
	ExtractDl(ctx context.Context, idfyrequest IdfyRequest, traceID string) (*IdfyDlResponse, error)
	ExtractVoter(ctx context.Context, idfyrequest IdfyRequest, traceID string) (*IdfyVoterIdResponse, error)
	ExtractPassport(ctx context.Context, idfyrequest IdfyRequest, traceID string) (*IdfyPassportResponse, error)
	FraudCheckPan(ctx context.Context, fraudCheckRequest FraudCheckRequest, traceID string) (*FraudCheckPanResponse, error)
	FraudCheckAadhar(ctx context.Context, fraudCheckRequest FraudCheckRequest, traceID string) (*FraudCheckAadharResponse, string, error)
	FraudCheckDl(ctx context.Context, fraudCheckRequest FraudCheckRequest, traceID string) (*FraudCheckDlResponse, error)
	FraudCheckVoter(ctx context.Context, fraudCheckRequest FraudCheckRequest, traceID string) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(ctx context.Context, fraudCheckRequest FraudCheckRequest, traceID string) (*FraudCheckPassportResponse, error)
	CheckTemperedImage(ctx context.Context, req CheckTemperedReq, traceID string) (bool, error)
	PostFruadValidationReq(ctx context.Context, documentType string, fraudCheckRequest FraudCheckRequest, traceID string) (*string, string, error)
	FetchPostedReq(requestID string, traceID string) (*FraudCheckAadharResponse, string, error)
	Healthcheck() (*HealthCheckRes, error)
	MaskAadharDoc(ctx context.Context, id string, maskAadharDocRequest MaskAadharDocRequest, traceID string) (*MaskAadharDocResponse, error)
}

type IdfyConfig interface {
	GetIdfyAccountId() string
	GetIdfyApiKey() string
	GetIdfyEndpoint() string
	GetIdfyHealthCheckEndpoint() string
	GetIdfyRetryAttemps() int
}
