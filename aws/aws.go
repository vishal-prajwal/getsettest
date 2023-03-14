package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

type AWS struct {
	sm *secretsmanager.SecretsManager
}

func New() (*AWS, error) {
	cred := aws.NewConfig()
	cred.WithRegion("ap-south-1")
	// Initialize a session in ap-south-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(cred)
	if err != nil {
		return &AWS{}, errors.Wrap(err, "getting aws session")
	}
	sm := secretsmanager.New(sess)
	return &AWS{sm}, nil
}
