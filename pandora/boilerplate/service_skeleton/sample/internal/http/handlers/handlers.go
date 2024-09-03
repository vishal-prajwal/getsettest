package handlers

import (
	"bitbucket.org/junglee_games/getsetgo/pandora/boilerplate/service_skeleton/sample/internal"
)

func NewSessionHandler(app *internal.Application) *SessionHandler {
	return &SessionHandler{
		app: app,
	}
}

type SessionHandler struct {
	app (*internal.Application)
}
