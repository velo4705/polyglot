package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type CommonLispHandler struct{}

func (h *CommonLispHandler) Name() string {
	return "Common Lisp"
}

func (h *CommonLispHandler) Extensions() []string {
	return []string{".lisp", ".lsp", ".cl", ".l"}
}

func (h *CommonLispHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *CommonLispHandler) NeedsCompilation() bool {
	return false
}

func (h *CommonLispHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *CommonLispHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"--script", file}, args...)
	cmd := exec.Command("sbcl", cmdArgs...)
	return cmd.CombinedOutput()
}