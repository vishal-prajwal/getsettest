package aws

import (
	"context"
	"encoding/json"
	"strings"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

type Secrets struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// TODO: return interface from here which has method with fetch with param key to make this generic :)
func (a AWS) GetFromSM(key string) (Secrets, error) {
	var secretsVals Secrets
	var ctx context.Context = context.Background()
	output, err := a.sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
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
