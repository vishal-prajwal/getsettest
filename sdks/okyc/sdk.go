package okyc

import (
	"context"
	"fmt"
)

type Vendor int

const (
	VendorDecentro Vendor = iota
	VendorCashfree
	VendorSurepass

	DecentroString = "Decentro"
	CashfreeString = "Cashfree"
	SurepassString = "Surepass"
)

func (v Vendor) String() string {
	switch v {
	case VendorDecentro:
		return DecentroString
	case VendorCashfree:
		return CashfreeString
	case VendorSurepass:
		return SurepassString
	default:
		return "Unknown"
	}
}

func VendorFromString(s string) (Vendor, error) {
	switch s {
	case DecentroString:
		return VendorDecentro, nil
	case CashfreeString:
		return VendorCashfree, nil
	case SurepassString:
		return VendorSurepass, nil
	default:
		return -1, fmt.Errorf("unknown OKYC Vendor: %s", s)
	}
}

type GenerateOTPRequest struct {
	AadharNumber string
}
type ValidateOTPRequest struct {
	ReferenceId   string
	TransactionId string
	OTP           string
}

type HealthCheckResponse struct {
	Percentage int
	Available  bool
}

type GenerateOTPResponse struct {
	ReferenceId   string
	TransactionId string
	MobileNumber  string // (Optional) this is the mobile number linked to the Aadhaar
}

type Gender string

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

// will be sending back our customized response and not the actual response comming from vendors
// we will be logging the actual response in APIUsageLogger (currently in dev)
type ValidateOTPResponse struct {
	PDFBytesStr   string
	ImageBytesStr string
	DOB           string // dd-mm-yyyy
	Gender        Gender
	FullName      string
	Country       string
	Pincode       string
	State         string
	Address       string
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
	ErrOTPAlreadySent            = fmt.Errorf("otp already sent") // in this case we will send success to user
	ErrInvalidAadhar             = fmt.Errorf("invalid Aadhaar provided")
	ErrAadharSuspended           = fmt.Errorf("UIDAI has suspended the Aadhaar")
	ErrNoMobileNumberAadhar      = fmt.Errorf("mobile is not lined to Aadhaar")
	ErrGenerateOTPFailure        = fmt.Errorf("error while generating OTP. Please try again")
	ErrUIDLockedAadhar           = fmt.Errorf("user has locked his/her Aadhaar. Cannot generate OTP")
	ErrInvalidResponseFromVendor = fmt.Errorf("invalid response from vendor")
	ErrSomethingWentWrong        = fmt.Errorf("something went wrong, please try again later")
	ErrInvalidOTP                = fmt.Errorf("invalid OTP provided")
	ErrInvalidReferenceId        = fmt.Errorf("invalid reference id provided")
	ErrInvalidCredentials        = fmt.Errorf("invalid credentials provided, please check your secrets and try again")
	ErrCreditsExhausted          = fmt.Errorf("your credits are exhausted, please contact support to add more credits")
)
