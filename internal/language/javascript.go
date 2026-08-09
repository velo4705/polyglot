package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type JavaScriptHandler struct{}

func (h *JavaScriptHandler) Name() string {
	return "JavaScript"
}

func (h *JavaScriptHandler) Extensions() []string {
	return []string{".js"}
}

func (h *JavaScriptHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *JavaScriptHandler) NeedsCompilation() bool {
	return false
}

func (h *JavaScriptHandler) Compile(source string, output string, extraFlags []string) error {
	// JavaScript doesn't need compilation
	return nil
}

func (h *JavaScriptHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("node", cmdArgs...)
	return cmd.CombinedOutput()
}
