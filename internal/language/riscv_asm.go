package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type RiscVAssemblyHandler struct{}

func (h *RiscVAssemblyHandler) Name() string {
	return "RISC-V Assembly"
}

func (h *RiscVAssemblyHandler) Extensions() []string {
	return []string{".riscv"}
}

func (h *RiscVAssemblyHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *RiscVAssemblyHandler) NeedsCompilation() bool {
	return true
}

func (h *RiscVAssemblyHandler) Compile(source string, output string, extraFlags []string) error {
	args := append([]string{source, "-o", output, "-nostdlib"}, extraFlags...)
	cmd := exec.Command("riscv64-linux-gnu-gcc", args...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("RISC-V assembly failed: %s", string(outputBytes))
	}
	return nil
}

func (h *RiscVAssemblyHandler) Run(file string, args []string) ([]byte, error) {
	ext := filepath.Ext(file)
	output := strings.TrimSuffix(file, ext)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
