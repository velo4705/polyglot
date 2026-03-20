package output

import (
	"encoding/json"
	"io"
)

// RunResult is the JSON payload for `polyglot run --json`.
type RunResult struct {
	Language   string `json:"language"`
	File       string `json:"file"`
	ExitCode   int    `json:"exit_code"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int64  `json:"duration_ms"`
}

// CheckEntry is one element of the `languages` array.
type CheckEntry struct {
	Language  string `json:"language"`
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
}

// CheckResult is the JSON payload for `polyglot check --json`.
type CheckResult struct {
	Languages []CheckEntry `json:"languages"`
}

// PrintRun serialises r as indented JSON and writes it to w.
func PrintRun(w io.Writer, r RunResult) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// PrintCheck serialises r as indented JSON and writes it to w.
func PrintCheck(w io.Writer, r CheckResult) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
