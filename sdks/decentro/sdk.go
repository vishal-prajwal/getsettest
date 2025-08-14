package decentro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/apilogger"
	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/monitoring"
	aadharlite "bitbucket.org/junglee_games/getsetgo/sdks/aadharlite"
	"bitbucket.org/junglee_games/getsetgo/sdks/constants"
	"bitbucket.org/junglee_games/getsetgo/sdks/okyc"
	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
)

const (
	// DefaultTimeout is the default timeout for API requests in seconds
	DefaultTimeout = 30
)

type Config struct {
	Endpoint     string `yaml:"Endpoint" json:"endpoint"`         // The base URL for the Decentro API
	Timeout      int    `yaml:"Timeout" json:"timeout"`           // Timeout for API requests in seconds
	ClientID     string `yaml:"ClientID" json:"clientId"`         // API key for authentication
	ClientSecret string `yaml:"ClientSecret" json:"clientSecret"` // API secret for authentication
}

func (c *Config) validate() error {
	if c.Endpoint == "" {
		return errors.New("endpoint cannot be empty")
	}
	if c.Timeout <= 0 {
		c.Timeout = DefaultTimeout // Set to default if not provided
	}
	if c.ClientID == "" {
		return errors.New("apikey cannot be empty")
	}
	if c.ClientSecret == "" {
		return errors.New("apisecret cannot be empty")
	}
	return nil
}

type SDK struct {
	config          *Config
	httpClient      *http.Client     // HTTP client for making API requests
	monitoringAgent monitoring.Agent // Optional monitoring agent
	apilogger       apilogger.ApiUsageLogger
}

type Options struct {
	// MonitoringAgent is an optional parameter to pass a monitoring agent for tracking time usage
	MonitoringAgent monitoring.Agent
}

func NewSDK(config *Config, ApiLogger apilogger.ApiUsageLogger, options ...Options) (*SDK, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}
	httpClient := httpclient.New(config.Timeout)
	sdk := &SDK{config: config,
		httpClient: &httpClient,
		apilogger:  ApiLogger,
	}
	if len(options) > 0 {
		sdk.monitoringAgent = options[0].MonitoringAgent
	}
	return sdk, nil
}

func (sdk *SDK) Name() okyc.Vendor {
	return okyc.VendorDecentro
}

