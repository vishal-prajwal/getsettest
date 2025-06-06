package apilogger

import (
	"context"
	"time"

	"bitbucket.org/junglee_games/getsetgo/web"
)

type ApiData struct {
	UserID    string       `json:"userId"`
	RequestID string       `json:"requestId,omitempty"` // Optional, can be used for tracking specific requests
	Reason    string       `json:"reason,omitempty"`    // Optional, can be used to provide additional context or reason for the API call
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

func NewApiDataBuilder(cfg Config) *ApiDataBuilder {
	return &ApiDataBuilder{
		apiData: ApiData{
			ProductID: cfg.ProductID,
			Source:    cfg.Source,
			EventType: cfg.EventType,
		},
	}
}

func (b *ApiDataBuilder) WithRequestContext(ctx context.Context) *ApiDataBuilder {
	// Extracting user ID and request ID from the conte
	if userID, ok := ctx.Value(web.WebKeyValues).(*web.WebValues); ok {
		b.apiData.UserID = userID.UserID
	}
	if requestID, ok := ctx.Value(web.WebKeyValues).(*web.WebValues); ok {
		b.apiData.RequestID = requestID.RequestID
	}
	return b
}

func (b *ApiDataBuilder) WithBasic(ctx context.Context, vendor, reason string) *ApiDataBuilder {
	b.apiData.Request.VendorName = vendor
	b.apiData.Reason = reason
	return b.WithRequestContext(ctx)
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
	b.apiData.EventTime = time.Now().UnixMilli()
	return &b.apiData
}
