package jungleegames_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func TestPopulateConfig(t *testing.T) {
	type testConfig struct {
		AppName      string `env:"APP_NAME" ssm:"/app_name"`
		DefaultValue string `env:"NOT_SET" ssm:"/not_set" default:"helloWorld"`
	}

	// I have no idea how else to do this
	jungleegames.AppEnv = "testing"

	require.NoError(t, os.Setenv("APP_NAME", "jungleegames"))

	defer func() {
		require.NoError(t, os.Unsetenv("APP_NAME"))
	}()

	cfg := &testConfig{}

	assert.NoError(t, jungleegames.PopulateConfig(context.Background(), cfg))

	assert.Equal(t, "jungleegames", cfg.AppName)
	assert.Equal(t, "helloWorld", cfg.DefaultValue)
}

func TestPopulateConfigWithUnsupportedFieldType(t *testing.T) {
	cfg := struct {
		Foo int
	}{}

	err := jungleegames.PopulateConfig(context.Background(), &cfg)

	assert.Error(t, err, "field: Foo has a none-string type which is not supported")
}
