package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type GasHandler struct{}

func (h *GasHandler) Name() string {
	return "GAS"
}

func (h *GasHandler) Extensions() []string {
	return []string{".s", ".S"}
}

func (h *GasHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *GasHandler) NeedsCompilation() bool {
	return true
}

func (h *GasHandler) Compile(source string, output string, extraFlags []string) error {
	// Use gcc as the frontend to handle linking
	args := append([]string{source, "-o", output, "-nostdlib"}, extraFlags...)
	cmd := exec.Command("gcc", args...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("GAS compilation failed: %s", string(outputBytes))
	}
	return nil
}

func (h *GasHandler) Run(file string, args []string) ([]byte, error) {
	ext := filepath.Ext(file)
	output := strings.TrimSuffix(file, ext)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
