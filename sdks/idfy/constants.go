package idfy

const (
	AADHAR_DOC_TYPE    = "/v3/tasks/sync/extract/ind_aadhaar"
	PAN_DOC_TYPE       = "/v3/tasks/sync/extract/ind_pan"
	DL_DOC_TYPE        = "/v3/tasks/sync/extract/ind_driving_license"
	VOTER_DOC_TYPE     = "/v3/tasks/sync/extract/ind_voter_id"
	PASSPORT_DOC_TYPE  = "/v3/tasks/sync/extract/ind_passport"
	FraudCheckAadhar   = "/v3/tasks/async/verify_with_source/aadhaar_lite"
	GetTaskStatus      = "/v3/tasks"
	HealthCheckTimeout = 5
	TemperedImage = "/v3/tasks/sync/check_tampering/document"
)
