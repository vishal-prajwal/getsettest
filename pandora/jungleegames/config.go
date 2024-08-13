package jungleegames

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/junglee_games/getsetgo/pandora/consul"
	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

// PopulateConfig will attempt to read the configuration values from the environment, ssm and read the default values
//
// The precedence of the configuration is the following and even if the key exists in one of them the value must not be
// empty for it to be considered as the appropriate value.
//
// 1. env
// 2. ssm
// 3. consul
// 3. default value if specified in tag
// 4. golang default value
//
// # Example usage
//
//	type config struct {
//		ProjectID string `env:"ACME_PROJECT_ID" ssm:"/acme/project_id" consul:"/jungleegames/svc/foo/bar" default:"THIS_IS_A_DEFAULT_VALUE"`
//	}
//
// cfg := config{}
// ctx := context.WithTimeout(context.Background(), time.Second * 3)
//
//	if err := jungleegames.PopulateConfig(ctx, &cfg); err != nil {
//	 return errors.Wrap(err, "failed to configure acme")
//	}
//
// todo: add support for different scalar values, currently only supports strings
func PopulateConfig(ctx context.Context, cfg interface{}) error {
	awsCfg, err := GetAwsConfig(ctx)
	if err != nil {
		return err
	}

	consulClient, err := consul.NewClient()
	if err != nil {
		return err
	}

	svc := ssm.NewFromConfig(awsCfg)

	t := reflect.TypeOf(cfg)

	if t.Kind() != reflect.Ptr {
		return errors.Errorf("config value passed to PopulateConfig must be a pointer to a struct")
	}

	// Since the cfg is passed as a pointer we need to access the actual struct
	// e.g ptr -> value
	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return errors.Errorf("config valued passed to PopulateConfig must be a pointer to a struct")
	}

	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i <= t.NumField()-1; i++ {
		tf := t.Field(i)
		vf := v.Field(i)
		tags := tf.Tag

		if vf.Kind() != reflect.String {
			return errors.Errorf("field: %s has a none-string type which is not supported", tf.Name)
		}

		l := log.WithFields(logrus.Fields{
			"config": fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()),
			"field":  tf.Name,
		})

		if defaultValue, ok := tags.Lookup("default"); ok {
			vf.SetString(defaultValue)
		}

		envKey, ok := tags.Lookup("env")
		if !ok {
			l.Info("missing env struct tag for field, skipping")
		} else if v := os.Getenv(envKey); v != "" {
			vf.SetString(v)

			continue
		}

		if consulKey, ok := tags.Lookup("consul"); ok {
			if val, err := consulClient.GetKVString(ctx, consulKey, ""); err != nil {
				return err
			} else if val != "" {
				vf.SetString(val)
				continue
			}
		} else {
			l.Info("Missing consul struct tag for field, skipping")
		}

		ssmKey, ok := tags.Lookup("ssm")

		switch {
		case !ok:
			l.Info("missing ssm struct tag for field, skipping")

		case IsOnAWS():
			// We need to prefix the ssm key with /prod or /staging
			// todo: this should be removed when we unique aws accounts for each env
			out, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
				Name:           aws.String(ssmKey),
				WithDecryption: aws.Bool(true),
			})

			if err != nil {
				notFoundErr := &types.ParameterNotFound{}

				if errors.As(err, &notFoundErr) {
					l.Warnf("configuration key: %s missing in ssm", ssmKey)

					continue
				}

				return errors.Wrapf(err, "invalid key %s", tf.Name)
			}

			vf.SetString(*out.Parameter.Value)

		default:
			l.Debug("not running on AWS ssm configuration skipped")
		}
	}

	return nil
}
