package aadharlite

import (
	"context"
)

// This SDK is the interface for AadharLite verification.

// Any vendor-specific SDK that performs AadharLite verification must implement this interface.
type aadharLiteSDK interface {
	ProcessAadharLite(ctx context.Context, aadharNumber string) (*FraudCheckAadharResponse, error)
	HealthCheckAadharLite(ctx context.Context) (*HealthCheckResponse, error)
}
