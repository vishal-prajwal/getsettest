package awss3

type Config interface {
	GetBucketName() string
	GetUploaderCount() int
	GetMaxRetries() int
}
