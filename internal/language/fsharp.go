package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type FSharpHandler struct{}

func (h *FSharpHandler) Name() string {
	return "F#"
}

func (h *FSharpHandler) Extensions() []string {
	return []string{".fs", ".fsx"}
}

func (h *FSharpHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *FSharpHandler) NeedsCompilation() bool {
	// F# scripts (.fsx) can be run directly
	return false
}

func (h *FSharpHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("fsharpc", cmdArgs...)
	return cmd.Run()
}

func (h *FSharpHandler) Run(file string, args []string) ([]byte, error) {
	// Use dotnet fsi for F# scripts
	cmdArgs := append([]string{"fsi", file}, args...)
	cmd := exec.Command("dotnet", cmdArgs...)
	return cmd.CombinedOutput()
}
