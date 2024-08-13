package database

import (
	"os"
)

func getWithDefault(name, defaultValue string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}

	return defaultValue
}
