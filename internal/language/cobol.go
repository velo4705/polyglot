package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type COBOLHandler struct{}

func (h *COBOLHandler) Name() string {
	return "COBOL"
}

func (h *COBOLHandler) Extensions() []string {
	return []string{".cob", ".cbl", ".COB", ".CBL"}
}

func (h *COBOLHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *COBOLHandler) NeedsCompilation() bool {
	return true
}

func (h *COBOLHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, "-x", source, "-o", output)
	cmd := exec.Command("cobc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *COBOLHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".cob")
	output = strings.TrimSuffix(output, ".cbl")
	output = strings.TrimSuffix(output, ".COB")
	output = strings.TrimSuffix(output, ".CBL")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
