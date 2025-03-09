package location

type LocationRequest struct {

	// Latitude of the location
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationResponse struct {
	StateName        string `json:"state"`
	StateShortCode   string `json:"state_short_code"`
	CountryName      string `json:"country"`
	CountryShortCode string `json:"country_short_code"`
	CityName         string `json:"city"`
	LocalityName     string `json:"locality"`
	PostalCode       string `json:"postal_code"`
}

type LocationServiceResponse struct {
	Country struct {
		LongName  string `json:"long_name"`
		ShortName string `json:"short_name"`
	} `json:"country"`
	State struct {
		LongName  string `json:"long_name"`
		ShortName string `json:"short_name"`
	} `json:"state"`
	City struct {
		LongName  string `json:"long_name"`
		ShortName string `json:"short_name"`
	} `json:"city"`
	Locality struct {
		LongName  string `json:"long_name"`
		ShortName string `json:"short_name"`
	} `json:"locality"`
	PostalCode struct {
		LongName  string `json:"long_name"`
		ShortName string `json:"short_name"`
	} `json:"postal_code"`
}

type Location interface {
	ExtractStateFromLatLong(locationrequest LocationRequest) (*LocationResponse, error)
}
