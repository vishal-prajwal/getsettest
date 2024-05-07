package hyperverge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"github.com/pkg/errors"
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
	req, err := newfileUploadRequest(url, hypervergeRequest.ImageFile,
		hypervergeImpl.config.GetHypervergeAppKey(), hypervergeImpl.config.GetHypervergeAppID(), hypervergeRequest.ImageName)
	if err != nil {
		return nil, errors.Wrap(ErrUploadError, err.Error())
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
		return nil, errors.Wrap(ErrHttpError, err.Error())
	}
	var hypervergePanResponse HypervergePanResponse
	err = json.Unmarshal(body.Bytes(), &hypervergePanResponse)
	if err != nil {
		return nil, err
	}
	if hypervergePanResponse.StatusCode != "200" {
		switch hypervergePanResponse.StatusCode {
		case "437":
			return nil, errors.Wrap(ErrBlurredImage, hypervergePanResponse.Error)
		case "432":
			return nil, errors.Wrap(ErrTemperedImage, hypervergePanResponse.Error)
		case "422":
			return nil, errors.Wrap(ErrInvalidDoc, hypervergePanResponse.Error)
		default:
			return nil, fmt.Errorf("status %v errorMessage %v", hypervergePanResponse.Status, hypervergePanResponse.Error)
		}
	}
	details := hypervergePanResponse.Result[0].Details
	panResponse := PanResponse{
		Date:        details.Date.Value,
		Father:      details.Father.Value,
		Name:        details.Name.Value,
		PanNo:       details.PanNo.Value,
		DateOfIssue: details.DateOfIssue.Value,
		RawResponse: body.String(),
	}

	return &panResponse, err

}

