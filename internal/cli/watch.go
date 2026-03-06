package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/config"
	"github.com/velo4705/polyglot/internal/detector"
	"github.com/velo4705/polyglot/internal/executor"
	"github.com/velo4705/polyglot/internal/language"
	"github.com/velo4705/polyglot/internal/ui"
	"github.com/velo4705/polyglot/internal/watcher"
	"github.com/velo4705/polyglot/pkg/types"
)

var (
	watchInterval time.Duration
	watchClear    bool
)

var watchCmd = &cobra.Command{
	Use:   "watch <file>",
	Short: "Watch a file and re-run on changes",
	Long: `Watch a source file and automatically re-run it when changes are detected.

This is useful during development to get immediate feedback.

Examples:
  polyglot watch hello.py
  polyglot watch --interval 500ms main.go
  polyglot watch --clear app.js`,
	Args: cobra.ExactArgs(1),
	RunE: watchFile,
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().DurationVar(&watchInterval, "interval", 1*time.Second, "Watch interval")
	watchCmd.Flags().BoolVar(&watchClear, "clear", false, "Clear screen before each run")
	watchCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	watchCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
	watchCmd.Flags().StringSliceVar(&args, "args", []string{}, "Arguments to pass to the program")
}

func watchFile(cmd *cobra.Command, cmdArgs []string) error {
	filename := cmdArgs[0]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Detect language
	lang, err := detector.DetectLanguageWithConfig(filename, cfg)
	if err != nil {
		return err
	}

	// Get language handler
	handler := language.GetHandler(lang)
	if handler == nil {
		return fmt.Errorf("no handler found for language: %s", lang)
	}

	ui.Header("Polyglot Watch Mode")
	ui.Info("File: %s", ui.File(filename))
	ui.Info("Language: %s", ui.Language(lang.Name()))
	fmt.Println()

	// Run once initially
	runWatchedFile(filename, handler, cfg)

	// Create watcher
	w := watcher.NewFileWatcher(filename, func() {
		if watchClear {
			clearScreen()
		}
		runWatchedFile(filename, handler, cfg)
	})

	if w == nil {
		return fmt.Errorf("failed to create file watcher")
	}

	// Start watching
	return w.Watch(watchInterval)
}

func runWatchedFile(filename string, handler types.LanguageHandler, cfg *config.Config) {
	startTime := time.Now()

	exec := executor.NewWithConfig(verbose, quiet, cfg)
	err := exec.Run(handler, filename, args)

	duration := time.Since(startTime)

	fmt.Println()
	if err != nil {
		ui.Error("Execution failed in %v", duration)
		if exitErr, ok := err.(*executor.ExecutionError); ok {
			ui.Dim("Exit code: %d", exitErr.ExitCode)
		}
	} else {
		ui.Success("Completed in %v", duration)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
