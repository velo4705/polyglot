package flags

import (
	"regexp"
)

// Rule defines a pattern-to-flags mapping for auto-detection.
type Rule struct {
	Pattern *regexp.Regexp
	Flags   []string
}

// Detector scans source code and returns auto-detected compiler flags.
type Detector struct {
	Rules    []Rule
	Defaults []string
}

// Detect scans source code against all rules and returns deduplicated flags.
func (d *Detector) Detect(source []byte) []string {
	flags := append([]string{}, d.Defaults...)
	seen := make(map[string]bool)
	for _, f := range flags {
		seen[f] = true
	}

	for _, rule := range d.Rules {
		if rule.Pattern.Match(source) {
			for _, f := range rule.Flags {
				if !seen[f] {
					flags = append(flags, f)
					seen[f] = true
				}
			}
		}
	}
	return flags
}

// detectors maps language names to their built-in detectors.
var detectors = map[string]*Detector{}

// Register adds a detector for a language.
func Register(lang string, d *Detector) {
	detectors[lang] = d
}

// Get returns the detector for a language, or nil if none exists.
func Get(lang string) *Detector {
	return detectors[lang]
}
