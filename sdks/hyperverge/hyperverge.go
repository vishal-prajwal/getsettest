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

func (hypervergeImpl *HypervergeImpl) ReadDocument(documentType string, hypervergeRequest HypervergeRequest) (*HypervergeResponse, error) {
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

	var result HypervergeResponse
	err = json.Unmarshal(body.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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
