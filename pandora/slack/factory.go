package slack

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

type config struct {
	Token string `ssm:"/slack/token" env:"SLACK_TOKEN"`
}

func NewClient(ctx context.Context) (Client, error) {
	cfg := config{}

	if err := jungleegames.PopulateConfig(ctx, &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to get slack config")
	}

	return &clientWrapper{
		api: slack.New(cfg.Token),
	}, nil
}
