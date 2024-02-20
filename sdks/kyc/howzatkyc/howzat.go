package howzatkyc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/monitoring"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/domain"
	"github.com/pkg/errors"
)

type HowzatKycServiceClient struct {
	endpoint        string
	monitoringAgent monitoring.Agent
	httpClient      *http.Client
}

func New(endpoint string, monitoringAgent monitoring.Agent, client *http.Client) *HowzatKycServiceClient {
	return &HowzatKycServiceClient{endpoint: endpoint, monitoringAgent: monitoringAgent, httpClient: client}
}

func (howzatImpl *HowzatKycServiceClient) FetchUserByPan(userByPanRequest domain.UserByPanRequest) (domain.UserByPanResponse, error) {
	userPanInfo := make([]domain.UserPanInfo, 0)
	ch := make(chan domain.UserPanInfo)
	for _, pan := range userByPanRequest.PanNumber {
		go func(pan string) {
			userPan := domain.UserPanInfo{
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

	userByPanResponse := domain.UserByPanResponse{
		UserPanInfo: userPanInfo,
	}
	return userByPanResponse, nil
}

func (howzatImpl *HowzatKycServiceClient) getHowzatKyc(panNumber string) (*KycResponse, error) {
	tr := howzatImpl.monitoringAgent.StartTransaction(HOWZAT_USER_BY_PAN_CALL)
	defer tr.End()
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

	body, err := io.ReadAll(res.Body)
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

func (howzatImpl *HowzatKycServiceClient) FetchPanByUserID(userID int) (*domain.PanByUserResponse, error) {
	tr := howzatImpl.monitoringAgent.StartTransaction(HOWZAT_USER_BY_PAN_CALL)
	defer tr.End()
	url := howzatImpl.endpoint + fmt.Sprintf(GET_PAN_BY_USER_URI, userID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}

	res, err := howzatImpl.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(ErrCallingHowzat, err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}
	var response PanByUserResponse

	var resp domain.PanByUserResponse
	switch res.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, errors.Wrap(ErrUnmarshlingResponse, err.Error())
		}
		if response.Data.DocumentType == "PAN" && response.Data.Status == "VERIFIED" {
			resp.PanNo = response.Data.DocumentNumber
			resp.UserID = userID
			return &resp, nil
		} else {
			return nil, ErrNotFound
		}
	case http.StatusBadRequest:
		return nil, ErrReqValidate
	case http.StatusInternalServerError:
		return nil, ErrHVServer
	default:
		return nil, ErrNotFound
	}

	return &resp, nil
}
