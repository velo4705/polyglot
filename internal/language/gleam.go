package language

import (
	"fmt"
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type GleamHandler struct{}

func (h *GleamHandler) Name() string {
	return "Gleam"
}

func (h *GleamHandler) Extensions() []string {
	return []string{".gleam"}
}

func (h *GleamHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *GleamHandler) NeedsCompilation() bool {
	return false
}

func (h *GleamHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *GleamHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"run"}, args...)
	cmd := exec.Command("gleam", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return output_bytes, fmt.Errorf("execution failed: %s", string(output_bytes))
	}
	return output_bytes, nil
}
