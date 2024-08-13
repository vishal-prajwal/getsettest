package internal

import (
	"github.com/uptrace/bun"

	"bitbucket.org/junglee_games/getsetgo/pandora/consul"
)

type Application struct {
	// DB Connection
	DB     *bun.DB
	Consul consul.Client
}
