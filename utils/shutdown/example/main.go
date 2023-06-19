package main

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/utils/shutdown"
)

func main() {

	shutdown.RunAndWait(
		func(ctx context.Context) error {
			return Task1()
		},
		func(ctx context.Context) error {
			return Task2()
		},
	)
}

func Task1() error {
	return nil
}

func Task2() error {
	return nil
}
