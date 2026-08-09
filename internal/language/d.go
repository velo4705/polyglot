package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type DHandler struct{}

func (h *DHandler) Name() string {
	return "D"
}

func (h *DHandler) Extensions() []string {
	return []string{".d"}
}

func (h *DHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *DHandler) NeedsCompilation() bool {
	return true
}

func (h *DHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-of"+output)
	cmd := exec.Command("dmd", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *DHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".d")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
