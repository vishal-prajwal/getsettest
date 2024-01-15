package main

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/howzat"
)

func main() {
	endpoint := "https://vkyc-inc-hzt-qa-1.howzatfantasy.com/"
	httpClient := http.Client{}
	userByPanRequest := howzat.UserByPanRequest{
		PanNumber: []string{"EVOPO1403E", "BQAPV3229H", "GXEPS1676F"},
	}

	t := time.Now()
	howzat := howzat.New(endpoint, newrelic.Agent{}, &httpClient)
	fmt.Println(howzat.FetchUserByPan(userByPanRequest))
	fmt.Println("Time Lapsed : ", time.Now().Sub(t).Seconds())
}
