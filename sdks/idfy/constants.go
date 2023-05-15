package idfy

const (
	AADHAR_DOC_TYPE   = "/ind_aadhaar"
	PAN_DOC_TYPE      = "/ind_pan"
	DL_DOC_TYPE       = "/ind_driving_license"
	VOTER_DOC_TYPE    = "/ind_voter_id"
	PASSPORT_DOC_TYPE = "/ind_passport"
	FraudCheckAadhar  = "/v3/tasks/async/verify_with_source/aadhaar_lite"
	GetTaskStatus     = "/v3/tasks"
)
