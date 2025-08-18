package idfy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/apilogger"
	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/sdks/constants"
	"bitbucket.org/junglee_games/getsetgo/sdks/hyperverge"
	"bitbucket.org/junglee_games/getsetgo/utils/aadharmasking"
	"github.com/google/uuid"
)

type IdfyImpl struct {
	config     IdfyConfig
	nr         newrelic.Agent
	httpClient httpclient.HTTPClient
	apilogger  apilogger.ApiUsageLogger
}

const (
	BaseDelay = 1 * time.Second
)

// New creates a new Idfy client
func New(config IdfyConfig, apiLogger apilogger.ApiUsageLogger, nr newrelic.Agent, client httpclient.HTTPClient) *IdfyImpl {
	idfy := IdfyImpl{
		config:     config,
		nr:         nr,
		httpClient: client,
		apilogger:  apiLogger,
	}
	return &idfy
}

func (idfyImpl *IdfyImpl) extract(documentType string, idfyrequest IdfyRequest, apiDataBuilder *apilogger.ApiDataBuilder) (*bytes.Buffer, int, error) {
	url := idfyImpl.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(idfyrequest)
	payload := strings.NewReader(string(reqObj))
	if apiDataBuilder != nil {
		apiDataBuilder.WithRequest(url, http.MethodPost, "", map[string]string{
			"task_id":  idfyrequest.TaskID,
			"group_id": idfyrequest.GroupID,
		})
	}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, 0, err
	}
	idfyImpl.addHeaders(req)

	res, err := idfyImpl.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	return body, res.StatusCode, err
}

func (idfyImpl *IdfyImpl) ExtractPan(ctx context.Context, idfyrequest IdfyRequest) (*IdfyPanResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadPan)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := PAN_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var idfyPanResp PanResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPanResp)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", statusCode), byteResp.String())
	err = idfyImpl.handleError(statusCode, idfyPanResp.Error)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "ExtractPan:: txnId : %s,response from Idfy %+v", idfyrequest.TaskID, idfyPanResp)
		return nil, err
	}

	return &idfyPanResp.Result.ExtractionOutput, nil
}

func (idfyImpl *IdfyImpl) ExtractAadhar(ctx context.Context, idfyrequest IdfyRequest) (*IdfyAadharResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadAadhar)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := AADHAR_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var idfyAadharResp AadharResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyAadharResp)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", statusCode), aadharmasking.MaskAddharInResponseJson(byteResp.String()))
	err = idfyImpl.handleError(statusCode, idfyAadharResp.Error)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "ExtractAadhar:: txnId : %s,response from Idfy %+v", idfyrequest.TaskID, idfyAadharResp)
		return nil, err
	}
	return &idfyAadharResp.Result.ExtractionOutput, nil
}

func (idfyImpl *IdfyImpl) ExtractDl(ctx context.Context, idfyrequest IdfyRequest) (*IdfyDlResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadDl)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := DL_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var idfyDlResp DlResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyDlResp)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", statusCode), byteResp.String())
	err = idfyImpl.handleError(statusCode, idfyDlResp.Error)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "ExtractDl:: txnId : %s,response from Idfy %+v", idfyrequest.TaskID, idfyDlResp)
		return nil, err
	}
	return &idfyDlResp.Result.ExtractionOutput, nil
}

func (idfyImpl *IdfyImpl) ExtractVoter(ctx context.Context, idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadVoterId)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := VOTER_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var idfyVoterResp VoterResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyVoterResp)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", statusCode), byteResp.String())
	err = idfyImpl.handleError(statusCode, idfyVoterResp.Error)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "ExtractVoter:: txnId : %s,response from Idfy %+v", idfyrequest.TaskID, idfyVoterResp)
		return nil, err
	}
	return &idfyVoterResp.Result.ExtractionOutput, nil
}

func (idfyImpl *IdfyImpl) ExtractPassport(ctx context.Context, idfyrequest IdfyRequest) (*IdfyPassportResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadPassport)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := PASSPORT_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var idfyPassportResp PassportResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPassportResp)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", statusCode), byteResp.String())
	err = idfyImpl.handleError(statusCode, idfyPassportResp.Error)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "ExtractPassport:: txnId : %s,response from Idfy %+v", idfyrequest.TaskID, idfyPassportResp)
		return nil, err
	}
	return &idfyPassportResp.Result.ExtractionOutput, nil
}

