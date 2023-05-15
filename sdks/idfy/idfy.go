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

func (idfyImpl *IdfyImpl) extract(documentType string, idfyrequest IdfyRequest) (*bytes.Buffer,int, error) {
	url := idfyImpl.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(idfyrequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil,0, err
	}
	idfyImpl.addHeaders(req)

	res, err := idfyImpl.httpClient.Do(req)
	if err != nil {
		return nil,0, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil,0, err
	}
	defer res.Body.Close()

	return body,res.StatusCode, err
}

func (idfyImpl *IdfyImpl) ExtractPan(idfyrequest IdfyRequest) (*IdfyPanResponse, error) {
	documentType := PAN_DOC_TYPE
	byteResp,statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyPanResp PanResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPanResp)
	if err != nil {
		return nil, err
	}
	return &idfyPanResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode,idfyPanResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractAadhar(idfyrequest IdfyRequest) (*IdfyAadharResponse, error) {
	documentType := AADHAR_DOC_TYPE
	byteResp,statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyAadharResp AadharResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyAadharResp)
	if err != nil {
		return nil, err
	}
	return &idfyAadharResp.Result.ExtractionOutput,  idfyImpl.handleError(statusCode,idfyAadharResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractDl(idfyrequest IdfyRequest) (*IdfyDlResponse, error) {
	documentType := DL_DOC_TYPE
	byteResp,statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyDlResp DlResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyDlResp)
	if err != nil {
		return nil, err
	}
	return &idfyDlResp.Result.ExtractionOutput,  idfyImpl.handleError(statusCode,idfyDlResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractVoter(idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error) {
	documentType := VOTER_DOC_TYPE
	byteResp,statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyVoterResp VoterResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyVoterResp)
	if err != nil {
		return nil, err
	}
	return &idfyVoterResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode,idfyVoterResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractPassport(idfyrequest IdfyRequest) (*IdfyPassportResponse, error) {
	documentType := PASSPORT_DOC_TYPE
	byteResp,statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyPassportResp PassportResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPassportResp)
	if err != nil {
		return nil, err
	}
	return &idfyPassportResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode,idfyPassportResp.Error)
}

func (idfyImpl *IdfyImpl) addHeaders(req *http.Request) {
	req.Header.Add("account-id", idfyImpl.config.GetIdfyAccountId())
	req.Header.Add("api-key", idfyImpl.config.GetIdfyApiKey())
	req.Header.Add("Content-Type", "application/json")
}

func (idfyImpl *IdfyImpl) fraudCheck(documentType string, fraudCheckRequest FraudCheckRequest) (*bytes.Buffer, error) {
	postUrl := idfyImpl.config.GetIdfyEndpoint() + documentType
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
	getUrl := idfyImpl.config.GetIdfyEndpoint() + GetTaskStatus
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
	return &fraudCheckPassportResponse, nil
}

func (idfyImpl *IdfyImpl) CheckTemperedImage(req CheckTemperedReq)(bool,error) {

	url := idfyImpl.config.GetIdfyEndpoint() + "/sync/check_tampering/document"
	reqObj, _ := json.Marshal(req)
	payload := strings.NewReader(string(reqObj))
	httpReq, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return false, err
	}
	idfyImpl.addHeaders(httpReq)

	res, err := idfyImpl.httpClient.Do(httpReq)
	if err != nil {
		return false, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	var httpRes CheckTemperedRes
	err = json.Unmarshal(body.Bytes(), &httpRes)
	if err != nil {
		return false, err
	}
	return httpRes.Result.IsTampered, nil
}

func (this *IdfyImpl) handleError(statusCode int,errMsg string) error {
	if statusCode == 200 {
		return nil
	}
	switch statusCode {
	case 422:
		if strings.Contains(errMsg,"INVALID_IMAGE") ||
			strings.Contains(errMsg,"PDF is non compliant to request/quality standard") ||
			strings.Contains(errMsg,"IMAGE_NOT_ACCESSIBLE"){
			return ErrImageNotAccessible
		}
	case 400:
		if strings.Contains(errMsg,"INVALID_IMAGE"){
			return ErrBadRequest
		}
	case 413:
		return ErrBadRequest
	case 500,502:
		return ErrInternalServerError
	case 504:
		return ErrTimeout
	default:
		return ErrInternalServerError
	}
	return nil
}