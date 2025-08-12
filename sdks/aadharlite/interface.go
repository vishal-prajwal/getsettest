package aadhaarlite

import (
	"context"
)

// This SDK is the interface for AadharLite verification.

// Any vendor-specific SDK that performs AadharLite verification must implement this interface.
type SDK interface {
	ProcessAadharLite(ctx context.Context, aadharNumber string) (*FraudCheckAadharResponse, string, error)
	HealthCheckAadhaarVerify(ctx context.Context) (*HealthCheckResponse, error)
}
