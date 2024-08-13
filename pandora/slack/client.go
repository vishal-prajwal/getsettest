package slack

import (
	"context"

	"github.com/sirupsen/logrus"
	s "github.com/slack-go/slack"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

// Client is essentially a wrapper around the official slack client but provides an interface so we can easily test
type Client interface {
	SendMessage(ctx context.Context, channel string, opt ...s.MsgOption) error
}

type clientWrapper struct {
	api *s.Client
}

func (c *clientWrapper) SendMessage(ctx context.Context, channel string, opt ...s.MsgOption) error {
	ch, timestamp, text, err := c.api.SendMessageContext(ctx, channel, opt...)

	log.WithFields(logrus.Fields{
		"channel":   ch,
		"timestamp": timestamp,
		"text":      text,
	}).Debug("Sent slack message")

	return err
}