func (idfyImpl *IdfyImpl) addHeaders(req *http.Request) {
	req.Header.Add("account-id", idfyImpl.config.GetIdfyAccountId())
	req.Header.Add("api-key", idfyImpl.config.GetIdfyApiKey())
	req.Header.Add("Content-Type", "application/json")
}

func (this *IdfyImpl) PostFruadValidationReq(ctx context.Context, documentType string, fraudCheckRequest FraudCheckRequest) (*string, string, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(this.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.FraudCheckAadharGetRequestID)
	defer func() {
		this.apilogger.Log(ctx, apiDataBuilder.Build())
	}()
	postUrl := this.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(fraudCheckRequest)
	payload := strings.NewReader(string(reqObj))
	apiDataBuilder.WithRequest(postUrl, http.MethodPost, aadharmasking.MaskAddharInResponseJson(fmt.Sprintf("%+v", fraudCheckRequest)), map[string]string{
		"task_id":  fraudCheckRequest.TaskID,
		"group_id": fraudCheckRequest.GroupID,
	})
	req, err := http.NewRequest(http.MethodPost, postUrl, payload)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_SEND_REQUEST), err
	}
	this.addHeaders(req)

	res, err := this.httpClient.Do(req)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_SEND_REQUEST), err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}
	defer res.Body.Close()
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())
	fmt.Printf("@@@@ debug %s", body.Bytes())
	var fraudCheckResponse FraudCheckResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}
	if fraudCheckResponse.RequestID == "" {
		apiDataBuilder.WithError("empty_requestid")
		return nil, fmt.Sprintf("%s:%s", IDFY, "empty_requestid"), fmt.Errorf("empty_requestid")
	}
	return &fraudCheckResponse.RequestID, fmt.Sprintf("%s:%+v", IDFY, fraudCheckResponse), nil
}

func (this *IdfyImpl) FetchPostedReq(ctx context.Context, requestID string) (*FraudCheckAadharResponse, string, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(this.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.FraudCheckAadhar)
	defer func() {
		this.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	var fraudCheckAadharResponse []FraudCheckAadharResponse
	getUrl := this.config.GetIdfyEndpoint() + GetTaskStatus
	params := url.Values{}
	params.Add("request_id", requestID)
	fullURL := fmt.Sprintf("%v?%v", getUrl, params.Encode())
	apiDataBuilder.WithRequest(fullURL, http.MethodGet, "", map[string]string{
		"task_id":  requestID,
		"group_id": requestID,
	})
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, "", err
	}
	this.addHeaders(request)
	res, err := this.httpClient.Do(request)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, UNABLE_TO_PARSE_VENDOR_RESPONSE, err
	}
	byteResp := &bytes.Buffer{}
	_, err = byteResp.ReadFrom(res.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, UNABLE_TO_PARSE_VENDOR_RESPONSE, err
	}
	defer res.Body.Close()
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), byteResp.String())
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckAadharResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, UNABLE_TO_PARSE_VENDOR_RESPONSE, fmt.Errorf("res %s error %v", byteResp.Bytes(), err)
	}
	if len(fraudCheckAadharResponse) == 0 {
		apiDataBuilder.WithError("unable to validate aadhar")
		return nil, NO_Vendor_Response, fmt.Errorf("unable to validate aadhar")
	}
	frRes := fraudCheckAadharResponse[0]
	if res.StatusCode == 422 || res.StatusCode == 403 || res.StatusCode == 401 {
		apiDataBuilder.WithError(fmt.Sprintf("%s:%+v", IDFY, frRes))
		return nil, fmt.Sprintf("%s:%+v", IDFY, frRes), ErrAddharLiteFetchError
	}
	if res.StatusCode != 200 {
		return nil, fmt.Sprintf("%s:%+v", IDFY, frRes), fmt.Errorf("statusCode %d body %s", res.StatusCode, res.Body)
	}

	return &frRes, fmt.Sprintf("%s:%+v", IDFY, frRes), err
}

