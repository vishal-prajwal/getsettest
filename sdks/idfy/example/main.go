package example

import (
	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/idfy"
)

type IdfyConfig struct {
	ApiKey    string
	Something string
}

func (this IdfyConfig) GetIdfyAccountId() string {
	return "blah blah"
}
func (this IdfyConfig) GetIdfyApiKey() string {
	return "blah blah"
}
func (this IdfyConfig) GetIdfyEndpoint() string {
	return "blah blah"
}
func (this IdfyConfig) GetIdfyFraudCheckPostEndpoint() string {
	return "blah blah"
}
func (this IdfyConfig) GetIdfyFraudCheckGetEndpoint() string {
	return "blah blah"
}

func (this IdfyConfig) GetIdfyHealthCheckEndpoint() string {
	return "blah blah"
}

func main() {
	d := idfy.New(IdfyConfig{}, newrelic.Agent{}, httpclient.NewHttpClient(30))
	d.FraudCheckAadhar(idfy.FraudCheckRequest{})
}
