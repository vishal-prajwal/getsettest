package hyperverge

type MatchOutput struct {
	Value string `json:"value"`
	Conf  int    `json:"conf"`
	Pin   string `json:"pin"`
}

type Details map[string]MatchOutput

type Result struct {
	Details Details `json:"details"`
	KycType string  `json:"type"`
}

type HypervergeResponse struct {
	Status     string   `json:"status,omitempty"`
	StatusCode string   `json:"statusCode,omitempty"`
	Err        string   `json:"error,omitempty"`
	Res        []Result `json:"result,omitempty"`
}

type HypervergeRequest struct {
	Path                string `json:"path"`
	EnableDashboard     string `json:"enableDashboard"`
	MaskAadhaarComplete string `json:"maskAadhaarComplete"`
	OutputImageUrl      string `json:"outputImageUrl"`
}