func (idfyImpl *IdfyImpl) fraudCheck(documentType string, fraudCheckRequest FraudCheckRequest, apiDataBuilder *apilogger.ApiDataBuilder) (*bytes.Buffer, error) {
	postUrl := idfyImpl.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(fraudCheckRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, postUrl, payload)
	if err != nil {
		return nil, err
	}
	idfyImpl.addHeaders(req)
	if apiDataBuilder != nil {
		requestPyload := fmt.Sprintf("%+v", fraudCheckRequest.Data)
		if documentType == AADHAR_DOC_TYPE {
			requestPyload = aadharmasking.MaskAddharInResponseJson(requestPyload)
		}
		apiDataBuilder.WithRequest(postUrl, http.MethodPost, requestPyload, map[string]string{
			"task_id":  fraudCheckRequest.TaskID,
			"group_id": fraudCheckRequest.GroupID,
		})
	}

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

func (idfyImpl *IdfyImpl) FraudCheckPan(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckPanResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadPan)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := PAN_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var fraudCheckPanResponse FraudCheckPanResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckPanResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	if fraudCheckPanResponse.Status != "completed" {
		apiDataBuilder.WithError(fmt.Sprintf("%v %v", fraudCheckPanResponse.Message, fraudCheckPanResponse.Error)).
			WithResponse("", byteResp.String())
		return nil, fmt.Errorf("%v %v", fraudCheckPanResponse.Message, fraudCheckPanResponse.Error)
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", 200), byteResp.String())
	return &fraudCheckPanResponse, err
}

func (idfyImpl *IdfyImpl) FraudCheckAadhar(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckAadharResponse, string, error) {
	documentType := FraudCheckAadhar
	requestID, vendorResp, err := idfyImpl.PostFruadValidationReq(ctx, documentType, fraudCheckRequest)
	if err != nil {
		return nil, vendorResp, err
	}
	frRes, res, err := idfyImpl.FetchPostedReq(ctx, *requestID)
	if err != nil {
		logger.Error(ctx, "Error in FetchPostedReq with res %s, error %v", res, err)
	}
	return frRes, res, err
}

func (idfyImpl *IdfyImpl) FraudCheckDl(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckDlResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadDl)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := DL_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var fraudCheckDlResponse FraudCheckDlResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckDlResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	if fraudCheckDlResponse.Status != "completed" {
		apiDataBuilder.WithError(fmt.Sprintf("%v %v", fraudCheckDlResponse.Message, fraudCheckDlResponse.Error)).
			WithResponse("", byteResp.String())
		return nil, fmt.Errorf("%v %v", fraudCheckDlResponse.Message, fraudCheckDlResponse.Error)
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", 200), byteResp.String())
	return &fraudCheckDlResponse, nil
}

func (idfyImpl *IdfyImpl) FraudCheckVoter(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckVoterResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadVoterId)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := VOTER_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var fraudCheckVoterResponse FraudCheckVoterResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckVoterResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	if fraudCheckVoterResponse.Status != "completed" {
		apiDataBuilder.WithError(fmt.Sprintf("%v %v", fraudCheckVoterResponse.Message, fraudCheckVoterResponse.Error)).
			WithResponse("", byteResp.String())
		return nil, fmt.Errorf("%v %v", fraudCheckVoterResponse.Message, fraudCheckVoterResponse.Error)
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", 200), byteResp.String())
	return &fraudCheckVoterResponse, nil
}

func (idfyImpl *IdfyImpl) FraudCheckPassport(ctx context.Context, fraudCheckRequest FraudCheckRequest) (*FraudCheckPassportResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.OCRReadPassport)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	documentType := PASSPORT_DOC_TYPE
	byteResp, err := idfyImpl.fraudCheck(documentType, fraudCheckRequest, apiDataBuilder)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	var fraudCheckPassportResponse FraudCheckPassportResponse
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckPassportResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", byteResp.String())
		return nil, err
	}
	if fraudCheckPassportResponse.Status != "completed" {
		apiDataBuilder.WithError(fmt.Sprintf("%v %v", fraudCheckPassportResponse.Message, fraudCheckPassportResponse.Error)).
			WithResponse("", byteResp.String())
		return nil, fmt.Errorf("%v %v", fraudCheckPassportResponse.Message, fraudCheckPassportResponse.Error)
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", 200), byteResp.String())
	return &fraudCheckPassportResponse, nil
}

