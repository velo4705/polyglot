# Implementation Plan: v1.1.0 UX Polish

## Overview

Implement five additive UX improvements to the Polyglot CLI: shell completions, stdin support, JSON output, improved error messages, and version numbers in `check`. Each feature is a discrete package addition or modification wired into the existing Cobra command tree.

## Tasks

- [x] 1. Add `pgregory.net/rapid` dependency
  - Run `go get pgregory.net/rapid` and commit the updated `go.mod` / `go.sum`
  - _Requirements: testing infrastructure for all property tests_

- [x] 2. Implement shell completions (`internal/cli/completion.go`)
  - [x] 2.1 Create `internal/cli/completion.go` with `completionCmd`
    - Define `var completionCmd` using `cobra.ExactArgs(1)` and `RunE: generateCompletion`
    - Write `completionLong` constant with per-shell install instructions (bash, zsh, fish)
    - Implement `generateCompletion`: delegate to `rootCmd.GenBashCompletion`, `GenZshCompletion`, `GenFishCompletion(os.Stdout, true)`; return error listing supported shells for any other argument
    - Register `completionCmd` on `rootCmd` in `root.go`
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6_

  - [ ]* 2.2 Write property test for completion output markers (Property 1)
    - **Property 1: Completion output contains shell-specific markers**
    - Use `rapid.SampledFrom([]string{"bash","zsh","fish"})` to drive shell selection
    - Assert stdout contains `complete -F` (bash), `#compdef` (zsh), `complete -c polyglot` (fish)
    - Tag: `// Feature: v1-1-ux-polish, Property 1`
    - **Validates: Requirements 1.2, 1.3, 1.4**

  - [ ]* 2.3 Write unit tests for completion edge cases
    - Test unsupported shell name returns non-zero exit and lists supported shells
    - Test help text contains install instructions for each shell
    - _Requirements: 1.5, 1.6_

- [x] 3. Implement stdin handler (`internal/stdin/handler.go`)
  - [x] 3.1 Create `internal/stdin/handler.go`
    - Implement `CanonicalExtension(langName string) (string, bool)` — map covering all 30 languages from `pkg/types/types.go`
    - Implement `ReadToTempFile(langName string) (string, error)` — read all bytes from `os.Stdin`, create temp file via `os.CreateTemp("", "polyglot-stdin-*<ext>")`, write bytes, return path
    - _Requirements: 2.1, 2.2_

  - [ ]* 3.2 Write property test for stdin temp file round-trip (Property 2)
    - **Property 2: Stdin temp file round-trip**
    - Use `rapid.SliceOf(rapid.Byte())` for content and `rapid.SampledFrom(supportedLangs)` for language
    - Assert file content equals input bytes and path suffix matches canonical extension
    - Tag: `// Feature: v1-1-ux-polish, Property 2`
    - **Validates: Requirements 2.1, 2.2**

  - [ ]* 3.3 Write property test for temp file cleanup invariant (Property 3)
    - **Property 3: Temp file cleanup invariant**
    - Use same generators; simulate error injection to cover all exit paths
    - Assert temp file path does not exist on disk after the run command returns
    - Tag: `// Feature: v1-1-ux-polish, Property 3`
    - **Validates: Requirements 2.4**

