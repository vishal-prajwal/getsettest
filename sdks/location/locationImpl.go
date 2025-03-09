package location

type LocationImpl struct {
	BaseURL           string
	DefaultAPITimeout int
}

type LocationConfig struct {
	BaseURL           string
	DefaultAPITimeout int
}

func New(cfg LocationConfig) Location {
	return &LocationImpl{
		BaseURL:           cfg.BaseURL,
		DefaultAPITimeout: cfg.DefaultAPITimeout,
	}
}

func (this *LocationImpl) ExtractStateFromLatLong(locationrequest LocationRequest) (*LocationResponse, error) {

	return &LocationResponse{StateName: "NY"}, nil
}
