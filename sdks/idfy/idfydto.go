package idfy

import "time"

type IdfyRequest struct {
	TaskID  string `json:"task_id"`
	GroupID string `json:"group_id"`
	Data    Data   `json:"data"`
}

type CheckTemperedReq struct {
	TaskID  string               `json:"task_id"`
	GroupID string               `json:"group_id"`
	Data    CheckTemperedReqData `json:"data"`
}

type CheckTemperedReqData struct {
	Document1 string `json:"document1"`
	DocType   string `json:"doc_type"`
}

type AdvancedDetails struct {
	ExtractQrInfo     bool `json:"extract_qr_info"`
	ExtractLast4Digit bool `json:"extract_last_4_digit"`
}
type Data struct {
	Document1       string          `json:"document1"`
	Document2       string          `json:"document2"`
	Consent         string          `json:"consent"`
	AdvancedDetails AdvancedDetails `json:"advanced_details"`
}

type IdfyPanResponse struct {
	Id_number     string `json:"id_number"`
	Name_on_card  string `json:"name_on_card"`
	Fathers_name  string `json:"fathers_name"`
	Date_of_birth string `json:"date_of_birth"`
	Date_of_issue string `json:"date_of_issues"`
	Age           int    `json:"age"`
	Minor         bool   `json:"minor"`
	Is_scanned    bool   `json:"is_scanned"`
	Pan_type      string `json:"pan_type"`
}

type IdfyAadharResponse struct {
	Id_number      string `json:"id_number"`
	Name_on_card   string `json:"name_on_card"`
	Fathers_name   string `json:"fathers_name"`
	Date_of_birth  string `json:"date_of_birth"`
	Year_of_birth  string `json:"year_of_birth"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	Street_address string `json:"street_address"`
	House_number   string `json:"house_number"`
	District       string `json:"district"`
	Pincode        string `json:"pincode"`
	State          string `json:"state"`
	Is_scanned     bool   `json:"is_scanned"`
}

type IdfyDlResponse struct {
	Id_number        string   `json:"id_number"`
	Name_on_card     string   `json:"name_on_card"`
	Fathers_name     string   `json:"fathers_name"`
	Date_of_birth    string   `json:"date_of_birth"`
	Date_of_validity string   `json:"date_of_validity"`
	Address          string   `json:"address"`
	District         string   `json:"district"`
	Street_address   string   `json:"street_address"`
	Pincode          string   `json:"pincode"`
	State            string   `json:"state"`
	Issue_dates      string   `json:"issue_dates"`
	Type             []string `json:"type"`
	Validity         string   `json:"validity"`
}

type IdfyVoterIdResponse struct {
	Id_number      string `json:"id_number"`
	Name_on_card   string `json:"name_on_card"`
	Fathers_name   string `json:"fathers_name"`
	Date_of_birth  string `json:"date_of_birth"`
	Year_of_birth  string `json:"year_of_birth"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	Street_address string `json:"street_address"`
	House_number   string `json:"house_number"`
	District       string `json:"district"`
	Pincode        string `json:"pincode"`
	State          string `json:"state"`
	Age            string `json:"age"`
}

type IdfyPassportResponse struct {
	Id_number       string `json:"id_number"`
	Is_scanned      bool   `json:"is_scanned"`
	First_name      string `json:"first_name"`
	Last_name       string `json:"last_name"`
	Name_on_card    string `json:"name_on_card"`
	Nationality     string `json:"nationality"`
	Pincode         string `json:"pincode"`
	Fathers_name    string `json:"fathers_name"`
	Mothers_name    string `json:"mothers_name"`
	Name_of_spouse  string `json:"name_of_spouse"`
	Date_of_birth   string `json:"date_of_birth"`
	Place_of_birth  string `json:"place_of_birth"`
	Date_of_issue   string `json:"date_of_issue"`
	District        string `json:"district"`
	State           string `json:"state"`
	Date_of_expirty string `json:"date_of_expirty"`
	Place_of_issue  string `json:"place_of_issue"`
	Address         string `json:"address"`
	Gender          string `json:"gender"`
	File_number     string `json:"file_number"`
}

type PanResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      ResultPan `json:"result"`
	Status      string    `json:"status"`
	TaskID      string    `json:"task_id"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Error       string    `json:"error"`
}

type ResultPan struct {
	ExtractionOutput IdfyPanResponse `json:"extraction_output"`
}

type AadharResponse struct {
	Action      string       `json:"action"`
	CompletedAt time.Time    `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at"`
	GroupID     string       `json:"group_id"`
	RequestID   string       `json:"request_id"`
	Result      ResultAadhar `json:"result"`
	Status      string       `json:"status"`
	TaskID      string       `json:"task_id"`
	Type        string       `json:"type"`
	Message     string       `json:"message"`
	Error       string       `json:"error"`
}

type ResultAadhar struct {
	ExtractionOutput IdfyAadharResponse `json:"extraction_output"`
}

