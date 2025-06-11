package hyperverge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/apilogger"
	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/sdks/constants"
	"bitbucket.org/junglee_games/getsetgo/utils/aadharmasking"
	"github.com/pkg/errors"
)

type HypervergeImpl struct {
	config     HypervergeConfig
	nr         newrelic.Agent
	httpClient httpclient.HTTPClient
	apilogger  apilogger.ApiUsageLogger
}

// New creates a new Hyperverge client
func New(config HypervergeConfig, apiLogger apilogger.ApiUsageLogger, nr newrelic.Agent, client httpclient.HTTPClient) *HypervergeImpl {
	hypervergeImpl := HypervergeImpl{
		config:     config,
		nr:         nr,
		httpClient: client,
		apilogger:  apiLogger,
	}
	return &hypervergeImpl
}

func (hypervergeImpl *HypervergeImpl) readDocument(documentType string, hypervergeRequest HypervergeRequest, apiloggerBuilder *apilogger.ApiDataBuilder) (*bytes.Buffer, error) {
	url := getURLFor(documentType, hypervergeImpl.config.GetHypervergeEndpoint())
	req, err := newfileUploadRequest(url, hypervergeRequest.ImageFile,
		hypervergeImpl.config.GetHypervergeAppKey(), hypervergeImpl.config.GetHypervergeAppID(), hypervergeRequest.ImageName, hypervergeRequest.TxnId, apiloggerBuilder)
	if err != nil {
		return nil, errors.Wrap(ErrUploadError, err.Error())
	}

	resp, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

func (hypervergeImpl *HypervergeImpl) ReadPan(ctx context.Context, hypervergeRequest HypervergeRequest) (*PanResponse, error) {
	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig())
	apiloggerBuilder.WithBasic(ctx, HYPERVERGE, constants.OCRReadPan)
	defer func() {
		hypervergeImpl.apilogger.Log(context.Background(), apiloggerBuilder.Build())
	}()
	body, err := hypervergeImpl.readDocument("pan", hypervergeRequest, apiloggerBuilder)
	if err != nil {
		apiloggerBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrHttpError, err.Error())
	}
	var hypervergePanResponse HypervergePanResponse
	err = json.Unmarshal(body.Bytes(), &hypervergePanResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", body.String())
		return nil, err
	}
	apiloggerBuilder.WithResponse(hypervergePanResponse.StatusCode, body.String())
	if hypervergePanResponse.StatusCode != "200" {
		logger.Error(ctx, "ReadPan:: txnId : %s,response from Hyeperverge %s", hypervergeRequest.TxnId, body.String())
		var mappedErr error
		switch hypervergePanResponse.StatusCode {
		case "437":
			mappedErr = errors.Wrap(ErrBlurredImage, hypervergePanResponse.Error)
		case "432":
			mappedErr = errors.Wrap(ErrTemperedImage, hypervergePanResponse.Error)
		case "422":
			mappedErr = errors.Wrap(ErrInvalidDoc, hypervergePanResponse.Error)
		default:
			mappedErr = fmt.Errorf("status %v errorMessage %v", hypervergePanResponse.Status, hypervergePanResponse.Error)
		}
		apiloggerBuilder.WithError(mappedErr.Error())
		return nil, mappedErr
	}

	details := hypervergePanResponse.Result[0].Details
	panResponse := PanResponse{
		Date:        details.Date.Value,
		Father:      details.Father.Value,
		Name:        details.Name.Value,
		PanNo:       details.PanNo.Value,
		DateOfIssue: details.DateOfIssue.Value,
		RawResponse: body.String(),
	}

	return &panResponse, err

}

