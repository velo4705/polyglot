package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type ErlangHandler struct{}

func (h *ErlangHandler) Name() string {
	return "Erlang"
}

func (h *ErlangHandler) Extensions() []string {
	return []string{".erl"}
}

func (h *ErlangHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *ErlangHandler) NeedsCompilation() bool {
	return true
}

func (h *ErlangHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source)
	cmd := exec.Command("erlc", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *ErlangHandler) Run(file string, args []string) ([]byte, error) {
	// Get module name from filename
	moduleName := strings.TrimSuffix(filepath.Base(file), ".erl")
	dir := filepath.Dir(file)

	// Run with erl
	cmdArgs := []string{"-noshell", "-pa", dir, "-s", moduleName, "main", "-s", "init", "stop"}
	cmd := exec.Command("erl", cmdArgs...)
	return cmd.CombinedOutput()
}
