package configs

type DefaultRedisConfig struct {
	Host     string
	Password string
}

func (c *DefaultRedisConfig) GetHost() string {
	return c.Host
}

func (c *DefaultRedisConfig) GetPassword() string {
	return c.Password
}
