package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type GroovyHandler struct{}

func (h *GroovyHandler) Name() string {
	return "Groovy"
}

func (h *GroovyHandler) Extensions() []string {
	return []string{".groovy", ".gvy"}
}

func (h *GroovyHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *GroovyHandler) NeedsCompilation() bool {
	return false
}

func (h *GroovyHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *GroovyHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("groovy", cmdArgs...)
	return cmd.CombinedOutput()
}
