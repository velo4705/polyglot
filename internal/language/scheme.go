package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type SchemeHandler struct{}

func (h *SchemeHandler) Name() string {
	return "Scheme"
}

func (h *SchemeHandler) Extensions() []string {
	return []string{".scm", ".ss"}
}

func (h *SchemeHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *SchemeHandler) NeedsCompilation() bool {
	return false
}

func (h *SchemeHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *SchemeHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("guile", cmdArgs...)
	return cmd.CombinedOutput()
}
