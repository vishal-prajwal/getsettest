package example

import (
	"fmt"
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/domain"
	"bitbucket.org/junglee_games/getsetgo/sdks/kyc/rummykyc"
)

func main() {
	endpoint := "http://localhost:3000"
	httpClient := http.Client{}
	userByPanRequest := domain.UserByPanRequest{
		XProductID: "RUMMY",
		PanNumber:  []string{"EVOPO1403E", "BQAPV3229H", "HDAUBSS4"},
	}
	kyc := rummykyc.New(endpoint, newrelic.Agent{}, &httpClient)
	fmt.Println(kyc.FetchUserByPan(userByPanRequest))
}