func (idfyImpl *IdfyImpl) CheckTemperedImage(ctx context.Context, req CheckTemperedReq) (bool, error) {

	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.CheckTemperedImage)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	url := idfyImpl.config.GetIdfyEndpoint() + TemperedImage
	reqObj, _ := json.Marshal(req)
	payload := strings.NewReader(string(reqObj))
	apiDataBuilder.WithRequest(url, http.MethodPost, "", map[string]string{
		"task_id":  req.TaskID,
		"group_id": req.GroupID,
	})
	httpReq, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return false, err
	}
	idfyImpl.addHeaders(httpReq)

	res, err := idfyImpl.httpClient.Do(httpReq)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "Error in CheckTemperedImage with res %+v, error %v", res, err)
		return false, err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return false, err
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())
	defer res.Body.Close()
	if res.StatusCode != 200 {
		apiDataBuilder.WithError(fmt.Sprintf("statusCode %d body %s", res.StatusCode, body.String()))
		logger.Error(ctx, "Error in CheckTemperedImage with res %+v, error %v", body.String(), err)
		return false, fmt.Errorf("return with error code %d res %v", res.StatusCode, res)
	}
	var httpRes CheckTemperedRes
	err = json.Unmarshal(body.Bytes(), &httpRes)
	if err != nil {
		return false, err
	}
	return httpRes.Result.IsTampered, nil
}

func (this *IdfyImpl) handleError(statusCode int, errMsg string) error {
	if statusCode == 200 {
		return nil
	}
	switch statusCode {
	case 422:
		// will have to confirm the scnario with idfy
		if strings.Contains(errMsg, "INVALID_IMAGE") ||
			strings.Contains(errMsg, "PDF is non compliant to request/quality standard") ||
			strings.Contains(errMsg, "IMAGE_NOT_ACCESSIBLE") {
			return ErrBadRequest
		}
	case 400:
		if strings.Contains(errMsg, "INVALID_IMAGE") {
			return ErrBadRequest
		}
	case 413:
		return ErrBadRequest
	case 500, 502:
		return ErrInternalServerError
	case 504:
		return ErrTimeout
	default:
		return ErrInternalServerError
	}
	return nil
}

func (idfyImpl *IdfyImpl) FraudCheckPanNSDL(ctx context.Context, NSDLPanRequest hyperverge.NSDLPanRequest, txnID string) (*hyperverge.NSDLPanResponse, error) {

	requestID, err := idfyImpl.postPanNSDLRequest(ctx, NSDLPanRequest, txnID)
	if err != nil {
		logger.Error(ctx, "Error in postPanNSDLRequest with res %v, error %v", requestID, err)
		return nil, err
	}
	idfyResponse, err := idfyImpl.fetchPanNSDLResponse(ctx, *requestID)
	if err != nil {
		logger.Error(ctx, "Error in fetchPanNSDLResponse with res %v, error %v", idfyResponse, err)
		return nil, err
	}

	nsdlPanResponse := &hyperverge.NSDLPanResponse{
		Status:     idfyResponse.Status,
		StatusCode: 200,
	}
	nsdlPanResponse.Result.Pan = idfyResponse.Result.SourceOutput.InputDetails.InputPanNumber
	nsdlPanResponse.Result.PanStatus = idfyResponse.Result.SourceOutput.PanStatus
	nsdlPanResponse.Result.Name = idfyResponse.Result.SourceOutput.InputDetails.InputName
	nsdlPanResponse.Result.Dob = idfyResponse.Result.SourceOutput.InputDetails.InputDOB
	if idfyResponse.Result.SourceOutput.AadhaarSeedingStatus {
		nsdlPanResponse.Result.AadharSeedingStatus = "true"
	} else {
		nsdlPanResponse.Result.AadharSeedingStatus = "false"
	}

	if !idfyResponse.Result.SourceOutput.DOBMatch {
		return nsdlPanResponse, ErrDobMismatch
	}

	return nsdlPanResponse, nil
}

