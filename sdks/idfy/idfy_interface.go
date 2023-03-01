package idfy

type Idfy interface {
	ExtractPan(idfyrequest IdfyRequest) (*IdfyPanResponse, error)
	ExtractAadhar(idfyrequest IdfyRequest) (*IdfyAadharResponse, error)
	ExtractDl(idfyrequest IdfyRequest) (*IdfyDlResponse, error)
	ExtractVoter(idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error)
	ExtractPassport(idfyrequest IdfyRequest) (*IdfyPassportResponse, error)
}

type IdfyConfig interface {
	GetIdfyAccountId() string
	GetIdfyApiKey() string
	GetIdfyEndpoint() string
}
