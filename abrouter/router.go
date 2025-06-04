package abrouter

import "math/rand"

type ABRouter[T any] interface {
	Choose(userId ...int) T
}

type Choice[T any] struct {
	Choice T
	Weight int
}

func NewChoice[T any](choice T, weight int) Choice[T] {
	return Choice[T]{
		Choice: choice,
		Weight: weight,
	}
}

type abRouterImpl[T any] struct {
	choices     []Choice[T]
	totalWeight int
}

func NewABRouter[T any](choices []Choice[T]) (ABRouter[T], error) {
	if len(choices) == 0 {
		return nil, ErrNoChoices
	}

	totalWeight := 0
	for _, choice := range choices {
		if choice.Weight < 0 {
			return nil, ErrInvalidWeight
		}
		totalWeight += choice.Weight
	}

	return &abRouterImpl[T]{
		choices:     choices,
		totalWeight: totalWeight,
	}, nil
}

func (r *abRouterImpl[T]) Choose(userId ...int) T {
	var selectionPoint int

	if len(userId) > 0 {
		// Use the first userId and mod it by the total weight to get a consistent selection point
		selectionPoint = userId[0] % r.totalWeight
	} else {
		// If no userId is provided, just get a random number
		selectionPoint = rand.Intn(r.totalWeight)
	}

	for _, choice := range r.choices {
		if selectionPoint < choice.Weight {
			return choice.Choice
		}
		selectionPoint -= choice.Weight
	}
	// fallback to the last choice if no other choice matched
	return r.choices[len(r.choices)-1].Choice
}
