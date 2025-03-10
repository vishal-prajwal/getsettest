package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LocationImpl struct {
	cfg                        *LocationConfig
	httpClient                 *http.Client
	nameOrShortToIndiaStateMap NameOrShortToIndiaStateMap
}

type LocationConfig struct {
	BaseURL           string
	DefaultAPITimeout int
}

func New(cfg *LocationConfig) Location {
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.DefaultAPITimeout) * time.Second,
	}

	return &LocationImpl{
		httpClient:                 httpClient,
		cfg:                        cfg,
		nameOrShortToIndiaStateMap: indiaStates.ToReverseMap(),
	}
}

func (locSDK *LocationImpl) ExtractStateFromLatLong(locationrequest LocationRequest) (*LocationResponse, error) {

	url := fmt.Sprintf("%s/state/%f,%f", locSDK.cfg.BaseURL, locationrequest.Latitude, locationrequest.Longitude)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := locSDK.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var locationResponse LocationServiceResponse
	err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		return nil, err
	}

	return &LocationResponse{
		StateName:        locationResponse.State.LongName,
		StateShortCode:   locationResponse.State.ShortName,
		CountryName:      locationResponse.Country.LongName,
		CountryShortCode: locationResponse.Country.ShortName,
		CityName:         locationResponse.City.LongName,
		LocalityName:     locationResponse.Locality.LongName,
		PostalCode:       locationResponse.PostalCode.LongName,
	}, nil
}

func (locSDK *LocationImpl) GetValidState(state string) (*IndiaState, error) {
	indiaState, exists := locSDK.nameOrShortToIndiaStateMap[state]
	if !exists {
		return nil, ErrInvalidState
	}
	return indiaState, nil
}
