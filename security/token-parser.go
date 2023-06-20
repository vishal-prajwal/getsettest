package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ExtractToken(request *http.Request) (token AuthToken) {
	claimsHeader := request.Header.Get("x-jg-claims")
	claimsHeader, _ = strconv.Unquote(string(claimsHeader))
	//fmt.Print("claimsHeader: ==", claimsHeader)
	token = AuthToken{}
	err := json.Unmarshal([]byte(claimsHeader), &token)
	if err != nil {
		fmt.Print("Error while extracting header: ==", err.Error())
	}
	//fmt.Print("Token: ==", token)
	return
}
