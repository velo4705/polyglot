package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "polyglot",
	Short:         "Universal compiler and interpreter wrapper",
	Long:          `Polyglot detects programming languages and executes them with the appropriate toolchain.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(listCmd)
	// Note: updateCmd, versionCmd, checkCmd, installCmd, configCmd, watchCmd
	// are registered in their respective init() functions
}
