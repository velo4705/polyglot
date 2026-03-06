package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/config"
	"github.com/velo4705/polyglot/internal/detector"
	"github.com/velo4705/polyglot/internal/executor"
	"github.com/velo4705/polyglot/internal/installer"
	"github.com/velo4705/polyglot/internal/language"
	"github.com/velo4705/polyglot/internal/ui"
	"github.com/velo4705/polyglot/pkg/types"
)

var (
	verbose bool
	quiet   bool
	args    []string
	dryRun  bool
)

var runCmd = &cobra.Command{
	Use:   "run <file>",
	Short: "Detect and run a source file",
	Args:  cobra.ExactArgs(1),
	RunE:  runFile,
}

func init() {
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	runCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (suppress polyglot messages)")
	runCmd.Flags().StringSliceVar(&args, "args", []string{}, "Arguments to pass to the program")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed without running")
}

func runFile(cmd *cobra.Command, cmdArgs []string) error {
	filename := cmdArgs[0]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		ui.Error("File not found: %s", filename)
		return fmt.Errorf("file not found: %s", filename)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		ui.Error("Failed to load configuration: %v", err)
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Detect language (with custom extension support)
	lang, err := detector.DetectLanguageWithConfig(filename, cfg)
	if err != nil {
		ui.Error("Failed to detect language: %v", err)
		return err
	}

	if !quiet {
		ui.Info("Detected: %s", ui.Language(lang.Name()))
	}

	// Check if language is enabled
	if !cfg.IsLanguageEnabled(lang.Name()) {
		ui.Error("Language %s is disabled in configuration", lang.Name())
		return fmt.Errorf("language %s is disabled in configuration", lang.Name())
	}

	// Get language handler
	handler := language.GetHandler(lang)
	if handler == nil {
		ui.Error("No handler found for language: %s", lang)
		return fmt.Errorf("no handler found for language: %s", lang)
	}

	// Dry run mode
	if dryRun {
		ui.Header("Dry Run Mode")
		ui.Info("File: %s", ui.File(filename))
		ui.Info("Language: %s", ui.Language(lang.Name()))
		ui.Info("Handler: %s", handler.Name())

		if handler.NeedsCompilation() {
			ui.Step("Would compile: %s", ui.Command(getCompileCommand(handler, filename)))
		}

		runCmd := getRunCommand(handler, filename, args)
		ui.Step("Would execute: %s", ui.Command(runCmd))

		if len(args) > 0 {
			ui.Info("Arguments: %s", strings.Join(args, ", "))
		}

		ui.Success("Dry run complete (no actual execution)")
		return nil
	}

	// Check if toolchain is available, install if needed
	if cfg.AutoInstall.Enabled {
		// Auto-install mode
		inst := installer.New(true, quiet)
		command := getCommandForLanguage(lang.Name())
		if command != "" && !inst.IsAvailable(command) {
			if !quiet {
				ui.Info("Auto-installing %s...", lang.Name())
			}
			if err := inst.Install(lang.Name(), command); err != nil {
				ui.Error("Failed to install %s: %v", lang.Name(), err)
				return fmt.Errorf("failed to install %s: %w", lang.Name(), err)
			}
		}
	} else {
		// Prompt mode (default)
		inst := installer.New(false, quiet)
		command := getCommandForLanguage(lang.Name())
		if command != "" && !inst.IsAvailable(command) {
			if err := inst.InstallIfMissing(lang.Name(), command); err != nil {
				return err
			}
		}
	}

	// Execute with configuration
	exec := executor.NewWithConfig(verbose, quiet, cfg)
	return exec.Run(handler, filename, args)
}

// getCommandForLanguage returns the command name for a language
func getCommandForLanguage(language string) string {
	commandMap := map[string]string{
		"Python":     "python3",
		"JavaScript": "node",
		"Go":         "go",
		"Ruby":       "ruby",
		"PHP":        "php",
		"Perl":       "perl",
		"Lua":        "lua",
		"Java":       "javac",
		"C":          "gcc",
		"C++":        "g++",
		"Rust":       "rustc",
		"Zig":        "zig",
		"Nim":        "nim",
		"Crystal":    "crystal",
		"D":          "dmd",
		"Shell":      "bash",
		"Haskell":    "ghc",
		"OCaml":      "ocamlc",
		"Elixir":     "elixir",
		"Erlang":     "erl",
		"F#":         "dotnet",
		"Kotlin":     "kotlin",
		"Scala":      "scala",
		"Groovy":     "groovy",
		"Swift":      "swift",
		"Dart":       "dart",
		"TypeScript": "ts-node",
		"R":          "Rscript",
		"Julia":      "julia",
	}
	return commandMap[language]
}

// getCompileCommand returns the compile command for a language
func getCompileCommand(handler types.LanguageHandler, filename string) string {
	switch handler.Name() {
	case "Java":
		return fmt.Sprintf("javac %s", filename)
	case "C":
		return fmt.Sprintf("gcc %s -o %s", filename, getOutputName(filename))
	case "C++":
		return fmt.Sprintf("g++ %s -o %s", filename, getOutputName(filename))
	case "Rust":
		return fmt.Sprintf("rustc %s -o %s", filename, getOutputName(filename))
	case "Zig":
		return fmt.Sprintf("zig build-exe %s", filename)
	case "Nim":
		return fmt.Sprintf("nim c %s", filename)
	case "Crystal":
		return fmt.Sprintf("crystal build %s", filename)
	case "D":
		return fmt.Sprintf("dmd %s", filename)
	default:
		return "unknown"
	}
}

// getRunCommand returns the run command for a language
func getRunCommand(handler types.LanguageHandler, filename string, args []string) string {
	var cmd string
	switch handler.Name() {
	case "Python":
		cmd = fmt.Sprintf("python3 %s", filename)
	case "Go":
		cmd = fmt.Sprintf("go run %s", filename)
	case "JavaScript":
		cmd = fmt.Sprintf("node %s", filename)
	case "Ruby":
		cmd = fmt.Sprintf("ruby %s", filename)
	case "PHP":
		cmd = fmt.Sprintf("php %s", filename)
	case "Perl":
		cmd = fmt.Sprintf("perl %s", filename)
	case "Lua":
		cmd = fmt.Sprintf("lua %s", filename)
	case "Shell":
		cmd = fmt.Sprintf("bash %s", filename)
	default:
		if handler.NeedsCompilation() {
			cmd = fmt.Sprintf("./%s", getOutputName(filename))
		} else {
			cmd = fmt.Sprintf("%s %s", strings.ToLower(handler.Name()), filename)
		}
	}

	if len(args) > 0 {
		cmd += " " + strings.Join(args, " ")
	}

	return cmd
}

// getOutputName returns the output binary name for a source file
func getOutputName(filename string) string {
	ext := ""
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i:]
			break
		}
	}
	return strings.TrimSuffix(filename, ext)
}
