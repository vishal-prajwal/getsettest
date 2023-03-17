package configs

type DefaultSlackConfig struct {
	URL string
}

func (c *DefaultSlackConfig) GetURL() string {
	return c.URL
}
