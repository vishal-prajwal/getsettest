package howzat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"github.com/pkg/errors"
)

type HowzatImpl struct {
	endpoint   string
	nr         newrelic.Agent
	httpClient *http.Client
}

func New(endpoint string, nr newrelic.Agent, client *http.Client) *HowzatImpl {
	return &HowzatImpl{endpoint: endpoint, nr: nr, httpClient: client}
}

func (howzatImpl *HowzatImpl) FetchUserByPan(userByPanRequest UserByPanRequest) (UserByPanResponse, error) {

	userPanInfo := make([]UserPanInfo, 0)
	ch := make(chan UserPanInfo)
	for _, pan := range userByPanRequest.PanNumber {
		go func(pan string) {
			userPan := UserPanInfo{
				PanNo: pan,
			}
			kycResponse, _ := howzatImpl.getHowzatKyc(pan)
			if kycResponse != nil {
				userPan.UserID = kycResponse.Data.UserID
			}
			ch <- userPan
		}(pan)
	}

	for i := 0; i < len(userByPanRequest.PanNumber); i++ {
		userInfo := <-ch
		if userInfo.UserID > 0 {
			userPanInfo = append(userPanInfo, userInfo)
		}
	}

	userByPanResponse := UserByPanResponse{
		UserPanInfo: userPanInfo,
	}
	if len(userPanInfo) == 0 {
		userByPanResponse.Error = ErrNotFound.Error()
	}
	return userByPanResponse, nil
}

func (howzatImpl *HowzatImpl) getHowzatKyc(panNumber string) (*KycResponse, error) {
	howzatImpl.nr.StartTransaction(HOWZAT_INITIATE_CALL)

	url := howzatImpl.endpoint + GET_USER_BY_PAN_URI + panNumber

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}

	res, err := howzatImpl.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(ErrCallingHowzat, err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}
	var response KycResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(ErrUnmarshlingResponse, err.Error())
	}

	return &response, nil
}
