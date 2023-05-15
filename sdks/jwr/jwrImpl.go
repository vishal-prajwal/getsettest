package jwr

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type JWRImpl struct {
	BaseURL           string
	Token             string
	DefaultAPITimeout int
}

func (this *JWRImpl) GetUserProfile(ctx context.Context, userID int, apiTimeOut int) (*UserProfile, error) {

	var result UserProfile
	request, err := http.NewRequest(http.MethodGet, this.BaseURL+GetUserProfilePath+"?id="+strconv.Itoa(userID), nil)
	if err != nil {
		return nil, err
	}

	timeout := this.DefaultAPITimeout
	if apiTimeOut > 0 {
		timeout = apiTimeOut
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", this.Token)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	timeoutDur := time.Duration(timeout) * time.Second
	httpClient := http.Client{
		Timeout:   timeoutDur,
		Transport: t,
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "while making api call to GET profile")
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "reading response from profile API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status + body.String())
	}
	err = json.Unmarshal(body.Bytes(), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (this JWRImpl) FullUpdateProfile(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int) error {

	timeout := this.DefaultAPITimeout
	if apiTimeOut > 0 {
		timeout = apiTimeOut
	}
	json, err := json.Marshal(userProfile)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, this.BaseURL+UpdateUserProfilePath+"?id="+strconv.Itoa(userID), bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", this.Token)
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	timeoutDur := time.Duration(timeout) * time.Second
	httpClient := http.Client{
		Timeout:   timeoutDur,
		Transport: t,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "while making api call to PUT profile")
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "reading response from profile API")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status + body.String())
		return err
	}
	return nil
}
