package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type NimHandler struct{}

func (h *NimHandler) Name() string {
	return "Nim"
}

func (h *NimHandler) Extensions() []string {
	return []string{".nim"}
}

func (h *NimHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *NimHandler) NeedsCompilation() bool {
	return true
}

func (h *NimHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append([]string{"c"}, extraFlags...)
	cmdArgs = append(cmdArgs, "-o:"+output, source)
	cmd := exec.Command("nim", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *NimHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".nim")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
