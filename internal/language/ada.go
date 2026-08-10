package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type AdaHandler struct{}

func (h *AdaHandler) Name() string {
	return "Ada"
}

func (h *AdaHandler) Extensions() []string {
	return []string{".adb", ".ads"}
}

func (h *AdaHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *AdaHandler) NeedsCompilation() bool {
	return true
}

func (h *AdaHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("gnatmake", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *AdaHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".adb")
	output = strings.TrimSuffix(output, ".ads")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
