package apilogger

import "time"

type ApiData struct {
	UserID    int64        `json:"userId"`
	ProductID int64        `json:"productId"`
	Source    string       `json:"source"`
	EventTime int64        `json:"eventTime"`
	EventType string       `json:"eventType"`
	Request   RequestData  `json:"request"`
	Response  ResponseData `json:"response"`
	Error     string       `json:"error,omitempty"`
}

type RequestData struct {
	VendorName string            `json:"vendorName"`
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Payload    string            `json:"payload"`
}

type ResponseData struct {
	StatusCode string `json:"statusCode"`
	Body       string `json:"body"`
}

type ApiDataBuilder struct {
	apiData ApiData
}

func NewApiDataBuilder() *ApiDataBuilder {
	return &ApiDataBuilder{
		apiData: ApiData{},
	}
}

func (b *ApiDataBuilder) WithBasic(userID int64, vendor string) *ApiDataBuilder {
	b.apiData.UserID = userID
	b.apiData.Request.VendorName = vendor
	return b
}

func (b *ApiDataBuilder) WithRequest(url, method, payload string, headers map[string]string) *ApiDataBuilder {
	b.apiData.Request = RequestData{
		URL:     url,
		Method:  method,
		Headers: headers,
		Payload: payload,
	}
	return b
}

func (b *ApiDataBuilder) WithResponse(statusCode string, body string) *ApiDataBuilder {
	b.apiData.Response = ResponseData{
		StatusCode: statusCode,
		Body:       body,
	}
	return b
}

func (b *ApiDataBuilder) WithError(err string) *ApiDataBuilder {
	b.apiData.Error = err
	return b
}

func (b *ApiDataBuilder) Build() *ApiData {
	b.apiData.ProductID = 3 // Assuming a fixed product ID for Rummy
	b.apiData.Source = RUMMY_GAME_AUDIT
	b.apiData.EventTime = time.Now().UnixMilli()
	b.apiData.EventType = KYC_API_USAGE_EVENT_TYPE
	return &b.apiData
}
