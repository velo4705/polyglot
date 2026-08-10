package language

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/velo4705/polyglot/pkg/types"
)

type ElmHandler struct{}

func (h *ElmHandler) Name() string {
	return "Elm"
}

func (h *ElmHandler) Extensions() []string {
	return []string{".elm"}
}

func (h *ElmHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *ElmHandler) NeedsCompilation() bool {
	return true
}

func (h *ElmHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, "make", source, "--output="+output)
	cmd := exec.Command("elm", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *ElmHandler) Run(file string, args []string) ([]byte, error) {
	output := filepath.Base(file)
	cmdArgs := append([]string{output}, args...)
	cmd := exec.Command("node", cmdArgs...)
	return cmd.CombinedOutput()
}
