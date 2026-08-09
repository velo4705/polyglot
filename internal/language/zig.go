package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type ZigHandler struct{}

func (h *ZigHandler) Name() string {
	return "Zig"
}

func (h *ZigHandler) Extensions() []string {
	return []string{".zig"}
}

func (h *ZigHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *ZigHandler) NeedsCompilation() bool {
	// Zig can use 'zig run' which compiles and runs
	return false
}

func (h *ZigHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append([]string{"build-exe"}, extraFlags...)
	cmdArgs = append(cmdArgs, source, "-femit-bin="+output)
	cmd := exec.Command("zig", cmdArgs...)
	return cmd.Run()
}

func (h *ZigHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"run", file}, args...)
	cmd := exec.Command("zig", cmdArgs...)
	return cmd.CombinedOutput()
}
