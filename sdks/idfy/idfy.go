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

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/google/uuid"
)

type IdfyImpl struct {
	config     IdfyConfig
	nr         newrelic.Agent
	httpClient httpclient.HTTPClient
}

const (
	BaseDelay = 1 * time.Second
)

// New creates a new Idfy client
func New(config IdfyConfig, nr newrelic.Agent, client httpclient.HTTPClient) *IdfyImpl {
	idfy := IdfyImpl{
		config:     config,
		nr:         nr,
		httpClient: client,
	}
	return &idfy
}

func (idfyImpl *IdfyImpl) extract(documentType string, idfyrequest IdfyRequest) (*bytes.Buffer, int, error) {
	url := idfyImpl.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(idfyrequest)
	payload := strings.NewReader(string(reqObj))
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

func (idfyImpl *IdfyImpl) ExtractPan(idfyrequest IdfyRequest) (*IdfyPanResponse, error) {
	documentType := PAN_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyPanResp PanResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPanResp)
	if err != nil {
		return nil, err
	}
	return &idfyPanResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode, idfyPanResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractAadhar(idfyrequest IdfyRequest) (*IdfyAadharResponse, error) {
	documentType := AADHAR_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyAadharResp AadharResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyAadharResp)
	if err != nil {
		return nil, err
	}
	return &idfyAadharResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode, idfyAadharResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractDl(idfyrequest IdfyRequest) (*IdfyDlResponse, error) {
	documentType := DL_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyDlResp DlResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyDlResp)
	if err != nil {
		return nil, err
	}
	return &idfyDlResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode, idfyDlResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractVoter(idfyrequest IdfyRequest) (*IdfyVoterIdResponse, error) {
	documentType := VOTER_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyVoterResp VoterResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyVoterResp)
	if err != nil {
		return nil, err
	}
	return &idfyVoterResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode, idfyVoterResp.Error)
}

func (idfyImpl *IdfyImpl) ExtractPassport(idfyrequest IdfyRequest) (*IdfyPassportResponse, error) {
	documentType := PASSPORT_DOC_TYPE
	byteResp, statusCode, err := idfyImpl.extract(documentType, idfyrequest)
	if err != nil {
		return nil, err
	}
	var idfyPassportResp PassportResponse
	err = json.Unmarshal(byteResp.Bytes(), &idfyPassportResp)
	if err != nil {
		return nil, err
	}
	return &idfyPassportResp.Result.ExtractionOutput, idfyImpl.handleError(statusCode, idfyPassportResp.Error)
}

func (idfyImpl *IdfyImpl) addHeaders(req *http.Request) {
	req.Header.Add("account-id", idfyImpl.config.GetIdfyAccountId())
	req.Header.Add("api-key", idfyImpl.config.GetIdfyApiKey())
	req.Header.Add("Content-Type", "application/json")
}

func (this *IdfyImpl) PostFruadValidationReq(documentType string, fraudCheckRequest FraudCheckRequest) (*string, string, error) {
	postUrl := this.config.GetIdfyEndpoint() + documentType
	reqObj, _ := json.Marshal(fraudCheckRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, postUrl, payload)
	if err != nil {
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_SEND_REQUEST), err
	}
	this.addHeaders(req)

	res, err := this.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_SEND_REQUEST), err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}
	defer res.Body.Close()
	fmt.Printf("@@@@ debug %s", body.Bytes())
	var fraudCheckResponse FraudCheckResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckResponse)
	if err != nil {
		return nil, fmt.Sprintf("%s:%s", IDFY, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}
	if fraudCheckResponse.RequestID == "" {
		return nil, fmt.Sprintf("%s:%s", IDFY, "empty_requestid"), fmt.Errorf("empty_requestid")
	}
	return &fraudCheckResponse.RequestID, fmt.Sprintf("%s:%+v", IDFY, fraudCheckResponse), nil
}

func (this *IdfyImpl) FetchPostedReq(requestID string) (*FraudCheckAadharResponse, string, error) {
	var fraudCheckAadharResponse []FraudCheckAadharResponse
	getUrl := this.config.GetIdfyEndpoint() + GetTaskStatus
	params := url.Values{}
	params.Add("request_id", requestID)
	fullURL := fmt.Sprintf("%v?%v", getUrl, params.Encode())
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, "", err
	}
	this.addHeaders(request)
	res, err := this.httpClient.Do(request)
	if err != nil {
		return nil, UNABLE_TO_PARSE_VENDOR_RESPONSE, err
	}
	byteResp := &bytes.Buffer{}
	_, err = byteResp.ReadFrom(res.Body)
	if err != nil {
		return nil, UNABLE_TO_PARSE_VENDOR_RESPONSE, err
	}
	defer res.Body.Close()
	err = json.Unmarshal(byteResp.Bytes(), &fraudCheckAadharResponse)
	if err != nil {
		return nil, UNABLE_TO_PARSE_VENDOR_RESPONSE, fmt.Errorf("res %s error %v", byteResp.Bytes(), err)
	}
	if len(fraudCheckAadharResponse) == 0 {
		return nil, NO_Vendor_Response, fmt.Errorf("unable to validate aadhar")
	}
	frRes := fraudCheckAadharResponse[0]
	if res.StatusCode == 422 || res.StatusCode == 403 || res.StatusCode == 401 {
		return nil, fmt.Sprintf("%s:%+v", IDFY, frRes), ErrAddharLiteFetchError
	}
	if res.StatusCode != 200 {
		return nil, fmt.Sprintf("%s:%+v", IDFY, frRes), fmt.Errorf("statusCode %d body %s", res.StatusCode, res.Body)
	}

	return &frRes, fmt.Sprintf("%s:%+v", IDFY, frRes), err
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

func (idfyImpl *IdfyImpl) FraudCheckAadhar(fraudCheckRequest FraudCheckRequest) (*FraudCheckAadharResponse, string, error) {
	documentType := FraudCheckAadhar
	requestID, vendorResp, err := idfyImpl.PostFruadValidationReq(documentType, fraudCheckRequest)
	if err != nil {
		return nil, vendorResp, err
	}
	return idfyImpl.FetchPostedReq(*requestID)
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

func (idfyImpl *IdfyImpl) CheckTemperedImage(req CheckTemperedReq) (bool, error) {

	url := idfyImpl.config.GetIdfyEndpoint() + TemperedImage
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
	if res.StatusCode != 200 {
		return false, fmt.Errorf("return with error code %d res %v", res.StatusCode, res)
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

func (idfyImpl *IdfyImpl) MaskAadharDoc(id string, maskAadharDocRequest MaskAadharDocRequest) (*MaskAadharDocResponse, error) {
	requestID, err := idfyImpl.getMaskAadharRequestId(id, maskAadharDocRequest)
	if err != nil {
		return nil, err
	}

	return idfyImpl.FetchMaskDoc(*requestID)
}

func (idfyImpl *IdfyImpl) getMaskAadharRequestId(id string, maskAadharDocRequest MaskAadharDocRequest) (*string, error) {
	postUrl := "https://run.mocky.io/v3/01e16030-82e6-42ea-8e82-1f55557f2160" + MASK_AADHAR_DOC
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
		if err == nil {
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
	getUrl := "https://run.mocky.io/v3/edf8318b-1c99-4e58-9631-f842f989353c" + GetTaskStatus
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
