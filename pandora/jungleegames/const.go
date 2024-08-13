package jungleegames

import (
	"os"
	"strings"
)

var (
	AppEnv     = strings.ToLower(os.Getenv("SENTRY_ENV"))
	AppVersion = os.Getenv("SENTRY_RELEASE")
	AppDebug   = os.Getenv("DEBUG") == "true"
)