func (hypervergeImpl *HypervergeImpl) ReadAadhar(ctx context.Context, hypervergeRequest HypervergeRequest) (*AadharResponse, error) {
	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig())
	apiloggerBuilder.WithBasic(ctx, HYPERVERGE, constants.OCRReadAadhar)
	defer func() {
		hypervergeImpl.apilogger.Log(context.Background(), apiloggerBuilder.Build())
	}()
	body, err := hypervergeImpl.readDocument("aadhar", hypervergeRequest, apiloggerBuilder)
	if err != nil {
		apiloggerBuilder.WithError(err.Error())
		return nil, err
	}
	var hypervergeAadharResponse HypervergeAadharResponse
	err = json.Unmarshal(body.Bytes(), &hypervergeAadharResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", body.String())
		return nil, err
	}
	apiloggerBuilder.WithResponse(hypervergeAadharResponse.StatusCode, aadharmasking.MaskAddharInResponseJson(body.String()))
	if hypervergeAadharResponse.StatusCode != "200" {
		logger.Error(ctx, "ReadAadhar:: txnId : %s,response from Hyeperverge %s", hypervergeRequest.TxnId, body.String())
		var mappedErr error
		switch hypervergeAadharResponse.StatusCode {
		case "437":
			mappedErr = errors.Wrap(ErrBlurredImage, hypervergeAadharResponse.Error)
		case "432":
			mappedErr = errors.Wrap(ErrTemperedImage, hypervergeAadharResponse.Error)
		case "422":
			mappedErr = errors.Wrap(ErrInvalidDoc, hypervergeAadharResponse.Error)
		default:
			mappedErr = fmt.Errorf("status %v errorMessage %v", hypervergeAadharResponse.Status, hypervergeAadharResponse.Error)
		}
		apiloggerBuilder.WithError(mappedErr.Error())
		return nil, mappedErr
	}
	if len(hypervergeAadharResponse.Result) > 1 {
		mergeAadharDetails(&hypervergeAadharResponse)
	}
	details := hypervergeAadharResponse.Result[0].Details
	aadharResponse := AadharResponse{
		Aadhaar:     details.Aadhaar.Value,
		Dob:         details.Dob.Value,
		Father:      details.Father.Value,
		Gender:      details.Gender.Value,
		Mother:      details.Mother.Value,
		Name:        details.Name.Value,
		Yob:         details.Yob.Value,
		Husband:     details.Husband.Value,
		Phone:       details.Phone.Value,
		Pin:         details.Pin.Value,
		CareOf:      details.Address.CareOf,
		District:    details.Address.District,
		City:        details.Address.City,
		Locality:    details.Address.Locality,
		Landmark:    details.Address.Landmark,
		Street:      details.Address.Street,
		Line1:       details.Address.Line1,
		Line2:       details.Address.Line2,
		HouseNumber: details.Address.HouseNumber,
		State:       details.Address.State,
		AddressPin:  details.Address.Pin,
		RawResponse: body.String(),
	}

	return &aadharResponse, err

}

