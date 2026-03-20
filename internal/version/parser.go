package version

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

// VersionSpec describes how to obtain and parse a version for one language.
type VersionSpec struct {
	Command string
	Args    []string
	// ParseFn extracts the version string from combined stdout+stderr output.
	ParseFn func(output string) string
}

// firstVersionMatch finds the first \d+\.\d+[\.\d]* match in output.
func firstVersionMatch(output string) string {
	re := regexp.MustCompile(`\d+\.\d+[\.\d]*`)
	return re.FindString(output)
}

// Specs is the registry of per-language version specs.
var Specs = map[string]VersionSpec{
	"Python": {
		Command: "python3",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// "Python 3.12.1" → "3.12.1"
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) == 0 {
				return ""
			}
			return parts[len(parts)-1]
		},
	},
	"JavaScript": {
		Command: "node",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// "v20.11.0" → "20.11.0"
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
	},
	"Go": {
		Command: "go",
		Args:    []string{"version"},
		ParseFn: func(output string) string {
			// "go version go1.22.0 linux/amd64" → "1.22.0"
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) < 3 {
				return ""
			}
			return strings.TrimPrefix(parts[2], "go")
		},
	},
	"Ruby": {
		Command: "ruby",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// "ruby 3.2.0 ..." → "3.2.0"
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		},
	},
	"Java": {
		Command: "java",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line: "java 21.0.1 ..." → "21.0.1"
			line := strings.SplitN(output, "\n", 2)[0]
			parts := strings.Fields(strings.TrimSpace(line))
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		},
	},
	"Rust": {
		Command: "rustc",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// "rustc 1.76.0 ..." → "1.76.0"
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		},
	},
	"PHP": {
		Command: "php",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line: "PHP 8.3.0 ..." → "8.3.0"
			line := strings.SplitN(output, "\n", 2)[0]
			parts := strings.Fields(strings.TrimSpace(line))
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		},
	},
	"Perl": {
		Command: "perl",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// "This is perl 5, version X" → extract \d+\.\d+[\.\d]*
			re := regexp.MustCompile(`\d+\.\d+[\.\d]*`)
			return re.FindString(output)
		},
	},
	"Lua": {
		Command: "lua",
		Args:    []string{"-v"},
		ParseFn: func(output string) string {
			// "Lua 5.4.6 ..." → "5.4.6"
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		},
	},
	"Zig": {
		Command: "zig",
		Args:    []string{"version"},
		ParseFn: func(output string) string {
			// output is just the version
			return strings.TrimSpace(output)
		},
	},
	"Nim": {
		Command: "nim",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line, extract first \d+\.\d+\.\d+ match
			line := strings.SplitN(output, "\n", 2)[0]
			return firstVersionMatch(line)
		},
	},
	"Crystal": {
		Command: "crystal",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line: "Crystal 1.11.0 ..." → "1.11.0"
			line := strings.SplitN(output, "\n", 2)[0]
			parts := strings.Fields(strings.TrimSpace(line))
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		},
	},
	"D": {
		Command: "dmd",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line, extract \d+\.\d+\.\d+ match
			line := strings.SplitN(output, "\n", 2)[0]
			return firstVersionMatch(line)
		},
	},
	"Shell": {
		Command: "bash",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line, extract \d+\.\d+\.\d+ match
			line := strings.SplitN(output, "\n", 2)[0]
			return firstVersionMatch(line)
		},
	},
	"Haskell": {
		Command: "ghc",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// last token of output
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) == 0 {
				return ""
			}
			return parts[len(parts)-1]
		},
	},
	"OCaml": {
		Command: "ocamlc",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// output is just the version
			return strings.TrimSpace(output)
		},
	},
	"Elixir": {
		Command: "elixir",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// line containing "Elixir", split on space, take index 1
			for _, line := range strings.Split(output, "\n") {
				if strings.Contains(line, "Elixir") {
					parts := strings.Fields(strings.TrimSpace(line))
					if len(parts) >= 2 {
						return parts[1]
					}
				}
			}
			return ""
		},
	},
	"Erlang": {
		Command: "erl",
		Args:    []string{"-eval", "erlang:display(erlang:system_info(otp_release)), halt().", "-noshell"},
		ParseFn: func(output string) string {
			// trim whitespace and quotes
			s := strings.TrimSpace(output)
			s = strings.Trim(s, `"`)
			return s
		},
	},
	"F#": {
		Command: "dotnet",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			return strings.TrimSpace(output)
		},
	},
	"Kotlin": {
		Command: "kotlin",
		Args:    []string{"-version"},
		ParseFn: func(output string) string {
			return firstVersionMatch(output)
		},
	},
	"Scala": {
		Command: "scala",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			return firstVersionMatch(output)
		},
	},
	"Groovy": {
		Command: "groovy",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			return firstVersionMatch(output)
		},
	},
	"Swift": {
		Command: "swift",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line, extract \d+\.\d+[\.\d]* match
			line := strings.SplitN(output, "\n", 2)[0]
			re := regexp.MustCompile(`\d+\.\d+[\.\d]*`)
			return re.FindString(line)
		},
	},
	"Dart": {
		Command: "dart",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			return firstVersionMatch(output)
		},
	},
	"TypeScript": {
		Command: "ts-node",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// trim leading "v"
			return strings.TrimPrefix(strings.TrimSpace(output), "v")
		},
	},
	"R": {
		Command: "Rscript",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			return firstVersionMatch(output)
		},
	},
	"Julia": {
		Command: "julia",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// split on space, take last token
			parts := strings.Fields(strings.TrimSpace(output))
			if len(parts) == 0 {
				return ""
			}
			return parts[len(parts)-1]
		},
	},
	"C": {
		Command: "gcc",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line, extract \d+\.\d+\.\d+ match
			line := strings.SplitN(output, "\n", 2)[0]
			return firstVersionMatch(line)
		},
	},
	"C++": {
		Command: "g++",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// first line, extract \d+\.\d+\.\d+ match
			line := strings.SplitN(output, "\n", 2)[0]
			return firstVersionMatch(line)
		},
	},
	"Brainfuck": {
		Command: "bf",
		Args:    []string{"--version"},
		ParseFn: func(output string) string {
			// best effort: trim whitespace
			return strings.TrimSpace(output)
		},
	},
}

// Get runs the version command for langName and returns the version string.
// Returns ("", false) if the language is not in Specs or the runtime is not installed.
// Returns ("unknown", true) if installed but output is unparseable or command exits non-zero.
// Returns (version, true) if the version is successfully extracted.
func Get(langName string) (version string, installed bool) {
	spec, ok := Specs[langName]
	if !ok {
		return "", false
	}

	out, err := exec.Command(spec.Command, spec.Args...).CombinedOutput()
	if err != nil {
		// Check if the executable was not found
		var execErr *exec.Error
		if errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound) {
			return "", false
		}
		// Also check for "executable file not found" in the error string
		if strings.Contains(err.Error(), "executable file not found") ||
			strings.Contains(err.Error(), "no such file or directory") {
			return "", false
		}
		// Command ran but exited non-zero — still installed
		v := spec.ParseFn(string(out))
		if v == "" {
			return "unknown", true
		}
		return v, true
	}

	v := spec.ParseFn(string(out))
	if v == "" {
		return "unknown", true
	}
	return v, true
}
