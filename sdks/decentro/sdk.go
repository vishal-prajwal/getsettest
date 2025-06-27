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

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/monitoring"
	"bitbucket.org/junglee_games/getsetgo/sdks/okyc"
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
}

type Options struct {
	// MonitoringAgent is an optional parameter to pass a monitoring agent for tracking time usage
	MonitoringAgent monitoring.Agent
}

func NewSDK(config *Config, options ...Options) (*SDK, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}
	httpClient := httpclient.New(config.Timeout)
	sdk := &SDK{config: config,
		httpClient: &httpClient,
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

	url := fmt.Sprintf("%s/v2/kyc/aadhaar/otp", sdk.config.Endpoint)
	sendOTPReq := NewSendOTPRequest(req)
	body, err := json.Marshal(sendOTPReq)
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to create request: %w", err)
	}

	sdk.addHeaders(httpReq)

	res, err := sdk.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to send request: %w", err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.GenerateOTP:: failed to read response body: %w", err)
	}

	var response DecentroResponse[any]

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error(ctx, "Critical::DecentroSDK.GenerateOTP:: failed to unmarshal response body: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}

	if res.StatusCode != http.StatusOK {
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

	url := fmt.Sprintf("%s/v2/kyc/aadhaar/otp/validate", sdk.config.Endpoint)

	body, err := json.Marshal(NewValidateOTPRequest(req))
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to create request: %w", err)
	}
	sdk.addHeaders(httpReq)
	res, err := sdk.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		logger.Error(ctx, "Critical::DecentroSDK.ValidateOTP:: failed to read response body: %v, body: %s", err, body)
		return nil, fmt.Errorf("DecentroSDK.ValidateOTP:: failed to read response body: %w", err)
	}
	var response DecentroResponse[*ValidateOTPResponseData]
	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error(ctx, "Critical::DecentroSDK.ValidateOTP:: failed to unmarshal response body: %v, body: %s", err, body)
		return nil, okyc.ErrInvalidResponseFromVendor
	}
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
