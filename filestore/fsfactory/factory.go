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

func GetFileStore(cfg Config) (filestore.FileStore, error) {
	switch cfg.GetName() {
	case AMAZON_S3:
		return awss3.NewS3Store(cfg.GetAmazonS3Config())
	case LOCAL_FILE_STORE:
		return localfilestore.NewLocalFileStore(cfg.GetLocalFileStoreConfig()), nil
	}

	return nil, ErrInvalidFileStoreName
}
