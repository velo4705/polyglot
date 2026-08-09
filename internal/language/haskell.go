package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type HaskellHandler struct{}

func (h *HaskellHandler) Name() string {
	return "Haskell"
}

func (h *HaskellHandler) Extensions() []string {
	return []string{".hs"}
}

func (h *HaskellHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *HaskellHandler) NeedsCompilation() bool {
	return true
}

func (h *HaskellHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("ghc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *HaskellHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".hs")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
