package userdetails

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
)

type userDetailsSDKImpl struct {
	config     UserDetailsConfig
	nr         newrelic.Agent
	httpClient httpclient.HTTPClient
}

// New creates a new UserDetails client
func New(config UserDetailsConfig, nr newrelic.Agent, client httpclient.HTTPClient) UserDetailsSDK {
	userDetails := userDetailsSDKImpl{
		config:     config,
		nr:         nr,
		httpClient: client,
	}
	return &userDetails
}

func (UserDetailsSDK *userDetailsSDKImpl) addHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

func (userDetailsSDK *userDetailsSDKImpl) GetUserDetailsFromMobileNumber(mobile string, platform Platform) (*UserDetailFromNumber, error) {
	nrTxn := userDetailsSDK.nr.StartTransaction(NR_USER_DETAILS_FETCH_DATA_FROM_NUMBER)
	defer nrTxn.End()

	var userDetailsResponse UserDetailsResponseFromNumber

	getBaseUrl := userDetailsSDK.config.GetUserDetailsBaseURL()
	endPoint := fmt.Sprintf(USER_DETAIL_FROM_NUMBER_END_POINT, mobile, platform)

	fullURL := getBaseUrl + endPoint
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	userDetailsSDK.addHeaders(request)
	res, err := userDetailsSDK.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("statusCode: %d, body: %s", res.StatusCode, string(body))
	}

	byteResp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteResp, &userDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("response: %s, error: %v", string(byteResp), err)
	}

	var ErrUserNotFound error = fmt.Errorf("User not found")
	if userDetailsResponse.Status == STATUS_FAILURE && userDetailsResponse.Message == ERROR_USER_DOES_NOT_EXIST {
		return nil, ErrUserNotFound
	}

	if userDetailsResponse.Status != STATUS_SUCCESS || userDetailsResponse.UserData == nil {
		return nil, fmt.Errorf("status : %s, message : %s", userDetailsResponse.Status, userDetailsResponse.Message)
	}

	return userDetailsResponse.UserData, nil
}

func (userDetailsSDK *userDetailsSDKImpl) GetUserDetailsFromUserId(userId int64, platform Platform) (*UserDetailFromUserId, error) {
	nrTxn := userDetailsSDK.nr.StartTransaction(NR_USER_DETAILS_FETCH_DATA_FROM_USER_ID)
	defer nrTxn.End()

	var userDetailsResponse UserDetailsResponseFromUserId

	getBaseUrl := userDetailsSDK.config.GetUserDetailsBaseURL()
	endPoint := fmt.Sprintf(USER_DETAIL_FROM_USER_ID_END_POINT, userId, platform)

	fullURL := getBaseUrl + endPoint
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	userDetailsSDK.addHeaders(request)
	res, err := userDetailsSDK.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("statusCode: %d, body: %s", res.StatusCode, string(body))
	}

	byteResp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteResp, &userDetailsResponse)
	if err != nil {
		return nil, fmt.Errorf("response: %s, error: %v", string(byteResp), err)
	}

	var ErrUserNotFound error = fmt.Errorf("User not found")
	if userDetailsResponse.Status == STATUS_FAILURE && userDetailsResponse.Message == ERROR_USER_DOES_NOT_EXIST {
		return nil, ErrUserNotFound
	}

	if userDetailsResponse.Status != STATUS_SUCCESS || userDetailsResponse.UserData == nil {
		return nil, fmt.Errorf("status : %s, message : %s", userDetailsResponse.Status, userDetailsResponse.Message)
	}

	return userDetailsResponse.UserData, nil
}