func (sdk *SDK) addHeaders(req *http.Request) {
	req.Header.Add("client_id", sdk.config.ClientID)
	req.Header.Add("client_secret", sdk.config.ClientSecret)
	// Add any other headers required by Decentro API
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

func (sdk *SDK) HealthCheck(ctx context.Context) (*okyc.HealthCheckResponse, error) {
	if sdk.monitoringAgent != nil {
		defer sdk.monitoringAgent.StartTransaction("DecentroSDK.HealthCheck").End()
	}

	url := fmt.Sprintf("%s/decentro/read/health/status/aadhar_xml_health_check_api", sdk.config.Endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.HealthCheck:: failed to create health check request: %w", err)
	}

	sdk.addHeaders(req)

	res, err := sdk.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.HealthCheck:: failed to read decentro healthcheck response body: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		logger.Error(ctx, "DecentroSDK.HealthCheck:: health check failed with status code: %d, body: %s", res.StatusCode, body)
		return nil, fmt.Errorf("DecentroSDK.HealthCheck:: health check failed with status code: %d", res.StatusCode)
	}

	var healthCheckResponse DecentroHealthCheckResponse
	if err = json.Unmarshal(body, &healthCheckResponse); err != nil {
		logger.Error(ctx, "DecentroSDK.HealthCheck:: failed to unmarshal decentro healthcheck response: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
	if len(healthCheckResponse) == 0 {
		return nil, fmt.Errorf("DecentroSDK.HealthCheck:: empty health check response from Decentro: %s", body)
	}
	// As discssed with decentro team, response will always have single entry
	response, ok := healthStatusToAvailabilityMap[healthCheckResponse[0].Status]
	if !ok {
		return nil, fmt.Errorf("DecentroSDK.HealthCheck:: unknown health status: %s", healthCheckResponse[0].Status)
	}
	return &response, nil
}

func (sdk *SDK) GenerateOTP(ctx context.Context, req *okyc.GenerateOTPRequest) (*okyc.GenerateOTPResponse, error) {
	if sdk.monitoringAgent != nil {
		defer sdk.monitoringAgent.StartTransaction("DecentroSDK.GenerateOTP").End()
	}
	apiDataBuilder := apilogger.NewApiDataBuilder(sdk.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, DECENTRO, constants.GenerateAadhaarOTP)

	defer func() {
		sdk.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()

	url := fmt.Sprintf("%s/v2/kyc/aadhaar/otp", sdk.config.Endpoint)
	sendOTPReq := NewSendOTPRequest(req)
	body, err := json.Marshal(sendOTPReq)
	if err != nil {
		apiDataBuilder.WithError(errMarshalRequestBody + err.Error())
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to marshal request body: %w", err)
	}

	apiDataBuilder.WithRequest(url, http.MethodPost, "", map[string]string{
		"reference_id": sendOTPReq.ReferenceId,
		"purpose":      sendOTPReq.Purpose,
		"consent":      "true",
	})

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to create request: %w", err)
	}

	sdk.addHeaders(httpReq)

	res, err := sdk.httpClient.Do(httpReq)
	if err != nil {
		apiDataBuilder.WithError(errSendRequest + err.Error())
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to send request: %w", err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		apiDataBuilder.WithError(errReadResponseBody + err.Error())
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to read response body: %w", err)
	}

	var response DecentroResponse[any]

	err = json.Unmarshal(body, &response)
	if err != nil {
		apiDataBuilder.WithError(errUnmarshalResponseBody + err.Error())
		logger.Error(ctx, "Critical::DecentroSDK.GenerateOTP:: failed to unmarshal response body: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), string(body))

	if res.StatusCode != http.StatusOK {
		apiDataBuilder.WithError(fmt.Sprintf("failed to send OTP, status code: %d, body: %s", res.StatusCode, body))
		logger.Error(ctx, "DecentroSDK.GenerateOTP:: failed to send OTP, status code: %d, body: %s", res.StatusCode, body)
		return nil, responseKeyToError(response.ResponseKey)
	}

	return &okyc.GenerateOTPResponse{
		ReferenceId:   sendOTPReq.ReferenceId,
		TransactionId: response.DecentroTxnId,
	}, nil
}

func (sdk *SDK) ValidateOTP(ctx context.Context, req *okyc.ValidateOTPRequest) (*okyc.ValidateOTPResponse, error) {
	if sdk.monitoringAgent != nil {
		defer sdk.monitoringAgent.StartTransaction("DecentroSDK.ValidateOTP").End()
	}
	apiDataBuilder := apilogger.NewApiDataBuilder(sdk.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, DECENTRO, constants.ValidateAadhaarOTP)

	defer func() {
		sdk.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()

	url := fmt.Sprintf("%s/v2/kyc/aadhaar/otp/validate", sdk.config.Endpoint)

	body, err := json.Marshal(NewValidateOTPRequest(req))
	if err != nil {
		apiDataBuilder.WithError(errMarshalRequestBody + err.Error())
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to marshal request body: %w", err)
	}

	apiDataBuilder.WithRequest(url, http.MethodPost, string(body), map[string]string{
		"reference_id":   req.ReferenceId,
		"transaction_id": req.TransactionId,
		"otp":            req.OTP,
		"consent":        "true",
		"generate_pdf":   "true",
		"generate_xml":   "false",
		"purpose":        "For Aadhaar Verification",
	})

	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to create request: %w", err)
	}
	sdk.addHeaders(httpReq)
	res, err := sdk.httpClient.Do(httpReq)
	if err != nil {
		apiDataBuilder.WithError(errSendRequest + err.Error())
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		apiDataBuilder.WithError(errReadResponseBody + err.Error())
		logger.Error(ctx, "Critical::DecentroSDK.ValidateOTP:: failed to read response body: %v, body: %s", err, body)
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to read response body: %w", err)
	}
	var response DecentroResponse[*ValidateOTPResponseData]
	err = json.Unmarshal(body, &response)
	if err != nil {
		apiDataBuilder.WithError(errUnmarshalResponseBody + err.Error())
		logger.Error(ctx, "Critical::DecentroSDK.ValidateOTP:: failed to unmarshal response body: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), "")
	if res.StatusCode != http.StatusOK {
		logger.Error(ctx, "DecentroSDK.ValidateOTP:: failed to validate OTP, status code: %d, body: %s", res.StatusCode, body)
		return nil, responseKeyToError(response.ResponseKey)
	}
	if response.Data == nil {
		logger.Error(ctx, "DecentroSDK.ValidateOTP:: response data is nil, body: %s", body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
	pdfByteStr, err := response.Data.GetDocumentPDF(ctx)
	if err != nil {
		logger.Error(ctx, "DecentroSDK.ValidateOTP:: failed to get document PDF: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}

	return &okyc.ValidateOTPResponse{
		DocumentBytes:     pdfByteStr,
		DocumentExtension: "pdf",
		DOB:               response.Data.ProofOfIdentity.DOB,
		Gender:            response.Data.ProofOfIdentity.Gender.ToOKYCGender(),
		FullName:          response.Data.ProofOfIdentity.Name,
		Country:           response.Data.ProofOfAddress.Country,
		Pincode:           response.Data.ProofOfAddress.Pincode,
		State:             response.Data.ProofOfAddress.State,
		Address:           response.Data.GetAddress(),
	}, nil
}

// ProcessAadharLite performs an AadharLite check using Decentro and conforms to the aadharlite.aadharLiteSDK interface.
func (sdk *SDK) ProcessAadharLite(ctx context.Context, aadharNumber string) (*aadharlite.FraudCheckAadharResponse, error) {
	if sdk.monitoringAgent != nil {
		defer sdk.monitoringAgent.StartTransaction("DecentroSDK.ProcessAadharLite").End()
	}
	apiDataBuilder := apilogger.NewApiDataBuilder(sdk.apilogger.GetConfig())
	apiDataBuilder.WithBasic(ctx, DECENTRO, aadharlite.AADHARLITE)
	defer func() {
		sdk.apilogger.Log(context.Background(), apiDataBuilder.Build())
	}()

	url := fmt.Sprintf("%s/v2/kyc/aadhaar/verify", sdk.config.Endpoint)
	reqPayload := VerifyAadharRequest{
		ReferenceID:  uuid.New().String(),
		Consent:      true,
		Purpose:      "For Aadhaar Verification",
		AadharNumber: aadharNumber,
	}
	body, err := json.Marshal(reqPayload)
	if err != nil {
		apiDataBuilder.WithError(errMarshalRequestBody + err.Error())
		return nil, pkgErrors.Wrap(aadharlite.MarshalErr, "DecentroSDK.ProcessAadharLite:: failed to marshal request body")
	}

	apiDataBuilder.WithRequest(url, http.MethodPost, "", map[string]string{
		"reference_id": reqPayload.ReferenceID,
		"purpose":      reqPayload.Purpose,
		"consent":      "true",
	})

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, pkgErrors.Wrap(aadharlite.CreateRequestErr, "DecentroSDK.ProcessAadharLite:: failed to create request")
	}

	sdk.addHeaders(httpReq)

	res, err := sdk.httpClient.Do(httpReq)
	if err != nil {
		apiDataBuilder.WithError(errSendRequest + err.Error())
		return nil, pkgErrors.Wrap(aadharlite.SendRequestErr, "DecentroSDK.ProcessAadharLite:: failed to send request")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		apiDataBuilder.WithError(errReadResponseBody + err.Error())
		return nil, pkgErrors.Wrap(aadharlite.ReadResponseBodyErr, "DecentroSDK.ProcessAadharLite:: failed to read response body")
	}

	var response DecentroResponse[AadharData]
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		apiDataBuilder.WithError(errUnmarshalResponseBody + err.Error())
		logger.Error(ctx, "Critical::DecentroSDK.ProcessAadharLite:: failed to unmarshal response body: %v, body: %s", err, resBody)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
	apiDataBuilder.WithResponse(fmt.Sprintf("%d", res.StatusCode), string(resBody))

	// Harden the response processing by ensuring the data payload is valid.
	if response.Data.Status == "" {
		apiDataBuilder.WithError("response data is missing or invalid")
		logger.Error(ctx, "DecentroSDK.ProcessAadharLite:: response data is missing or invalid, body: %s", resBody)
		return nil, okyc.ErrInvalidResponseFromVendor
	}

	if res.StatusCode != http.StatusOK {
		apiDataBuilder.WithError(fmt.Sprintf("failed to process aadhar lite, status code: %d, body: %s", res.StatusCode, resBody))
		logger.Error(ctx, "DecentroSDK.ProcessAadharLite:: failed to process aadhar lite, status code: %d, body: %s", res.StatusCode, resBody)
		return nil, responseKeyToError(response.ResponseKey)
	}

	return toFraudCheckAadharResponse(&response), nil
}

// HealthCheckAadharLite checks the health of the Decentro Aadhar Verify service.
func (sdk *SDK) HealthCheckAadharLite(ctx context.Context) (*aadharlite.HealthCheckResponse, error) {
	if sdk.monitoringAgent != nil {
		defer sdk.monitoringAgent.StartTransaction("DecentroSDK.HealthCheckAadharLite").End()
	}

	url := fmt.Sprintf("%s/decentro/read/health/status/aadhaar_verify", sdk.config.Endpoint)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader("{}"))
	if err != nil {
		return nil, pkgErrors.Wrap(aadharlite.CreateRequestErr, "DecentroSDK.HealthCheckAadharLite:: failed to create health check request")
	}

	sdk.addHeaders(req)

	res, err := sdk.httpClient.Do(req)
	if err != nil {
		return nil, pkgErrors.Wrap(aadharlite.SendRequestErr, "DecentroSDK.HealthCheckAadharLite:: failed to send request")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, pkgErrors.Wrap(aadharlite.ReadResponseBodyErr, "DecentroSDK.HealthCheckAadharLite:: failed to read decentro healthcheck response body")
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		logger.Error(ctx, "DecentroSDK.HealthCheckAadharLite:: health check failed with status code: %d, body: %s", res.StatusCode, body)
		return nil, pkgErrors.Wrap(aadharlite.HealthCheckFailedErr, fmt.Sprintf("DecentroSDK.HealthCheckAadharLite:: health check failed with status code: %d", res.StatusCode))
	}

	var healthCheckResponse DecentroHealthCheckResponse
	if err = json.Unmarshal(body, &healthCheckResponse); err != nil {
		logger.Error(ctx, "DecentroSDK.HealthCheckAadharLite:: failed to unmarshal decentro healthcheck response: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
	if len(healthCheckResponse) == 0 {
		return nil, pkgErrors.Wrap(aadharlite.EmptyHealthCheckResponseErr, fmt.Sprintf("DecentroSDK.HealthCheckAadharLite:: empty health check response from Decentro: %s", body))
	}

	response, ok := healthStatusToAvailabilityMap[healthCheckResponse[0].Status]
	if !ok {
		return nil, pkgErrors.Wrap(aadharlite.UnknownHealthStatusErr, fmt.Sprintf("DecentroSDK.HealthCheckAadharLite:: unknown health status: %s", healthCheckResponse[0].Status))
	}
	return &aadharlite.HealthCheckResponse{
		Percentage: response.Percentage,
		Available:  response.Available,
	}, nil
}

// toFraudCheckAadharResponse translates the Decentro-specific response to the standard FraudCheckAadharResponse.
func toFraudCheckAadharResponse(resp *DecentroResponse[AadharData]) *aadharlite.FraudCheckAadharResponse {
	ageBandParts := strings.Split(resp.Data.AgeBand, "-")
	var ageBand aadharlite.AgeBand
	if len(ageBandParts) == 2 {
		ageBand.LowerLimit = ageBandParts[0]
		ageBand.UpperLimit = ageBandParts[1]
	}

	return &aadharlite.FraudCheckAadharResponse{
		Action:      "AADHAAR_LITE_VERIFICATION",
		CompletedAt: time.Now(),
		CreatedAt:   time.Now(),
		RequestID:   resp.DecentroTxnId,
		Result: aadharlite.FraudCheckResult{
			Data: aadharlite.FraudCheckAadharData{
				AgeBand:      ageBand,
				Gender:       resp.Data.Gender,
				MobileNumber: resp.Data.MaskedMobileNumber,
				State:        resp.Data.Address,
				Status:       resp.Data.AadharStatus,
			},
		},
		Status:  resp.Status,
		Type:    "AADHAR_LITE",
		Message: resp.Message,
	}
}