func (hypervergeImpl *HypervergeImpl) ReadAadhar(hypervergeRequest HypervergeRequest) (*AadharResponse, error) {
	body, err := hypervergeImpl.readDocument("aadhar", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergeAadharResponse HypervergeAadharResponse
	err = json.Unmarshal(body.Bytes(), &hypervergeAadharResponse)
	if err != nil {
		return nil, err
	}
	if hypervergeAadharResponse.StatusCode != "200" {
		switch hypervergeAadharResponse.StatusCode {
		case "437":
			return nil, errors.Wrap(ErrBlurredImage, hypervergeAadharResponse.Error)
		case "432":
			return nil, errors.Wrap(ErrTemperedImage, hypervergeAadharResponse.Error)
		case "422":
			return nil, errors.Wrap(ErrInvalidDoc, hypervergeAadharResponse.Error)
		default:
			return nil, fmt.Errorf("status %v errorMessage %v", hypervergeAadharResponse.Status, hypervergeAadharResponse.Error)
		}
	}
	if len(hypervergeAadharResponse.Result) > 1 {
		mergeAadharDetails(&hypervergeAadharResponse)
	}
	details := hypervergeAadharResponse.Result[0].Details
	aadharResponse := AadharResponse{
		Aadhaar:     details.Aadhaar.Value,
		Dob:         details.Dob.Value,
		Father:      details.Father.Value,
		Gender:      details.Gender.Value,
		Mother:      details.Mother.Value,
		Name:        details.Name.Value,
		Yob:         details.Yob.Value,
		Husband:     details.Husband.Value,
		Phone:       details.Phone.Value,
		Pin:         details.Pin.Value,
		CareOf:      details.Address.CareOf,
		District:    details.Address.District,
		City:        details.Address.City,
		Locality:    details.Address.Locality,
		Landmark:    details.Address.Landmark,
		Street:      details.Address.Street,
		Line1:       details.Address.Line1,
		Line2:       details.Address.Line2,
		HouseNumber: details.Address.HouseNumber,
		State:       details.Address.State,
		AddressPin:  details.Address.Pin,
		RawResponse: body.String(),
	}

	return &aadharResponse, err

}

func (hypervergeImpl *HypervergeImpl) ReadPassport(hypervergeRequest HypervergeRequest) (*PassportResponse, error) {
	body, err := hypervergeImpl.readDocument("passport", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergePassportResponse HypervergePassportResponse
	err = json.Unmarshal(body.Bytes(), &hypervergePassportResponse)
	if err != nil {
		return nil, err
	}
	if hypervergePassportResponse.StatusCode != "200" {
		switch hypervergePassportResponse.StatusCode {
		case "437":
			return nil, errors.Wrap(ErrBlurredImage, hypervergePassportResponse.Error)
		case "432":
			return nil, errors.Wrap(ErrTemperedImage, hypervergePassportResponse.Error)
		case "422":
			return nil, errors.Wrap(ErrInvalidDoc, hypervergePassportResponse.Error)
		default:
			return nil, fmt.Errorf("status %v errorMessage %v", hypervergePassportResponse.Status, hypervergePassportResponse.Error)
		}
	}
	details := hypervergePassportResponse.Result[0].Details
	passportResponse := PassportResponse{
		CountryCode:     details.CountryCode.Value,
		Dob:             details.Dob.Value,
		Doe:             details.Doe.Value,
		Doi:             details.Doi.Value,
		Gender:          details.Gender.Value,
		GivenName:       details.GivenName.Value,
		Nationality:     details.Nationality.Value,
		PassportNum:     details.PassportNum.Value,
		PlaceOfBirth:    details.PlaceOfBirth.Value,
		PlaceOfIssue:    details.PlaceOfIssue.Value,
		Surname:         details.Surname.Value,
		Mrz:             details.Mrz.Line1,
		Type:            details.Type.Value,
		District:        details.Address.District,
		City:            details.Address.City,
		Locality:        details.Address.Locality,
		Landmark:        details.Address.Landmark,
		Street:          details.Address.Street,
		Line1:           details.Address.Line1,
		Line2:           details.Address.Line2,
		HouseNumber:     details.Address.HouseNumber,
		State:           details.Address.State,
		Father:          details.Father.Value,
		Mother:          details.Mother.Value,
		FileNum:         details.FileNum.Value,
		OldDoi:          details.OldDoi.Value,
		OldPassportNum:  details.OldPassportNum.Value,
		OldPlaceOfIssue: details.OldPlaceOfIssue.Value,
		Pin:             details.Pin.Value,
		Spouse:          details.Spouse.Value,
		AddressPin:      details.Address.Pin,
		RawResponse:     body.String(),
	}

	return &passportResponse, err
}

func (hypervergeImpl *HypervergeImpl) ReadVotedID(hypervergeRequest HypervergeRequest) (*VoterIdResponse, error) {
	body, err := hypervergeImpl.readDocument("voter", hypervergeRequest)
	if err != nil {
		return nil, err
	}
	var hypervergeVoterIdResponse HypervergeVoterIdResponse
	err = json.Unmarshal(body.Bytes(), &hypervergeVoterIdResponse)
	if err != nil {
		return nil, err
	}
	if hypervergeVoterIdResponse.StatusCode != "200" {
		switch hypervergeVoterIdResponse.StatusCode {
		case "437":
			return nil, errors.Wrap(ErrBlurredImage, hypervergeVoterIdResponse.Error)
		case "432":
			return nil, errors.Wrap(ErrTemperedImage, hypervergeVoterIdResponse.Error)
		case "422":
			return nil, errors.Wrap(ErrInvalidDoc, hypervergeVoterIdResponse.Error)
		default:
			return nil, fmt.Errorf("status %v errorMessage %v", hypervergeVoterIdResponse.Status, hypervergeVoterIdResponse.Error)
		}
	}
	details := hypervergeVoterIdResponse.Result[0].Details
	voterIdResponse := VoterIdResponse{
		Voterid:     details.Voterid.Value,
		Name:        details.Name.Value,
		Gender:      details.Gender.Value,
		Relation:    details.Relation.Value,
		Dob:         details.Dob.Value,
		Doc:         details.Doc.Value,
		Age:         details.Age.Value,
		Pin:         details.Pin.Value,
		Date:        details.Date.Value,
		Type:        details.Type.Value,
		District:    details.Address.District,
		City:        details.Address.City,
		Locality:    details.Address.Locality,
		Street:      details.Address.Street,
		Line1:       details.Address.Line1,
		Line2:       details.Address.Line2,
		HouseNumber: details.Address.HouseNumber,
		State:       details.Address.State,
		AddressPin:  details.Address.Pin,
		RawResponse: body.String(),
	}

	return &voterIdResponse, err
}

func newfileUploadRequest(uri, file, appKey, appID, fileName string) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return nil, err
	}
	src := bytes.NewReader([]byte(file))
	_, err = io.Copy(part, src)
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
	return req, err
}

func getURLFor(docType, baseURL string) string {
	var path string
	switch docType {
	case "aadhar":
		path = "/readAadhaar"
	case "passport":
		path = "/readPassport"
	case "pan":
		path = "/readPAN"
	case "voter":
		path = "/readVoterID"
	default:
		path = "/readKYC"
	}
	return baseURL + path
}

