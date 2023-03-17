package configs

import (
	"bitbucket.org/junglee_games/getsetgo/filestore/impls/awss3"
	"bitbucket.org/junglee_games/getsetgo/filestore/impls/localfilestore"
)

type DefaultFilestoreConfig struct {
	Name  string
	AWSS3 DefaultS3Config
	Local DefaultLocalFileStoreConfig
}

func (c *DefaultFilestoreConfig) GetName() string {
	return c.Name
}

func (c *DefaultFilestoreConfig) GetAmazonS3Config() awss3.Config {
	return &c.AWSS3
}
func (c *DefaultFilestoreConfig) GetLocalFileStoreConfig() localfilestore.Config {
	return &c.Local
}