func (hypervergeImpl *HypervergeImpl) ReadPassport(ctx context.Context, hypervergeRequest HypervergeRequest) (*PassportResponse, error) {
	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig())
	apiloggerBuilder.WithBasic(ctx, HYPERVERGE, constants.OCRReadPassport)
	defer func() {
		hypervergeImpl.apilogger.Log(context.Background(), apiloggerBuilder.Build())
	}()
	body, err := hypervergeImpl.readDocument("passport", hypervergeRequest, apiloggerBuilder)
	if err != nil {
		apiloggerBuilder.WithError(err.Error())
		return nil, err
	}
	var hypervergePassportResponse HypervergePassportResponse
	err = json.Unmarshal(body.Bytes(), &hypervergePassportResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", body.String())
		return nil, err
	}
	apiloggerBuilder.WithResponse(hypervergePassportResponse.StatusCode, body.String())
	if hypervergePassportResponse.StatusCode != "200" {
		logger.Error(ctx, "ReadPassport:: txnId : %s,response from Hyeperverge %s", hypervergeRequest.TxnId, body.String())
		var mappedErr error
		switch hypervergePassportResponse.StatusCode {
		case "437":
			mappedErr = errors.Wrap(ErrBlurredImage, hypervergePassportResponse.Error)
		case "432":
			mappedErr = errors.Wrap(ErrTemperedImage, hypervergePassportResponse.Error)
		case "422":
			mappedErr = errors.Wrap(ErrInvalidDoc, hypervergePassportResponse.Error)
		default:
			mappedErr = fmt.Errorf("status %v errorMessage %v", hypervergePassportResponse.Status, hypervergePassportResponse.Error)
		}
		apiloggerBuilder.WithError(mappedErr.Error())
		return nil, mappedErr
	}
	details := hypervergePassportResponse.Result[0].Details
	passportResponse := PassportResponse{
		CountryCode:     details.CountryCode.Value,
		Dob:             details.Dob.Value,
		Doe:             details.Doe.Value,
		Doi:             details.Doi.Value,
		Gender:          details.Gender.Value,
		GivenName:       details.GivenName.Value,
		Nationality:     details.Nationality.Value,
		PassportNum:     details.PassportNum.Value,
		PlaceOfBirth:    details.PlaceOfBirth.Value,
		PlaceOfIssue:    details.PlaceOfIssue.Value,
		Surname:         details.Surname.Value,
		Mrz:             details.Mrz.Line1,
		Type:            details.Type.Value,
		District:        details.Address.District,
		City:            details.Address.City,
		Locality:        details.Address.Locality,
		Landmark:        details.Address.Landmark,
		Street:          details.Address.Street,
		Line1:           details.Address.Line1,
		Line2:           details.Address.Line2,
		HouseNumber:     details.Address.HouseNumber,
		State:           details.Address.State,
		Father:          details.Father.Value,
		Mother:          details.Mother.Value,
		FileNum:         details.FileNum.Value,
		OldDoi:          details.OldDoi.Value,
		OldPassportNum:  details.OldPassportNum.Value,
		OldPlaceOfIssue: details.OldPlaceOfIssue.Value,
		Pin:             details.Pin.Value,
		Spouse:          details.Spouse.Value,
		AddressPin:      details.Address.Pin,
		RawResponse:     body.String(),
	}

	return &passportResponse, err
}

func (hypervergeImpl *HypervergeImpl) ReadVotedID(ctx context.Context, hypervergeRequest HypervergeRequest) (*VoterIdResponse, error) {
	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig())
	apiloggerBuilder.WithBasic(ctx, HYPERVERGE, constants.OCRReadVoterId)
	defer func() {
		hypervergeImpl.apilogger.Log(context.Background(), apiloggerBuilder.Build())
	}()
	body, err := hypervergeImpl.readDocument("voter", hypervergeRequest, apiloggerBuilder)
	if err != nil {
		apiloggerBuilder.WithError(err.Error())
		return nil, err
	}
	var hypervergeVoterIdResponse HypervergeVoterIdResponse
	err = json.Unmarshal(body.Bytes(), &hypervergeVoterIdResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: "+err.Error()).
			WithResponse("", body.String())
		return nil, err
	}
	apiloggerBuilder.WithResponse(hypervergeVoterIdResponse.StatusCode, body.String())
	if hypervergeVoterIdResponse.StatusCode != "200" {
		logger.Error(ctx, "ReadVotedID:: txnId : %s,response from Hyeperverge %s", hypervergeRequest.TxnId, body.String())
		var mappedErr error
		switch hypervergeVoterIdResponse.StatusCode {
		case "437":
			mappedErr = errors.Wrap(ErrBlurredImage, hypervergeVoterIdResponse.Error)
		case "432":
			mappedErr = errors.Wrap(ErrTemperedImage, hypervergeVoterIdResponse.Error)
		case "422":
			mappedErr = errors.Wrap(ErrInvalidDoc, hypervergeVoterIdResponse.Error)
		default:
			mappedErr = fmt.Errorf("status %v errorMessage %v", hypervergeVoterIdResponse.Status, hypervergeVoterIdResponse.Error)
		}
		apiloggerBuilder.WithError(mappedErr.Error())
		return nil, mappedErr
	}
	details := hypervergeVoterIdResponse.Result[0].Details
	voterIdResponse := VoterIdResponse{
		Voterid:     details.Voterid.Value,
		Name:        details.Name.Value,
		Gender:      details.Gender.Value,
		Relation:    details.Relation.Value,
		Dob:         details.Dob.Value,
		Doc:         details.Doc.Value,
		Age:         details.Age.Value,
		Pin:         details.Pin.Value,
		Date:        details.Date.Value,
		Type:        details.Type.Value,
		District:    details.Address.District,
		City:        details.Address.City,
		Locality:    details.Address.Locality,
		Street:      details.Address.Street,
		Line1:       details.Address.Line1,
		Line2:       details.Address.Line2,
		HouseNumber: details.Address.HouseNumber,
		State:       details.Address.State,
		AddressPin:  details.Address.Pin,
		RawResponse: body.String(),
	}

	return &voterIdResponse, err
}

