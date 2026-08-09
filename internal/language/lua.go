package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type LuaHandler struct{}

func (h *LuaHandler) Name() string {
	return "Lua"
}

func (h *LuaHandler) Extensions() []string {
	return []string{".lua"}
}

func (h *LuaHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *LuaHandler) NeedsCompilation() bool {
	return false
}

func (h *LuaHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *LuaHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("lua", cmdArgs...)
	return cmd.CombinedOutput()
}
