package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type ShellHandler struct{}

func (h *ShellHandler) Name() string {
	return "Shell"
}

func (h *ShellHandler) Extensions() []string {
	return []string{".sh", ".bash"}
}

func (h *ShellHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *ShellHandler) NeedsCompilation() bool {
	return false
}

func (h *ShellHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *ShellHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("bash", cmdArgs...)
	return cmd.CombinedOutput()
}
