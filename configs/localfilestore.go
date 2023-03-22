package configs

type DefaultLocalFileStoreConfig struct {
	DirectoryPath string
}

func (c *DefaultLocalFileStoreConfig) GetDirectoryPath() string {
	return c.DirectoryPath
}