func newfileUploadRequest(uri, file, appKey, appID, fileName, txnID string, apiloggerBuilder *apilogger.ApiDataBuilder) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return nil, err
	}
	src := bytes.NewReader([]byte(file))
	_, err = io.Copy(part, src)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("appId", appID)
	req.Header.Set("appkey", appKey)
	req.Header.Set("transactionId", txnID)
	if apiloggerBuilder != nil {
		apiloggerBuilder.WithRequest(uri, http.MethodPost, "", map[string]string{
			"transactionId": txnID,
		})
	}
	return req, err
}

func getURLFor(docType, baseURL string) string {
	var path string
	switch docType {
	case "aadhar":
		path = "/readAadhaar"
	case "passport":
		path = "/readPassport"
	case "pan":
		path = "/readPAN"
	case "voter":
		path = "/readVoterID"
	default:
		path = "/readKYC"
	}
	return baseURL + path
}

func mergeAadharDetails(hypervergeAadharResult *HypervergeAadharResponse) {
	size := len(hypervergeAadharResult.Result)
	for i := size - 2; i >= 0; i-- {
		if hypervergeAadharResult.Result[i].Details.Aadhaar.Value == "" {
			hypervergeAadharResult.Result[i].Details.Aadhaar.Value = hypervergeAadharResult.Result[i+1].Details.Aadhaar.Value
		}
		if hypervergeAadharResult.Result[i].Details.Dob.Value == "" {
			hypervergeAadharResult.Result[i].Details.Dob.Value = hypervergeAadharResult.Result[i+1].Details.Dob.Value
		}
		if hypervergeAadharResult.Result[i].Details.Pin.Value == "" {
			hypervergeAadharResult.Result[i].Details.Pin.Value = hypervergeAadharResult.Result[i+1].Details.Pin.Value
		}
		if hypervergeAadharResult.Result[i].Details.Name.Value == "" {
			hypervergeAadharResult.Result[i].Details.Name.Value = hypervergeAadharResult.Result[i+1].Details.Name.Value
		}
	}
}

func (hypervergeImpl *HypervergeImpl) addHeaders(req *http.Request, txnID string) {
	req.Header.Add("appId", hypervergeImpl.config.GetHypervergeAppID())
	req.Header.Add("appKey", hypervergeImpl.config.GetHypervergeAppKey())
	req.Header.Add("transactionId", txnID)
	req.Header.Add("Content-Type", "application/json")
}

func (hypervergeImpl *HypervergeImpl) FraudCheckPan(ctx context.Context, fraudCheckPanRequest FraudCheckPanRequest, txnID string) (*FraudCheckPanResponse, error) {
	url := hypervergeImpl.config.GetHypervergeFraudCheckEndpoint() + "/verifyPAN"
	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig()).
		WithBasic(ctx, HYPERVERGE, constants.FraudCheckPan).
		WithRequest(
			url,
			http.MethodPost,
			fmt.Sprintf("%+v", fraudCheckPanRequest),
			map[string]string{"transactionId": txnID},
		)

	// Ensure logging happens at the end using original context
	defer func() {
		hypervergeImpl.apilogger.Log(ctx, apiloggerBuilder.Build())
	}()

	reqObj, err := json.Marshal(fraudCheckPanRequest)
	if err != nil {
		apiloggerBuilder.WithError("failed to marshal request: " + err.Error())
		return nil, err
	}

	requestBody := bytes.NewBuffer(reqObj)
	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		apiloggerBuilder.WithError("failed to create request: " + err.Error())
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)

	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		apiloggerBuilder.WithError("HTTP call failed: " + err.Error())
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	defer res.Body.Close()
	if err != nil {
		apiloggerBuilder.WithError("failed to read response body: " + err.Error())
		return nil, err
	}

	// Log raw response for traceability
	apiloggerBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())

	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		logger.Error(ctx, "FraudCheckPan:: txnId : %s, response from Hyperverge: %s", txnID, body.String())
		apiloggerBuilder.WithError("non-200 response: " + err.Error())
		return nil, err
	}

	var fraudCheckPanResponse FraudCheckPanResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckPanResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	return &fraudCheckPanResponse, nil
}

