package slack

type Config interface {
	GetURL() string
}

type SlackConfig struct {
	URL string
}

func (sc *SlackConfig) GetURL() string {
	return sc.URL
}
