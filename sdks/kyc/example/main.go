package example

import (
	"fmt"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc"
)

func main() {
	endpoint := "http://localhost:3000"
	httpClient := http.Client{}
	userByPanRequest := kyc.UserByPanRequest{
		XProductID: "RUMMY",
		PanNumber:  []string{"EVOPO1403E", "BQAPV3229H", "HDAUBSS4"},
	}
	kyc := kyc.New(endpoint, newrelic.Agent{}, &httpClient)
	fmt.Println(kyc.FetchUserByPan(userByPanRequest))
}
