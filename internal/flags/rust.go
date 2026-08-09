package flags

import "regexp"

func init() {
	Register("Rust", &Detector{
		Defaults: []string{"--edition", "2021"},
		Rules: []Rule{
			// Async/await requires 2018+
			{regexp.MustCompile(`async\s+fn|\.await\s|tokio::|async-std::`), []string{"--edition", "2021"}},
			// Gen blocks require 2024
			{regexp.MustCompile(`gen\s*\{|gen\s+move\s*\{`), []string{"--edition", "2024"}},
			// Closures with precise capturing (2021+)
			{regexp.MustCompile(`capture_precise_closing|precise_capturing`), []string{"--edition", "2021"}},
			// let-else (2021+)
			{regexp.MustCompile(`let\s+else\s`), []string{"--edition", "2021"}},
			// RPIT (Return Position Impl Trait) - stable since 2018 but common
			{regexp.MustCompile(`fn\s+\w+.*->\s*impl\s+\w+`), []string{}},
		},
	})
}
