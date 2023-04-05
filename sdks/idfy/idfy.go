package idfy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (idfyImpl *IdfyImpl) extract(documentType string, idfyrequest IdfyRequest) (*bytes.Buffer, error) {
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

	return body, err
}

func (idfyImpl *IdfyImpl) ExtractPan(idfyrequest IdfyRequest) (*IdfyPanResponse, error) {
	documentType := PAN_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyPanResp PanResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPanResp)
	if err != nil {
		return nil, err
	}
	if idfyPanResp.Status != "completed" {
		return nil, fmt.Errorf("%v %v", idfyPanResp.Message, idfyPanResp.Error)
	}
	return &idfyPanResp.Result.ExtractionOutput, err
}

func (idfyImpl *IdfyImpl) ExtractAadhar(idfyrequest IdfyRequest) (*IdfyAadharResponse, error) {
	documentType := AADHAR_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyAadharResp AadharResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyAadharResp)
	if err != nil {
		return nil, err
	}
	if idfyAadharResp.Status != "completed" {
		return nil, fmt.Errorf("%v %v", idfyAadharResp.Message, idfyAadharResp.Error)
	}
	return &idfyAadharResp.Result.ExtractionOutput, err
}

func (idfyImpl *IdfyImpl) ExtractDl(idfyrequest IdfyRequest) (*IdfyDlResponse, error) {
	documentType := DL_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyDlResp DlResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyDlResp)
	if err != nil {
		return nil, err
	}
	if idfyDlResp.Status != "completed" {
		return nil, fmt.Errorf("%v %v", idfyDlResp.Message, idfyDlResp.Error)
	}
	return &idfyDlResp.Result.ExtractionOutput, err
}

func (idfyImpl *IdfyImpl) ExtractVoter(idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error) {
	documentType := VOTER_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyVoterResp VoterResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyVoterResp)
	if err != nil {
		return nil, err
	}
	if idfyVoterResp.Status != "completed" {
		return nil, fmt.Errorf("%v %v", idfyVoterResp.Message, idfyVoterResp.Error)
	}
	return &idfyVoterResp.Result.ExtractionOutput, err
}

func (idfyImpl *IdfyImpl) ExtractPassport(idfyrequest IdfyRequest) (*IdfyPassportResponse, error) {
	documentType := PASSPORT_DOC_TYPE
	byteResp, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyPassportResp PassportResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPassportResp)
	if err != nil {
		return nil, err
	}
	if idfyPassportResp.Status != "completed" {
		return nil, fmt.Errorf("%v %v", idfyPassportResp.Message, idfyPassportResp.Error)
	}
	return &idfyPassportResp.Result.ExtractionOutput, err
}

func (idfyImpl *IdfyImpl) addHeaders(req *http.Request) {
	req.Header.Add("account-id", idfyImpl.config.GetIdfyAccountId())
	req.Header.Add("api-key", idfyImpl.config.GetIdfyApiKey())
	req.Header.Add("Content-Type", "application/json")
}

func (idfyImpl *IdfyImpl) fraudCheck(documentType string, fraudCheckRequest FraudCheckRequest) (*bytes.Buffer, error) {
	postUrl := idfyImpl.config.GetIdfyFraudCheckPostEndpoint() + documentType
	reqObj, _ := json.Marshal(fraudCheckRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, postUrl, payload)
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
	var fraudCheckResponse FraudCheckResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckResponse)
	if err != nil {
		return nil, err
	}
	if fraudCheckResponse.RequestID == "" {
		return nil, fmt.Errorf("empty_requestid")
	}
	getUrl := idfyImpl.config.GetIdfyFraudCheckGetEndpoint()
	params := url.Values{}
	params.Add("request_id", fraudCheckResponse.RequestID)
	fullURL := fmt.Sprintf("%v?%v", getUrl, params.Encode())
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	idfyImpl.addHeaders(request)
	res, err = idfyImpl.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	body = &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return body, err
}

func (idfyImpl *IdfyImpl) FraudCheckPan(fraudCheckRequest FraudCheckRequest) (*FraudCheckPanResponse, error) {
	documentType := PAN_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest)
	if err != nil {
		return nil, err
	}
	var fraudCheckPanResponse FraudCheckPanResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckPanResponse)
	if err != nil {
		return nil, err
	}
	if fraudCheckPanResponse.Status != "completed" {
		return nil, fmt.Errorf("%v %v", fraudCheckPanResponse.Message, fraudCheckPanResponse.Error)
	}
	return &fraudCheckPanResponse, err
}

func (idfyImpl *IdfyImpl) FraudCheckAadhar(fraudCheckRequest FraudCheckRequest) (*FraudCheckAadharResponse, error) {
	documentType := FraudCheckAadhar
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest)
	if err != nil {
		return nil, err
	}
	var fraudCheckAadharResponse FraudCheckAadharResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckAadharResponse)
	if err != nil {
		return nil, err
	}
	if fraudCheckAadharResponse.Status != "completed" {
		return nil, fmt.Errorf("%v %v", fraudCheckAadharResponse.Message, fraudCheckAadharResponse.Error)
	}
	return &fraudCheckAadharResponse, err
}

func (idfyImpl *IdfyImpl) FraudCheckDl(fraudCheckRequest FraudCheckRequest) (*FraudCheckDlResponse, error) {
	documentType := DL_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest)
	if err != nil {
		return nil, err
	}
	var fraudCheckDlResponse FraudCheckDlResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckDlResponse)
	if err != nil {
		return nil, err
	}
	if fraudCheckDlResponse.Status != "completed" {
		return nil, fmt.Errorf("%v %v", fraudCheckDlResponse.Message, fraudCheckDlResponse.Error)
	}
	return &fraudCheckDlResponse, err
}

func (idfyImpl *IdfyImpl) FraudCheckVoter(fraudCheckRequest FraudCheckRequest) (*FraudCheckVoterResponse, error) {
	documentType := VOTER_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest)
	if err != nil {
		return nil, err
	}
	var fraudCheckVoterResponse FraudCheckVoterResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckVoterResponse)
	if err != nil {
		return nil, err
	}
	if fraudCheckVoterResponse.Status != "completed" {
		return nil, fmt.Errorf("%v %v", fraudCheckVoterResponse.Message, fraudCheckVoterResponse.Error)
	}
	return &fraudCheckVoterResponse, err
}

func (idfyImpl *IdfyImpl) FraudCheckPassport(fraudCheckRequest FraudCheckRequest) (*FraudCheckPassportResponse, error) {
	documentType := PASSPORT_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest)
	if err != nil {
		return nil, err
	}
	var fraudCheckPassportResponse FraudCheckPassportResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckPassportResponse)
	if err != nil {
		return nil, err
	}
	if fraudCheckPassportResponse.Status != "completed" {
		return nil, fmt.Errorf("%v %v", fraudCheckPassportResponse.Message, fraudCheckPassportResponse.Error)
	}
	return &fraudCheckPassportResponse, err
}
