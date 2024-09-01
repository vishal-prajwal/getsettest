package config

type Config struct {
	Server     ServerConfigurations
	Repository RepositoryConfig
	AppConfig  *ApplicationConfig
	Cache      CacheConfig
}

type CacheConfig struct {
	Use          string
	RedisCluster RedisClusterConfig
	RedisSimple  RedisSimpleConfig
}

type RedisSimpleConfig struct {
	Addrs    string
	PoolSize int
}

type RedisClusterConfig struct {
	Addrs    string
	PoolSize int
}

// ServerConfigurations exported
type ServerConfigurations struct {
	APIHost         string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type RepositoryConfig struct {
	Type        string
	RedisConfig RedisConfig
}

type RedisConfig struct {
	Host string
}

type ApplicationConfig struct {
}
