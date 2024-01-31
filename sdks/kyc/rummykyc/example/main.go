package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/rummykyc"
)

func main() {
	endpoint := "https://kyc-qa-3.jungleerummyqa.com"
	httpClient := http.Client{}
	kyc := rummykyc.New(endpoint, newrelic.Agent{}, &httpClient)
	fmt.Println(kyc.FetchPanByUserID(68717))
}
