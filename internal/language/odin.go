package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type OdinHandler struct{}

func (h *OdinHandler) Name() string {
	return "Odin"
}

func (h *OdinHandler) Extensions() []string {
	return []string{".odin"}
}

func (h *OdinHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *OdinHandler) NeedsCompilation() bool {
	return true
}

func (h *OdinHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, "build", source, "-out:"+output)
	cmd := exec.Command("odin", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *OdinHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".odin")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}