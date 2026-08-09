# Changelog

All notable changes to Polyglot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-08-09

### Added
- 21 new languages: Fortran, Pascal, Ada, COBOL, Scheme, Common Lisp, Forth, Prolog, Tcl, Clojure, Gleam, Elm, PureScript, Roc, V, Odin, NASM, GAS, ARM Assembly, MIPS Assembly, RISC-V Assembly
- Assembly language support (NASM, GAS, ARM, MIPS, RISC-V)
- `--compile-flags` CLI flag for manual compiler flag override
- Auto-detection of compiler flags from source code patterns (C/C++: -std=c++17/20/23, -pthread, -lm, -fopenmp, etc.)
- `--self-correct` flag properly registered (was dead code)

### Fixed
- Version constant single source of truth (`pkg/version/version.go`) — was duplicated in `main.go` (1.0.1) and `cli/version.go` (1.2.0)
- `--self-correct` flag now registered with Cobra (was declared but never wired up)
- HTTP client timeouts added (30s default) — `http.DefaultClient` had no timeout
- Spinner data race fixed — `active` field changed to `atomic.Bool`
- Spinner double-stop deadlock fixed — uses `select` with `default`
- Duplicate error output fixed — `SilenceErrors: true` on root command

### Changed
- `LanguageHandler.Compile()` signature now accepts `extraFlags []string`
- All compiled language handlers inject extraFlags into build commands
- `pkg/types/types.go` refactored from giant switch statements to table-driven registry
- Extension detection centralized in `types.FromExtension()` (was duplicated in `detector.go`)
- `Language` type methods: `Name()`, `Extensions()`, `Type()` now use registry look up
- Added `AllLanguages()`, `AllLanguageNames()`, `AllExtensions()` helpers to `pkg/types`
- Added `Language.NeedsCompilation()` convenience method

## [1.0.1] - 2026-02-22

### Added - Auto-Update Feature

#### Update Command
- New `polyglot update` command for self-updating
- `--check` flag to check for updates without installing
- Automatic version checking against GitHub releases
- Platform-specific binary downloads (Linux, macOS, Windows)
- Support for amd64 and arm64 architectures

#### Background Update Checks
- Automatic update checks every 24 hours
- Non-blocking background checks
- Update notifications in terminal
- Configurable check interval via `~/.polyglot/last_update_check`

#### Package Manager Integration
- Automatic detection of Homebrew installations
- Automatic detection of APT installations
- Falls back to direct binary update if not installed via package manager
- Respects package manager update workflows

#### Safety Features
- Automatic backup before updating
- Rollback on update failure
- Verification of downloaded binaries
- HTTPS-only downloads from GitHub releases

#### Documentation
- New `docs/AUTO_UPDATE.md` with comprehensive guide
- Updated `docs/QUICKSTART.md` with update instructions
- Updated `README.md` to mention update command
- Examples and troubleshooting guide

### Changed
- Version bumped to 1.0.1
- Main CLI now includes background update checking
- Version command updated

