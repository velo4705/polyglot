package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type MipsAssemblyHandler struct{}

func (h *MipsAssemblyHandler) Name() string {
	return "MIPS Assembly"
}

func (h *MipsAssemblyHandler) Extensions() []string {
	return []string{".mips"}
}

func (h *MipsAssemblyHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *MipsAssemblyHandler) NeedsCompilation() bool {
	return true
}

func (h *MipsAssemblyHandler) Compile(source string, output string, extraFlags []string) error {
	args := append([]string{source, "-o", output, "-nostdlib"}, extraFlags...)
	cmd := exec.Command("mips-linux-gnu-gcc", args...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("MIPS assembly failed: %s", string(outputBytes))
	}
	return nil
}

func (h *MipsAssemblyHandler) Run(file string, args []string) ([]byte, error) {
	ext := filepath.Ext(file)
	output := strings.TrimSuffix(file, ext)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
