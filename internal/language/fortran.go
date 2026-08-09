package language

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/velo4705/polyglot/pkg/types"
)

type FortranHandler struct{}

func (h *FortranHandler) Name() string {
	return "Fortran"
}

func (h *FortranHandler) Extensions() []string {
	return []string{".f", ".f90", ".f95", ".f03", ".f08", ".for", ".fpp"}
}

func (h *FortranHandler) Type() types.LanguageType {
	return types.Compiled
}

func (h *FortranHandler) NeedsCompilation() bool {
	return true
}

func (h *FortranHandler) Compile(source string, output string, extraFlags []string) error {
	cmdArgs := append(extraFlags, source, "-o", output)
	cmd := exec.Command("gfortran", cmdArgs...)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output_bytes))
	}
	return nil
}

func (h *FortranHandler) Run(file string, args []string) ([]byte, error) {
	output := strings.TrimSuffix(file, ".f")
	output = strings.TrimSuffix(output, ".f90")
	output = strings.TrimSuffix(output, ".f95")
	output = strings.TrimSuffix(output, ".f03")
	output = strings.TrimSuffix(output, ".f08")
	output = strings.TrimSuffix(output, ".for")
	output = strings.TrimSuffix(output, ".fpp")
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(output, cmdArgs...)
	return cmd.CombinedOutput()
}