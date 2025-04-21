package userdetails

const (
	USER_DETAIL_FROM_NUMBER_END_POINT       string   = "/v2/users/mobile/%s/%s"
	USER_DETAIL_FROM_USER_ID_END_POINT      string   = "/v2/users/%d/products/%s/mobile"
	RDC                                     Platform = "RUMMYCOM"
	JR                                      Platform = "JUNGLEERUMMY"
	NR_USER_DETAILS_FETCH_DATA_FROM_NUMBER  string   = "USER_DETAILS_FETCH_DATA_FROM_NUMBER"
	NR_USER_DETAILS_FETCH_DATA_FROM_USER_ID string   = "USER_DETAILS_FETCH_DATA_FROM_USER_ID"
	ERROR_USER_DOES_NOT_EXIST               string   = "user does not exist"
	STATUS_FAILURE                          string   = "FAILURE"
	STATUS_SUCCESS                          string   = "SUCCESS"
)

type Platform string