- [x] 4. Implement JSON formatter (`internal/output/json.go`)
  - [x] 4.1 Create `internal/output/json.go`
    - Define `RunResult` struct with fields: `Language`, `File`, `ExitCode`, `Stdout`, `Stderr`, `DurationMs` and correct `json:` tags
    - Define `CheckEntry` struct with fields: `Language`, `Installed`, `Version`
    - Define `CheckResult` struct with `Languages []CheckEntry`
    - Implement `PrintRun(w io.Writer, r RunResult) error` and `PrintCheck(w io.Writer, r CheckResult) error` using `encoding/json`
    - _Requirements: 3.2, 3.3, 3.7_

  - [ ]* 4.2 Write property test for run JSON output structure (Property 4)
    - **Property 4: Run JSON output structure**
    - Use `rapid.String()` and `rapid.Int()` generators for all `RunResult` fields
    - Assert output parses as valid JSON and all six fields are present with correct types
    - Tag: `// Feature: v1-1-ux-polish, Property 4`
    - **Validates: Requirements 3.2, 3.7**

  - [ ]* 4.3 Write property test for check JSON output structure (Property 5)
    - **Property 5: Check JSON output structure**
    - Use `rapid.SliceOf(checkEntryGen)` where `checkEntryGen` generates `CheckEntry` values
    - Assert output parses as valid JSON with a `languages` array; every element has correct field types
    - Tag: `// Feature: v1-1-ux-polish, Property 5`
    - **Validates: Requirements 3.3, 5.7**

  - [ ]* 4.4 Write property test for JSON error encoding (Property 6)
    - **Property 6: JSON error encoding**
    - Use `rapid.String()` for error messages and `rapid.IntRange(1, 255)` for exit codes
    - Assert `exit_code` is non-zero, `stderr` is non-empty, and no plain-text error appears on stdout
    - Tag: `// Feature: v1-1-ux-polish, Property 6`
    - **Validates: Requirements 3.4, 4.6**

- [x] 5. Implement error advisor (`internal/errors/advisor.go`)
  - [x] 5.1 Create `internal/errors/advisor.go`
    - Define `executableToLanguage` map covering `python3→python`, `node→javascript`, `ruby→ruby`, `go→go`, `java→java`, `rustc→rust`, and all other supported executables
    - Implement `NotFound(langName string) string` — format: `<Language> is not installed. Run: polyglot install <language>`
    - Implement `FileNotFound(path string) string` — format: `File not found: <path>. Check the path and try again.`
    - Implement `UnknownExtension(ext string, supported []string) string` — mention supported extensions and `--lang` flag
    - Implement `Generic(context string, err error) string` — prefix original error with context string
    - Implement `IsNotFoundError(err error) (executable string, ok bool)` — detect `exec.ErrNotFound` or "no such file" errors
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

  - [ ]* 5.2 Write property test for missing runtime error format (Property 7)
    - **Property 7: Missing runtime error message format**
    - Use `rapid.SampledFrom(supportedLangs)` to drive language selection
    - Assert `NotFound` output matches `<Language> is not installed. Run: polyglot install <language>`
    - Tag: `// Feature: v1-1-ux-polish, Property 7`
    - **Validates: Requirements 4.1**

  - [ ]* 5.3 Write property test for file-not-found error format (Property 8)
    - **Property 8: File-not-found error message format**
    - Use `rapid.String()` for arbitrary file paths
    - Assert `FileNotFound` output matches `File not found: <path>. Check the path and try again.`
    - Tag: `// Feature: v1-1-ux-polish, Property 8`
    - **Validates: Requirements 4.2**

  - [ ]* 5.4 Write property test for unknown extension error format (Property 9)
    - **Property 9: Unknown extension error message format**
    - Use `rapid.String()` filtered to extensions not in the known set
    - Assert `UnknownExtension` output mentions the supported extensions list and the `--lang` flag
    - Tag: `// Feature: v1-1-ux-polish, Property 9`
    - **Validates: Requirements 4.3**

  - [ ]* 5.5 Write property test for executable-to-language mapping completeness (Property 10)
    - **Property 10: Executable-to-language mapping completeness**
    - Use `rapid.SampledFrom(executableNames)` over all keys defined in Requirement 4.5
    - Assert `executableToLanguage[exe]` returns the correct canonical Polyglot language name for every entry
    - Tag: `// Feature: v1-1-ux-polish, Property 10`
    - **Validates: Requirements 4.5**

  - [ ]* 5.6 Write unit tests for generic error wrapping
    - Test `Generic` prefixes the original error message with the context string
    - _Requirements: 4.4_

