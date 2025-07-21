package jwr

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	nrf "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker/v2"
)

type JWRImpl struct {
	BaseURL           string
	InternalURL       string
	Token             string
	DefaultAPITimeout int
	cb                *gobreaker.CircuitBreaker[[]byte]
}

func (this *JWRImpl) GetUserProfile(ctx context.Context, userID int, readFromDB bool, apiTimeOut int) (*UserProfile, error) {

	defer nrf.FromContext(ctx).StartSegment("GetUserProfile").End()
	var result UserProfile
	url := this.BaseURL + GetUserProfilePath + "?id=" + strconv.Itoa(userID)
	if readFromDB {
		url += "&rdcsyncuser=1"
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

func (this JWRImpl) FullUpdateProfileV2(ctx context.Context, userID int, userProfile UpdateUserProfileRequest, apiTimeOut int, retries int) error {
	defer nrf.FromContext(ctx).StartSegment("FullUpdateProfile").End()
	timeout := this.DefaultAPITimeout
	if userProfile.Gender != nil && *userProfile.Gender == "" {
		*userProfile.Gender = "UNDEFINED"
	}
	if apiTimeOut > 0 {
		timeout = apiTimeOut
	}
	json, err := json.Marshal(userProfile)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, this.InternalURL+UpdateUserProfilePathV2+strconv.Itoa(userID), bytes.NewBuffer(json))
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
	request := func() ([]byte, error) {
		resp, err := httpClient.Do(req)
		if err != nil {
			err = errors.Wrapf(err, "while making api call to PUT profile")
			return nil, err
		}
		body := &bytes.Buffer{}
		_, err = body.ReadFrom(resp.Body)
		if err != nil {
			err = errors.Wrapf(err, "reading response from profile API")
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = errors.New(resp.Status + body.String())
			return nil, err
		}
		return nil, nil
	}

	if this.cb != nil {
		_, err = this.cb.Execute(request)
	} else {
		_, err = request()
	}

	if err != nil {
		return errors.Wrapf(err, "while making api call to PUT profile")
	}
	return nil
}

func (this JWRImpl) FullUpdateProfile(ctx context.Context, userID int, userProfile UserProfile, apiTimeOut int) error {

	defer nrf.FromContext(ctx).StartSegment("FullUpdateProfile").End()
	timeout := this.DefaultAPITimeout
	if userProfile.Gender == "" {
		userProfile.Gender = "UNDEFINED"
	}
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
