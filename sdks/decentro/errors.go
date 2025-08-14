package decentro

import "bitbucket.org/junglee_games/getsetgo/sdks/okyc"

type ResponseKey string

const (
	ErrorKeyDuplicateReferenceID       ResponseKey = "error_duplicate_reference_id"
	ErrorKeyNoMobileFound              ResponseKey = "error_no_mobile_found"
	ErrorKeyCancelledAadhaar           ResponseKey = "error_cancelled_aadhaar"
	ErrorKeyInvalidAadhaar             ResponseKey = "error_invalid_aadhaar"        // Aadhaar number is invalid or not registered
	ErrorKeyInvalidAadhaarNumber       ResponseKey = "error_invalid_aadhaar_number" // Aadhaar number is not valid or does not match the expected format
	ErrorKeyOTPLimitExceeded           ResponseKey = "error_otp_limit_exceeded"
	ErrorKeyResendOTP                  ResponseKey = "error_resend_otp" // wait for 30 seconds before retrying
	ErrorKeyInvalidOTP                 ResponseKey = "error_invalid_otp"
	ErrorKeyInvalidOTPRequest          ResponseKey = "error_invalid_otp_request"
	ErrorKeyInvalidSession             ResponseKey = "error_invalid_session"
	ErrorKeyProviderError              ResponseKey = "error_provider_error"       // User Aadhaar is locked
	ErrorKeyOTPFlooding                ResponseKey = "error_aadhaar_otp_flooding" // too many OTP requests in a short time wait for 120 seconds
	ErrorKeyMissingConfiguration       ResponseKey = "error_missing_configuration"
	ErrorKeyEmptyModuleSecret          ResponseKey = "error_empty_module_secret"
	ErrorKeyNoSubscriptionFound        ResponseKey = "error_no_subscription_found"
	ErrorKeyUnauthorizedModule         ResponseKey = "error_unauthorized_module"
	ErrorKeyModuleAccessExpired        ResponseKey = "error_module_access_expired"
	ErrorKeyModuleCreditsExhausted     ResponseKey = "error_module_credits_exhausted"
	ErrorKeyInsufficientAccountBalance ResponseKey = "error_insufficient_account_balance"
)

const (
	errMarshalRequestBody    = "failed to marshal request body: "
	errSendRequest           = "failed to send request: "
	errReadResponseBody      = "failed to read response body: "
	errUnmarshalResponseBody = "failed to unmarshal response body: "
)

var (
	responseKeyToErrorMap = map[ResponseKey]error{
		ErrorKeyDuplicateReferenceID:       okyc.ErrInvalidReferenceId,
		ErrorKeyNoMobileFound:              okyc.ErrNoMobileNumberAadhar,
		ErrorKeyCancelledAadhaar:           okyc.ErrAadharSuspended,
		ErrorKeyInvalidAadhaar:             okyc.ErrInvalidAadhar,
		ErrorKeyInvalidAadhaarNumber:       okyc.ErrInvalidAadhar,
		ErrorKeyOTPLimitExceeded:           okyc.ErrOTPAlreadySent,     // in this case we will send success to user
		ErrorKeyResendOTP:                  okyc.ErrGenerateOTPFailure, // wait for 30 seconds before retrying
		ErrorKeyInvalidOTP:                 okyc.ErrInvalidOTP,
		ErrorKeyInvalidSession:             okyc.ErrInvalidReferenceId, // invalid reference id provided
		ErrorKeyProviderError:              okyc.ErrUIDLockedAadhar,    // User Aadhaar is locked
		ErrorKeyOTPFlooding:                okyc.ErrGenerateOTPFailure, // too many OTP requests in a short time wait for 120 seconds
		ErrorKeyMissingConfiguration:       okyc.ErrInvalidCredentials,
		ErrorKeyEmptyModuleSecret:          okyc.ErrInvalidCredentials,
		ErrorKeyNoSubscriptionFound:        okyc.ErrInvalidCredentials,
		ErrorKeyUnauthorizedModule:         okyc.ErrInvalidCredentials,
		ErrorKeyModuleAccessExpired:        okyc.ErrInvalidCredentials,
		ErrorKeyModuleCreditsExhausted:     okyc.ErrCreditsExhausted,
		ErrorKeyInsufficientAccountBalance: okyc.ErrCreditsExhausted,
		ErrorKeyInvalidOTPRequest:          okyc.ErrInvalidOTP,
	}
)

func responseKeyToError(statusKey ResponseKey) error {
	if err, exists := responseKeyToErrorMap[statusKey]; exists {
		return err
	}
	return okyc.ErrInvalidResponseFromVendor
}
