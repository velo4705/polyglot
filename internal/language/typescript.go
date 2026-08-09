package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type TypeScriptHandler struct{}

func (h *TypeScriptHandler) Name() string {
	return "TypeScript"
}

func (h *TypeScriptHandler) Extensions() []string {
	return []string{".ts"}
}

func (h *TypeScriptHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *TypeScriptHandler) NeedsCompilation() bool {
	// Use ts-node for direct execution
	return false
}

func (h *TypeScriptHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source)
	cmd := exec.Command("tsc", cmdArgs...)
	return cmd.Run()
}

func (h *TypeScriptHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("ts-node", cmdArgs...)
	return cmd.CombinedOutput()
}
