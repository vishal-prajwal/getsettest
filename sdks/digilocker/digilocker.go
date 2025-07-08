package digilocker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/apilogger"
	httpclient "bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/sdks/constants"
	nrf "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

type DigilockerImpl struct {
	appId       string
	appKey      string
	hvEndpoint  string
	httpClient  httpclient.HTTPClient
	redirectURL string
	apilogger   apilogger.ApiUsageLogger
}

// New creates a new digilocker client
func New(appId, appKey string, hvEndpoint string, client httpclient.HTTPClient, redirectURL string, apiLogger apilogger.ApiUsageLogger) *DigilockerImpl {
	dl := DigilockerImpl{
		appId:       appId,
		appKey:      appKey,
		hvEndpoint:  hvEndpoint,
		httpClient:  client,
		redirectURL: redirectURL,
		apilogger:   apiLogger,
	}
	return &dl
}

func (dl *DigilockerImpl) GetRedirectURL(ctx context.Context) string {
	return dl.redirectURL
}

func (dl *DigilockerImpl) addHeaders(req *http.Request) {
	req.Header.Add("appId", dl.appId)
	req.Header.Add("appKey", dl.appKey)
	req.Header.Add("Content-Type", "application/json")

}
func (dl *DigilockerImpl) StartKYC(ctx context.Context, transactionId, referenceId, redirectURL string) (*KYCStartDetails, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(dl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, HYPERVERGE, constants.DigilockerStartKYC)
	defer nrf.FromContext(ctx).StartSegment("StartKYC").End()
	defer func() {
		dl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	url := dl.hvEndpoint + "/api/digilocker/start"
	method := "POST"
	apiDataBuilder.WithRequest(url, method, "", map[string]string{
		"transactionId": transactionId,
		"referenceId":   referenceId,
	})

	reqObj, _ := json.Marshal(&KYCStartRequest{ReferenceId: referenceId, RedirectURL: redirectURL})
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}
	dl.addHeaders(req)
	req.Header.Add("transactionId", transactionId)

	res, err := dl.httpClient.Do(req)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrCallingHyperverge, err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrReadResponseBody, err.Error())
	}
	logger.Info(ctx, fmt.Sprintf("/api/digilocker/start response =>  %s", string(body)))
	var result KYCStartResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrUnmarshalJson, err.Error())
	}
	apiDataBuilder.WithResponse(result.StatusCode, string(body))
	switch result.StatusCode {
	case "200":
		return &result.Result, nil
	case "400":
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrReqValidate, result.Error.Message)
	case "500":
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	}
	return &result.Result, nil
}

func (dl *DigilockerImpl) CheckAccountstatus(ctx context.Context, mobile, aadhaar string) (*AccountStatusDetails, error) {
	defer nrf.FromContext(ctx).StartSegment("CheckAccountstatus").End()
	url := dl.hvEndpoint + "/api/digilocker/accountStatus"
	method := "POST"

	reqObj, _ := json.Marshal(&AccountStatusRequest{Mobile: mobile, Aadhaar: aadhaar})
	payload := strings.NewReader(string(reqObj))

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}
	dl.addHeaders(req)

	res, err := dl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrCallingHyperverge, err.Error())

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadResponseBody, err.Error())
	}
	logger.Info(ctx, fmt.Sprintf("/api/digilocker/accountStatus response =>  %s", string(body)))

	var result AccountStatusResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrap(ErrUnmarshalJson, err.Error())
	}
	switch result.StatusCode {
	case "200":
		return &result.Result, nil
	case "400":
		return nil, errors.Wrap(ErrReqValidate, result.Error.Message)
	case "500":
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	}
	return &result.Result, nil
}

func (dl *DigilockerImpl) GetAddharDetails(ctx context.Context, transactionId, referenceId string) (*AadhaarDetails, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(dl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, HYPERVERGE, constants.DigilockerGetAadhaarDetails)
	defer func() {
		dl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	defer nrf.FromContext(ctx).StartSegment("GetAddharDetails").End()
	url := dl.hvEndpoint + "/api/digilocker/eAadhaarDetails"
	method := "POST"
	apiDataBuilder.WithRequest(url, method, "", map[string]string{
		"transactionId": transactionId,
		"referenceId":   referenceId,
	})

	reqObj, _ := json.Marshal(&EAadhaarDetailsRequest{ReferenceId: referenceId, AadhaarFile: "yes"})
	payload := strings.NewReader(string(reqObj))

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}
	dl.addHeaders(req)
	req.Header.Add("transactionId", transactionId)

	res, err := dl.httpClient.Do(req)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrCallingHyperverge, err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrReadResponseBody, err.Error())
	}
	logger.Info(ctx, fmt.Sprintf("/api/digilocker/eAadhaarDetails response =>  %s", string(body)))

	var result AadhaarDetailsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrUnmarshalJson, err.Error())
	}
	apiDataBuilder.WithResponse(result.StatusCode, "")

	if result.Error.Code == "ER_CONSENT_MISSING" {
		return nil, errors.Wrap(ErrConsentNotProvided, result.Error.Message)
	}

	switch result.StatusCode {
	case "200":
		err = result.Result.SethPinCodeFromXmlFile()
		if err != nil {
			apiDataBuilder.WithError(fmt.Sprintf("error setting pincode from xml file: %v", err))
			return nil, errors.Wrap(ErrExtractingXML, err.Error())
		}
		return &result.Result, nil
	case "504":
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	case "500":
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	default:
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	}

}