func (idfyImpl *IdfyImpl) postPanNSDLRequest(ctx context.Context, NSDLPanRequest hyperverge.NSDLPanRequest, txnID string) (*string, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, NSDL_PAN_DOC_TYPE)
	defer func() {
		idfyImpl.apilogger.Log(ctx, apiDataBuilder.Build())
	}()

	postUrl := idfyImpl.config.GetIdfyEndpoint() + NSDL_PAN_DOC_TYPE
	idfyRequest := IDfyNSDLPanRequest{
		TaskID:  uuid.New().String(),
		GroupID: uuid.New().String(),
		Data: IDfyNSDLPanData{
			IDNumber: NSDLPanRequest.Pan,
			FullName: NSDLPanRequest.Name,
			DOB:      NSDLPanRequest.Dob,
		},
	}
	reqObj, _ := json.Marshal(idfyRequest)
	payload := strings.NewReader(string(reqObj))
	apiDataBuilder.WithRequest(postUrl, http.MethodPost, string(reqObj), map[string]string{
		"task_id":  idfyRequest.TaskID,
		"group_id": idfyRequest.GroupID,
	})
	req, err := http.NewRequest(http.MethodPost, postUrl, payload)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	idfyImpl.addHeaders(req)

	res, err := idfyImpl.httpClient.Do(req)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, err
	}
	defer res.Body.Close()
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())

	var idfyPostResponse IDfyNSDLPanPostResponse
	err = json.Unmarshal(body.Bytes(), &idfyPostResponse)
	if err != nil {
		apiDataBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	if idfyPostResponse.RequestID == "" {
		apiDataBuilder.WithError("empty_requestid")
		return nil, fmt.Errorf("empty_requestid")
	}

	return &idfyPostResponse.RequestID, nil
}

func (idfyImpl *IdfyImpl) fetchPanNSDLResponse(ctx context.Context, requestID string) (*IDfyNSDLPanGetResponse, error) {
	var err error
	var res *IDfyNSDLPanGetResponse
	initialTime := time.Now()

	for attempt := 1; attempt <= idfyImpl.config.GetIdfyRetryAttemps(); attempt++ {
		// Calculate the delay with exponential backoff and jitter
		delay := BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay += jitter

		res, err = idfyImpl.callGetPanNSDLAPI(requestID)
		if err == nil && res != nil && res.Status == "completed" {
			successTime := time.Now()
			logger.Info(context.Background(), "PAN NSDL request with request_id %s completed in time %v second, took %d attempt", requestID, successTime.Sub(initialTime).Seconds(), attempt)
			return res, nil
		}

		time.Sleep(delay)
	}

	logger.Error(context.Background(), "PAN NSDL request with request_id %s failed in time %v second, took %d attempt", requestID, time.Now().Sub(initialTime).Seconds(), idfyImpl.config.GetIdfyRetryAttemps())
	return nil, err
}

func (idfyImpl *IdfyImpl) callGetPanNSDLAPI(requestID string) (*IDfyNSDLPanGetResponse, error) {
	var panNSDLResponse []IDfyNSDLPanGetResponse
	getUrl := idfyImpl.config.GetIdfyEndpoint() + GetTaskStatus
	params := url.Values{}
	params.Add("request_id", requestID)
	fullURL := fmt.Sprintf("%v?%v", getUrl, params.Encode())

	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	idfyImpl.addHeaders(request)

	res, err := idfyImpl.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("statusCode %d body %s", res.StatusCode, res.Body)
	}

	byteResp := &bytes.Buffer{}
	_, err = byteResp.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(byteResp.Bytes(), &panNSDLResponse)
	if err != nil {
		return nil, fmt.Errorf("res %s error %v", byteResp.Bytes(), err)
	}
	if len(panNSDLResponse) == 0 {
		return nil, fmt.Errorf("unable to validate pan")
	}

	frRes := panNSDLResponse[0]
	return &frRes, err
}

func (idfyImpl *IdfyImpl) Healthcheck() (*HealthCheckRes, error) {

	var healthCheckRes HealthCheckRes
	url := idfyImpl.config.GetIdfyHealthCheckEndpoint() + HealthCheck
	req := HealthCheckReq{
		TaskID:  uuid.NewString(),
		GroupID: uuid.NewString(),
		Data: HealthCheckReqData{
			TaskType: "verify_with_source_aadhaar_lite",
		},
	}
	reqObj, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(string(reqObj))
	httpReq, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	idfyImpl.addHeaders(httpReq)
	idfyImpl.httpClient = httpclient.NewHttpClient(HealthCheckTimeout)
	res, err := idfyImpl.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &healthCheckRes)
	if err != nil {
		return nil, err
	}
	return &healthCheckRes, nil
}