type PassportResponse struct {
	Action      string         `json:"action"`
	CompletedAt time.Time      `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	GroupID     string         `json:"group_id"`
	RequestID   string         `json:"request_id"`
	Result      ResultPassport `json:"result"`
	Status      string         `json:"status"`
	TaskID      string         `json:"task_id"`
	Type        string         `json:"type"`
	Message     string         `json:"message"`
	Error       string         `json:"error"`
}

type ResultPassport struct {
	ExtractionOutput IdfyPassportResponse `json:"extraction_output"`
}

type VoterResponse struct {
	Action      string      `json:"action"`
	CompletedAt time.Time   `json:"completed_at"`
	CreatedAt   time.Time   `json:"created_at"`
	GroupID     string      `json:"group_id"`
	RequestID   string      `json:"request_id"`
	Result      ResultVoter `json:"result"`
	Status      string      `json:"status"`
	TaskID      string      `json:"task_id"`
	Type        string      `json:"type"`
	Message     string      `json:"message"`
	Error       string      `json:"error"`
}

type ResultVoter struct {
	ExtractionOutput IdfyVoterIdResponse `json:"extraction_output"`
}

type DlResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      ResultDl  `json:"result"`
	Status      string    `json:"status"`
	TaskID      string    `json:"task_id"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Error       string    `json:"error"`
}

type ResultDl struct {
	ExtractionOutput IdfyDlResponse `json:"extraction_output"`
}

type FraudCheckRequest struct {
	TaskID  string         `json:"task_id"`
	GroupID string         `json:"group_id"`
	Data    FraudCheckData `json:"data"`
}
type FraudCheckData struct {
	IdNumber           string `json:"id_number,omitempty"`
	PassportFileNumber string `json:"passport_file_number,omitempty"`
	DateOfBirth        string `json:"date_of_birth,omitempty"`
	AadhaarNumber      string `json:"aadhaar_number,omitempty"`
}

type FraudCheckResponse struct {
	RequestID string `json:"request_id"`
}

type FraudCheckDlResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		SourceOutput struct {
			Address      interface{} `json:"address"`
			BadgeDetails interface{} `json:"badge_details"`
			CardSerialNo interface{} `json:"card_serial_no"`
			City         interface{} `json:"city"`
			CovDetails   []struct {
				Category  string `json:"category"`
				Cov       string `json:"cov"`
				IssueDate string `json:"issue_date"`
			} `json:"cov_details"`
			DateOfIssue           string      `json:"date_of_issue"`
			DateOfLastTransaction interface{} `json:"date_of_last_transaction"`
			DlStatus              string      `json:"dl_status"`
			Dob                   string      `json:"dob"`
			FaceImage             interface{} `json:"face_image"`
			Gender                interface{} `json:"gender"`
			HazardousValidTill    interface{} `json:"hazardous_valid_till"`
			HillValidTill         interface{} `json:"hill_valid_till"`
			IDNumber              string      `json:"id_number"`
			IssuingRtoName        string      `json:"issuing_rto_name"`
			LastTransactedAt      interface{} `json:"last_transacted_at"`
			Name                  string      `json:"name"`
			NtValidityFrom        string      `json:"nt_validity_from"`
			NtValidityTo          string      `json:"nt_validity_to"`
			RelativesName         interface{} `json:"relatives_name"`
			Source                string      `json:"source"`
			Status                string      `json:"status"`
			TValidityFrom         string      `json:"t_validity_from"`
			TValidityTo           string      `json:"t_validity_to"`
			State                 string      `json:"state"`
			IsMinor               bool        `json:"is_minor"`
		} `json:"source_output"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type IDfyNSDLPanRequest struct {
	TaskID  string            `json:"task_id"`
	GroupID string            `json:"group_id"`
	Data    IDfyNSDLPanData `json:"data"`
}

type IDfyNSDLPanData struct {
	IDNumber string `json:"id_number"`
	FullName string `json:"full_name"`
	DOB      string `json:"dob"`
}

type IDfyNSDLPanPostResponse struct {
	RequestID string `json:"request_id"`
}

type IDfyNSDLPanGetResponse struct {
	Action     string    `json:"action"`
	CompletedAt string    `json:"completed_at"`
	CreatedAt   string    `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		SourceOutput struct {
			AadhaarSeedingStatus bool   `json:"aadhaar_seeding_status"`
			PanStatus            string `json:"pan_status"`
			NameMatch            bool   `json:"name_match"`
			DOBMatch             bool   `json:"dob_match"`
			InputDetails         struct {
				InputPanNumber string `json:"input_pan_number"`
				InputName      string `json:"input_name"`
				InputDOB       string `json:"input_dob"`
			} `json:"input_details"`
			Status string `json:"status"`
		} `json:"source_output"`
	} `json:"result"`
	Status string `json:"status"`
	TaskID string `json:"task_id"`
	Type   string `json:"type"`
}

type FraudCheckPanResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		SourceOutput struct {
			AadhaarSeedingStatus bool        `json:"aadhaar_seeding_status"`
			FirstName            string      `json:"first_name"`
			Gender               interface{} `json:"gender"`
			IDNumber             string      `json:"id_number"`
			LastName             string      `json:"last_name"`
			MiddleName           string      `json:"middle_name"`
			NameOnCard           string      `json:"name_on_card"`
			Source               string      `json:"source"`
			Status               string      `json:"status"`
		} `json:"source_output"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type FraudCheckVoterResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		SourceOutput struct {
			AcNo        string      `json:"ac_no"`
			DateOfBirth interface{} `json:"date_of_birth"`
			District    string      `json:"district"`
			Gender      string      `json:"gender"`
			HouseNo     interface{} `json:"house_no"`
			IDNumber    string      `json:"id_number"`
			LastUpdate  string      `json:"last_update"`
			NameOnCard  string      `json:"name_on_card"`
			PartNo      string      `json:"part_no"`
			PsLatLong   string      `json:"ps_lat_long"`
			PsName      string      `json:"ps_name"`
			RlnName     string      `json:"rln_name"`
			SectionNo   string      `json:"section_no"`
			Source      string      `json:"source"`
			StCode      string      `json:"st_code"`
			State       string      `json:"state"`
			Status      string      `json:"status"`
		} `json:"source_output"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type CheckTemperedRes struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		ConfidenceScore interface{} `json:"confidence_score"`
		Details         struct {
			IsAppCreated          bool        `json:"is_app_created"`
			IsDataValidated       interface{} `json:"is_data_validated"`
			IsLayoutIncorrect     interface{} `json:"is_layout_incorrect"`
			IsLive                interface{} `json:"is_live"`
			IsMrzBarcodeValidated interface{} `json:"is_mrz_barcode_validated"`
			IsScanned             interface{} `json:"is_scanned"`
		} `json:"details"`
		IsTampered bool `json:"is_tampered"`
	} `json:"result"`
	Status string `json:"status"`
	TaskID string `json:"task_id"`
	Type   string `json:"type"`
}

type FraudCheckPassportResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		SourceOutput struct {
			ApplicationDate string `json:"application_date"`
			DateOfBirth     string `json:"date_of_birth"`
			FileNumber      string `json:"file_number"`
			Name            string `json:"name"`
			PassportStatus  string `json:"passport_status"`
			Status          string `json:"status"`
			Surname         string `json:"surname"`
		} `json:"source_output"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type FraudCheckAadharResponse struct {
	Action      string    `json:"action"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	GroupID     string    `json:"group_id"`
	RequestID   string    `json:"request_id"`
	Result      struct {
		SourceOutput struct {
			AgeBand struct {
				LowerLimit string `json:"lower_limit"`
				UpperLimit string `json:"upper_limit"`
			} `json:"age_band"`
			Gender       string `json:"gender"`
			MobileNumber string `json:"mobile_number"`
			State        string `json:"state"`
			Status       string `json:"status"`
		} `json:"source_output"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HealthCheckRes struct {
	RequestID string            `json:"request_id"`
	Result    HealthCheckResult `json:"result"`
}

type HealthCheckResult struct {
	LastUpdated string `json:"last_updated"`
	Metadata    string `json:"metadata"`
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
}

type HealthCheckReq struct {
	TaskID  string             `json:"task_id"`
	GroupID string             `json:"group_id"`
	Data    HealthCheckReqData `json:"data"`
}

type HealthCheckReqData struct {
	TaskType string `json:"task_type"`
}

type MaskAadharDocRequest struct {
	GroupID string `json:"group_id"`
	TaskID  string `json:"task_id"`
	Data    struct {
		Consent        string `json:"consent"`
		Document1      string `json:"document1"`
		LastFourDigits bool   `json:"last_four_digits"`
		MaskAllDigits  bool   `json:"mask_all_digits"`
		MaskQr         bool   `json:"mask_qr"`
	} `json:"data"`
}

type MaskAadharRequestID struct {
	RequestID string `json:"request_id"`
}

type MaskAadharDocResponse struct {
	Action      string `json:"action"`
	CompletedAt string `json:"completed_at"`
	CreatedAt   string `json:"created_at"`
	GroupID     string `json:"group_id"`
	RequestID   string `json:"request_id"`
	Result      struct {
		DocumentURL         string `json:"document_url"`
		BaseImage           string `json:"base64_image"`
		IDNumber            string `json:"id_number"`
		IDNumberFound       bool   `json:"id_number_found"`
		OriginalDocumentURL string `json:"original_document_url"`
		SelfLink            string `json:"self_link"`
	} `json:"result"`
	Status  string `json:"status"`
	TaskID  string `json:"task_id"`
	Type    string `json:"type"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
