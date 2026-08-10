package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type PrologHandler struct{}

func (h *PrologHandler) Name() string {
	return "Prolog"
}

func (h *PrologHandler) Extensions() []string {
	return []string{".pl", ".pro", ".prolog"}
}

func (h *PrologHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *PrologHandler) NeedsCompilation() bool {
	return false
}

func (h *PrologHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *PrologHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"-s", file, "--"}, args...)
	cmd := exec.Command("swipl", cmdArgs...)
	return cmd.CombinedOutput()
}
