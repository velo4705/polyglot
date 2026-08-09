package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type RHandler struct{}

func (h *RHandler) Name() string {
	return "R"
}

func (h *RHandler) Extensions() []string {
	return []string{".r", ".R"}
}

func (h *RHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *RHandler) NeedsCompilation() bool {
	return false
}

func (h *RHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *RHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"--vanilla", "--slave", "-f", file}, args...)
	cmd := exec.Command("Rscript", cmdArgs...)
	return cmd.CombinedOutput()
}
