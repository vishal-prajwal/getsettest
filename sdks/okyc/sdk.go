package okyc

import (
	"context"
	"fmt"
)

type GenerateOTPRequest struct {
	AadharNumber string
}
type ValidateOTPRequest struct {
	ReferenceId string
	RequestId   string
	OTP         string
}

type HealthCheckResponse struct {
	Percentage int
	Available  bool
}

type GenerateOTPResponse struct {
	ReferenceId  string
	RequestId    string
	MobileNumber string // this is the mobile number linked to the Aadhaar
}

// will be sending back our customized response and not the actual response comming from vendors
// we will be logging the actual response in APIUsageLogger (currently in dev)
type ValidateOTPResponse struct {
	Image    string
	DOB      string // dd-mm-yyyy
	Gender   string
	FullName string
	Country  string
	Pincode  string
	State    string
	Address  string
}

type OKYCSDK interface {
	// Name returns the name of the SDK
	Name() string
	HealthCheck(context.Context) (*HealthCheckResponse, error)
	GenerateOTP(context.Context, *GenerateOTPRequest) (*GenerateOTPResponse, error)
	ValidateOTP(context.Context, *ValidateOTPRequest) (*ValidateOTPResponse, error)
}

// we will be sending typed error from here based on the status and respective service will take care.
var (
	//following error will be sent when 429 is returned from surepass
	ErrOTPAlreadySent       = fmt.Errorf("otp already sent") // in this case we will send success to user
	ErrInvalidAadhar        = fmt.Errorf("invalid Aadhaar provided")
	ErrAadharSuspended      = fmt.Errorf("UIDAI has suspended the Aadhaar")
	ErrNoMobileNumberAadhar = fmt.Errorf("mobile is not lined to Aadhaar")
	ErrGenerateOTPFailure   = fmt.Errorf("error while generating OTP. Please try again")
	ErrUIDLockedAadhar      = fmt.Errorf("user has locked his/her Aadhaar. Cannot generate OTP")
)
