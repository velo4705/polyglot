package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/velo4705/polyglot/internal/config"
	"github.com/velo4705/polyglot/internal/ui"
	"github.com/velo4705/polyglot/pkg/types"
)

type ExecutionError struct {
	Message  string
	ExitCode int
	Stderr   string
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

	stderrStr, err := e.runStreamed(handler, filename, args)

	if err != nil {
		// Check if it's an exit error to preserve exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &ExecutionError{
				Message:  "program exited with error",
				ExitCode: exitErr.ExitCode(),
				Stderr:   stderrStr,
			}
		}
		return &ExecutionError{
			Message:  err.Error(),
			ExitCode: 1,
			Stderr:   stderrStr,
		}
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

	sandboxEnabled := e.config != nil && e.config.Sandbox.Enabled
	sandboxTimeout := 0
	if sandboxEnabled {
		sandboxTimeout = e.config.Sandbox.Timeout
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if sandboxEnabled && sandboxTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(sandboxTimeout)*time.Second)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	cmd := e.buildCmd(ctx, handler, filename, args)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("execution timed out (limit: %ds)", sandboxTimeout)
	}
	return output, err
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
func (e *Executor) runStreamed(handler types.LanguageHandler, filename string, args []string) (string, error) {
	sandboxEnabled := e.config != nil && e.config.Sandbox.Enabled
	sandboxTimeout := 0
	if sandboxEnabled {
		sandboxTimeout = e.config.Sandbox.Timeout
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if sandboxEnabled && sandboxTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(sandboxTimeout)*time.Second)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	cmd := e.buildCmd(ctx, handler, filename, args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stderrBuf.String(), fmt.Errorf("execution timed out (limit: %ds)", sandboxTimeout)
	}
	return stderrBuf.String(), err
}

// buildCmd builds the command to execute, applying sandboxing resource limits (memory and CPU) via prlimit if enabled.
func (e *Executor) buildCmd(ctx context.Context, handler types.LanguageHandler, filename string, args []string) *exec.Cmd {
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

	sandboxEnabled := false
	sandboxMemory := int64(0)
	sandboxCPU := 0

	if e.config != nil {
		sandboxEnabled = e.config.Sandbox.Enabled
		sandboxMemory = e.config.Sandbox.MemoryLimitMB
		sandboxCPU = e.config.Sandbox.CPULimit
	}

	hasPrlimit := false
	if _, err := exec.LookPath("prlimit"); err == nil {
		hasPrlimit = true
	}

	if sandboxEnabled {
		prlimitArgs := []string{}
		if sandboxMemory > 0 {
			// Limit virtual memory size (address space) in bytes
			prlimitArgs = append(prlimitArgs, fmt.Sprintf("--as=%d", sandboxMemory*1024*1024))
		}
		if sandboxCPU > 0 {
			// Limit CPU time in seconds
			prlimitArgs = append(prlimitArgs, fmt.Sprintf("--cpu=%d", sandboxCPU))
		}

		if len(prlimitArgs) > 0 && hasPrlimit {
			prlimitArgs = append(prlimitArgs, name)
			prlimitArgs = append(prlimitArgs, cmdArgs...)
			name = "prlimit"
			cmdArgs = prlimitArgs
		} else if len(prlimitArgs) > 0 && !hasPrlimit && !e.quiet {
			ui.Warning("prlimit not found. CPU and memory limits will not be enforced.")
		}
	}

	cmd := exec.CommandContext(ctx, name, cmdArgs...)

	// For Java, run from the file's directory so the class is found
	if handler.Name() == "Java" {
		cmd.Dir = filepath.Dir(filename)
	}

	return cmd
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
