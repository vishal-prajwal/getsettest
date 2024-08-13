package jungleegames_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func TestUuidStringSliceToAndFrom(t *testing.T) {
	data := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	stringSlice := jungleegames.UuidSliceToStringSlice(data)
	assert.Len(t, stringSlice, len(data))

	uuidSlice := jungleegames.StringSliceToUUID(stringSlice)
	assert.Len(t, uuidSlice, len(data))

	assert.Equal(t, data[0], uuidSlice[0])
	assert.Equal(t, data[1], uuidSlice[1])
	assert.Equal(t, data[2], uuidSlice[2])
	assert.Equal(t, data[3], uuidSlice[3])
}

func TestRemoveIdsFromUUICSlice(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()
	uuid4 := uuid.New()
	uuid5 := uuid.New()
	uuid6 := uuid.New()
	data := []uuid.UUID{
		uuid1,
		uuid2,
		uuid3,
		uuid4,
		uuid5,
		uuid6,
	}

	removeUUIDs := []uuid.UUID{
		uuid2,
		uuid.New(),
	}

	result := jungleegames.RemoveIdsFromUUICSlice(data, removeUUIDs)
	assert.Len(t, result, 5)
	assert.Equal(t, data[0], result[0])
	assert.Equal(t, data[5], result[4])
}
