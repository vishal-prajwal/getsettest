package jungleegames

import (
	"github.com/google/uuid"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

func Uuid(id uuid.UUID) *uuid.UUID {
	return &id
}

func UuidPtrToString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}

	return id.String()
}

func UuidPtrFromString(id string) *uuid.UUID {
	if id == "" {
		return nil
	}

	return Uuid(uuid.MustParse(id))
}

func UuidSliceToStringSlice(uuids []uuid.UUID) []string {
	var result []string

	for _, id := range uuids {
		result = append(result, id.String())
	}

	return result
}

func StringSliceToUUID(slice []string) []uuid.UUID {
	var items []uuid.UUID

	for _, id := range slice {
		if parsedId, err := uuid.Parse(id); err == nil {
			items = append(items, parsedId)
		} else {
			log.WithError(err).Errorf("error parsing string id %s to uuid ", id)
		}
	}

	return items
}

func RemoveIdsFromUUICSlice(slice []uuid.UUID, removeIds []uuid.UUID) []uuid.UUID {
	var result []uuid.UUID

	for _, id := range slice {
		remove := false

		for _, removeId := range removeIds {
			if id == removeId {
				remove = true
				break
			}
		}

		if !remove {
			result = append(result, id)
		}
	}

	return result
}