- [x] 6. Checkpoint — ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. Implement version parser (`internal/version/parser.go`)
  - [x] 7.1 Create `internal/version/parser.go`
    - Define `VersionSpec` struct with `Command string`, `Args []string`, `ParseFn func(string) string`
    - Populate `var Specs map[string]VersionSpec` for Python (`python3 --version`, split+last token), Node.js (`node --version`, trim `v`), Go (`go version`, split+3rd token, trim `go`), Ruby (`ruby --version`, split+2nd token), Java (`java --version`, first line split+2nd token), Rust (`rustc --version`, split+2nd token), and all other supported languages (first `\d+\.\d+` match)
    - Implement `Get(langName string) (version string, installed bool)` — run version command, call `ParseFn`, return `("unknown", true)` on non-zero exit or unparseable output, `("", false)` if not installed
    - _Requirements: 5.1, 5.2, 5.5, 5.6_

  - [ ]* 7.2 Write property test for version parser with known outputs (Property 11)
    - **Property 11: Version parser returns non-empty string for known output**
    - Use `rapid.SampledFrom(supportedLangs)` with fixture version command outputs per language
    - Assert `Get` returns a non-empty, non-`"unknown"` version string for each fixture
    - Tag: `// Feature: v1-1-ux-polish, Property 11`
    - **Validates: Requirements 5.1, 5.2, 5.6**

- [x] 8. Modify `internal/cli/run.go` — wire stdin, JSON, and error advisor
  - [x] 8.1 Add `--json` flag registration to `runCmd`
    - Register `--json` boolean flag on `runCmd`
    - _Requirements: 3.1_

  - [x] 8.2 Wire stdin handler into `runCmd`
    - Detect stdin via `!term.IsTerminal(int(os.Stdin.Fd()))` or equivalent
    - If stdin detected and `--lang` not set, call `ErrorAdvisor.Generic` and exit 1
    - If stdin detected and `--lang` set, call `stdin.ReadToTempFile`; defer `os.Remove` of the temp path; pass temp path to executor
    - If `stdin.ReadToTempFile` returns an error, call `ErrorAdvisor.Generic` and exit 1
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

  - [x] 8.3 Wire JSON formatter and error advisor into `runCmd`
    - When `--json` is set: capture executor stdout/stderr into `bytes.Buffer`; populate `RunResult`; call `output.PrintRun`; suppress all `ui.Error` calls in favour of encoding errors into `RunResult.Stderr` and setting `RunResult.ExitCode`
    - Replace bare `ui.Error` calls for runtime-not-found, file-not-found, and unknown-extension with `ErrorAdvisor.NotFound`, `ErrorAdvisor.FileNotFound`, `ErrorAdvisor.UnknownExtension`, and `ErrorAdvisor.Generic`
    - Honour `--json` + `--quiet` precedence: JSON wins
    - _Requirements: 3.2, 3.4, 3.5, 3.6, 4.1, 4.2, 4.3, 4.4, 4.6_

  - [ ]* 8.4 Write unit tests for `run.go` wiring
    - Test `--lang` required when stdin is detected
    - Test `--json` + `--quiet` produces JSON output
    - Test `--json` flag is registered on `runCmd`
    - _Requirements: 2.5, 3.1, 3.6_

- [x] 9. Modify `internal/cli/check.go` — wire version parser and JSON
  - [x] 9.1 Add `--json` flag registration to `checkCmd`
    - Register `--json` boolean flag on `checkCmd`
    - _Requirements: 3.1_

  - [x] 9.2 Wire version parser into `checkCmd`
    - For each installed language, call `version.Get(langName)`
    - Display version string alongside `INSTALLED` status; display `INSTALLED (version unknown)` when `Get` returns `"unknown"`; display `NOT FOUND` with no version when not installed
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

  - [x] 9.3 Wire JSON formatter into `checkCmd`
    - When `--json` is set: build `CheckResult` slice from checker results (including version strings); call `output.PrintCheck`; suppress plain-text table output
    - _Requirements: 3.3, 3.5, 5.7_

  - [ ]* 9.4 Write unit tests for `check.go` wiring
    - Test NOT FOUND row has no version string
    - Test `--json` flag is registered on `checkCmd`
    - _Requirements: 3.1, 5.4_

- [x] 10. Final checkpoint — ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use `pgregory.net/rapid` with a minimum of 100 iterations each
- All property tests are tagged with `// Feature: v1-1-ux-polish, Property <N>`
- Checkpoints at tasks 6 and 10 ensure incremental validation
