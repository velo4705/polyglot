package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type PythonHandler struct{}

func (h *PythonHandler) Name() string {
	return "Python"
}

func (h *PythonHandler) Extensions() []string {
	return []string{".py"}
}

func (h *PythonHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *PythonHandler) NeedsCompilation() bool {
	return false
}

func (h *PythonHandler) Compile(source string, output string, extraFlags []string) error {
	// Python doesn't need compilation
	return nil
}

func (h *PythonHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("python3", cmdArgs...)
	return cmd.CombinedOutput()
}
