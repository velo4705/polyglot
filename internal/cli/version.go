package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Polyglot v%s\n", version.Version)
		fmt.Printf("%s\n", version.Milestone)
		fmt.Println("Universal compiler and interpreter wrapper")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