func (hypervergeImpl *HypervergeImpl) FraudCheckPanV2(ctx context.Context, NSDLPanRequest NSDLPanRequest, txnID string) (*NSDLPanResponse, error) {
	url := hypervergeImpl.config.GetHypervergeNSDLUrl() + "/NSDLPanVerification"

	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig()).
		WithBasic(ctx, HYPERVERGE, constants.FraudCheckPanV2). // Assuming UserID is part of NSDLPanRequest
		WithRequest(
			url,
			http.MethodPost,
			fmt.Sprintf("%+v", NSDLPanRequest),
			map[string]string{"transactionId": txnID},
		)

	defer func() {
		hypervergeImpl.apilogger.Log(ctx, apiloggerBuilder.Build())
	}()

	reqObj, err := json.Marshal(NSDLPanRequest)
	if err != nil {
		apiloggerBuilder.WithError("failed to marshal request: " + err.Error())
		return nil, err
	}

	requestBody := bytes.NewBuffer(reqObj)
	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		apiloggerBuilder.WithError("failed to create request: " + err.Error())
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)

	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		apiloggerBuilder.WithError("HTTP call failed: " + err.Error())
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	defer res.Body.Close()
	if err != nil {
		apiloggerBuilder.WithError("failed to read response body: " + err.Error())
		return nil, err
	}

	apiloggerBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())

	err = hypervergeImpl.handlNSDLErrorStatusCode(res)
	if err != nil {
		logger.Error(ctx, "FraudCheckPanV2:: txnId : %s, response from Hyperverge: %s", txnID, body.String())
		apiloggerBuilder.WithError("non-200 response: " + err.Error())
		return nil, err
	}

	var nsdlPanResp NSDLPanResponse
	err = json.Unmarshal(body.Bytes(), &nsdlPanResp)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	return &nsdlPanResp, nil
}

func (hypervergeImpl HypervergeImpl) handlFruadCheckErrorStatusCode(res *http.Response) error {
	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 422:
			return errors.Wrap(ErrInvalidDocID, fmt.Sprintf("%d  %v", res.StatusCode, res))
		case 400:
			return errors.Wrap(ErrBadRequest, fmt.Sprintf("%d %v", res.StatusCode, res))
		case 500:
			return errors.Wrap(ErrSomethingWentWrong, fmt.Sprintf("%d %v", res.StatusCode, res))
		case 401:
			return errors.Wrap(ErrFruadCheckUnauthrised, fmt.Errorf("%d %v", res.StatusCode, res).Error())
		default:
			return errors.Wrap(ErrSomethingWentWrong, fmt.Sprintf("%d %v", res.StatusCode, res))
		}
	}
	return nil
}

func (hypervergeImpl HypervergeImpl) handlNSDLErrorStatusCode(res *http.Response) error {
	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 404:
			return errors.Wrap(ErrInvalidDocID, fmt.Sprintf("%d  %v", res.StatusCode, res))
		case 400:
			return errors.Wrap(ErrBadRequest, fmt.Sprintf("%d %v", res.StatusCode, res))
		case 500:
			return errors.Wrap(ErrSomethingWentWrong, fmt.Sprintf("%d %v", res.StatusCode, res))
		case 401:
			return errors.Wrap(ErrFruadCheckUnauthrised, fmt.Errorf("%d %v", res.StatusCode, res).Error())
		default:
			return errors.Wrap(ErrSomethingWentWrong, fmt.Sprintf("%d %v", res.StatusCode, res))
		}
	}
	return nil
}

