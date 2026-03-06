package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/detector"
	"github.com/velo4705/polyglot/internal/executor"
	"github.com/velo4705/polyglot/internal/language"
)

var compileCmd = &cobra.Command{
	Use:   "compile <file>",
	Short: "Compile a source file without running it",
	Args:  cobra.ExactArgs(1),
	RunE:  compileFile,
}

func compileFile(cmd *cobra.Command, cmdArgs []string) error {
	filename := cmdArgs[0]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Detect language
	lang, err := detector.DetectLanguage(filename)
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Printf("Detected: %s\n", lang.Name())
	}

	// Get language handler
	handler := language.GetHandler(lang)
	if handler == nil {
		return fmt.Errorf("no handler found for language: %s", lang)
	}

	// Compile
	exec := executor.New(verbose, quiet)
	return exec.Compile(handler, filename)
}
