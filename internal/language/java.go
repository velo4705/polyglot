package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type JavaHandler struct{}

func (h *JavaHandler) Name() string {
	return "Java"
}

func (h *JavaHandler) Extensions() []string {
	return []string{".java"}
}

func (h *JavaHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *JavaHandler) NeedsCompilation() bool {
	return true
}

func (h *JavaHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source)
	cmd := exec.Command("javac", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *JavaHandler) Run(file string, args []string) ([]byte, error) {
	// Extract class name from file (e.g., HelloWorld.java -> HelloWorld)
	className := strings.TrimSuffix(filepath.Base(file), ".java")
	dir := filepath.Dir(file)

	// Change to directory and run
	cmdArgs := append([]string{className}, args...)
	cmd := exec.Command("java", cmdArgs...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}
