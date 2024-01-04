package digilocker

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	httpclientmocks "bitbucket.org/junglee_games/getsetgo/httpclient/mocks"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDigilockerSuite(t *testing.T) {
	suite.Run(t, new(digilockerSuite))
}

type digilockerSuite struct {
	suite.Suite
	httpClient httpclientmocks.HTTPClient
	nr         newrelic.Agent
	srv        Digilocker
}

func (suite *digilockerSuite) SetupTest() {
	suite.httpClient = *httpclientmocks.NewHTTPClient(suite.T())
	//	suite.nr = *newrelicmocks.NewAgent(suite.T())
	//	suite.nr.On("StartTransaction", mock.Anything).Return(nil)
	suite.srv = New("test", "testing", "url.com", &suite.httpClient, "")
}

func (suite *digilockerSuite) TestNew() {
	type args struct {
		appId  string
		appKey string
	}
	tests := []struct {
		name string
		args args
		want *DigilockerImpl
	}{
		{name: "1",
			args: args{
				appId:  "123",
				appKey: "appKey",
			},
			want: &DigilockerImpl{
				appId:  "123",
				appKey: "appKey",
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if got := New(tt.args.appId, tt.args.appKey, "", nil, ""); !reflect.DeepEqual(got, tt.want) {
				suite.Equal(tt.want, got)
			}
		})
	}
}

func (suite *digilockerSuite) TestDigilocker_addHeaders() {
	type fields struct {
		AppId  string
		AppKey string
	}
	type args struct {
		req *http.Request
	}
	reqArgs, _ := http.NewRequest(http.MethodPost, "", nil)
	req, _ := http.NewRequest(http.MethodPost, "", nil)
	req.Header.Add("appId", "id")
	req.Header.Add("appKey", "key")
	req.Header.Add("Content-Type", "application/json")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *http.Request
	}{
		{
			name: "1",
			fields: fields{AppId: "id",
				AppKey: "key",
			},
			args: args{req: reqArgs},
			want: req,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			dl := &DigilockerImpl{
				appId:  tt.fields.AppId,
				appKey: tt.fields.AppKey,
			}
			dl.addHeaders(tt.args.req)
			suite.Equal(tt.want, tt.args.req)
		})
	}
}

func (suite *digilockerSuite) TestDigilockerImpl_GetAddharDetails() {

	type args struct {
		transactionId string
		referenceId   string
	}
	tests := []struct {
		name    string
		args    args
		want    *AadhaarDetails
		doErr   error
		wantErr bool
	}{

		{
			name:    "when digilocker returns error",
			want:    nil,
			doErr:   errors.New("some error"),
			wantErr: true,
		},
		{
			name: "when digilocker works fine",
			want: &AadhaarDetails{
				Name:    "Jon Doe",
				DOB:     "30-08-1996",
				Address: "Somewhere on the earth",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var str []byte
			if tt.want != nil {
				str, _ = json.Marshal(&AadhaarDetailsResponse{Result: *tt.want})
			} else {
				str = nil
			}
			suite.httpClient.ExpectedCalls = []*mock.Call{}
			suite.httpClient.Calls = []mock.Call{}
			suite.httpClient.On("Do", mock.Anything).Return(&http.Response{Status: "Okay", StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(string(str)))}, tt.doErr)
			got, err := suite.srv.GetAddharDetails(context.TODO(), tt.args.transactionId, tt.args.referenceId)
			suite.Equal(tt.want, got)
			suite.Equal(tt.wantErr, err != nil)

		})
	}
}

func (suite *digilockerSuite) TestDigilockerImpl_CheckAccountstatus() {

	type fields struct {
		AppId  string
		AppKey string
	}
	type args struct {
		mobile  string
		aadhaar string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccountStatusDetails
		doErr   error
		wantErr bool
	}{

		{
			name:    "when digilocker returns error",
			want:    nil,
			doErr:   errors.New("some error"),
			wantErr: true,
		},
		{
			name: "when digilocker works fine",
			want: &AccountStatusDetails{
				Code:    "200",
				Mobile:  "704030****",
				Aadhaar: "yes",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var str []byte
			if tt.want != nil {
				str, _ = json.Marshal(&AccountStatusResponse{Result: *tt.want})
			} else {
				str = nil
			}
			httpRes := &http.Response{Status: "Okay", StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(string(str)))}
			suite.httpClient.ExpectedCalls = []*mock.Call{}
			suite.httpClient.Calls = []mock.Call{}
			suite.httpClient.On("Do", mock.Anything).Return(httpRes, tt.doErr)
			got, err := suite.srv.CheckAccountstatus(context.TODO(), tt.args.mobile, tt.args.aadhaar)
			suite.Equal(tt.want, got)
			suite.Equal(tt.wantErr, err != nil)
		})
	}
}

func (suite *digilockerSuite) TestDigilockerImpl_StartKYC() {
	type fields struct {
		AppId  string
		AppKey string
	}
	type args struct {
		transactionId string
		referenceId   string
		redirectURL   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *KYCStartDetails
		wantErr bool
		doErr   error
	}{

		{
			name:    "when digilocker returns error",
			want:    nil,
			doErr:   errors.New("some error"),
			wantErr: true,
		},
		{
			name: "when digilocker works fine",

			want: &KYCStartDetails{
				URL: "www.google.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {

			var str []byte
			if tt.want != nil {
				str, _ = json.Marshal(&KYCStartResponse{Result: *tt.want})
			} else {
				str = nil
			}
			httpRes := &http.Response{Status: "Okay", StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(string(str)))}
			suite.httpClient.ExpectedCalls = []*mock.Call{}
			suite.httpClient.Calls = []mock.Call{}
			suite.httpClient.On("Do", mock.Anything).Return(httpRes, tt.doErr)
			got, err := suite.srv.StartKYC(context.TODO(), tt.args.transactionId, tt.args.referenceId, tt.args.redirectURL)
			suite.Equal(tt.want, got)
			suite.Equal(tt.wantErr, err != nil)
		})
	}
}

func (suite *digilockerSuite) TestDigilockerImpl_Healthcheck() {
	type fields struct {
	}
	tests := []struct {
		name    string
		fields  fields
		want    *HealthcheckResult
		wantErr bool
		doErr   error
	}{
		{
			name:    "when digilocker returns error",
			want:    nil,
			doErr:   errors.New("some error"),
			wantErr: true,
		},
		{
			name: "when digilocker works fine",

			want: &HealthcheckResult{
				Severity: "Available",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {

			var str []byte
			if tt.want != nil {
				str, _ = json.Marshal(&HypervergeHealthcheckResponse{Result: *tt.want})
			} else {
				str = nil
			}
			httpRes := &http.Response{Status: "Okay", StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(string(str)))}
			suite.httpClient.ExpectedCalls = []*mock.Call{}
			suite.httpClient.Calls = []mock.Call{}

			suite.httpClient.On("Do", mock.Anything).Return(httpRes, tt.doErr)

			got, err := suite.srv.Healthcheck(context.TODO())
			suite.Equal(tt.want, got)
			suite.Equal(tt.wantErr, err != nil)
		})
	}
}
