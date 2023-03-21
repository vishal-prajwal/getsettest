package fsfactory

import (
	"bitbucket.org/junglee_games/getsetgo/filestore"
	"bitbucket.org/junglee_games/getsetgo/filestore/impls/awss3"
	"bitbucket.org/junglee_games/getsetgo/filestore/impls/localfilestore"
)

type Config interface {
	GetName() string
	GetAmazonS3Config() awss3.Config
	GetLocalFileStoreConfig() localfilestore.Config
}

type Factory struct {
	config Config
}

func NewFileStoreFactory(config Config) *Factory {
	return &Factory{config: config}
}

func (fsf *Factory) GetFileStore(name string) (filestore.FileStore, error) {
	switch name {
	case AMAZON_S3:
		return awss3.NewS3Store(fsf.config.GetAmazonS3Config())
	case LOCAL_FILE_STORE:
		return localfilestore.NewLocalFileStore(fsf.config.GetLocalFileStoreConfig()), nil
	}

	return nil, ErrInvalidFileStoreName
}
