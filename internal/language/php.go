package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type PHPHandler struct{}

func (h *PHPHandler) Name() string {
	return "PHP"
}

func (h *PHPHandler) Extensions() []string {
	return []string{".php"}
}

func (h *PHPHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *PHPHandler) NeedsCompilation() bool {
	return false
}

func (h *PHPHandler) Compile(source string, output string, extraFlags []string) error {
	// PHP doesn't need compilation
	return nil
}

func (h *PHPHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("php", cmdArgs...)
	return cmd.CombinedOutput()
}
