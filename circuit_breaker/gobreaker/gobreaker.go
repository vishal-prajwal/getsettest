package gobreaker

import (
	"io"
	"net/http"
	"time"

	"github.com/sony/gobreaker/v2"
)

var cb *gobreaker.CircuitBreaker[[]byte]

type GobreakerCfg struct {
	Name          string
	Requests      uint32
	FailiureRatio float64
	Timeout       int
}

func GetCircutBreaker(cfg *GobreakerCfg) (*gobreaker.CircuitBreaker[[]byte], error) {
	var st gobreaker.Settings
	st.Name = cfg.Name
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= cfg.Requests && failureRatio >= cfg.FailiureRatio
	}
	st.Timeout = time.Second * time.Duration(cfg.Timeout)
	return gobreaker.NewCircuitBreaker[[]byte](st), nil
}

// Get wraps http.Get in CircuitBreaker.
func Do(url string) ([]byte, error) {
	body, err := cb.Execute(func() ([]byte, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	if err != nil {
		return nil, err
	}

	return body, nil
}
