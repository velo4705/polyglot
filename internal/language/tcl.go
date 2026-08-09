package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type TclHandler struct{}

func (h *TclHandler) Name() string {
	return "Tcl"
}

func (h *TclHandler) Extensions() []string {
	return []string{".tcl", ".tk"}
}

func (h *TclHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *TclHandler) NeedsCompilation() bool {
	return false
}

func (h *TclHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *TclHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("tclsh", cmdArgs...)
	return cmd.CombinedOutput()
}