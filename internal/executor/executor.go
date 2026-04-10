package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/internal/config"
	"github.com/velo4705/polyglot/internal/ui"
	"github.com/velo4705/polyglot/pkg/types"
)

type ExecutionError struct {
	Message  string
	ExitCode int
}

func (e *ExecutionError) Error() string {
	return e.Message
}

type Executor struct {
	verbose bool
	quiet   bool
	config  *config.Config
}

func New(verbose, quiet bool) *Executor {
	return &Executor{
		verbose: verbose,
		quiet:   quiet,
	}
}

func NewWithConfig(verbose, quiet bool, cfg *config.Config) *Executor {
	return &Executor{
		verbose: verbose,
		quiet:   quiet,
		config:  cfg,
	}
}

func (e *Executor) Run(handler types.LanguageHandler, filename string, args []string) error {
	// Check if language is enabled
	if e.config != nil && !e.config.IsLanguageEnabled(string(handler.Name())) {
		return fmt.Errorf("language %s is disabled in configuration", handler.Name())
	}

	if handler.NeedsCompilation() {
		// Compile first
		output := e.getOutputName(filename)

		if !e.quiet {
			ui.Step("Compiling %s...", ui.File(filename))
		}

		spinner := ui.NewSpinner(fmt.Sprintf("Compiling %s", handler.Name()))
		if !e.quiet && !e.verbose {
			spinner.Start()
		}

		err := handler.Compile(filename, output)

		if !e.quiet && !e.verbose {
			spinner.Stop()
		}

		if err != nil {
			if !e.quiet {
				ui.Error("Compilation failed")
			}
			return err
		}

		if !e.quiet {
			ui.Success("Compilation successful")
		}

		// Clean up compiled binary after execution
		defer func() {
			os.Remove(output)
			// Also clean up Java .class files
			if handler.Name() == "Java" {
				dir := filepath.Dir(filename)
				className := strings.TrimSuffix(filepath.Base(filename), ".java")
				classFile := filepath.Join(dir, className+".class")
				os.Remove(classFile)
			}
		}()
	}

	if !e.quiet {
		cmdStr := e.getCommand(handler, filename)
		if len(args) > 0 {
			cmdStr += " " + strings.Join(args, " ")
		}
		ui.Step("Executing: %s", ui.Command(cmdStr))
	}

	err := e.runStreamed(handler, filename, args)

	if err != nil {
		// Check if it's an exit error to preserve exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &ExecutionError{
				Message:  "program exited with error",
				ExitCode: exitErr.ExitCode(),
			}
		}
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (e *Executor) Compile(handler types.LanguageHandler, filename string) error {
	if !handler.NeedsCompilation() {
		if !e.quiet {
			fmt.Println("No compilation needed")
		}
		return nil
	}

	output := e.getOutputName(filename)

	if !e.quiet {
		fmt.Printf("Compiling: %s %s -o %s\n", handler.Name(), filename, output)
	}

	if err := handler.Compile(filename, output); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	if !e.quiet {
		fmt.Printf("Compilation successful: %s\n", output)
	}

	return nil
}

// RunBuffered runs the program and captures output (used for --json mode).
func (e *Executor) RunBuffered(handler types.LanguageHandler, filename string, args []string) ([]byte, error) {
	if handler.NeedsCompilation() {
		outputBin := e.getOutputName(filename)
		spinner := ui.NewSpinner(fmt.Sprintf("Compiling %s", handler.Name()))
		if !e.quiet && !e.verbose {
			spinner.Start()
		}
		err := handler.Compile(filename, outputBin)
		if !e.quiet && !e.verbose {
			spinner.Stop()
		}
		if err != nil {
			return nil, err
		}
		defer func() {
			os.Remove(outputBin)
			if handler.Name() == "Java" {
				dir := filepath.Dir(filename)
				className := strings.TrimSuffix(filepath.Base(filename), ".java")
				os.Remove(filepath.Join(dir, className+".class"))
			}
		}()
	}
	return handler.Run(filename, args)
}

func (e *Executor) getOutputName(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// getAbsOutputName returns an absolute path for the compiled binary,
// ensuring it can be executed regardless of working directory.
func (e *Executor) getAbsOutputName(filename string) string {
	abs, err := filepath.Abs(e.getOutputName(filename))
	if err != nil {
		return e.getOutputName(filename)
	}
	return abs
}

// runStreamed executes the program with stdin/stdout/stderr wired to the terminal.
func (e *Executor) runStreamed(handler types.LanguageHandler, filename string, args []string) error {
	var name string
	var cmdArgs []string

	switch handler.Name() {
	case "C", "C++":
		binary := e.getAbsOutputName(filename)
		name = binary
		cmdArgs = args
	case "Java":
		className := strings.TrimSuffix(filepath.Base(filename), ".java")
		name = "java"
		cmdArgs = append([]string{className}, args...)
	case "Go":
		name = "go"
		cmdArgs = append([]string{"run", filename}, args...)
	case "Python":
		name = "python3"
		cmdArgs = append([]string{filename}, args...)
	case "JavaScript":
		name = "node"
		cmdArgs = append([]string{filename}, args...)
	case "Ruby":
		name = "ruby"
		cmdArgs = append([]string{filename}, args...)
	case "PHP":
		name = "php"
		cmdArgs = append([]string{filename}, args...)
	case "Perl":
		name = "perl"
		cmdArgs = append([]string{filename}, args...)
	case "Lua":
		name = "lua"
		cmdArgs = append([]string{filename}, args...)
	case "Shell":
		name = "bash"
		cmdArgs = append([]string{filename}, args...)
	case "TypeScript":
		name = "ts-node"
		cmdArgs = append([]string{filename}, args...)
	case "Rust":
		binary := e.getAbsOutputName(filename)
		name = binary
		cmdArgs = args
	default:
		// For other compiled languages, run the output binary
		if handler.NeedsCompilation() {
			name = e.getAbsOutputName(filename)
			cmdArgs = args
		} else {
			name = strings.ToLower(handler.Name())
			cmdArgs = append([]string{filename}, args...)
		}
	}

	cmd := exec.Command(name, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// For Java, run from the file's directory so the class is found
	if handler.Name() == "Java" {
		cmd.Dir = filepath.Dir(filename)
	}

	return cmd.Run()
}

func (e *Executor) getCommand(handler types.LanguageHandler, filename string) string {
	switch handler.Name() {
	case "Python":
		return "python3"
	case "Go":
		return "go run"
	case "JavaScript":
		return "node"
	default:
		return handler.Name()
	}
}
