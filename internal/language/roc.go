package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type RocHandler struct{}

func (h *RocHandler) Name() string {
	return "Roc"
}

func (h *RocHandler) Extensions() []string {
	return []string{".roc"}
}

func (h *RocHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *RocHandler) NeedsCompilation() bool {
	return false
}

func (h *RocHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *RocHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"run", file}, args...)
	cmd := exec.Command("roc", cmdArgs...)
	return cmd.CombinedOutput()
}