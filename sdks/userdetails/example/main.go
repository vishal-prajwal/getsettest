package main

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/junglee_games/getsetgo/httpclient"
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/userdetails"
)

type UserDetailsConfig struct{}

func (udc UserDetailsConfig) GetUserDetailsBaseURL() string {
	return "http://userdetails.jwr-qa-4.jwrnonprod.int"
}

func main() {
	userDetailsSDK := userdetails.New(UserDetailsConfig{}, newrelic.Agent{}, httpclient.NewHttpClient(30))
	res, err := userDetailsSDK.GetUserDetailsFromMobileNumber("8860351487", userdetails.JR)
	if err != nil {
		fmt.Println("Error in calling user details API : ", err)
		return
	}
	prettyJSON(res)
}

func prettyJSON(res interface{}) {
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
