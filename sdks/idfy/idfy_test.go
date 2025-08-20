package idfy

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"bitbucket.org/junglee_games/getsetgo/apilogger"
	"bitbucket.org/junglee_games/getsetgo/httpclient/mocks"
	"bitbucket.org/junglee_games/getsetgo/sdks/hyperverge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIdfyImpl_FraudCheckPanNSDL(t *testing.T) {
	type args struct {
		ctx            context.Context
		NSDLPanRequest hyperverge.NSDLPanRequest
		txnID          string
	}
	tests := []struct {
		name          string
		args          args
		mockDoFunc    func() *mocks.HTTPClient
		want      *hyperverge.NSDLPanResponse
		assertion assert.ErrorAssertionFunc
	}{
		{
			"Testcase for successful PAN validation",
			args{
				context.Background(),
				hyperverge.NSDLPanRequest{
					Pan:  "ABCDE1234F",
					Name: "JOHN DOE",
					Dob:  "1990-01-01",
				},
				"test-txn-id",
			},
			func() *mocks.HTTPClient {
				mockClient := &mocks.HTTPClient{}
				// Mock POST call
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
					StatusCode: 200,
					Body:       NewMockHTTPResponseBody(`{"request_id": "test-request-id"}`),
				}, nil).Once()

				// Mock GET call
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
					StatusCode: 200,
					Body:       NewMockHTTPResponseBody(`[{"status": "completed", "result": {"source_output": {"dob_match": true, "pan_status": "VALID", "aadhaar_seeding_status": true, "input_details": {"input_pan_number": "ABCDE1234F", "input_name": "JOHN DOE", "input_dob": "1990-01-01"}}}}`),
				}, nil).Once()
				return mockClient
			},
			&hyperverge.NSDLPanResponse{
				Status:     "completed",
				StatusCode: 200,
				Result: struct {
					Pan                 string `json:"pan"`
					PanStatus           string `json:"panStatus"`
					Name                string `json:"name"`
					Dob                 string `json:"dateOfBirth"`
					AadharSeedingStatus string `json:"aadhaarSeedingStatus"`
				}{
					Pan:                 "ABCDE1234F",
					PanStatus:           "VALID",
					Name:                "JOHN DOE",
					Dob:                 "1990-01-01",
					AadharSeedingStatus: "true",
				},
			},
			assert.NoError,
		},
		{
			"Testcase for DOB mismatch",
			args{
				context.Background(),
				hyperverge.NSDLPanRequest{
					Pan:  "ABCDE1234F",
					Name: "JOHN DOE",
					Dob:  "1990-01-01",
				},
				"test-txn-id",
			},
			func() *mocks.HTTPClient {
				mockClient := &mocks.HTTPClient{}
				// Mock POST call
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
					StatusCode: 200,
					Body:       NewMockHTTPResponseBody(`{"request_id": "test-request-id"}`),
				}, nil).Once()

				// Mock GET call
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
					StatusCode: 200,
					Body:       NewMockHTTPResponseBody(`[{"status": "completed", "result": {"source_output": {"dob_match": false, "pan_status": "VALID", "aadhaar_seeding_status": true, "input_details": {"input_pan_number": "ABCDE1234F", "input_name": "JOHN DOE", "input_dob": "1990-01-01"}}}}`),
				}, nil).Once()
				return mockClient
			},
			&hyperverge.NSDLPanResponse{
				Status:     "completed",
				StatusCode: 200,
				Result: struct {
					Pan                 string `json:"pan"`
					PanStatus           string `json:"panStatus"`
					Name                string `json:"name"`
					Dob                 string `json:"dateOfBirth"`
					AadharSeedingStatus string `json:"aadhaarSeedingStatus"`
				}{
					Pan:                 "ABCDE1234F",
					PanStatus:           "VALID",
					Name:                "JOHN DOE",
					Dob:                 "1990-01-01",
					AadharSeedingStatus: "true",
				},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, ErrDobMismatch.Error())
			},
		},
		{
			"Testcase for Post request failed",
			args{
				context.Background(),
				hyperverge.NSDLPanRequest{
					Pan:  "ABCDE1234F",
					Name: "JOHN DOE",
					Dob:  "1990-01-01",
				},
				"test-txn-id",
			},
			func() *mocks.HTTPClient {
				mockClient := &mocks.HTTPClient{}
				// Mock POST call
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("post failed")).Once()
				return mockClient
			},
			nil,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "post failed")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idfyImpl := &IdfyImpl{
				httpClient: tt.mockDoFunc(),
				config:     &mockIdfyConfig{},
				apilogger:  &mockApiLogger{},
			}
			got, err := idfyImpl.FraudCheckPanNSDL(tt.args.ctx, tt.args.NSDLPanRequest, tt.args.txnID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type mockIdfyConfig struct{}

func (m *mockIdfyConfig) GetIdfyAccountId() string {
	return "test-account-id"
}

func (m *mockIdfyConfig) GetIdfyApiKey() string {
	return "test-api-key"
}

func (m *mockIdfyConfig) GetIdfyEndpoint() string {
	return "https://test.idfy.com"
}

func (m *mockIdfyConfig) GetIdfyHealthCheckEndpoint() string {
	return "https://test.idfy.com/health"
}

func (m *mockIdfyConfig) GetIdfyRetryAttemps() int {
	return 1
}

type mockHTTPResponseBody struct {
	data string
	read int
}

func NewMockHTTPResponseBody(data string) *mockHTTPResponseBody {
	return &mockHTTPResponseBody{
		data: data,
	}
}

func (m *mockHTTPResponseBody) Read(p []byte) (n int, err error) {
	if m.read >= len(m.data) {
		return 0, io.EOF
	}
	n = copy(p, m.data[m.read:])
	m.read += n
	return n, nil
}

func (m *mockHTTPResponseBody) Close() error {
	return nil
}

type mockApiLogger struct{}

func (m *mockApiLogger) Log(ctx context.Context, apiData *apilogger.ApiData) error {
	return nil
}

func (m *mockApiLogger) GetConfig() *apilogger.Config {
	return &apilogger.Config{}
}
