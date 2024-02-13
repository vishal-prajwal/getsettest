package digilocker

import "errors"

var (
	ErrNotAvailable        = errors.New("ERROR : hyperverge not available")
	ErrUnmarshalJson       = errors.New("ERROR : unmarshling response")
	ErrReadResponseBody    = errors.New("ERROR : reading response body")
	ErrCallingHyperverge   = errors.New("ERROR : calling hyperverge")
	ErrCreatingRequest     = errors.New("ERROR : creating request")
	ErrReqValidate         = errors.New("ERROR : hyperverge request validate")
	ErrHVServer            = errors.New("ERROR : Hyperverge Server")
	ErrConsentNotProvided  = errors.New("ERROR : User did not provide the consent")
	ErrHVServerMissingData = errors.New("ERROR : Hyperverge Sent incomplete data")
	ErrDocumentNotFound    = errors.New("ERROR : document not found")
	ErrDownloadFile        = errors.New("ERROR : Download addhar xml file")
	ErrExtractingXML       = errors.New("ERROR : Unable to extract data from xml file addhar")
)
