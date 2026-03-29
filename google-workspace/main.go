package main

import (
	"os"

	"github.com/panthrocorp/openclaw-skills/google-workspace/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
