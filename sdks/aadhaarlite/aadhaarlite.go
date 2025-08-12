package aadhaarlite

import (
	"context"
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/sdks/decentro"
	"github.com/google/uuid"
)

type sdk struct {
	decentroSDK *decentro.SDK
}

func New(decentroSDK *decentro.SDK) SDK {
	return &sdk{
		decentroSDK: decentroSDK,
	}
}

func (s *sdk) ProcessAadharLite(ctx context.Context, aadharNumber string) (*FraudCheckAadharResponse, error) {
	req := &decentro.VerifyAadhaarRequest{
		ReferenceID:   uuid.New().String(),
		Consent:       true,
		Purpose:       "For Aadhaar Verification",
		AadhaarNumber: aadharNumber,
	}

	resp, err := s.decentroSDK.VerifyAadhaar(ctx, req)
	if err != nil {
		return nil, err
	}

	return toFraudCheckAadharResponse(resp, aadharNumber), nil
}

func (s *sdk) HealthCheck(ctx context.Context) (*HealthCheckResponse, error) {
	resp, err := s.decentroSDK.HealthCheckAadhaarVerify(ctx)
	if err != nil {
		return nil, err
	}

	return &HealthCheckResponse{
		Percentage: resp.Percentage,
		Available:  resp.Available,
	}, nil
}

func toFraudCheckAadharResponse(resp *decentro.DecentroResponse[decentro.AadhaarData], aadharNumber string) *FraudCheckAadharResponse {
	ageBandParts := strings.Split(resp.Data.AgeBand, "-")
	var ageBand AgeBand
	if len(ageBandParts) == 2 {
		ageBand.LowerLimit = ageBandParts[0]
		ageBand.UpperLimit = ageBandParts[1]
	}

	return &FraudCheckAadharResponse{
		Action:      "AADHAAR_LITE_VERIFICATION",
		CompletedAt: time.Now(),
		CreatedAt:   time.Now(),
		RequestID:   resp.DecentroTxnId,
		Result: struct {
			Data FraudCheckAadharData `json:"source_output"`
		}{
			Data: FraudCheckAadharData{
				AgeBand:      ageBand,
				Gender:       resp.Data.Gender,
				MobileNumber: resp.Data.MaskedMobileNumber,
				State:        resp.Data.Address,
				Status:       resp.Data.AadhaarStatus,
			},
		},
		Status:  resp.Status,
		Type:    "AADHAAR_LITE",
		Message: resp.Message,
	}
}
