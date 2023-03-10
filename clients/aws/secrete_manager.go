package aws

import (
	"context"
	"encoding/json"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

type SecreteManager struct {
	sm *secretsmanager.SecretsManager
}

func NewSecreteManager() (*SecreteManager, error) {
	cred := aws.NewConfig()
	cred.WithRegion("ap-south-1")
	// Initialize a session in ap-south-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(cred)
	if err != nil {
		return nil, errors.Wrap(err, "getting aws session")
	}
	return &SecreteManager{secretsmanager.New(sess)}, nil
}

type Secrets struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (sm *SecreteManager) GetFromSM(ctx context.Context, key string) (Secrets, error) {
	var secretsVals Secrets
	output, err := sm.sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &key,
	})
	if err != nil {
		logger.Error(ctx, err.Error())
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				logger.Error(ctx, aerr.Error())
				logger.Error(ctx, "Secrets Manager could not decrypt the secret.")
			case secretsmanager.ErrCodeInternalServiceError:
				logger.Error(ctx, aerr.Error())
				logger.Error(ctx, "Server side error.")
			case secretsmanager.ErrCodeInvalidParameterException:
				logger.Error(ctx, aerr.Error())
				logger.Error(ctx, "Invalid parameter. Check inputs.")
			case secretsmanager.ErrCodeInvalidRequestException:
				logger.Error(ctx, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				logger.Error(ctx, aerr.Error())
				logger.Error(ctx, "Is your secret name correct?")
			}
		}
		return secretsVals, err
	}
	dec := json.NewDecoder(strings.NewReader(*output.SecretString))
	if err := dec.Decode(&secretsVals); err != nil {
		return secretsVals, errors.Wrapf(err, "getting key %s from secret manager", key)
	}
	return secretsVals, nil
}
