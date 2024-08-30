package main

import (
	"github.com/joho/godotenv"

	"bitbucket.org/junglee_games/getsetgo/pandora/boilerplate/service_skeleton/sample/cmd"
)

func main() {
	// Load the local configuration first because gotdotenv will not overide existing environment variables
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")

	cmd.Execute()
}
