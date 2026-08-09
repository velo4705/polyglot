package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type BrainfuckHandler struct{}

func (h *BrainfuckHandler) Name() string {
	return "Brainfuck"
}

func (h *BrainfuckHandler) Extensions() []string {
	return []string{".bf", ".b"}
}

func (h *BrainfuckHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *BrainfuckHandler) NeedsCompilation() bool {
	return false
}

func (h *BrainfuckHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *BrainfuckHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("bf", cmdArgs...)
	return cmd.CombinedOutput()
}