func mergeAadharDetails(hypervergeAadharResult *HypervergeAadharResponse) {
	size := len(hypervergeAadharResult.Result)
	for i := size - 2; i >= 0; i-- {
		if hypervergeAadharResult.Result[i].Details.Aadhaar.Value == "" {
			hypervergeAadharResult.Result[i].Details.Aadhaar.Value = hypervergeAadharResult.Result[i+1].Details.Aadhaar.Value
		}
		if hypervergeAadharResult.Result[i].Details.Dob.Value == "" {
			hypervergeAadharResult.Result[i].Details.Dob.Value = hypervergeAadharResult.Result[i+1].Details.Dob.Value
		}
		if hypervergeAadharResult.Result[i].Details.Pin.Value == "" {
			hypervergeAadharResult.Result[i].Details.Pin.Value = hypervergeAadharResult.Result[i+1].Details.Pin.Value
		}
		if hypervergeAadharResult.Result[i].Details.Name.Value == "" {
			hypervergeAadharResult.Result[i].Details.Name.Value = hypervergeAadharResult.Result[i+1].Details.Name.Value
		}
	}
}

func (hypervergeImpl *HypervergeImpl) addHeaders(req *http.Request, txnID string) {
	req.Header.Add("appId", hypervergeImpl.config.GetHypervergeAppID())
	req.Header.Add("appKey", hypervergeImpl.config.GetHypervergeAppKey())
	req.Header.Add("transactionId", txnID)
	req.Header.Add("Content-Type", "application/json")
}

func (hypervergeImpl *HypervergeImpl) FraudCheckPan(fraudCheckPanRequest FraudCheckPanRequest, txnID string) (*FraudCheckPanResponse, error) {
	url := hypervergeImpl.config.GetHypervergeFraudCheckEndpoint() + "/verifyPAN"
	reqObj, err := json.Marshal(fraudCheckPanRequest)
	if err != nil {
		return nil, err
	}
	requestBody := bytes.NewBuffer(reqObj)
	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)
	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var fraudCheckPanResponse FraudCheckPanResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckPanResponse)
	return &fraudCheckPanResponse, err
}

func (hypervergeImpl HypervergeImpl) handlFruadCheckErrorStatusCode(res *http.Response) error {
	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 422:
			return errors.Wrap(ErrInvalidDocID, fmt.Sprintf("%d  %v", res.StatusCode, res))
		case 400:
			return errors.Wrap(ErrBadRequest, fmt.Sprintf("%d %v", res.StatusCode, res))
		case 500:
			return errors.Wrap(ErrSomethingWentWrong, fmt.Sprintf("%d %v", res.StatusCode, res))
		case 401:
			return errors.Wrap(ErrFruadCheckUnauthrised, fmt.Errorf("%d %v", res.StatusCode, res).Error())
		default:
			return errors.Wrap(ErrSomethingWentWrong, fmt.Sprintf("%d %v", res.StatusCode, res))
		}
	}
	return nil
}

func (hypervergeImpl *HypervergeImpl) FraudCheckDl(fraudCheckDlRequest FraudCheckDlRequest, txnID string) (*FraudCheckDlResponse, error) {
	url := hypervergeImpl.config.GetHypervergeEndpoint() + "/api/checkDL"
	reqObj, _ := json.Marshal(fraudCheckDlRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)
	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var fraudCheckDlResponse FraudCheckDlResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckDlResponse)
	return &fraudCheckDlResponse, err
}

func (hypervergeImpl *HypervergeImpl) FraudCheckVoter(fraudCheckVoterRequest FraudCheckVoterRequest, txnID string) (*FraudCheckVoterResponse, error) {
	url := hypervergeImpl.config.GetHypervergeEndpoint() + "/api/checkVoterId"
	reqObj, _ := json.Marshal(fraudCheckVoterRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)
	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var fraudCheckVoterResponse FraudCheckVoterResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckVoterResponse)
	return &fraudCheckVoterResponse, err
}

func (hypervergeImpl *HypervergeImpl) FraudCheckPassport(fraudCheckPassportRequest FraudCheckPassportRequest, txnID string) (*FraudCheckPassportResponse, error) {
	url := hypervergeImpl.config.GetHypervergeEndpoint() + "/api/verifyPassport"
	reqObj, _ := json.Marshal(fraudCheckPassportRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)
	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var fraudCheckPassportResponse FraudCheckPassportResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckPassportResponse)
	return &fraudCheckPassportResponse, err
}

func (hypervergeImpl *HypervergeImpl) FraudCheckAadhar(fraudCheckAadharRequest FraudCheckAadharRequest, txnID string) (*FraudCheckAadharResponse, error) {
	url := hypervergeImpl.config.GetHypervergeEndpoint() + "/api/verifyAadhaar"
	reqObj, _ := json.Marshal(fraudCheckAadharRequest)
	payload := strings.NewReader(string(reqObj))
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	hypervergeImpl.addHeaders(req, txnID)
	res, err := hypervergeImpl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = hypervergeImpl.handlFruadCheckErrorStatusCode(res)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var fraudCheckAadharResponse FraudCheckAadharResponse
	err = json.Unmarshal(body.Bytes(), &fraudCheckAadharResponse)
	return &fraudCheckAadharResponse, err
}
