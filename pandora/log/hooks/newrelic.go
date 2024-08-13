package hooks

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type newrelicEntry struct {
	Timestamp  int64                  `json:"timestamp"`
	Attributes map[string]interface{} `json:"attributes"`
	Message    string                 `json:"message"`
}

type common struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type payload struct {
	Common common          `json:"common"`
	Logs   []newrelicEntry `json:"logs"`
}

func Newrelic(license string) *nrHook {
	client := retryablehttp.NewClient()
	client.Logger = nil

	return &nrHook{
		license:    license,
		client:     client,
		stack:      []newrelicEntry{},
		stackMtx:   &sync.Mutex{},
		lastUpdate: time.Now(),
	}
}

type nrHook struct {
	license string
	client  *retryablehttp.Client

	stack      []newrelicEntry
	stackMtx   *sync.Mutex
	lastUpdate time.Time
}

func (n nrHook) Levels() []logrus.Level {
	return defaultLevels
}

func (n *nrHook) Fire(entry *logrus.Entry) error {
	e := newrelicEntry{
		Timestamp: entry.Time.Unix(),
		Attributes: map[string]interface{}{
			"level": entry.Level,
		},
		Message: entry.Message,
	}

	for k, v := range entry.Data {
		e.Attributes[k] = formatData(v) // use default formatter
	}

	if ctx := entry.Context; nil != ctx {
		if txn := newrelic.FromContext(ctx); nil != txn {
			logcontext.AddLinkingMetadata(e.Attributes, txn.GetLinkingMetadata())
		}
	}

	n.stackMtx.Lock()
	n.stack = append(n.stack, e)

	if time.Since(n.lastUpdate).Seconds() >= 10 || len(n.stack) > 100 {
		stack := n.stack

		go n.send(stack)

		n.stack = []newrelicEntry{}
		n.lastUpdate = time.Now()
	}

	n.stackMtx.Unlock()

	return nil
}

func (n *nrHook) Flush() {
	n.stackMtx.Lock()
	defer n.stackMtx.Unlock()

	n.send(n.stack)
	n.stack = []newrelicEntry{}
}

func (n nrHook) send(stack []newrelicEntry) {
	if len(stack) == 0 {
		return
	}

	hostname, _ := os.Hostname()

	body, err := json.Marshal([]payload{{
		Common: common{
			Attributes: map[string]interface{}{
				"hostname":    hostname,
				"environment": os.Getenv("SENTRY_ENV"),
				"dist":        os.Getenv("SENTRY_DIST"),
			},
		},
		Logs: stack,
	}})

	if err != nil {
		sentry.CaptureException(err)
		return
	}

	buf := &bytes.Buffer{}

	gz := gzip.NewWriter(buf)

	if _, err := gz.Write(body); err != nil {
		sentry.CaptureException(err)
		return
	}

	if err := gz.Close(); err != nil {
		sentry.CaptureException(err)
		return
	}

	req, err := retryablehttp.NewRequest(http.MethodPost, "https://log-api.eu.newrelic.com/log/v1", buf)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/gzip")
	req.Header.Set("X-License-Key", n.license)

	resp, err := n.client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusAccepted {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			data, _ := ioutil.ReadAll(resp.Body)

			scope.SetExtra("responseBody", string(data))

			sentry.CaptureException(err)
		})
	}
}
