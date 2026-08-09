package language

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type NasmHandler struct{}

func (h *NasmHandler) Name() string {
	return "NASM"
}

func (h *NasmHandler) Extensions() []string {
	return []string{".asm", ".nasm"}
}

func (h *NasmHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *NasmHandler) NeedsCompilation() bool {
	return true
}

func (h *NasmHandler) Compile(source string, output string, extraFlags []string) error {
	objFile := output + ".o"
	objArgs := append([]string{"-f", "elf64", source, "-o", objFile}, extraFlags...)
	cmd := exec.Command("nasm", objArgs...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nasm assembly failed: %s", string(outputBytes))
	}
	// Link with ld
	linkCmd := exec.Command("ld", objFile, "-o", output)
	linkOut, err := linkCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("linking failed: %s", string(linkOut))
	}
	return nil
}

func (h *NasmHandler) Run(file string, args []string) ([]byte, error) {
	ext := filepath.Ext(file)
	output := strings.TrimSuffix(file, ext)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}
