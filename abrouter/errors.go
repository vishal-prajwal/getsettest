package abrouter

import "errors"

var (
	// ErrNoChoices is returned when no choices are provided to the AB router.
	ErrNoChoices = errors.New("no choices provided")
	// ErrInvalidWeight is returned when a choice has a non-positive weight.
	ErrInvalidWeight = errors.New("invalid weight for choice, must be greater than zero")
)
