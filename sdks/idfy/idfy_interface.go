package idfy

type Idfy interface {
	ExtractPan(idfyrequest IdfyRequest) (*IdfyPanResponse, error)
	ExtractAadhar(idfyrequest IdfyRequest) (*IdfyAadharResponse, error)
	ExtractDl(idfyrequest IdfyRequest) (*IdfyDlResponse, error)
	ExtractVoter(idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error)
	ExtractPassport(idfyrequest IdfyRequest) (*IdfyPassportResponse, error)
	FraudCheckPan(fraudCheckRequest FraudCheckRequest) (*FraudCheckPanResponse, error)
	FraudCheckAadhar(fraudCheckRequest FraudCheckRequest) (*FraudCheckAadharResponse, error)
	FraudCheckDl(fraudCheckRequest FraudCheckRequest) (*FraudCheckDlResponse, error)
	FraudCheckVoter(fraudCheckRequest FraudCheckRequest) (*FraudCheckVoterResponse, error)
	FraudCheckPassport(fraudCheckRequest FraudCheckRequest) (*FraudCheckPassportResponse, error)
	CheckTemperedImage(req CheckTemperedReq)(bool,error)
}

type IdfyConfig interface {
	GetIdfyAccountId() string
	GetIdfyApiKey() string
	GetIdfyEndpoint() string
}
