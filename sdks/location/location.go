package location

type LocationRequest struct {

	// Latitude of the location
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationResponse struct {
	StateName      string `json:"state"`
	StateShortCode string `json:"state_short_code"`
	StateCode      int    `json:"state_code"`
}

type Location interface {
	ExtractStateFromLatLong(locationrequest LocationRequest) (*LocationResponse, error)
}
