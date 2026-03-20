package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.1.0"
const milestone = "UX Polish - Shell Completions, Stdin Support, JSON Output"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Polyglot v%s\n", version)
		fmt.Printf("%s\n", milestone)
		fmt.Println("Universal compiler and interpreter wrapper")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
