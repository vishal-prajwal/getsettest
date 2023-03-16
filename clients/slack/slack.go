package slack

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type SlackClient struct {
	config Config
}

func NewSlackClient(c Config) *SlackClient {
	return &SlackClient{config: c}
}

func (s *SlackClient) SendMessage(msg string) error {

	method := "POST"
	nmsg := fmt.Sprintf(`{"text":"%s"}`, msg)
	payload := strings.NewReader(nmsg)

	client := &http.Client{}
	req, err := http.NewRequest(method, s.config.GetURL(), payload)

	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {

		return errors.Wrap(err, "calling slack api")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("status code other than 200")
	}
	return nil
}
