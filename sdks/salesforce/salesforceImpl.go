package salesforce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"
	"mime/multipart"

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"github.com/pkg/errors"
)

type SalesforceImpl struct {
	BaseURL           string
	Token             string
	DefaultAPITimeout int
	httpClient        httpclient.HTTPClient
}

func (salesforceImpl *SalesforceImpl) RequestAccessToken(ctx context.Context, accessTokenRequest AccessTokenRequest, apiTimeout int) (*AccessTokenResponse, error) {
	url := "https://test.salesforce.com/services/oauth2/token"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", accessTokenRequest.Username)
	_ = writer.WriteField("password", accessTokenRequest.Password)
	_ = writer.WriteField("grant_type", "password")
	_ = writer.WriteField("client_id", accessTokenRequest.ClientId)
	_ = writer.WriteField("client_secret", accessTokenRequest.ClientSecret)
	_ = writer.WriteField("assertion", "tmp")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var accessTokenResponse AccessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return nil, err
	}

	return &accessTokenResponse, nil
}

func (salesforceImpl *SalesforceImpl) CreateTask(ctx context.Context, createTaskRequest CreateTaskRequest, apiTimeout int) (*CreateTaskResponse, error) {

	httpReq := SaleForceCreateTaskHTTPRequest{
		UserID:        createTaskRequest.UserID,
		SocialNetwork: createTaskRequest.SocialNetwork,
	}

	createTaskRequestBytes, err := json.Marshal(httpReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, createTaskRequest.BaseURL+CreateTaskPath, bytes.NewReader(createTaskRequestBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set(contentType, applicationJson)
	request.Header.Set("Authorization", "Bearer "+createTaskRequest.AccessToken)

	resp, err := salesforceImpl.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "while making api call to Create Task")
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "reading response from Create Task")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status + body.String())
	}

	var createTaskResponse CreateTaskResponse
	err = json.Unmarshal(body.Bytes(), &createTaskResponse)
	if err != nil {
		return nil, err
	}

	if len(createTaskResponse.Errors) > 0 {
		return nil, fmt.Errorf("something went wrong")
	}
	return &createTaskResponse, nil
}
