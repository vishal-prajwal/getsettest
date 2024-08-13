package database

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type IdentifiableContent interface {
	GetId() uuid.UUID
}

func GenerateSlug(tx bun.Tx, model IdentifiableContent, rawMaterial string) (string, error) {
	index := 1
	slugged := slug.Make(rawMaterial)

	for {
		var exists int

		exists, err := tx.NewSelect().
			Model(model).
			Where("slug = ? AND id != ?", slugged, model.GetId()).
			Count(context.Background())

		if err != nil {
			return "", err
		}

		if exists == 0 {
			return slugged, nil
		}

		// Generate a new slug
		slugged = slug.Make(rawMaterial + " " + strconv.Itoa(index))

		// Increment index
		index++

		if index == 1000 {
			return "", errors.Errorf("Slugify should not run for more than 1000 iterations")
		}
	}
}
