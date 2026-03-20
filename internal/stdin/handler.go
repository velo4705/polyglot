package stdin

import (
	"fmt"
	"io"
	"os"
)

// canonicalExtensions maps language names to their primary file extension.
var canonicalExtensions = map[string]string{
	"Python":     ".py",
	"Go":         ".go",
	"JavaScript": ".js",
	"Java":       ".java",
	"C":          ".c",
	"C++":        ".cpp",
	"Rust":       ".rs",
	"Ruby":       ".rb",
	"PHP":        ".php",
	"Perl":       ".pl",
	"Lua":        ".lua",
	"Zig":        ".zig",
	"Nim":        ".nim",
	"Crystal":    ".cr",
	"D":          ".d",
	"Shell":      ".sh",
	"Haskell":    ".hs",
	"OCaml":      ".ml",
	"Elixir":     ".ex",
	"Erlang":     ".erl",
	"F#":         ".fs",
	"Kotlin":     ".kt",
	"Scala":      ".scala",
	"Groovy":     ".groovy",
	"Swift":      ".swift",
	"Dart":       ".dart",
	"TypeScript": ".ts",
	"R":          ".r",
	"Julia":      ".jl",
	"Brainfuck":  ".bf",
}

// CanonicalExtension returns the primary file extension for a language name.
// The second return value is false if the language is not recognised.
func CanonicalExtension(langName string) (string, bool) {
	ext, ok := canonicalExtensions[langName]
	return ext, ok
}

// ReadToTempFile reads all bytes from os.Stdin and writes them to a temporary
// file whose name ends with the canonical extension for langName.
// The caller is responsible for removing the file when done (e.g. defer os.Remove(path)).
func ReadToTempFile(langName string) (string, error) {
	ext, ok := CanonicalExtension(langName)
	if !ok {
		return "", fmt.Errorf("unsupported language: %q", langName)
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}

	f, err := os.CreateTemp("", "polyglot-stdin-*"+ext)
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}

	path := f.Name()

	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(path)
		return "", fmt.Errorf("writing temp file: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("closing temp file: %w", err)
	}

	return path, nil
}
