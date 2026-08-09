package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type RustHandler struct{}

func (h *RustHandler) Name() string {
	return "Rust"
}

func (h *RustHandler) Extensions() []string {
	return []string{".rs"}
}

func (h *RustHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *RustHandler) NeedsCompilation() bool {
	return true
}

func (h *RustHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("rustc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *RustHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".rs")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