func (hypervergeImpl *HypervergeImpl) FraudCheckDl(ctx context.Context, fraudCheckDlRequest FraudCheckDlRequest, txnID string) (*FraudCheckDlResponse, error) {
	url := hypervergeImpl.config.GetHypervergeFraudCheckEndpoint() + "/checkDL"

	// Create apilogger entry
	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig()).
		WithBasic(ctx, HYPERVERGE, constants.FraudCheckDl).
		WithRequest(
			url,
			http.MethodPost,
			fmt.Sprintf("%+v", fraudCheckDlRequest),
			map[string]string{"transactionId": txnID},
		)

	// Ensure logger logs at the end
	defer func() {
		hypervergeImpl.apilogger.Log(ctx, apiloggerBuilder.Build())
	}()

	// Marshal request
	reqObj, err := json.Marshal(fraudCheckDlRequest)
	if err != nil {
		apiloggerBuilder.WithError("failed to marshal request: " + err.Error())
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqObj))
	if err != nil {
		apiloggerBuilder.WithError("failed to create HTTP request: " + err.Error())
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)

	// Execute HTTP request
	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		apiloggerBuilder.WithError("HTTP request failed: " + err.Error())
		return nil, err
	}

	// Read response body
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	defer res.Body.Close()
	if err != nil {
		apiloggerBuilder.WithError("failed to read response body: " + err.Error())
		return nil, err
	}

	// Log raw response
	apiloggerBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())

	// Handle non-200 status codes
	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		logger.Error(ctx, "FraudCheckDl:: txnId : %s, response from Hyperverge: %s", txnID, body.String())
		apiloggerBuilder.WithError("non-200 response: " + err.Error())
		return nil, err
	}

	// Unmarshal final response
	var fraudCheckDlResponse FraudCheckDlResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckDlResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	return &fraudCheckDlResponse, nil
}

func (hypervergeImpl *HypervergeImpl) FraudCheckVoter(ctx context.Context, fraudCheckVoterRequest FraudCheckVoterRequest, txnID string) (*FraudCheckVoterResponse, error) {
	url := hypervergeImpl.config.GetHypervergeFraudCheckEndpoint() + "/checkVoterId"

	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig()).
		WithBasic(ctx, HYPERVERGE, constants.FraudCheckVoter).
		WithRequest(
			url,
			http.MethodPost,
			fmt.Sprintf("%+v", fraudCheckVoterRequest),
			map[string]string{"transactionId": txnID},
		)

	defer func() {
		hypervergeImpl.apilogger.Log(ctx, apiloggerBuilder.Build())
	}()

	reqObj, err := json.Marshal(fraudCheckVoterRequest)
	if err != nil {
		apiloggerBuilder.WithError("failed to marshal request: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqObj))
	if err != nil {
		apiloggerBuilder.WithError("failed to create HTTP request: " + err.Error())
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)

	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		apiloggerBuilder.WithError("HTTP request failed: " + err.Error())
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	defer res.Body.Close()
	if err != nil {
		apiloggerBuilder.WithError("failed to read response body: " + err.Error())
		return nil, err
	}

	apiloggerBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())

	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		logger.Error(ctx, "FraudCheckVoter:: txnId : %s, response from Hyperverge: %s", txnID, body.String())
		apiloggerBuilder.WithError("non-200 response: " + err.Error())
		return nil, err
	}

	var fraudCheckVoterResponse FraudCheckVoterResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckVoterResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	return &fraudCheckVoterResponse, nil
}

