package jungleegames

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	nraws "github.com/newrelic/go-agent/v3/integrations/nrawssdk-v2"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

var (
	ec2InstanceDataOnce         sync.Once
	ec2InstanceIdentityDocument imds.InstanceIdentityDocument
)

func GetEc2InstanceData() imds.InstanceIdentityDocument {
	ec2InstanceDataOnce.Do(func() {
		// We do this because it otherwise we have to wait 3 seconds for the detection to fail
		if os.Getenv("CI_SERVER") == "yes" {
			log.Infof("env CI_SERVER set to true, disabled ec2 instance data")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		cfg, err := GetAwsConfig(ctx)
		if err != nil {
			log.WithError(err).Error("failed to load aws configuration")
			return
		}

		svc := imds.NewFromConfig(cfg)

		out, err := svc.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})

		if err != nil {
			var oe *smithy.OperationError

			if errors.As(err, &oe); errors.Unwrap(oe.Err) == context.DeadlineExceeded {
				log.Infof("Not running on AWS")
			} else {
				log.WithError(err).Info("Failed to get instance information")
			}

			return
		}

		ec2InstanceIdentityDocument = out.InstanceIdentityDocument

		log.WithFields(logrus.Fields{
			"instanceId":       out.InstanceID,
			"instanceType":     out.InstanceType,
			"availabilityZone": out.AvailabilityZone,
		}).Info("Running on ec2 instance")
	})

	return ec2InstanceIdentityDocument
}

func GetAwsConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion("ap-south-1"))
	if err != nil {
		return aws.Config{}, errors.Wrap(err, "failed to load default config")
	}

	// This ensures all the AWS calls are monitored by New Relic
	nraws.AppendMiddlewares(&cfg.APIOptions, nil)

	return cfg, nil
}

func IsOnAWS() bool {
	return GetEc2InstanceData().InstanceID != ""
}
