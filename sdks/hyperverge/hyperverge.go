package hyperverge

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
)

type HypervergeImpl struct {
	config     HypervergeConfig
	nr         newrelic.Agent
	httpClient httpclient.HTTPClient
}

// New creates a new Hyperverge client
func New(config HypervergeConfig, nr newrelic.Agent, client httpclient.HTTPClient) *HypervergeImpl {
	hypervergeImpl := HypervergeImpl{
		config:     config,
		nr:         nr,
		httpClient: client,
	}
	return &hypervergeImpl
}

func (hypervergeImpl *HypervergeImpl) readDocument(documentType string, hypervergeRequest HypervergeRequest) (*bytes.Buffer, error) {
	url := getURLFor(documentType, hypervergeImpl.config.GetHypervergeEndpoint())
	req, err := newfileUploadRequest(url, hypervergeRequest.Path, hypervergeImpl.config.GetHypervergeAppKey(), hypervergeImpl.config.GetHypervergeAppID(), hypervergeImpl.config.GetHypervergeTransactionId())
	if err != nil {
		return nil, err
	}

	resp, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil

}

func (hypervergeImpl *HypervergeImpl) ReadPan(hypervergeRequest HypervergeRequest) (*PanResponse, error) {
	body, err := hypervergeImpl.readDocument("pan", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergePanResponseStruct HypervergePanResponse
	err = json.Unmarshal(body.Bytes(), &hypervergePanResponseStruct)
	if err != nil {
		return nil, err
	}
	hypervergePanResponse := hypervergePanResponseStruct.Result[0].Details
	panResponse := PanResponse{
		Date:        hypervergePanResponse.Date.Value,
		Father:      hypervergePanResponse.Father.Value,
		Name:        hypervergePanResponse.Name.Value,
		PanNo:       hypervergePanResponse.PanNo.Value,
		DateOfIssue: hypervergePanResponse.DateOfIssue.Value,
	}

	return &panResponse, err

}

func (hypervergeImpl *HypervergeImpl) ReadAadhar(hypervergeRequest HypervergeRequest) (*AadharResponse, error) {
	body, err := hypervergeImpl.readDocument("aadhar", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergeAadharResponseStruct HypervergeAadharResponse
	err = json.Unmarshal(body.Bytes(), &hypervergeAadharResponseStruct)
	if err != nil {
		return nil, err
	}
	hypervergeAadharResponse := hypervergeAadharResponseStruct.Result[0].Details
	aadharResponse := AadharResponse{
		Aadhaar:     hypervergeAadharResponse.Aadhaar.Value,
		Dob:         hypervergeAadharResponse.Dob.Value,
		Father:      hypervergeAadharResponse.Father.Value,
		Gender:      hypervergeAadharResponse.Gender.Value,
		Mother:      hypervergeAadharResponse.Mother.Value,
		Name:        hypervergeAadharResponse.Name.Value,
		Yob:         hypervergeAadharResponse.Yob.Value,
		Husband:     hypervergeAadharResponse.Husband.Value,
		Phone:       hypervergeAadharResponse.Phone.Value,
		Pin:         hypervergeAadharResponse.Pin.Value,
		CareOf:      hypervergeAadharResponse.Address.CareOf,
		District:    hypervergeAadharResponse.Address.District,
		City:        hypervergeAadharResponse.Address.City,
		Locality:    hypervergeAadharResponse.Address.Locality,
		Landmark:    hypervergeAadharResponse.Address.Landmark,
		Street:      hypervergeAadharResponse.Address.Street,
		Line1:       hypervergeAadharResponse.Address.Line1,
		Line2:       hypervergeAadharResponse.Address.Line2,
		HouseNumber: hypervergeAadharResponse.Address.HouseNumber,
		State:       hypervergeAadharResponse.Address.State,
		AddressPin:  hypervergeAadharResponse.Address.Pin,
	}

	return &aadharResponse, err

}

func (hypervergeImpl *HypervergeImpl) ReadPassport(hypervergeRequest HypervergeRequest) (*PassportResponse, error) {
	body, err := hypervergeImpl.readDocument("passport", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergePassportResponseStruct HypervergePassportResponse
	err = json.Unmarshal(body.Bytes(), &hypervergePassportResponseStruct)
	if err != nil {
		return nil, err
	}
	hypervergePassportResponse := hypervergePassportResponseStruct.Result[0].Details
	passportResponse := PassportResponse{
		CountryCode:     hypervergePassportResponse.CountryCode.Value,
		Dob:             hypervergePassportResponse.Dob.Value,
		Doe:             hypervergePassportResponse.Doe.Value,
		Doi:             hypervergePassportResponse.Doi.Value,
		Gender:          hypervergePassportResponse.Gender.Value,
		GivenName:       hypervergePassportResponse.GivenName.Value,
		Nationality:     hypervergePassportResponse.Nationality.Value,
		PassportNum:     hypervergePassportResponse.PassportNum.Value,
		PlaceOfBirth:    hypervergePassportResponse.PlaceOfBirth.Value,
		PlaceOfIssue:    hypervergePassportResponse.PlaceOfIssue.Value,
		Surname:         hypervergePassportResponse.Surname.Value,
		Mrz:             hypervergePassportResponse.Mrz.Line1,
		Type:            hypervergePassportResponse.Type.Value,
		District:        hypervergePassportResponse.Address.District,
		City:            hypervergePassportResponse.Address.City,
		Locality:        hypervergePassportResponse.Address.Locality,
		Landmark:        hypervergePassportResponse.Address.Landmark,
		Street:          hypervergePassportResponse.Address.Street,
		Line1:           hypervergePassportResponse.Address.Line1,
		Line2:           hypervergePassportResponse.Address.Line2,
		HouseNumber:     hypervergePassportResponse.Address.HouseNumber,
		State:           hypervergePassportResponse.Address.State,
		Father:          hypervergePassportResponse.Father.Value,
		Mother:          hypervergePassportResponse.Mother.Value,
		FileNum:         hypervergePassportResponse.FileNum.Value,
		OldDoi:          hypervergePassportResponse.OldDoi.Value,
		OldPassportNum:  hypervergePassportResponse.OldPassportNum.Value,
		OldPlaceOfIssue: hypervergePassportResponse.OldPlaceOfIssue.Value,
		Pin:             hypervergePassportResponse.Pin.Value,
		Spouse:          hypervergePassportResponse.Spouse.Value,
		AddressPin:      hypervergePassportResponse.Address.Pin,
	}

	return &passportResponse, err
}

func (hypervergeImpl *HypervergeImpl) ReadVotedID(hypervergeRequest HypervergeRequest) (*VoterIdResponse, error) {
	body, err := hypervergeImpl.readDocument("voter", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergeVoterIdResponseStruct HypervergeVoterIdResponse
	err = json.Unmarshal(body.Bytes(), &hypervergeVoterIdResponseStruct)
	if err != nil {
		return nil, err
	}
	hypervergeVoterIdResponse := hypervergeVoterIdResponseStruct.Result[0].Details
	voterIdResponse := VoterIdResponse{
		Voterid:     hypervergeVoterIdResponse.Voterid.Value,
		Name:        hypervergeVoterIdResponse.Name.Value,
		Gender:      hypervergeVoterIdResponse.Gender.Value,
		Relation:    hypervergeVoterIdResponse.Relation.Value,
		Dob:         hypervergeVoterIdResponse.Dob.Value,
		Doc:         hypervergeVoterIdResponse.Doc.Value,
		Age:         hypervergeVoterIdResponse.Age.Value,
		Pin:         hypervergeVoterIdResponse.Pin.Value,
		Date:        hypervergeVoterIdResponse.Date.Value,
		Type:        hypervergeVoterIdResponse.Type.Value,
		District:    hypervergeVoterIdResponse.Address.District,
		City:        hypervergeVoterIdResponse.Address.City,
		Locality:    hypervergeVoterIdResponse.Address.Locality,
		Street:      hypervergeVoterIdResponse.Address.Street,
		Line1:       hypervergeVoterIdResponse.Address.Line1,
		Line2:       hypervergeVoterIdResponse.Address.Line2,
		HouseNumber: hypervergeVoterIdResponse.Address.HouseNumber,
		State:       hypervergeVoterIdResponse.Address.State,
		AddressPin:  hypervergeVoterIdResponse.Address.Pin,
	}

	return &voterIdResponse, err
}

func newfileUploadRequest(uri, path, appKey, appID, transactionId string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("appId", appID)
	req.Header.Set("appkey", appKey)
	req.Header.Set("transactionId", transactionId)
	return req, err
}

func getURLFor(docType, baseURL string) string {
	var path string
	switch docType {
	case "aadhar":
		path = "readAadhaar"
	case "passport":
		path = "readPassport"
	case "pan":
		path = "readPAN"
	case "voter":
		path = "readVoterID"
	default:
		path = "readKYC"
	}
	return baseURL + path
}