### Technical Details
- New `internal/updater/updater.go` package
- New `internal/cli/update.go` command
- Integration with main CLI execution flow
- Cross-platform binary replacement logic
- Windows-specific update handling (can't replace running executable)

## [1.0.0] - 2026-02-22

### Added - Production Ready Release

Polyglot is now production-ready with:
- 30+ programming languages
- Auto-install system
- Configuration management
- Watch mode
- Colored output
- Comprehensive testing (85%+ coverage)
- Distribution packages
- Complete documentation

## [0.5.0] - 2026-02-22

### Added - Milestone 7 (Developer Experience)

#### Colored Output
- ANSI color support for all terminal output
- Success messages in green (✓)
- Error messages in red (✗)
- Warnings in yellow (⚠)
- Info messages in blue (ℹ)
- Step indicators in cyan (→)
- Language names highlighted in magenta
- Commands highlighted in cyan
- Filenames highlighted in blue
- Automatic color detection (respects NO_COLOR env var)

#### Watch Mode
- New `polyglot watch` command
- Automatic file change detection
- Configurable watch interval (`--interval` flag)
- Optional screen clearing (`--clear` flag)
- Execution timing display
- Colored status messages
- Support for program arguments

#### Dry Run Mode
- New `--dry-run` flag for `polyglot run`
- Preview execution without running
- Shows compilation and execution commands
- Useful for learning and safety

#### Progress Indicators
- Animated spinner for compilation
- Progress bars for multi-step operations
- Automatic hiding when complete
- Respects NO_COLOR environment variable

#### Better Error Messages
- Colored error messages
- Clear, actionable feedback
- Helpful suggestions
- Improved formatting

#### UI Package
- New `internal/ui/` package
- `colors.go` - Color utilities
- `progress.go` - Spinner and progress bar
- Reusable UI components

#### Watcher Package
- New `internal/watcher/` package
- File watching implementation
- Efficient polling mechanism

### Changed
- Version bumped to 0.5.0
- All commands now use colored output
- Enhanced error messages across all commands
- Improved visual hierarchy in output
- Better formatting for all messages

### Technical Changes
- Added `internal/ui/colors.go` - Color system
- Added `internal/ui/progress.go` - Progress indicators
- Added `internal/watcher/watcher.go` - File watching
- Added `internal/cli/watch.go` - Watch command
- Enhanced `internal/cli/run.go` - Dry-run mode, colored output
- Enhanced `internal/executor/executor.go` - Progress indicators

## [0.4.0] - 2026-02-22

### Added - Milestone 6 (Configuration & Customization)

#### Configuration System
- **Global Configuration** (`~/.polyglot/config.yaml`) - System-wide settings
- **Project Configuration** (`.polyglot.yaml`) - Project-specific settings
- **Configuration Merging** - Project config overrides global config
- **Nested Key Access** - Dot notation for configuration keys

#### Custom Language Definitions
- Per-language enable/disable
- Custom command paths
- Custom compiler/interpreter flags
- Version preferences
- Installation method configuration
- Language-specific environment variables

#### Custom Extension Mapping
- Map any file extension to any language
- Support for framework-specific extensions (.jsx, .tsx, .vue, etc.)
- Dialect support (.pyx for Cython, etc.)
- Module system extensions (.mjs, .cjs, etc.)

#### Environment Variables
- Global environment variable configuration
- Language-specific environment variables
- Automatic merging of global and language-specific vars

#### Enhanced Config Command
- `polyglot config init` - Initialize configuration file
- `polyglot config show` - Show current configuration
- `polyglot config get <key>` - Get configuration value
- `polyglot config set <key> <value>` - Set configuration value
- `polyglot config path` - Show configuration file paths

#### Documentation
- New comprehensive configuration guide (docs/CONFIGURATION.md)
- Milestone 6 completion document (MILESTONE6_COMPLETE.md)
- Updated STATUS.md with progress tracking

### Changed
- Version bumped to 0.4.0
- Enhanced detector to support custom extensions
- Enhanced executor to support configuration
- Enhanced run command to load and use configuration
- Updated all documentation to reflect new features

### Technical Changes
- Added `LanguageConfig` struct for per-language settings
- Added `CustomExtensions` map for extension mapping
- Added `Environment` map for environment variables
- Added `FromString()` function to convert language names
- Added `DetectLanguageWithConfig()` for custom extension support
- Added configuration merging logic
- Added nested key access helpers

## [0.3.0] - 2026-02-21

### Added - Milestone 5 (Auto-Install System)

#### Auto-Install Features
- Automatic toolchain detection
- System package manager integration (dnf, apt, brew, pacman, zypper, apk)
- Interactive installation prompts
- Configuration system (~/.polyglot/config.yaml)
- `polyglot install <language>` command
- `polyglot config` command
- Auto-install mode (skip prompts)
- System-wide installation

#### Package Manager Support
- Fedora/RHEL (dnf)
- Ubuntu/Debian (apt)
- macOS (brew)
- Arch Linux (pacman)
- openSUSE (zypper)
- Alpine Linux (apk)

#### Configuration
- YAML-based configuration file
- Auto-install enable/disable
- Installation method selection
- Prompt enable/disable

### Changed
- Version bumped to 0.3.0
- Enhanced run command with auto-install integration
- Updated documentation with auto-install information

## [0.2.0] - 2026-02-21

### Added - Milestone 2 & 3 (Compiled + Extended Languages)

#### New Languages (11 added!)

**Milestone 2 - Compiled Languages:**
- Java (.java) - Full compile and run support with automatic cleanup
- C (.c) - Compiled with gcc
- C++ (.cpp, .cc, .cxx, .C) - Compiled with g++
- Rust (.rs) - Compiled with rustc

**Milestone 3 - Extended Languages:**
- Perl (.pl) - Interpreted scripting language
- Lua (.lua) - Lightweight scripting language
- Zig (.zig) - Modern systems programming language
- Nim (.nim) - Efficient compiled language
- Crystal (.cr) - Ruby-like compiled language
- D (.d) - Systems programming language
- Shell (.sh, .bash) - Bash shell scripts

#### Improvements
- Enhanced compilation workflow with automatic cleanup
- Better handling of compiled binaries
- Automatic removal of Java .class files
- Improved error messages for compilation failures
- Support for multiple file extensions per language

#### Total Language Count
- **16 languages** now supported (up from 5)
- **7 interpreted** languages
- **9 compiled** languages

### Technical Changes
- Updated language registry with 11 new handlers
- Enhanced executor to handle compilation cleanup
- Improved detector to support multiple extensions
- Updated check command to show all 16 toolchains

## [0.1.0] - 2026-02-21

### Added - Milestone 1 (MVP)

#### Core Features
- Automatic language detection from file extensions
- Process execution with proper output handling
- Exit code propagation from executed programs
- Error handling and user-friendly error messages

#### Language Support
- Python (.py) - interpreted via python3
- Go (.go) - compiled and run via go run
- JavaScript (.js) - interpreted via node
- Ruby (.rb) - interpreted via ruby
- PHP (.php) - interpreted via php

#### CLI Commands
- `polyglot run <file>` - Detect and run any supported file
- `polyglot compile <file>` - Compile without running (for compiled languages)
- `polyglot list` - Show all supported languages and their extensions
- `polyglot check` - Verify which language toolchains are installed
- `polyglot version` - Display version and milestone information

#### CLI Flags
- `-v, --verbose` - Show detailed execution information
- `-q, --quiet` - Suppress polyglot messages (show only program output)
- `--args` - Pass comma-separated arguments to the executed program

#### Documentation
- Comprehensive README with installation and usage instructions
- Quick Start Guide for getting started in 5 minutes
- Language Support Matrix showing all supported languages
- Examples document with real-world usage scenarios
- Contributing Guide for adding new languages
- 10-milestone Roadmap for future development
- Project Structure documentation
- Visual Demo of all features
- Milestone 1 Summary and Checklist

#### Build & Distribution
- Makefile with targets: build, install, test, check, fmt, clean, help
- Installation script (scripts/install.sh)
- Comprehensive test script (scripts/test-all.sh)
- GitHub Actions CI workflow for automated testing
- GitHub Actions Release workflow for multi-platform binaries
- Cross-platform support (Linux, macOS, Windows)

#### Testing
- Unit tests for language detector
- Test fixtures for all 5 supported languages
- Integration test suite
- Code formatting checks
- Go vet checks

#### Project Infrastructure
- Clean modular architecture with clear separation of concerns
- Interface-based design for easy extensibility
- MIT License
- .gitignore for Go projects
- go.mod with dependencies (cobra for CLI)

### Technical Details

#### Architecture
- `cmd/polyglot/` - Application entry point
- `internal/cli/` - CLI command implementations
- `internal/detector/` - Language detection logic
- `internal/executor/` - Process execution wrapper
- `internal/language/` - Language-specific handlers
- `pkg/types/` - Shared types and interfaces

#### Dependencies
- github.com/spf13/cobra v1.8.0 - CLI framework
- Go 1.21+ required

### Metrics
- ~1,500 lines of code
- 30+ files
- 5 languages supported
- 5 CLI commands
- 7 documentation files
- ~8MB binary size
- ~2 second build time

---

