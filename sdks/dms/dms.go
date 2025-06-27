package dms

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"

	dmsService "bitbucket.org/junglee_games/japi-dms-service/models/dto"
	"github.com/pkg/errors"
)

type DMSImpl struct {
	endpoint   string
	nr         newrelic.Agent
	httpClient *http.Client
}

func New(endpoint string, nr newrelic.Agent, client *http.Client) *DMSImpl {
	return &DMSImpl{endpoint: endpoint, nr: nr, httpClient: client}
}

func (dms *DMSImpl) Initiate(req IntiateRequest) (*IntiateResponse, error) {
	dms.nr.StartTransaction(DMS_INITIATE_CALL)
	url := dms.endpoint + "/v1/documents"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("userID", strconv.Itoa(req.UserID))
	_ = writer.WriteField("docType", req.DocType)
	file, errFile3 := os.Open(req.FrontImage)
	if errFile3 != nil {
		return nil, errors.Wrap(ErrFileNotFound, errFile3.Error())
	}
	defer file.Close()
	part3, errFile3 := writer.CreateFormFile("front", filepath.Base(req.FrontImage))
	if errFile3 != nil {
		return nil, errors.Wrap(ErrFileNotFound, errFile3.Error())
	}
	_, errFile3 = io.Copy(part3, file)

	if errFile3 != nil {
		return nil, errors.Wrap(ErrFileNotFound, errFile3.Error())
	}
	file, errFile4 := os.Open(req.BackImage)
	if errFile4 != nil {
		return nil, errors.Wrap(ErrFileNotFound, errFile4.Error())
	}
	defer file.Close()
	part4, errFile4 := writer.CreateFormFile("back", filepath.Base(req.BackImage))
	if errFile4 != nil {
		return nil, errors.Wrap(ErrFileNotFound, errFile4.Error())
	}
	_, errFile4 = io.Copy(part4, file)
	if errFile4 != nil {
		return nil, errors.Wrap(ErrFileNotFound, errFile4.Error())
	}
	err := writer.Close()
	if err != nil {
		return nil, errors.Wrap(ErrFileNotFound, err.Error())
	}

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}
	request.Header.Add("X-PRODUCT-ID", req.XProductID)

	request.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := dms.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(ErrCallingDMS, err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}
	var result IntiateResponse
	result.ID = string(body)
	// err = json.Unmarshal(body, &result)
	// if err != nil {
	// 	return nil, errors.Wrap(ErrUnmarshlingResponse, err.Error())
	// }
	return &result, nil
}

func (dms *DMSImpl) SaveKycInfo(req dmsService.KycInfo) (*dmsService.KycInfoResponse, error) {
	dms.nr.StartTransaction(DMS_SAVE_KYC_INFO_CALL)
	url := dms.endpoint + "/v1/internal/KycInfo"
	method := "POST"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-PRODUCT-ID", req.ProductID)

	res, err := dms.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(ErrCallingDMS, err.Error())
	}
	if res.Body == nil {
		return nil, errors.New("response body is nil")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("DMS call failed with status code: %d, response: %s", res.StatusCode, body)
	}

	var result dmsService.KycInfoResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}
	return &result, nil
}
