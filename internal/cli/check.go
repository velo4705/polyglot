package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/output"
	versionpkg "github.com/velo4705/polyglot/internal/version"
)

var checkJSON bool

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check which language toolchains are installed",
	RunE:  checkToolchains,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&checkJSON, "json", false, "Output results as JSON")
}

func checkToolchains(cmd *cobra.Command, args []string) error {
	toolchains := map[string]string{
		"Python":          "python3",
		"Go":              "go",
		"JavaScript":      "node",
		"Ruby":            "ruby",
		"PHP":             "php",
		"Java":            "javac",
		"C":               "gcc",
		"C++":             "g++",
		"Rust":            "rustc",
		"Perl":            "perl",
		"Lua":             "lua",
		"Zig":             "zig",
		"Nim":             "nim",
		"Crystal":         "crystal",
		"D":               "dmd",
		"Shell":           "bash",
		"Haskell":         "ghc",
		"OCaml":           "ocamlc",
		"Elixir":          "elixir",
		"Erlang":          "erl",
		"F#":              "dotnet",
		"Kotlin":          "kotlin",
		"Scala":           "scala",
		"Groovy":          "groovy",
		"Swift":           "swift",
		"Dart":            "dart",
		"TypeScript":      "ts-node",
		"R":               "Rscript",
		"Julia":           "julia",
		"Brainfuck":       "bf",
		"Fortran":         "gfortran",
		"Pascal":          "fpc",
		"Ada":             "gnatmake",
		"COBOL":           "cobc",
		"Scheme":          "guile",
		"Common Lisp":     "sbcl",
		"Forth":           "gforth",
		"Prolog":          "swipl",
		"Tcl":             "tclsh",
		"Clojure":         "clojure",
		"Gleam":           "gleam",
		"Elm":             "elm",
		"PureScript":      "spago",
		"Roc":             "roc",
		"V":               "v",
		"Odin":            "odin",
		"NASM":            "nasm",
		"GAS":             "gcc",
		"ARM Assembly":    "arm-linux-gnueabi-gcc",
		"MIPS Assembly":   "mips-linux-gnu-gcc",
		"RISC-V Assembly": "riscv64-linux-gnu-gcc",
	}

	type langResult struct {
		lang        string
		command     string
		isInstalled bool
		versionStr  string
	}

	results := make([]langResult, 0, len(toolchains))
	installed := 0

	for lang, command := range toolchains {
		_, err := exec.LookPath(command)
		isInstalled := err == nil
		var versionStr string
		if isInstalled {
			installed++
			v, ok := versionpkg.Get(lang)
			if !ok {
				versionStr = ""
			} else {
				versionStr = v
			}
		}
		results = append(results, langResult{
			lang:        lang,
			command:     command,
			isInstalled: isInstalled,
			versionStr:  versionStr,
		})
	}

	if checkJSON {
		entries := make([]output.CheckEntry, 0, len(results))
		for _, r := range results {
			entries = append(entries, output.CheckEntry{
				Language:  r.lang,
				Installed: r.isInstalled,
				Version:   r.versionStr,
			})
		}
		return output.PrintCheck(os.Stdout, output.CheckResult{Languages: entries})
	}

	// Plain-text output
	fmt.Println("Checking installed toolchains...")
	fmt.Println("================================")

	total := len(toolchains)
	for _, r := range results {
		var status string
		if !r.isInstalled {
			status = "✗ NOT FOUND"
		} else if r.versionStr == "" || r.versionStr == "unknown" || strings.Contains(r.versionStr, "Error") {
			status = "✓ (version unknown)"
		} else {
			status = "✓ " + r.versionStr
		}
		fmt.Printf("%-15s %-15s %s\n", r.lang, r.command, status)
	}

	fmt.Println("================================")
	fmt.Printf("Summary: %d/%d toolchains installed\n", installed, total)

	if installed < total {
		fmt.Println("\nTo use all languages, install the missing toolchains.")
	}

	return nil
}
