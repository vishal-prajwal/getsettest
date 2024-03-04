package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/howzatkyc"
)

func main() {
	endpoint := "https://vkyc-inc-hzt-qa-1.howzatfantasy.com/"
	httpClient := http.Client{}

	howzat := howzatkyc.New(endpoint, newrelic.Agent{}, &httpClient)
	fmt.Println(howzat.FetchPanByUserID(998632))
}
