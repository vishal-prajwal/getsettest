package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/pkg/errors"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

type Client interface {
	Publish(ctx context.Context, payload interface{}) error
}

type Opt func(c *client) error

func WithStaticTopic(topic string) Opt {
	return func(c *client) error {
		c.topicArn = topic

		return nil
	}
}

func WithTopicFromEnv(topicName string) Opt {
	return func(c *client) error {
		topicARN := os.Getenv(fmt.Sprintf("SNS_%s_ARN", topicName))

		if topicARN == "" {
			log.Warnf("Missing env key %s for sns client", topicName)
		}

		c.topicArn = topicARN

		return nil
	}
}

func NewClient(cfg aws.Config, opts ...Opt) (Client, error) {
	c := &client{
		sns: sns.NewFromConfig(cfg),
	}

	for index, opt := range opts {
		if err := opt(c); err != nil {
			return nil, errors.Wrapf(err, "failed to apply option at index %d", index)
		}
	}

	return c, nil
}

type client struct {
	sns      *sns.Client
	topicArn string
}

func (c client) Publish(ctx context.Context, payload interface{}) error {
	// todo: remove this and anything that relies on publish should use the mock when testing
	if c.topicArn == "" {
		return nil
	}

	d, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	out, err := c.sns.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(d)),
		TopicArn: aws.String(c.topicArn),
	})

	if err != nil {
		return err
	}

	log.
		WithField("msgId", out.MessageId).
		Debug("Published message to SNS")

	return nil
}
