package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type CppHandler struct{}

func (h *CppHandler) Name() string {
	return "C++"
}

func (h *CppHandler) Extensions() []string {
	return []string{".cpp", ".cc", ".cxx", ".C"}
}

func (h *CppHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *CppHandler) NeedsCompilation() bool {
	return true
}

func (h *CppHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("g++", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *CppHandler) Run(file string, args []string) ([]byte, error) {
	// Get output name without extension
	ext := filepath.Ext(file)
	output := strings.TrimSuffix(file, ext)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
