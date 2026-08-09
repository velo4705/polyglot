package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type JuliaHandler struct{}

func (h *JuliaHandler) Name() string {
	return "Julia"
}

func (h *JuliaHandler) Extensions() []string {
	return []string{".jl"}
}

func (h *JuliaHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *JuliaHandler) NeedsCompilation() bool {
	return false
}

func (h *JuliaHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *JuliaHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("julia", cmdArgs...)
	return cmd.CombinedOutput()
}
