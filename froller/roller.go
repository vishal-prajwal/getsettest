package froller

import "math/rand"

type FeatureRoller[T comparable] interface {
	// IsEnabled checks if a feature is enabled for a given user.
	// userId is optional and can be used to ensure consistent feature enablement for the same user.
	// If no userId is provided, a random selection will be made based on the feature's percentage.
	IsEnabled(feature T, userId ...int) bool
}

type Feature[T comparable] struct {
	Feature    T
	Percentage int
}

func NewFeature[T comparable](feature T, percentage int) Feature[T] {
	return Feature[T]{
		Feature:    feature,
		Percentage: percentage,
	}
}

type featureRollerImpl[T comparable] struct {
	features map[T]int
}

func NewFeatureRoller[T comparable](features []Feature[T]) (FeatureRoller[T], error) {
	// validations
	if len(features) == 0 {
		return nil, ErrNoFeatures
	}
	featureMap := make(map[T]int)
	for _, feature := range features {
		if feature.Percentage < 0 || feature.Percentage > 100 {
			return nil, ErrInvalidPercentage
		}
		featureMap[feature.Feature] = feature.Percentage
	}

	return &featureRollerImpl[T]{
		features: featureMap,
	}, nil
}

func (r *featureRollerImpl[T]) IsEnabled(feature T, userId ...int) bool {
	var selectionPoint int

	if len(userId) > 0 {
		// Use the first userId and mod it by 100 to get a consistent selection point
		selectionPoint = userId[0] % 100
	} else {
		// If no userId is provided, just get a random number between 0 and 99
		selectionPoint = rand.Intn(100)
	}

	if percentage, exists := r.features[feature]; exists {
		return selectionPoint < percentage
	}
	// If the feature does not exist, we consider it disabled
	return false
}
