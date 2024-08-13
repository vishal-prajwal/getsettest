package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

func TestSetLevel(t *testing.T) {
	// Assert before calling setName
	assert.Nil(t, log.SetLevel("debug"))

	log.SetAppName("test")

	assert.Nil(t, log.SetLevel("info"))
}