func (idfyImpl *IdfyImpl) MaskAadharDoc(ctx context.Context, id string, maskAadharDocRequest MaskAadharDocRequest) (*MaskAadharDocResponse, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(idfyImpl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, IDFY, constants.MaskAadhar)
	defer func() {
		idfyImpl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	apiDataBuilder.WithRequest(idfyImpl.config.GetIdfyEndpoint()+MASK_AADHAR_DOC, http.MethodPost, "", map[string]string{
		"task_id":  maskAadharDocRequest.TaskID,
		"group_id": maskAadharDocRequest.GroupID,
	})
	requestID, err := idfyImpl.getMaskAadharRequestId(id, maskAadharDocRequest)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "Error in MaskAadharDoc with res %+v, error %v", requestID, err)
		return nil, err
	}
	res, err := idfyImpl.FetchMaskDoc(*requestID)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		logger.Error(ctx, "Error in MaskAadharDoc with res %+v, error %v", res, err)
	}
	return res, err
}

func (idfyImpl *IdfyImpl) getMaskAadharRequestId(id string, maskAadharDocRequest MaskAadharDocRequest) (*string, error) {
	postUrl := idfyImpl.config.GetIdfyEndpoint() + MASK_AADHAR_DOC
	reqObj, err := json.Marshal(maskAadharDocRequest)
	if err != nil {
		return nil, err
	}
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
	var maskAadharRequestID MaskAadharRequestID
	err = json.Unmarshal(body.Bytes(), &maskAadharRequestID)
	if err != nil {
		return nil, err
	}
	logger.Info(context.Background(), "Mask Aadhar RequestID Response for id %v is %v", id, body.String())

	if maskAadharRequestID.RequestID == "" {
		return nil, fmt.Errorf("empty_requestid")
	}
	return &maskAadharRequestID.RequestID, nil
}

func (idfyImpl *IdfyImpl) FetchMaskDoc(requestID string) (*MaskAadharDocResponse, error) {
	var err error
	var res *MaskAadharDocResponse
	initialTime := time.Now()
	for attempt := 1; attempt <= idfyImpl.config.GetIdfyRetryAttemps(); attempt++ {

		// Calculate the delay with exponential backoff and jitter
		delay := BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay += jitter

		res, err = idfyImpl.callMaskingApi(requestID)
		if err == nil && res != nil && res.Status == "completed" {
			successTime := time.Now()
			logger.Info(context.Background(), "Mask Aadhar request with request_id %s completed in time %v second, took %d attempt", requestID, successTime.Sub(initialTime).Seconds(), attempt)
			return res, nil
		}

		time.Sleep(delay)
	}
	logger.Error(context.Background(), "Mask Aadhar request with request_id %s failed in time %v second, took %d attempt", requestID, time.Now().Sub(initialTime).Seconds(), idfyImpl.config.GetIdfyRetryAttemps())
	return nil, err
}

func (idfyImpl *IdfyImpl) callMaskingApi(requestID string) (*MaskAadharDocResponse, error) {
	var maskAadharDocResponse []MaskAadharDocResponse
	getUrl := idfyImpl.config.GetIdfyEndpoint() + GetTaskStatus
	params := url.Values{}
	params.Add("request_id", requestID)
	fullURL := fmt.Sprintf("%v?%v", getUrl, params.Encode())
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	idfyImpl.addHeaders(request)
	res, err := idfyImpl.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("statusCode %d body %s", res.StatusCode, res.Body)
	}
	byteResp := &bytes.Buffer{}
	_, err = byteResp.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.Unmarshal(byteResp.Bytes(), &maskAadharDocResponse)
	if err != nil {
		return nil, fmt.Errorf("res %s error %v", byteResp.Bytes(), err)
	}
	if len(maskAadharDocResponse) == 0 {
		return nil, fmt.Errorf("unable to validate aadhar")
	}
	frRes := maskAadharDocResponse[0]
	return &frRes, err
}
