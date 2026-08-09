package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type CHandler struct{}

func (h *CHandler) Name() string {
	return "C"
}

func (h *CHandler) Extensions() []string {
	return []string{".c"}
}

func (h *CHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *CHandler) NeedsCompilation() bool {
	return true
}

func (h *CHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("gcc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *CHandler) Run(file string, args []string) ([]byte, error) {
	// For C, we need to get the compiled binary name
	output := strings.TrimSuffix(file, ".c")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
