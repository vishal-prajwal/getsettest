package jungleegames_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func TestGetAwsConfigLoadDefaultRegion(t *testing.T) {
	require.NoError(t, os.Setenv("AWS_REGION", ""))

	cfg, err := jungleegames.GetAwsConfig(context.Background())

	require.NoError(t, err)
	assert.Equal(t, cfg.Region, "eu-west-1")
}

func TestGetAwsConfigOverrideDefaultRegionFromEnv(t *testing.T) {
	require.NoError(t, os.Setenv("AWS_REGION", "eu-central-1"))

	cfg, err := jungleegames.GetAwsConfig(context.Background())

	require.NoError(t, err)
	assert.Equal(t, cfg.Region, "eu-central-1")
}
