package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/config"
	"github.com/velo4705/polyglot/internal/detector"
	advisorpkg "github.com/velo4705/polyglot/internal/errors"
	"github.com/velo4705/polyglot/internal/executor"
	"github.com/velo4705/polyglot/internal/installer"
	"github.com/velo4705/polyglot/internal/language"
	"github.com/velo4705/polyglot/internal/output"
	stdinpkg "github.com/velo4705/polyglot/internal/stdin"
	"github.com/velo4705/polyglot/internal/ui"
	"github.com/velo4705/polyglot/pkg/types"
)

var (
	verbose    bool
	quiet      bool
	args       []string
	dryRun     bool
	jsonOutput bool
	lang       string
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Detect and run a source file",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  runFile,
}

func init() {
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	runCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (suppress polyglot messages)")
	runCmd.Flags().StringSliceVar(&args, "args", []string{}, "Arguments to pass to the program")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed without running")
	runCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	runCmd.Flags().StringVar(&lang, "lang", "", "Language name (required when reading from stdin)")
}

func runFile(cmd *cobra.Command, cmdArgs []string) error {
	var filename string

	// Detect stdin mode (no file argument provided)
	if len(cmdArgs) == 0 {
		if lang == "" {
			fmt.Fprintln(os.Stderr, "stdin requires --lang flag (e.g. --lang Python)")
			return fmt.Errorf("stdin requires --lang flag")
		}
		tempPath, err := stdinpkg.ReadToTempFile(lang)
		if err != nil {
			fmt.Fprintln(os.Stderr, advisorpkg.Generic("failed to read stdin", err))
			return err
		}
		defer os.Remove(tempPath)
		filename = tempPath
	} else {
		filename = cmdArgs[0]
	}

	// When --json is active, suppress normal ui output
	if jsonOutput {
		quiet = true
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		msg := advisorpkg.FileNotFound(filename)
		if jsonOutput {
			result := output.RunResult{
				Language: lang,
				File:     filename,
				ExitCode: 1,
				Stderr:   msg,
			}
			_ = output.PrintRun(os.Stdout, result)
			return nil
		}
		ui.Error("%s", msg)
		return fmt.Errorf("%s", msg)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		msg := advisorpkg.Generic("failed to load configuration", err)
		if jsonOutput {
			result := output.RunResult{
				Language: lang,
				File:     filename,
				ExitCode: 1,
				Stderr:   msg,
			}
			_ = output.PrintRun(os.Stdout, result)
			return nil
		}
		ui.Error("%s", msg)
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Detect language (with custom extension support)
	detectedLang, err := detector.DetectLanguageWithConfig(filename, cfg)
	if err != nil {
		msg := advisorpkg.Generic("failed to detect language", err)
		if jsonOutput {
			result := output.RunResult{
				Language: lang,
				File:     filename,
				ExitCode: 1,
				Stderr:   msg,
			}
			_ = output.PrintRun(os.Stdout, result)
			return nil
		}
		ui.Error("%s", msg)
		return err
	}

	if !quiet {
		ui.Info("Detected: %s", ui.Language(detectedLang.Name()))
	}

	// Check if language is enabled
	if !cfg.IsLanguageEnabled(detectedLang.Name()) {
		msg := fmt.Sprintf("language %s is disabled in configuration", detectedLang.Name())
		if jsonOutput {
			result := output.RunResult{
				Language: detectedLang.Name(),
				File:     filename,
				ExitCode: 1,
				Stderr:   msg,
			}
			_ = output.PrintRun(os.Stdout, result)
			return nil
		}
		ui.Error("Language %s is disabled in configuration", detectedLang.Name())
		return fmt.Errorf("%s", msg)
	}

	// Get language handler
	handler := language.GetHandler(detectedLang)
	if handler == nil {
		msg := fmt.Sprintf("no handler found for language: %s", detectedLang)
		if jsonOutput {
			result := output.RunResult{
				Language: detectedLang.Name(),
				File:     filename,
				ExitCode: 1,
				Stderr:   msg,
			}
			_ = output.PrintRun(os.Stdout, result)
			return nil
		}
		ui.Error("No handler found for language: %s", detectedLang)
		return fmt.Errorf("%s", msg)
	}

	// Dry run mode
	if dryRun {
		ui.Header("Dry Run Mode")
		ui.Info("File: %s", ui.File(filename))
		ui.Info("Language: %s", ui.Language(detectedLang.Name()))
		ui.Info("Handler: %s", handler.Name())

		if handler.NeedsCompilation() {
			ui.Step("Would compile: %s", ui.Command(getCompileCommand(handler, filename)))
		}

		cmdStr := getRunCommand(handler, filename, args)
		ui.Step("Would execute: %s", ui.Command(cmdStr))

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
		command := getCommandForLanguage(detectedLang.Name())
		if command != "" && !inst.IsAvailable(command) {
			if !quiet {
				ui.Info("Auto-installing %s...", detectedLang.Name())
			}
			if err := inst.Install(detectedLang.Name(), command); err != nil {
				ui.Error("Failed to install %s: %v", detectedLang.Name(), err)
				return fmt.Errorf("failed to install %s: %w", detectedLang.Name(), err)
			}
		}
	} else {
		// Prompt mode (default)
		inst := installer.New(false, quiet)
		command := getCommandForLanguage(detectedLang.Name())
		if command != "" && !inst.IsAvailable(command) {
			if err := inst.InstallIfMissing(detectedLang.Name(), command); err != nil {
				return err
			}
		}
	}

	// Execute with configuration
	exec := executor.NewWithConfig(verbose, quiet, cfg)

	if jsonOutput {
		startTime := time.Now()
		outBytes, runErr := exec.RunBuffered(handler, filename, args)
		durationMs := time.Since(startTime).Milliseconds()

		exitCode := 0
		stderrMsg := ""
		stdoutStr := string(outBytes)
		if runErr != nil {
			exitCode = 1
			stderrMsg = runErr.Error()
			if execErr, ok := runErr.(*executor.ExecutionError); ok {
				exitCode = execErr.ExitCode
			}
		}

		result := output.RunResult{
			Language:   detectedLang.Name(),
			File:       filename,
			ExitCode:   exitCode,
			Stdout:     stdoutStr,
			Stderr:     stderrMsg,
			DurationMs: durationMs,
		}
		return output.PrintRun(os.Stdout, result)
	}

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
	var cmdStr string
	switch handler.Name() {
	case "Python":
		cmdStr = fmt.Sprintf("python3 %s", filename)
	case "Go":
		cmdStr = fmt.Sprintf("go run %s", filename)
	case "JavaScript":
		cmdStr = fmt.Sprintf("node %s", filename)
	case "Ruby":
		cmdStr = fmt.Sprintf("ruby %s", filename)
	case "PHP":
		cmdStr = fmt.Sprintf("php %s", filename)
	case "Perl":
		cmdStr = fmt.Sprintf("perl %s", filename)
	case "Lua":
		cmdStr = fmt.Sprintf("lua %s", filename)
	case "Shell":
		cmdStr = fmt.Sprintf("bash %s", filename)
	default:
		if handler.NeedsCompilation() {
			cmdStr = fmt.Sprintf("./%s", getOutputName(filename))
		} else {
			cmdStr = fmt.Sprintf("%s %s", strings.ToLower(handler.Name()), filename)
		}
	}

	if len(args) > 0 {
		cmdStr += " " + strings.Join(args, " ")
	}

	return cmdStr
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
