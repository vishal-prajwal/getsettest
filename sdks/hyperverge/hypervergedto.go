package hyperverge

type HypervergeRequest struct {
	ImageFile                string `json:"imageFile"`
	EnableDashboard     string `json:"enableDashboard"`
	MaskAadhaarComplete string `json:"maskAadhaarComplete"`
	OutputImageUrl      string `json:"outputImageUrl"`
	ImageName string `json:"imageName"`
}

type HypervergePanResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     []struct {
		Type    string `json:"type"`
		Details struct {
			Date struct {
				Value string `json:"value"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"date"`
			Father struct {
				Value string `json:"value"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"father"`
			Name struct {
				Value string `json:"value"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"name"`
			PanNo struct {
				Value string `json:"value"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"pan_no"`
			DateOfIssue struct {
				Value string `json:"value"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"date_of_issue"`
		} `json:"details"`
	} `json:"result"`
}

type HypervergeAadharResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     []struct {
		Type    string `json:"type"`
		Details struct {
			Aadhaar struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"aadhaar,omitempty"`
			Dob struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"dob,omitempty"`
			Father struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"father,omitempty"`
			Gender struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"gender,omitempty"`
			Mother struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"mother,omitempty"`
			Name struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"name,omitempty"`
			Yob struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"yob,omitempty"`
			Address struct {
				CareOf      string `json:"care_of,omitempty"`
				District    string `json:"district,omitempty"`
				City        string `json:"city,omitempty"`
				Locality    string `json:"locality,omitempty"`
				Landmark    string `json:"landmark,omitempty"`
				Street      string `json:"street,omitempty"`
				Line1       string `json:"line1,omitempty"`
				Line2       string `json:"line2,omitempty"`
				HouseNumber string `json:"house_number,omitempty"`
				Pin         string `json:"pin,omitempty"`
				State       string `json:"state,omitempty"`
				Value       string `json:"value,omitempty"`
				Conf        int    `json:"conf,omitempty"`
			} `json:"address,omitempty"`
			Husband struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"husband,omitempty"`
			Phone struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"phone,omitempty"`
			Pin struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"pin,omitempty"`
		} `json:"details"`
	} `json:"result"`
}

type HypervergePassportResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     []struct {
		Type    string `json:"type"`
		Details struct {
			CountryCode struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"country_code,omitempty"`
			Dob struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"dob,omitempty"`
			Doe struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"doe,omitempty"`
			Doi struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"doi,omitempty"`
			Gender struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"gender,omitempty"`
			GivenName struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"given_name,omitempty"`
			Nationality struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"nationality,omitempty"`
			PassportNum struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"passport_num,omitempty"`
			PlaceOfBirth struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"place_of_birth,omitempty"`
			PlaceOfIssue struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"place_of_issue,omitempty"`
			Surname struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"surname,omitempty"`
			Mrz struct {
				Line1 string `json:"line1,omitempty"`
				Line2 string `json:"line2,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"mrz,omitempty"`
			Type struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"type,omitempty"`
			Address struct {
				District    string `json:"district,omitempty"`
				City        string `json:"city,omitempty"`
				Locality    string `json:"locality,omitempty"`
				Landmark    string `json:"landmark,omitempty"`
				Street      string `json:"street,omitempty"`
				Line1       string `json:"line1,omitempty"`
				Line2       string `json:"line2,omitempty"`
				HouseNumber string `json:"house_number,omitempty"`
				Pin         string `json:"pin,omitempty"`
				State       string `json:"state,omitempty"`
				Value       string `json:"value,omitempty"`
				Conf        int    `json:"conf,omitempty"`
			} `json:"address,omitempty"`
			Father struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"father,omitempty"`
			Mother struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"mother,omitempty"`
			FileNum struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"file_num,omitempty"`
			OldDoi struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"old_doi,omitempty"`
			OldPassportNum struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"old_passport_num,omitempty"`
			OldPlaceOfIssue struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"old_place_of_issue,omitempty"`
			Pin struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"pin,omitempty"`
			Spouse struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"spouse,omitempty"`
		} `json:"details"`
	} `json:"result"`
}

type HypervergeVoterIdResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     []struct {
		Type    string `json:"type"`
		Details struct {
			Voterid struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"voterid,omitempty"`
			Name struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"name,omitempty"`
			Gender struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"gender,omitempty"`
			Relation struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"relation,omitempty"`
			Dob struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"dob,omitempty"`
			Doc struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"doc,omitempty"`
			Age struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"age,omitempty"`

			Pin struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"pin,omitempty"`

			Date struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"date,omitempty"`
			Type struct {
				Value string `json:"value,omitempty"`
				Conf  int    `json:"conf,omitempty"`
			} `json:"type,omitempty"`
			Address struct {
				District    string `json:"district,omitempty"`
				City        string `json:"city,omitempty"`
				Locality    string `json:"locality,omitempty"`
				Landmark    string `json:"landmark,omitempty"`
				Street      string `json:"street,omitempty"`
				Line1       string `json:"line1,omitempty"`
				Line2       string `json:"line2,omitempty"`
				HouseNumber string `json:"house_number,omitempty"`
				Pin         string `json:"pin,omitempty"`
				State       string `json:"state,omitempty"`
				Value       string `json:"value,omitempty"`
				Conf        int    `json:"conf,omitempty"`
			} `json:"address,omitempty"`
		} `json:"details"`
	} `json:"result"`
}

type PanResponse struct {
	Date        string `json:"date,omitempty"`
	Father      string `json:"father,omitempty"`
	Name        string `json:"name,omitempty"`
	PanNo       string `json:"pan_no,omitempty"`
	DateOfIssue string `json:"date_of_issue,omitempty"`
}

type AadharResponse struct {
	Aadhaar     string `json:"aadhaar,omitempty"`
	Dob         string `json:"dob,omitempty"`
	Father      string `json:"father,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Mother      string `json:"mother,omitempty"`
	Name        string `json:"name,omitempty"`
	Yob         string `json:"yob,omitempty"`
	CareOf      string `json:"care_of,omitempty"`
	District    string `json:"district,omitempty"`
	City        string `json:"city,omitempty"`
	Locality    string `json:"locality,omitempty"`
	Landmark    string `json:"landmark,omitempty"`
	Street      string `json:"street,omitempty"`
	Line1       string `json:"line1,omitempty"`
	Line2       string `json:"line2,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
	State       string `json:"state,omitempty"`
	Value       string `json:"value,omitempty"`
	Husband     string `json:"husband,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Pin         string `json:"pin,omitempty"`
	AddressPin  string `json:"address_pin,omitempty"`
}

type PassportResponse struct {
	CountryCode     string `json:"country_code,omitempty"`
	Dob             string `json:"dob,omitempty"`
	Doe             string `json:"doe,omitempty"`
	Doi             string `json:"doi,omitempty"`
	Gender          string `json:"gender,omitempty"`
	GivenName       string `json:"given_name,omitempty"`
	Nationality     string `json:"nationality,omitempty"`
	PassportNum     string `json:"passport_num,omitempty"`
	PlaceOfBirth    string `json:"place_of_birth,omitempty"`
	PlaceOfIssue    string `json:"place_of_issue,omitempty"`
	Surname         string `json:"surname,omitempty"`
	Mrz             string `json:"mrz,omitempty"`
	Type            string `json:"type,omitempty"`
	District        string `json:"district,omitempty"`
	City            string `json:"city,omitempty"`
	Locality        string `json:"locality,omitempty"`
	Landmark        string `json:"landmark,omitempty"`
	Street          string `json:"street,omitempty"`
	Line1           string `json:"line1,omitempty"`
	Line2           string `json:"line2,omitempty"`
	HouseNumber     string `json:"house_number,omitempty"`
	State           string `json:"state,omitempty"`
	Value           string `json:"value,omitempty"`
	Conf            string `json:"conf,omitempty"`
	Father          string `json:"father,omitempty"`
	Mother          string `json:"mother,omitempty"`
	FileNum         string `json:"file_num,omitempty"`
	OldDoi          string `json:"old_doi,omitempty"`
	OldPassportNum  string `json:"old_passport_num,omitempty"`
	OldPlaceOfIssue string `json:"old_place_of_issue,omitempty"`
	Pin             string `json:"pin,omitempty"`
	Spouse          string `json:"spouse,omitempty"`
	AddressPin      string `json:"address_pin,omitempty"`
}

type VoterIdResponse struct {
	Voterid     string `json:"voterid,omitempty"`
	Name        string `json:"name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Dob         string `json:"dob,omitempty"`
	Doc         string `json:"doc,omitempty"`
	Age         string `json:"age,omitempty"`
	Pin         string `json:"pin,omitempty"`
	Date        string `json:"date,omitempty"`
	Type        string `json:"type,omitempty"`
	District    string `json:"district,omitempty"`
	City        string `json:"city,omitempty"`
	Locality    string `json:"locality,omitempty"`
	Landmark    string `json:"landmark,omitempty"`
	Street      string `json:"street,omitempty"`
	Line1       string `json:"line1,omitempty"`
	Line2       string `json:"line2,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
	State       string `json:"state,omitempty"`
	Value       string `json:"value,omitempty"`
	Conf        string `json:"conf,omitempty"`
	AddressPin  string `json:"address_pin,omitempty"`
}

type FraudCheckPanRequest struct {
	Pan                       string `json:"pan"`
	Name                      string `json:"name"`
	Dob                       string `json:"dob"`
	StrictlyUseGetNameFromPan string `json:"strictlyUseGetNameFromPan"`
	MatchDob                  string `json:"matchDob"`
}

type FraudCheckPanResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     struct {
		DobMatch  bool        `json:"dobMatch"`
		NameMatch bool        `json:"nameMatch"`
		Status    string      `json:"status"`
		Duplicate interface{} `json:"duplicate"`
	} `json:"result"`
}

type FraudCheckDlRequest struct {
	DlNumber string `json:"dlNumber"`
	Dob      string `json:"dob"`
}

type FraudCheckDlResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     struct {
		IssueDate  string `json:"issue_date"`
		FatherName string `json:"father/husband"`
		Name       string `json:"name"`
		Image      string `json:"img"`
		BloodGroup string `json:"blood_group"`
		DOB        string `json:"dob"`
		Validity   struct {
			NonTransport string `json:"non-transport"`
			Transport    string `json:"transport"`
		} `json:"validity"`
		COVDetails []struct {
			IssueDate string `json:"issue_date"`
			COV       string `json:"cov"`
		} `json:"cov_details"`
		Address string `json:"address"`
	} `json:"result"`
}

type FraudCheckPassportRequest struct {
	FileNo     string `json:"fileNo"`
	Dob        string `json:"dob"`
	Doi        string `json:"doi"`
	PassportNo string `json:"passportNo"`
	Name       string `json:"name"`
}

type FraudCheckPassportResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     struct {
		ApplicationDate string `json:"applicationDate"`
		DateOfIssue     struct {
			DispatchedOnFromSource string `json:"dispatchedOnFromSource"`
			DateOfIssueMatch       bool   `json:"dateOfIssueMatch"`
		} `json:"dateOfIssue"`
		PassportNumber struct {
			PassportNumberFromSource string `json:"passportNumberFromSource"`
			PassportNumberMatch      bool   `json:"passportNumberMatch"`
		} `json:"passportNumber"`
		Name struct {
			NameMatch           bool   `json:"nameMatch"`
			SurnameFromPassport string `json:"surnameFromPassport"`
			NameScore           int    `json:"nameScore"`
			NameFromPassport    string `json:"nameFromPassport"`
		} `json:"name"`
		TypeOfApplication string `json:"typeOfApplication"`
	} `json:"result"`
}

type FraudCheckVoterRequest struct {
	EpicNumber string `json:"epicNumber"`
}

type FraudCheckVoterResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
	Result     struct {
		PSLatLong  string `json:"ps_lat_long"`
		RLNNameV1  string `json:"rln_name_v1"`
		RLNNameV2  string `json:"rln_name_v2"`
		RLNNameV3  string `json:"rln_name_v3"`
		PartNo     string `json:"part_no"`
		RLNType    string `json:"rln_type"`
		SectionNo  string `json:"section_no"`
		ID         string `json:"id"`
		EpicNo     string `json:"epic_no"`
		RLNName    string `json:"rln_name"`
		District   string `json:"district"`
		LastUpdate string `json:"last_update"`
		State      string `json:"state"`
		ACNo       string `json:"ac_no"`
		HouseNo    string `json:"house_no"`
		PSName     string `json:"ps_name"`
		PCName     string `json:"pc_name"`
		SlnoInPart string `json:"slno_inpart"`
		Name       string `json:"name"`
		PartName   string `json:"part_name"`
		DOB        string `json:"dob"`
		Gender     string `json:"gender"`
		Age        int    `json:"age"`
		ACName     string `json:"ac_name"`
		NameV1     string `json:"name_v1"`
		StCode     string `json:"st_code"`
		NameV3     string `json:"name_v3"`
		NameV2     string `json:"name_v2"`
	} `json:"result"`
}

type FraudCheckAadharResponse struct {
	Message            string            `json:"message"`
	Error              string            `json:"error"`
	StatusCode         int               `json:"statusCode"`
	Code               int               `json:"code"`
	DetailsFromUIDAI   map[string]string `json:"detailsFromUIDAI"`
	AadhaarNumberCheck fraudMatch        `json:"aadhaarNumberCheck"`
	DobMatch           fraudMatch        `json:"dobMatch"`
	StateMatch         fraudMatch        `json:"stateMatch"`
}

type fraudMatch struct {
	Match   bool   `json:"match"`
	Channel string `json:"channel"`
}

type FraudCheckAadharRequest struct{
	SessionID     string `json:"sessionId"`
	AadhaarNoUser string `json:"aadhaarNoUser"`
	SecurityCode  string `json:"securityCode"`
	Match         bool   `json:"match"`
	DOB           string `json:"dob"`
}
