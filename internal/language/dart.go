package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type DartHandler struct{}

func (h *DartHandler) Name() string {
	return "Dart"
}

func (h *DartHandler) Extensions() []string {
	return []string{".dart"}
}

func (h *DartHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *DartHandler) NeedsCompilation() bool {
	// Dart can run directly with JIT
	return false
}

func (h *DartHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append([]string{"compile", "exe"}, extraFlags...)
	cmdArgs = append(cmdArgs, source, "-o", output)
	cmd := exec.Command("dart", cmdArgs...)
	return cmd.Run()
}

func (h *DartHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"run", file}, args...)
	cmd := exec.Command("dart", cmdArgs...)
	return cmd.CombinedOutput()
}
