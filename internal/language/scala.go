package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type ScalaHandler struct{}

func (h *ScalaHandler) Name() string {
	return "Scala"
}

func (h *ScalaHandler) Extensions() []string {
	return []string{".scala", ".sc"}
}

func (h *ScalaHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *ScalaHandler) NeedsCompilation() bool {
	// Scala can use scala-cli for quick execution
	return false
}

func (h *ScalaHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source)
	cmd := exec.Command("scalac", cmdArgs...)
	return cmd.Run()
}

func (h *ScalaHandler) Run(file string, args []string) ([]byte, error) {
	// Use scala command for direct execution
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("scala", cmdArgs...)
	return cmd.CombinedOutput()
}
