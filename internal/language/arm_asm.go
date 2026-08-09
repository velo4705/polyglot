package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type ArmAssemblyHandler struct{}

func (h *ArmAssemblyHandler) Name() string {
	return "ARM Assembly"
}

func (h *ArmAssemblyHandler) Extensions() []string {
	return []string{".arm"}
}

func (h *ArmAssemblyHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *ArmAssemblyHandler) NeedsCompilation() bool {
	return true
}

func (h *ArmAssemblyHandler) Compile(source string, output string, extraFlags []string) error {
	args := append([]string{source, "-o", output, "-nostdlib"}, extraFlags...)
	cmd := exec.Command("arm-linux-gnueabi-gcc", args...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ARM assembly failed: %s", string(outputBytes))
	}
	return nil
}

func (h *ArmAssemblyHandler) Run(file string, args []string) ([]byte, error) {
	ext := filepath.Ext(file)
	output := strings.TrimSuffix(file, ext)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
