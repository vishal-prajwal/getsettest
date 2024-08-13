package jungleegames

import (
	"github.com/Unleash/unleash-client-go/v3"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

// nolint: gochecknoinits
func init() {
	if AppEnv == "" {
		return
	}

	err := unleash.Initialize(
		unleash.WithUrl("https://gitlab.com/api/v4/feature_flags/unleash/19693690"),
		unleash.WithInstanceId("aNP2H8ymuVSsx595d_Yu"),

		// This needs to match the app environment that we deploy to gitlab
		// https://gitlab.com/jungleegames/backend/core/-/feature_flags
		unleash.WithAppName(AppEnv),
	)

	if err != nil {
		log.WithError(err).Errorf("failed to initiate unleash")
	}
}
