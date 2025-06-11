package surepass

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/sdks/okyc"
)

const (
	NAME = "surepass"
)

type Config struct {
	Endpoint string
	Timeout  int // in seconds
	APIKey   string
}

type SurepassSDK struct {
	config *Config
	okyc.OKYCSDK
}

func New(config *Config) *SurepassSDK {
	return &SurepassSDK{
		config: config,
	}
}

// Name returns the name of the SDK
func (sdk *SurepassSDK) Name() string {
	return NAME
}

func (sdk *SurepassSDK) HealthCheck(ctx context.Context) (*okyc.HealthCheckResponse, error) {
	// Implement the health check logic here
	return &okyc.HealthCheckResponse{
		Percentage: 100,
		Available:  true,
	}, nil
}
func (sdk *SurepassSDK) GenerateOTP(ctx context.Context, req *okyc.GenerateOTPRequest) (*okyc.GenerateOTPResponse, error) {
	// Implement the logic to generate OTP here
	return &okyc.GenerateOTPResponse{
		ReferenceId:  "ref123",
		RequestId:    "req123",
		MobileNumber: "1234567890", // Example mobile number
	}, nil
}
func (sdk *SurepassSDK) ValidateOTP(ctx context.Context, req *okyc.ValidateOTPRequest) (*okyc.ValidateOTPResponse, error) {
	// Implement the logic to validate OTP here
	return &okyc.ValidateOTPResponse{
		Image:    "base64ImageString",
		DOB:      "01-01-1990",
		FullName: "John Doe",
	}, nil
}
