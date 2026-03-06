package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/installer"
)

var installCmd = &cobra.Command{
	Use:   "install <language>",
	Short: "Install a language toolchain",
	Long: `Install a language toolchain system-wide using your system's package manager.

Examples:
  polyglot install python
  polyglot install rust
  polyglot install node`,
	Args: cobra.MinimumNArgs(1),
	RunE: installLanguage,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installLanguage(cmd *cobra.Command, args []string) error {
	inst := installer.New(false, false)

	// Map language names to commands
	langMap := map[string]struct {
		name    string
		command string
	}{
		"python":     {"Python", "python3"},
		"node":       {"Node.js", "node"},
		"javascript": {"Node.js", "node"},
		"js":         {"Node.js", "node"},
		"go":         {"Go", "go"},
		"rust":       {"Rust", "rustc"},
		"java":       {"Java", "javac"},
		"c":          {"C", "gcc"},
		"cpp":        {"C++", "g++"},
		"c++":        {"C++", "g++"},
		"ruby":       {"Ruby", "ruby"},
		"php":        {"PHP", "php"},
		"perl":       {"Perl", "perl"},
		"lua":        {"Lua", "lua"},
		"haskell":    {"Haskell", "ghc"},
		"elixir":     {"Elixir", "elixir"},
		"julia":      {"Julia", "julia"},
		"r":          {"R", "Rscript"},
		"dart":       {"Dart", "dart"},
		"kotlin":     {"Kotlin", "kotlin"},
		"scala":      {"Scala", "scala"},
		"groovy":     {"Groovy", "groovy"},
		"swift":      {"Swift", "swift"},
		"typescript": {"TypeScript", "ts-node"},
		"ts":         {"TypeScript", "ts-node"},
	}

	for _, lang := range args {
		if info, ok := langMap[lang]; ok {
			fmt.Printf("\n=== Installing %s ===\n", info.name)

			// Check if already installed
			if inst.IsAvailable(info.command) {
				fmt.Printf("✓ %s is already installed\n", info.name)
				continue
			}

			// Install
			if err := inst.Install(info.name, info.command); err != nil {
				fmt.Printf("❌ Failed to install %s: %v\n", info.name, err)
				continue
			}
		} else {
			fmt.Printf("❌ Unknown language: %s\n", lang)
			fmt.Println("   Run 'polyglot list' to see supported languages")
		}
	}

	return nil
}
