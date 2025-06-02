package hyperverge

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
)

type HypvergeTest struct {
}

func (this HypvergeTest) GetHypervergeAppID() string {
	return "376412"
}
func (this HypvergeTest) GetHypervergeAppKey() string {
	return "ee80e13a789929c70dd0"
}
func (this HypvergeTest) GetHypervergeEndpoint() string {
	return "https://ind-docs.hyperverge.co/v2.0"
}
func (this HypvergeTest) GetHypervergeFraudCheckEndpoint() string {
	return ""
}
func (this HypvergeTest) GetHypervergeFraudCheckAadharEndpoint() string {
	return ""
}

func (this HypvergeTest) GetHypervergeNSDLUrl() string {
	return ""
}
func TestHypverge(t *testing.T) {
	var hvConfig HypvergeTest
	hvClient := New(hvConfig, newrelic.Agent{}, httpclient.NewHttpClient(10), &configs.DefaultKafkaConfig{})

	bytes, err := ioutil.ReadFile("tt_front.jpeg")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	res, err := hvClient.ReadAadhar(context.Background(), HypervergeRequest{
		ImageFile: string(bytes),
	})
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// fmt.Println(toBase64(bytes))
	t.Log(res)
}
