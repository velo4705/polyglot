package errors

import (
	stderrors "errors"
	"fmt"
	"os/exec"
	"strings"
)

// executableToLanguage maps runtime binary names to Polyglot language names.
var executableToLanguage = map[string]string{
	"python3": "python",
	"node":    "javascript",
	"ruby":    "ruby",
	"go":      "go",
	"javac":   "java",
	"java":    "java",
	"rustc":   "rust",
	"gcc":     "c",
	"g++":     "c++",
	"perl":    "perl",
	"lua":     "lua",
	"zig":     "zig",
	"nim":     "nim",
	"crystal": "crystal",
	"dmd":     "d",
	"bash":    "shell",
	"ghc":     "haskell",
	"ocamlc":  "ocaml",
	"elixir":  "elixir",
	"erl":     "erlang",
	"dotnet":  "fsharp",
	"kotlin":  "kotlin",
	"scala":   "scala",
	"groovy":  "groovy",
	"swift":   "swift",
	"dart":    "dart",
	"ts-node": "typescript",
	"Rscript": "r",
	"julia":   "julia",
	"bf":      "brainfuck",
}

// NotFound returns the actionable message for a missing runtime.
func NotFound(langName string) string {
	return fmt.Sprintf("%s is not installed. Run: polyglot install %s", langName, strings.ToLower(langName))
}

// FileNotFound returns the actionable message for a missing file.
func FileNotFound(path string) string {
	return fmt.Sprintf("File not found: %s. Check the path and try again.", path)
}

// UnknownExtension returns the actionable message for an unrecognised extension.
func UnknownExtension(ext string, supported []string) string {
	return fmt.Sprintf(
		"Unknown file extension %q. Supported extensions: %s. Use --lang to override detection.",
		ext,
		strings.Join(supported, ", "),
	)
}

// Generic wraps an arbitrary error with a context prefix.
func Generic(context string, err error) string {
	return fmt.Sprintf("%s: %s", context, err.Error())
}

// IsNotFoundError inspects an error and returns the executable name if it is
// an "executable not found" error (exec.ErrNotFound or "no such file").
func IsNotFoundError(err error) (executable string, ok bool) {
	if err == nil {
		return "", false
	}

	if stderrors.Is(err, exec.ErrNotFound) {
		executable = extractExecutable(err.Error())
		return executable, true
	}

	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "executable file not found") || strings.Contains(msg, "no such file or directory") {
		executable = extractExecutable(err.Error())
		return executable, true
	}

	return "", false
}

// extractExecutable tries to pull the executable name out of an exec error message.
// exec errors typically look like: exec: "python3": executable file not found in $PATH
func extractExecutable(msg string) string {
	// Look for the pattern: exec: "name":
	if idx := strings.Index(msg, `exec: "`); idx != -1 {
		rest := msg[idx+7:] // skip 'exec: "'
		if end := strings.Index(rest, `"`); end != -1 {
			return rest[:end]
		}
	}
	return ""
}

// ExecutableToLanguage looks up the canonical Polyglot language name for an executable.
func ExecutableToLanguage(exe string) (string, bool) {
	lang, ok := executableToLanguage[exe]
	return lang, ok
}
