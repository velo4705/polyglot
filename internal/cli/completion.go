package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const completionLong = `Generate shell completion scripts for Polyglot.

Bash:
  polyglot completion bash > /etc/bash_completion.d/polyglot
  # or for a single user:
  polyglot completion bash > ~/.bash_completion

Zsh:
  polyglot completion zsh > "${fpath[1]}/_polyglot"
  # Then start a new shell or run: autoload -U compinit && compinit

Fish:
  polyglot completion fish > ~/.config/fish/completions/polyglot.fish
`

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long:  completionLong,
	Args:  cobra.ExactArgs(1),
	RunE:  generateCompletion,
}

func generateCompletion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	default:
		return fmt.Errorf("unsupported shell %q: supported shells are bash, zsh, fish", args[0])
	}
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
