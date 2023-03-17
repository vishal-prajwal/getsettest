package idfy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
)

type IdfyImpl struct {
	config     IdfyConfig
	nr         newrelic.Agent
	httpClient httpclient.HTTPClient
}

// New creates a new Idfy client
func New(config IdfyConfig, nr newrelic.Agent, client httpclient.HTTPClient) *IdfyImpl {
	idfy := IdfyImpl{
		config:     config,
		nr:         nr,
		httpClient: client,
	}
	return &idfy
}

func (idfyImpl *IdfyImpl) extract(documentType string, idfyrequest IdfyRequest) (*IdfyResponse, error) {
	url := idfyImpl.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(idfyrequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	idfyImpl.addHeaders(req)

	res, err := idfyImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var idfyResp IdfyResponse
	err = json.Unmarshal(body.Bytes(), &idfyResp)
	return &idfyResp, err
}

func (idfyImpl *IdfyImpl) ExtractPan(idfyrequest IdfyRequest) (*IdfyPanResponse, error) {
	documentType := PAN_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	return &byteResp.Result.IdfyPanResponse, err
}

func (idfyImpl *IdfyImpl) ExtractAadhar(idfyrequest IdfyRequest) (*IdfyAadharResponse, error) {
	documentType := AADHAR_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	return &byteResp.Result.IdfyAadharResponse, err
}

func (idfyImpl *IdfyImpl) ExtractDl(idfyrequest IdfyRequest) (*IdfyDlResponse, error) {
	documentType := DL_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	return &byteResp.Result.IdfyDlResponse, err
}

func (idfyImpl *IdfyImpl) ExtractVoter(idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error) {
	documentType := VOTER_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	return &byteResp.Result.IdfyVoterIdResponse, err
}

func (idfyImpl *IdfyImpl) ExtractPassport(idfyrequest IdfyRequest) (*IdfyPassportResponse, error) {
	documentType := PASSPORT_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	return &byteResp.Result.IdfyPassportResponse, err
}

func (idfyImpl *IdfyImpl) addHeaders(req *http.Request) {
	req.Header.Add("account-id", idfyImpl.config.GetIdfyAccountId())
	req.Header.Add("api-key", idfyImpl.config.GetIdfyApiKey())
	req.Header.Add("Content-Type", "application/json")
}
