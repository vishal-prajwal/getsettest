package aadharlite

import "time"

// HealthCheckResponse represents a standardized health check response.
type HealthCheckResponse struct {
	Percentage int
	Available  bool
}

// AgeBand represents the age range from a verification response.
type AgeBand struct {
	LowerLimit string `json:"lower_limit"`
	UpperLimit string `json:"upper_limit"`
}

// FraudCheckAadharData contains the specific data points from an AadharLite check.
type FraudCheckAadharData struct {
	AgeBand      AgeBand `json:"age_band"`
	Gender       string  `json:"gender"`
	MobileNumber string  `json:"mobile_number"`
	State        string  `json:"state"`
	Status       string  `json:"status"`
}

// FraudCheckResult wraps the data from a fraud check.
type FraudCheckResult struct {
	Data FraudCheckAadharData `json:"source_output"`
}

// FraudCheckAadharResponse is the standardized internal response for an AadharLite check.
type FraudCheckAadharResponse struct {
	Action      string           `json:"action"`
	CompletedAt time.Time        `json:"completed_at"`
	CreatedAt   time.Time        `json:"created_at"`
	GroupID     string           `json:"group_id"`
	RequestID   string           `json:"request_id"`
	Result      FraudCheckResult `json:"result"`
	Status      string           `json:"status"`
	TaskID      string           `json:"task_id"`
	Type        string           `json:"type"`
	Error       string           `json:"error"`
	Message     string           `json:"message"`
}
