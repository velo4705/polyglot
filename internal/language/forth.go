package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type ForthHandler struct{}

func (h *ForthHandler) Name() string {
	return "Forth"
}

func (h *ForthHandler) Extensions() []string {
	return []string{".fth", ".forth"}
}

func (h *ForthHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *ForthHandler) NeedsCompilation() bool {
	return false
}

func (h *ForthHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *ForthHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("gforth", cmdArgs...)
	return cmd.CombinedOutput()
}