package cli

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/ui"
	"github.com/velo4705/polyglot/internal/updater"
)

var (
	checkOnly bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Polyglot to the latest version",
	Long: `Check for and install the latest version of Polyglot.

This command will:
  - Check GitHub for the latest release
  - Download the appropriate binary for your platform
  - Backup the current version
  - Install the new version

If installed via a package manager (Homebrew, APT), it will use that instead.`,
	Example: `  # Check for updates and install if available
  polyglot update

  # Only check for updates without installing
  polyglot update --check`,
	Run: func(cmd *cobra.Command, args []string) {
		u := updater.New(version, false)

		// Check for updates
		release, hasUpdate, err := u.CheckForUpdates()
		if err != nil {
			ui.Error("Failed to check for updates: %v", err)
			os.Exit(1)
		}

		if !hasUpdate {
			// Already on latest version
			return
		}

		// If check-only mode, just report and exit
		if checkOnly {
			ui.Info("New version available: %s", release.TagName)
			ui.Dim("Run 'polyglot update' to upgrade")
			return
		}

		// Attempt package manager update first
		if err := u.UpdateViaPackageManager(); err == nil {
			// Successfully updated via package manager
			return
		}

		// Fall back to direct binary update
		if err := u.Update(release); err != nil {
			if strings.Contains(err.Error(), "no binary found") {
				ui.Error("No pre-built binary found for your platform (%s)", err)
				ui.Dim("You can build from source: https://github.com/velo4705/polyglot")
			} else {
				ui.Error("Update failed: %v", err)
			}
			os.Exit(1)
		}
	},
}

func init() {
	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates without installing")
	rootCmd.AddCommand(updateCmd)
}
