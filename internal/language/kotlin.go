package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type KotlinHandler struct{}

func (h *KotlinHandler) Name() string {
	return "Kotlin"
}

func (h *KotlinHandler) Extensions() []string {
	return []string{".kt", ".kts"}
}

func (h *KotlinHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *KotlinHandler) NeedsCompilation() bool {
	// Kotlin scripts (.kts) can be run directly
	return false
}

func (h *KotlinHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-include-runtime", "-d", output+".jar")
	cmd := exec.Command("kotlinc", cmdArgs...)
	return cmd.Run()
}

func (h *KotlinHandler) Run(file string, args []string) ([]byte, error) {
	// Use kotlin command for scripts
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("kotlin", cmdArgs...)
	return cmd.CombinedOutput()
}
