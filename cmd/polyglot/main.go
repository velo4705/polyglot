package main

import (
	"fmt"
	"os"

	"github.com/velo4705/polyglot/internal/cli"
	"github.com/velo4705/polyglot/internal/updater"
	"github.com/velo4705/polyglot/pkg/version"
)

func main() {
	// Check for updates in background (non-blocking)
	u := updater.New(version.Version, true)
	u.CheckForUpdatesInBackground()

	// Execute CLI
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
