package froller

import "errors"

var (
	ErrNoFeatures        = errors.New("no features provided")
	ErrInvalidPercentage = errors.New("invalid percentage for feature")
)
