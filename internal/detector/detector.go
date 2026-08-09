package detector

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/velo4705/polyglot/internal/config"
	"github.com/velo4705/polyglot/pkg/types"
)

// DetectLanguage detects the programming language based on file extension
func DetectLanguage(filename string) (types.Language, error) {
	return DetectLanguageWithConfig(filename, nil)
}

// DetectLanguageWithConfig detects language with custom configuration support
func DetectLanguageWithConfig(filename string, cfg *config.Config) (types.Language, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	// Check custom extensions first
	if cfg != nil {
		if customLang := cfg.GetCustomExtension(ext); customLang != "" {
			lang := types.FromString(customLang)
			if lang != types.Unknown {
				return lang, nil
			}
		}
	}

	lang := types.FromExtension(ext)
	if lang == types.Unknown {
		return types.Unknown, fmt.Errorf("unsupported file extension: %s", ext)
	}
	return lang, nil
}