func (hypervergeImpl *HypervergeImpl) FraudCheckPassport(ctx context.Context, fraudCheckPassportRequest FraudCheckPassportRequest, txnID string) (*FraudCheckPassportResponse, error) {
	url := hypervergeImpl.config.GetHypervergeFraudCheckEndpoint() + "/verifyPassport"

	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig()).
		WithBasic(ctx, HYPERVERGE, constants.FraudCheckPassport).
		WithRequest(
			url,
			http.MethodPost,
			fmt.Sprintf("%+v", fraudCheckPassportRequest),
			map[string]string{"transactionId": txnID},
		)

	defer func() {
		hypervergeImpl.apilogger.Log(ctx, apiloggerBuilder.Build())
	}()

	reqObj, err := json.Marshal(fraudCheckPassportRequest)
	if err != nil {
		apiloggerBuilder.WithError("failed to marshal request: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqObj))
	if err != nil {
		apiloggerBuilder.WithError("failed to create HTTP request: " + err.Error())
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)

	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		apiloggerBuilder.WithError("HTTP request failed: " + err.Error())
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	defer res.Body.Close()
	if err != nil {
		apiloggerBuilder.WithError("failed to read response body: " + err.Error())
		return nil, err
	}

	apiloggerBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), body.String())

	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		logger.Error(ctx, "FraudCheckPassport:: txnId : %s, response from Hyperverge: %s", txnID, body.String())
		apiloggerBuilder.WithError("non-200 response: " + err.Error())
		return nil, err
	}

	var fraudCheckPassportResponse FraudCheckPassportResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckPassportResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, err
	}

	return &fraudCheckPassportResponse, nil
}

func (hypervergeImpl *HypervergeImpl) FraudCheckAadhar(ctx context.Context, fraudCheckAadharRequest FraudCheckAadharRequest, txnID string) (*FraudCheckAadharResponse, string, error) {
	url := hypervergeImpl.config.GetHypervergeFraudCheckEndpoint() + "/verifyAadhaar"

	apiloggerBuilder := apilogger.NewApiDataBuilder(hypervergeImpl.apilogger.GetConfig()).
		WithBasic(ctx, HYPERVERGE, "FraudCheckAadhar").
		WithRequest(
			url,
			http.MethodPost,
			aadharmasking.MaskAddharInResponseJson(fmt.Sprintf("%+v", fraudCheckAadharRequest)),
			map[string]string{"transactionId": txnID},
		)

	defer func() {
		hypervergeImpl.apilogger.Log(ctx, apiloggerBuilder.Build())
	}()

	reqObj, err := json.Marshal(fraudCheckAadharRequest)
	if err != nil {
		apiloggerBuilder.WithError("failed to marshal request: " + err.Error())
		return nil, fmt.Sprintf("%s,%+v", HYPERVERGE, UNABLE_TO_SEND_REQUEST), err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqObj))
	if err != nil {
		apiloggerBuilder.WithError("failed to create HTTP request: " + err.Error())
		return nil, fmt.Sprintf("%s,%+v", HYPERVERGE, UNABLE_TO_SEND_REQUEST), err
	}

	hypervergeImpl.addHeaders(req, txnID)

	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		apiloggerBuilder.WithError("HTTP request failed: " + err.Error())
		return nil, fmt.Sprintf("%s,%+v", HYPERVERGE, UNABLE_TO_SEND_REQUEST), err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	defer res.Body.Close()
	if err != nil {
		apiloggerBuilder.WithError("failed to read response body: " + err.Error())
		return nil, fmt.Sprintf("%s,%+v", HYPERVERGE, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}

	apiloggerBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), aadharmasking.MaskAddharInResponseJson(body.String()))

	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		logger.Error(ctx, "FraudCheckAadhar:: txnId : %s, response from Hyperverge: %s", txnID, body.String())
		apiloggerBuilder.WithError("non-200 response: " + err.Error())
		return nil, fmt.Sprintf("%s,%+v", HYPERVERGE, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}

	var fraudCheckAadharResponse FraudCheckAadharResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckAadharResponse)
	if err != nil {
		apiloggerBuilder.WithError("failed to unmarshal response: " + err.Error())
		return nil, fmt.Sprintf("%s,%+v", HYPERVERGE, UNABLE_TO_PARSE_VENDOR_RESPONSE), err
	}

	return &fraudCheckAadharResponse, fmt.Sprintf("%s,%+v", HYPERVERGE, fraudCheckAadharResponse), nil
}
