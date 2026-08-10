package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type PureScriptHandler struct{}

func (h *PureScriptHandler) Name() string {
	return "PureScript"
}

func (h *PureScriptHandler) Extensions() []string {
	return []string{".purs"}
}

func (h *PureScriptHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *PureScriptHandler) NeedsCompilation() bool {
	return false
}

func (h *PureScriptHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *PureScriptHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{"run"}, args...)
	cmd := exec.Command("spago", cmdArgs...)
	return cmd.CombinedOutput()
}
