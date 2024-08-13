package uuidutils_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/junglee_games/getsetgo/pandora/uuidutils"
)

func TestUUIDFromStringOrNil(t *testing.T) {
	expected := uuid.New()

	u := uuidutils.UUIDFromStringOrNil(expected.String() + "\n  ")

	assert.Equal(t, expected.String(), u.String())
}
