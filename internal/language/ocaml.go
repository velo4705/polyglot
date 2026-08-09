package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type OCamlHandler struct{}

func (h *OCamlHandler) Name() string {
	return "OCaml"
}

func (h *OCamlHandler) Extensions() []string {
	return []string{".ml", ".mli"}
}

func (h *OCamlHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *OCamlHandler) NeedsCompilation() bool {
	return true
}

func (h *OCamlHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("ocamlc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *OCamlHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(strings.TrimSuffix(file, ".ml"), ".mli")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
