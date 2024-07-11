package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/rummykyc"
)

func main() {
	endpoint := "http://kyc.jwr-qa-2.jwrnonprod.int"
	httpClient := http.Client{}
	kyc := rummykyc.New(endpoint, newrelic.Agent{}, &httpClient)
	fmt.Println(kyc.FetchPanByUserID(1976566, "RUMMY"))
}
