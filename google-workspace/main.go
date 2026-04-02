package main

import (
	"os"
	"runtime/debug"

	"github.com/PanthroCorp-Limited/openclaw-skills/google-workspace/cmd"
)

func getVersion() string {
	// Try to read build info (works when built with `go build` or `go install`)
	if info, ok := debug.ReadBuildInfo(); ok {
		// Use the module version from build info
		if info.Main.Version != "" {
			return info.Main.Version
		}
	}
	// Fallback to dev if not available
	return "dev"
}

func main() {
	version := getVersion()
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
