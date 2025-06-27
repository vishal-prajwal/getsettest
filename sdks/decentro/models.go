package decentro

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/encryption"
	"bitbucket.org/junglee_games/getsetgo/sdks/okyc"
	"github.com/google/uuid"
)

type HealthStatus string

const (
	HealthStatusFluctuations      HealthStatus = "FLUCTUATIONS"
	HealthStatusMajorFluctuations HealthStatus = "MAJORFLUCTUATIONS"
	HealthStatusDowntime          HealthStatus = "DOWNTIME"
	HealthStatusInconclusive      HealthStatus = "INCONCLUSIVE"
	HealthStatusOperational       HealthStatus = "OPERATIONAL"
)

type DecentroHealthCheckResponse []struct {
	ApiName       string       `json:"ApiName"`
	Status        HealthStatus `json:"Status"`
	FromTimestamp string       `json:"from_timestamp"`
	ToTimestamp   string       `json:"to_timestamp"`
}

var (
	healthStatusToAvailabilityMap = map[HealthStatus]okyc.HealthCheckResponse{
		HealthStatusOperational: {
			Percentage: 100,
			Available:  true,
		},
		HealthStatusFluctuations: {
			Percentage: 70,
			Available:  true,
		},
		HealthStatusMajorFluctuations: {
			Percentage: 50,
			Available:  false,
		},
		HealthStatusInconclusive: {
			Percentage: 30,
			Available:  false,
		},
		HealthStatusDowntime: {
			Percentage: 0,
			Available:  false,
		},
	}
)

type SendOTPRequest struct {
	AadhaarNumber string `json:"aadhaar_number"`
	ReferenceId   string `json:"reference_id"`
	Purpose       string `json:"purpose"` // e.g., "For Aadhaar Verification"
	Consent       bool   `json:"consent"` // true if user has given consent
}

func NewSendOTPRequest(req *okyc.GenerateOTPRequest) *SendOTPRequest {
	return &SendOTPRequest{
		AadhaarNumber: req.AadharNumber,
		ReferenceId:   uuid.New().String(),        // random uuid, we need always a unique reference id
		Purpose:       "For Aadhaar Verification", // this wil not change
		Consent:       true,                       // concent will always be true else decentro will give error
	}
}

type DecentroResponse[T any] struct {
	DecentroTxnId string      `json:"decentroTxnId"`
	Status        string      `json:"status"`
	Data          T           `json:"data"` // Generic type for data
	ResponseCode  string      `json:"responseCode"`
	Message       string      `json:"message"`
	ResponseKey   ResponseKey `json:"responseKey"`
}

type ValidateOTPRequest struct {
	ReferenceId   string `json:"reference_id"`              // Unique reference ID for the request
	TransactionId string `json:"initiation_transaction_id"` // Unique transaction ID for the request
	OTP           string `json:"otp"`                       // OTP received by the user
	Consent       bool   `json:"consent"`                   // true if user has given consent
	GeneratePDF   bool   `json:"generate_pdf"`              // true if PDF needs to be generated
	GenerateXML   bool   `json:"generate_xml"`              // true if XML needs to be generated
	Purpose       string `json:"purpose"`                   // e.g., "For Aadhaar Verification"
}

func NewValidateOTPRequest(req *okyc.ValidateOTPRequest) *ValidateOTPRequest {
	return &ValidateOTPRequest{
		ReferenceId:   req.ReferenceId,
		TransactionId: req.TransactionId,
		OTP:           req.OTP,
		Consent:       true,                       // concent will always be true else decentro will give error
		GeneratePDF:   true,                       // this will always be true as per discussion with decentro team
		GenerateXML:   false,                      // no need to generate XML as per current requirements
		Purpose:       "For Aadhaar Verification", // this wil not change
	}
}

type ValidateOTPResponseData struct {
	AadhaarReferenceNumber string `json:"aadhaarReferenceNumber"`
	ProofOfIdentity        struct {
		DOB          string `json:"dob"`          // Date of Birth in dd-mm-yyyy format
		HashedEmail  string `json:"hashedEmail"`  // Hashed email address
		Gender       Gender `json:"gender"`       // Gender (M/F)
		Name         string `json:"name"`         // Name of the individual
		MobileNumber string `json:"mobileNumber"` // Mobile number
	} `json:"proofOfIdentity"`
	ProofOfAddress struct {
		CareOf      string `json:"careOf"`      // Care of
		Country     string `json:"country"`     // Country
		District    string `json:"district"`    // District
		House       string `json:"house"`       // House number
		Landmark    string `json:"landmark"`    // Landmark
		Locality    string `json:"locality"`    // Locality
		Pincode     string `json:"pincode"`     // Pincode
		PostOffice  string `json:"postOffice"`  // Post office
		State       string `json:"state"`       // State
		Street      string `json:"street"`      // Street
		SubDistrict string `json:"subDistrict"` // Sub-district
		VTC         string `json:"vtc"`         // Village/Town/City
	} `json:"proofOfAddress"`
	Image    string `json:"image"`    // Base64 encoded image string
	PDF      string `json:"pdf"`      // Base64 encoded PDF string
	Password string `json:"password"` // Password
}

func (r ValidateOTPResponseData) GetAddress() string {
	address := []string{}
	addressProof := r.ProofOfAddress
	if addressProof.House != "" {
		address = append(address, addressProof.House)
	}
	if addressProof.Street != "" {
		address = append(address, addressProof.Street)
	}
	if addressProof.Locality != "" {
		address = append(address, addressProof.Locality)
	}
	if addressProof.VTC != "" {
		address = append(address, addressProof.VTC)
	}
	if addressProof.SubDistrict != "" {
		address = append(address, addressProof.SubDistrict)
	}
	if addressProof.District != "" {
		address = append(address, addressProof.District)
	}
	if addressProof.State != "" {
		address = append(address, addressProof.State)
	}
	if addressProof.Country != "" {
		address = append(address, addressProof.Country)
	}
	if addressProof.Pincode != "" {
		address = append(address, addressProof.Pincode)
	}
	return strings.Join(address, ", ")
}

func (r ValidateOTPResponseData) GetDocumentPDF(ctx context.Context) (string, error) {
	if r.PDF == "" {
		return "", fmt.Errorf("PDF is empty in response data")
	}
	// convert the PDF base64 to bytes

	base64Data, err := base64.StdEncoding.DecodeString(r.PDF)
	if err != nil {
		return "", fmt.Errorf("failed to decode PDF base64 data: %w", err)
	}
	// unprotect the PDF
	bytes, err := encryption.UnprotectPdfBytes(base64Data, r.Password)
	if err != nil {
		return "", fmt.Errorf("failed to unprotect PDF bytes: %w", err)
	}
	// encode the bytes to base64 string
	return base64.StdEncoding.EncodeToString(bytes), nil
}

type Gender string

const (
	GenderMale   = "M"
	GenderFemale = "F"
)

func (g Gender) ToOKYCGender() okyc.Gender {
	switch g {
	case GenderMale:
		return okyc.GenderMale
	case GenderFemale:
		return okyc.GenderFemale
	default:
		return okyc.GenderOther
	}
}
