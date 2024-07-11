package rummykyc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/monitoring"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/domain"
	"github.com/pkg/errors"
)

type KYCImpl struct {
	endpoint        string
	monitoringAgent monitoring.Agent
	httpClient      *http.Client
}

func New(endpoint string, agent monitoring.Agent, client *http.Client) *KYCImpl {
	return &KYCImpl{endpoint: endpoint, monitoringAgent: agent, httpClient: client}
}

func (kyc *KYCImpl) FetchUserByPan(userByPanRequest domain.UserByPanRequest) (*domain.UserByPanResponse, error) {
	kyc.monitoringAgent.StartTransaction(KYC_INITIATE_CALL)

	url := kyc.endpoint + GET_USER_BY_PAN_URI + encodeQueryParams("?panNos=", userByPanRequest.PanNumber)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set(X_PRODUCT_ID, userByPanRequest.XProductID)

	res, err := kyc.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(ErrCallingKYC, err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}
	var response domain.UserByPanResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(ErrUnmarshlingResponse, err.Error())
	}
	switch res.StatusCode {
	case http.StatusOK:
		return &response, nil
	case http.StatusBadRequest:
		return nil, errors.Wrap(ErrReqValidate, response.Error)
	case http.StatusInternalServerError:
		return nil, errors.Wrap(ErrHVServer, response.Error)
	}

	return &response, nil
}

func (kyc *KYCImpl) FetchPanByUserID(userID int, productID string) (*domain.PanByUserResponse, error) {
	kyc.monitoringAgent.StartTransaction(KYC_INITIATE_CALL)

	url := kyc.endpoint + fmt.Sprintf(GET_PAN_BY_USER, userID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(ErrCreatingRequest, err.Error())
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set(X_PRODUCT_ID, productID)

	res, err := kyc.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(ErrCallingKYC, err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(ErrReadingResponseBody, err.Error())
	}
	var response Kyc

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(ErrUnmarshlingResponse, err.Error())
	}
	var resp domain.PanByUserResponse
	switch res.StatusCode {
	case http.StatusOK:
		if response.PanProof.Status == "APPROVED" {
			resp.PanNo = response.PanProof.DocumentID
			resp.UserID = userID
			return &resp, nil
		} else {
			return nil, ErrNotFound
		}
	case http.StatusBadRequest:
		return nil, ErrReqValidate
	case http.StatusInternalServerError:
		return nil, ErrHVServer
	}

	return &resp, nil
}

// encodeQueryParams encodes query parameters for a URL.
func encodeQueryParams(key string, values []string) string {
	return key + strings.Join(values, ",")
}
