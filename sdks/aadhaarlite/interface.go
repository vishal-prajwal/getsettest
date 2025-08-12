package aadhaarlite

import (
	"context"
	"time"
)

type SDK interface {
	ProcessAadharLite(ctx context.Context, aadharNumber string) (*FraudCheckAadharResponse, error)
	HealthCheck(ctx context.Context) (*HealthCheckResponse, error)
}

type HealthCheckResponse struct {
	Percentage int
	Available  bool
}

type AgeBand struct {
	LowerLimit string `json:"lower_limit"`
	UpperLimit string `json:"upper_limit"`
}

type FraudCheckAadharData struct {
	AgeBand      AgeBand `json:"age_band"`
	Gender       string  `json:"gender"`
	MobileNumber string  `json:"mobile_number"`
	State        string  `json:"state"`
	Status       string  `json:"status"`
}

type FraudCheckAadharResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		Data FraudCheckAadharData `json:"source_output"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
