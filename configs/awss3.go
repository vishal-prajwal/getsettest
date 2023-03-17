package configs

type DefaultS3Config struct {
	Bucket        string
	UploaderCount int
	MaxRetries    int
}

func (c *DefaultS3Config) GetBucketName() string {
	return c.Bucket
}

func (c *DefaultS3Config) GetUploaderCount() int {
	return c.UploaderCount
}

func (c *DefaultS3Config) GetMaxRetries() int {
	return c.MaxRetries
}