func (dl *DigilockerImpl) Healthcheck(ctx context.Context) (*HealthcheckResult, error) {
	defer nrf.FromContext(ctx).StartSegment(HV_HEALTHCHECK_CALL).End()
	url := dl.hvEndpoint + "/api/health/digilocker/eAadhaarDetails"

	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}
	dl.addHeaders(req)
	res, err := dl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrCallingHyperverge, err.Error())

	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.Wrap(ErrHVServer, fmt.Sprintf("status code %d", res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadResponseBody, err.Error())
	}
	logger.Info(ctx, fmt.Sprintf("/api/health/digilocker/eAadhaarDetails response =>  %s", string(body)))

	var result HypervergeHealthcheckResponse
	logger.Info(ctx, string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrap(ErrUnmarshalJson, err.Error())
	}
	return &result.Result, nil
}

func (dl *DigilockerImpl) GetPanDigilockerDoc(ctx context.Context, refId string) (*PanDetails, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(dl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, HYPERVERGE, constants.GetPanDigilockerDoc)
	defer func() {
		dl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	defer nrf.FromContext(ctx).StartSegment(HV_HEALTHCHECK_CALL).End()
	url := dl.hvEndpoint + "/api/digilocker/docDetails"

	appID := dl.appId
	appKey := dl.appKey
	payload := DigilockerDocumentsRequetsPayload{
		ReferenceID: refId,
		PAN:         "yes",
		PANFile:     "yes",
		EnableRetry: "yes",
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		apiDataBuilder.WithError(fmt.Sprintf("error marshaling request payload: %v", err))
		return nil, fmt.Errorf("error marshaling request payload: %v", err)
	}

	apiDataBuilder.WithRequest(url, "POST", string(jsonPayload), nil)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		apiDataBuilder.WithError(fmt.Sprintf("error creating request: %v", err))
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set request headers
	req.Header.Set("appid", appID)
	req.Header.Set("appKey", appKey)
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	resp, err := dl.httpClient.Do(req)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrNotAvailable, err.Error())
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrReadResponseBody, err.Error())
	}

	// Unmarshal the response body into KYCResult
	var result DigilockerPanResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrUnmarshalJson, err.Error())
	}
	apiDataBuilder.WithResponse(result.StatusCode, "")
	if result.StatusCode != "200" {
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	}
	if len(result.Result.DocsFound) == 0 {
		apiDataBuilder.WithError("no documents found")
		return nil, ErrDocumentNotFound
	}
	if len(result.Result.Details[0].PAN) == 0 || len(result.Result.Details[0].Name) == 0 || len(result.Result.Details[0].DOB) == 0 || len(result.Result.Details[0].FileUrl) == 0 {
		apiDataBuilder.WithError(fmt.Sprintf("missing data in response: %+v", result.Result.Details[0]))
		return nil, errors.Wrap(ErrHVServerMissingData, fmt.Sprintf("Data recived %+v", result.Result.Details[0]))
	}
	logger.Info(ctx, "%v", result)
	panDetails := PanDetails{
		PanNumber:   result.Result.Details[0].PAN,
		Name:        result.Result.Details[0].Name,
		DOB:         result.Result.Details[0].DOB,
		Address:     "",
		DocImageUrl: result.Result.Details[0].FileUrl,
	}

	return &panDetails, nil

}

func (dl *DigilockerImpl) GetPanDetails(ctx context.Context, refId string, panNumber string, fullName string) (*PanDetails, error) {
	apiDataBuilder := apilogger.NewApiDataBuilder(dl.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, HYPERVERGE, constants.GetPanDetails)
	defer func() {
		dl.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()
	defer nrf.FromContext(ctx).StartSegment(HV_HEALTHCHECK_CALL).End()
	url := dl.hvEndpoint + "/api/digilocker/fetchDocuments"

	appID := dl.appId
	appKey := dl.appKey

	// Create the request payload
	payload := PanDetailsRequestPayload{
		ReferenceID: refId,
		Docs: []DocumentInfo{
			{
				DocID: "001891_PANCR",
				File:  "yes",
				Params: DocumentParams{
					PANNo:       panNumber,
					PANFullName: fullName,
				},
			},
		},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		apiDataBuilder.WithError(fmt.Sprintf("error marshaling request payload: %v", err))
		return nil, fmt.Errorf("error marshaling request payload: %v", err)
	}
	apiDataBuilder.WithRequest(url, "POST", string(jsonPayload), map[string]string{
		"referenceId": refId,
	})

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		apiDataBuilder.WithError(fmt.Sprintf("error creating request: %v", err))
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set request headers
	req.Header.Set("appid", appID)
	req.Header.Set("appKey", appKey)
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	resp, err := dl.httpClient.Do(req)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrNotAvailable, err.Error())
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		apiDataBuilder.WithError(err.Error())
		return nil, errors.Wrap(ErrReadResponseBody, err.Error())
	}

	// Unmarshal the response body into KYCResult
	var result KYCResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrap(ErrUnmarshalJson, err.Error())
	}
	apiDataBuilder.WithResponse(result.StatusCode, "")
	if result.StatusCode != "200" || len(result.Result) == 0 {
		apiDataBuilder.WithError(result.Error.Message)
		return nil, errors.Wrap(ErrHVServer, result.Error.Message)
	}
	logger.Info(ctx, "%v", result)
	details := result.Result[0].Data
	var panDetails PanDetails
	panDetails.Address = details.Address
	panDetails.DOB = details.DOB
	panDetails.DocImageUrl = details.File
	panDetails.Name = details.Name
	panDetails.PanNumber = details.Number

	return &panDetails, nil
}
