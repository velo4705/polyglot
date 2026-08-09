package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type PerlHandler struct{}

func (h *PerlHandler) Name() string {
	return "Perl"
}

func (h *PerlHandler) Extensions() []string {
	return []string{".pl"}
}

func (h *PerlHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *PerlHandler) NeedsCompilation() bool {
	return false
}

func (h *PerlHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *PerlHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("perl", cmdArgs...)
	return cmd.CombinedOutput()
}
