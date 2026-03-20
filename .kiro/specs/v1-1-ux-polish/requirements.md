# Requirements Document

## Introduction

Polyglot v1.1.0 introduces five UX polish improvements to the existing CLI tool that detects and runs source files in 30+ programming languages. The release focuses on developer experience: shell completions for discoverability, stdin support for pipeline workflows, JSON output for CI/CD integration, improved error messages with actionable guidance, and version numbers in the `check` command.

## Glossary

- **CLI**: The `polyglot` command-line interface built with Go and Cobra.
- **Runner**: The internal component that executes source files in a detected or specified language.
- **Detector**: The internal component that infers a programming language from a file extension.
- **Checker**: The internal component that probes whether language runtimes are installed.
- **Version_Parser**: The per-language component that extracts a version string from runtime version output.
- **Completion_Generator**: The Cobra-based component that emits shell completion scripts.
- **Stdin_Handler**: The component that reads from standard input, writes to a temp file, and passes it to the Runner.
- **JSON_Formatter**: The component that serializes command output to machine-readable JSON.
- **Error_Advisor**: The component that maps low-level system errors to human-readable, actionable messages.

---

## Requirements

### Requirement 1: Shell Completions

**User Story:** As a developer, I want to run `polyglot completion bash|zsh|fish` so that my shell can tab-complete Polyglot subcommands and flags without manual configuration.

#### Acceptance Criteria

1. THE CLI SHALL expose a `completion` subcommand that accepts exactly one positional argument: `bash`, `zsh`, or `fish`.
2. WHEN the user runs `polyglot completion bash`, THE Completion_Generator SHALL write a valid Bash completion script to stdout.
3. WHEN the user runs `polyglot completion zsh`, THE Completion_Generator SHALL write a valid Zsh completion script to stdout.
4. WHEN the user runs `polyglot completion fish`, THE Completion_Generator SHALL write a valid Fish completion script to stdout.
5. WHEN the user runs `polyglot completion` with an unsupported shell name, THE CLI SHALL print an error message listing the supported shells and exit with a non-zero status code.
6. THE CLI help text for the `completion` subcommand SHALL include installation instructions for each supported shell.

---

### Requirement 2: Stdin Support

**User Story:** As a developer, I want to pipe source code into `polyglot run --lang <language>` so that I can run code snippets inline without creating a file.

#### Acceptance Criteria

1. WHEN the user pipes data to `polyglot run` and provides the `--lang` flag, THE Stdin_Handler SHALL read all bytes from stdin until EOF.
2. WHEN stdin input is received, THE Stdin_Handler SHALL write the content to a temporary file with the canonical extension for the specified language.
3. WHEN the temporary file has been created, THE Runner SHALL execute it using the same logic applied to regular file arguments.
4. WHEN execution completes or is interrupted, THE Stdin_Handler SHALL delete the temporary file.
5. IF the user pipes data to `polyglot run` without providing the `--lang` flag, THEN THE CLI SHALL print an error message stating that `--lang` is required when reading from stdin and exit with a non-zero status code.
6. IF writing the temporary file fails, THEN THE Stdin_Handler SHALL print an error message and exit with a non-zero status code.
7. WHILE stdin input is being read, THE CLI SHALL accept the existing `--verbose`, `--quiet`, `--args`, and `--dry-run` flags with their current semantics.

---

### Requirement 3: JSON Output Flag

**User Story:** As a CI/CD pipeline author, I want `polyglot run` and `polyglot check` to support a `--json` flag so that I can parse structured output programmatically.

#### Acceptance Criteria

1. THE CLI SHALL expose a `--json` boolean flag on both the `run` and `check` subcommands.
2. WHEN `polyglot run --json` is invoked, THE JSON_Formatter SHALL write a single JSON object to stdout containing the fields: `language` (string), `file` (string), `exit_code` (integer), `stdout` (string), `stderr` (string), and `duration_ms` (integer).
3. WHEN `polyglot check --json` is invoked, THE JSON_Formatter SHALL write a single JSON object to stdout containing a `languages` array where each element has the fields: `language` (string), `installed` (boolean), and `version` (string, empty string when not installed).
4. WHEN `--json` is active and the Runner encounters an error, THE JSON_Formatter SHALL include the error detail in the `stderr` field and set `exit_code` to the appropriate non-zero value rather than printing a plain-text error.
5. WHEN `--json` is active, THE CLI SHALL write no plain-text output to stdout; all diagnostic information SHALL be encoded in the JSON payload.
6. IF `--json` and `--quiet` are both provided to `run`, THEN THE CLI SHALL treat `--json` as taking precedence and produce JSON output.
7. THE JSON output SHALL be valid JSON as defined by RFC 8259 for all inputs.

---

### Requirement 4: Improved Error Messages

**User Story:** As a developer, I want Polyglot to show actionable error messages so that I can resolve problems without searching documentation.

#### Acceptance Criteria

1. WHEN the Runner fails because a language runtime executable is not found on PATH, THE Error_Advisor SHALL display a message in the format: `<Language> is not installed. Run: polyglot install <language>`.
2. WHEN the user provides a file path that does not exist, THE Error_Advisor SHALL display a message in the format: `File not found: <path>. Check the path and try again.`
3. WHEN the user provides a file with an extension that the Detector does not recognise, THE Error_Advisor SHALL display a message listing the supported extensions and suggesting the `--lang` flag to override detection.
4. WHEN the Runner fails for a reason other than the cases above, THE Error_Advisor SHALL display the original system error prefixed with a human-readable context string.
5. THE Error_Advisor SHALL map the following runtime executables to their install hint language names: `python3` → `python`, `node` → `javascript`, `ruby` → `ruby`, `go` → `go`, `java` → `java`, `rustc` → `rust`, and all other supported language executables to their canonical Polyglot language names.
6. IF `--json` is active, THEN THE Error_Advisor SHALL encode all error messages in the JSON payload rather than printing plain text.

---

### Requirement 5: `polyglot check` Version Numbers

**User Story:** As a developer, I want `polyglot check` to display the installed version of each language runtime so that I can verify my environment at a glance.

#### Acceptance Criteria

1. WHEN `polyglot check` is run, THE Checker SHALL invoke the version command for each language runtime that is detected as installed.
2. WHEN a version command succeeds, THE Version_Parser SHALL extract the version string from the command output using a per-language parsing rule.
3. WHEN a version string is successfully extracted, THE Checker SHALL display it alongside the `INSTALLED` status in the check output.
4. WHEN a language runtime is not installed, THE Checker SHALL display `NOT FOUND` with no version string, preserving the existing behaviour.
5. IF a version command exits with a non-zero status or produces unparseable output, THEN THE Checker SHALL display `INSTALLED (version unknown)` rather than failing the entire check.
6. THE Version_Parser SHALL support per-language version extraction for at minimum: Python (`python3 --version`), Node.js (`node --version`), Go (`go version`), Ruby (`ruby --version`), Java (`java --version`), Rust (`rustc --version`), and all other languages supported by Polyglot.
7. WHEN `polyglot check --json` is invoked, THE JSON_Formatter SHALL include the extracted version string in the `version` field of each language entry as specified in Requirement 3.
