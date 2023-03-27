package configs

type DefaultNewrelicConfig struct {
	ServiceName string
	Key         string
}

func (c *DefaultNewrelicConfig) GetKey() string {
	return c.Key
}

func (c *DefaultNewrelicConfig) GetServiceName() string {
	return c.ServiceName
}
