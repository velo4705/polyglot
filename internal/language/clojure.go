package language

import (
	"os/exec"

	"github.com/velo4705/polyglot/pkg/types"
)

type ClojureHandler struct{}

func (h *ClojureHandler) Name() string {
	return "Clojure"
}

func (h *ClojureHandler) Extensions() []string {
	return []string{".clj", ".cljs", ".cljc", ".edn"}
}

func (h *ClojureHandler) Type() types.LanguageType {
	return types.Interpreted
}

func (h *ClojureHandler) NeedsCompilation() bool {
	return false
}

func (h *ClojureHandler) Compile(source string, output string, extraFlags []string) error {
	return nil
}

func (h *ClojureHandler) Run(file string, args []string) ([]byte, error) {
	cmdArgs := append([]string{file}, args...)
	cmd := exec.Command("clojure", cmdArgs...)
	return cmd.CombinedOutput()
}