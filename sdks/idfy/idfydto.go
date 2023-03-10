package idfy

type IdfyRequest struct {
	TaskID  string `json:"task_id"`
	GroupID string `json:"group_id"`
	Data    Data   `json:"data"`
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
	Minor         int    `json:"minor"`
	Is_scanned    int    `json:"is_scanned"`
	Pan_type      int    `json:"pan_type"`
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
	Is_scanned     int    `json:"is_scanned"`
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
	Age            int    `json:"age"`
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
