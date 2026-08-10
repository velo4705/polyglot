package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type PascalHandler struct{}

func (h *PascalHandler) Name() string {
	return "Pascal"
}

func (h *PascalHandler) Extensions() []string {
	return []string{".pas", ".pp", ".lpr", ".dpr"}
}

func (h *PascalHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *PascalHandler) NeedsCompilation() bool {
	return true
}

func (h *PascalHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source)
	cmd := exec.Command("fpc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *PascalHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".pas")
	output = strings.TrimSuffix(output, ".pp")
	output = strings.TrimSuffix(output, ".lpr")
	output = strings.TrimSuffix(output, ".dpr")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
