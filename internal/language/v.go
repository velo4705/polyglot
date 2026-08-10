package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type VHandler struct{}

func (h *VHandler) Name() string {
	return "V"
}

func (h *VHandler) Extensions() []string {
	return []string{".v"}
}

func (h *VHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *VHandler) NeedsCompilation() bool {
	return true
}

func (h *VHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, "-o", output, source)
	cmd := exec.Command("v", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *VHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".v")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
