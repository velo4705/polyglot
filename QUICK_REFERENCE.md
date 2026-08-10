# Polyglot Quick Reference

One-page reference for all Polyglot commands and features.

## Installation

```bash
# Quick install
./scripts/install.sh

# Or with make
make install

# Or manually
go build -o polyglot ./cmd/polyglot
sudo mv polyglot /usr/local/bin/
```

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `polyglot run <file>` | Run a file | `polyglot run hello.py` |
| `polyglot compile <file>` | Compile only | `polyglot compile main.rs` |
| `polyglot list` | List languages | `polyglot list` |
| `polyglot check` | Check toolchains | `polyglot check` |
| `polyglot update` | Update Polyglot | `polyglot update` |
| `polyglot version` | Show version | `polyglot version` |
| `polyglot --help` | Show help | `polyglot --help` |

## Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Detailed output | `polyglot run -v hello.py` |
| `--quiet` | `-q` | Silent mode | `polyglot run -q hello.py` |
| `--args` | | Pass arguments | `polyglot run script.py --args a,b,c` |
| `--compile-flags` | | Compiler flags | `polyglot run hello.cpp --compile-flags "-O3"` |
| `--self-correct` | | LLM error correction | `polyglot run broken.py --self-correct` |
| `--dry-run` | | Preview only | `polyglot run --dry-run hello.c` |
| `--lang` | | Force language (stdin) | `polyglot run --lang Python` |
| `--sandbox` | | Sandbox execution | `polyglot run --sandbox hello.sh` |
| `--json` | | JSON output | `polyglot run --json hello.py` |
| `--help` | `-h` | Show help | `polyglot run --help` |

## Supported Languages

| Language | Extension | Command | Type |
|----------|-----------|---------|------|
| Python | `.py` | `python3` | Interpreted |
| JavaScript | `.js` | `node` | Interpreted |
| Ruby | `.rb` | `ruby` | Interpreted |
| PHP | `.php` | `php` | Interpreted |
| Perl | `.pl` | `perl` | Interpreted |
| Lua | `.lua` | `lua` | Interpreted |
| Shell | `.sh` `.bash` | `bash` | Interpreted |
| Go | `.go` | `go run` | Compiled* |
| Java | `.java` | `javac` + `java` | Compiled |
| C | `.c` | `gcc` | Compiled |
| C++ | `.cpp` `.cc` `.cxx` `.C` | `g++` | Compiled |
| Rust | `.rs` | `rustc` | Compiled |
| Zig | `.zig` | `zig run` | Compiled* |
| Nim | `.nim` | `nim c` | Compiled |
| Crystal | `.cr` | `crystal run` | Compiled* |
| D | `.d` | `dmd` | Compiled |
| Swift | `.swift` | `swift` | Compiled* |
| Kotlin | `.kt` `.kts` | `kotlinc` | Compiled* |
| Scala | `.scala` | `scala` | Compiled* |
| Haskell | `.hs` | `ghc` | Compiled* |
| OCaml | `.ml` `.mli` | `ocaml` | Compiled* |
| TypeScript | `.ts` | `ts-node` | Compiled* |
| Dart | `.dart` | `dart` | Compiled* |
| Elixir | `.ex` `.exs` | `elixir` | Interpreted |
| Erlang | `.erl` | `escript` | Interpreted |
| Groovy | `.groovy` `.gvy` | `groovy` | Interpreted |
| R | `.r` `.R` | `Rscript` | Interpreted |
| Julia | `.jl` | `julia` | Interpreted |
| Brainfuck | `.bf` `.b` | interpreter | Interpreted |
| Fortran | `.f` `.f90` `.f95` `.f03` `.f08` `.for` `.fpp` | `gfortran` | Compiled |
| Pascal | `.pas` `.pp` `.lpr` `.dpr` | `fpc` | Compiled |
| Ada | `.adb` `.ads` | `gnatmake` | Compiled |
| COBOL | `.cob` `.cbl` | `cobc` | Compiled |
| Scheme | `.scm` `.ss` | `guile` | Interpreted |
| Common Lisp | `.lisp` `.lsp` `.cl` `.l` | `sbcl --script` | Interpreted |
| Forth | `.fth` `.forth` | `gforth` | Interpreted |
| Prolog | `.pro` `.prolog` | `swipl` | Interpreted |
| Tcl | `.tcl` `.tk` | `tclsh` | Interpreted |
| Clojure | `.clj` `.cljs` `.cljc` `.edn` | `clojure` | Interpreted |
| Gleam | `.gleam` | `gleam run` | Compiled* |
| Elm | `.elm` | `elm make` | Compiled |
| PureScript | `.purs` | `spago run` | Compiled* |
| Roc | `.roc` | `roc run` | Compiled* |
| V | `.v` | `v` | Compiled |
| Odin | `.odin` | `odin build` | Compiled |
| NASM | `.asm` `.nasm` | `nasm + ld` | Compiled |
| GAS | `.s` `.S` | `gcc` | Compiled |
| ARM Assembly | `.arm` | `arm-linux-gnueabi-gcc` | Compiled |
| MIPS Assembly | `.mips` | `mips-linux-gnu-gcc` | Compiled |
| RISC-V Assembly | `.riscv` | `riscv64-linux-gnu-gcc` | Compiled |

**Total: 50 languages supported!**

## Common Usage

### Run a file
```bash
polyglot run hello.py
polyglot run main.go
polyglot run app.js
```

### Update Polyglot
```bash
# Check and install updates
polyglot update

# Only check for updates
polyglot update --check
```

### Pass arguments
```bash
polyglot run script.py --args arg1,arg2,arg3
```

### Quiet mode (for piping)
```bash
polyglot run -q data.py | jq .
```

### Verbose mode
```bash
polyglot run -v hello.py
```

### Check what's installed
```bash
polyglot check
```

### List supported languages
```bash
polyglot list
```

## Error Messages

| Error | Meaning | Solution |
|-------|---------|----------|
| `file not found` | File doesn't exist | Check file path |
| `unsupported file extension` | Language not supported | Check `polyglot list` |
| `executable file not found` | Toolchain not installed | Install language toolchain |
| `program exited with error` | Program failed | Check program code |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error (file not found, unsupported extension, etc.) |
| `N` | Program exit code (propagated from executed program) |

## File Structure

```
polyglot/
├── cmd/polyglot/          # Entry point
├── internal/
│   ├── cli/               # Commands
│   ├── detector/          # Detection
│   ├── executor/          # Execution
│   └── language/          # Handlers
├── pkg/types/             # Types
└── test/fixtures/         # Tests
```

## Development

### Build
```bash
make build
```

### Test
```bash
make test
```

### Install
```bash
make install
```

### Clean
```bash
make clean
```

### Format
```bash
make fmt
```

### Check
```bash
make check
```

## Adding a Language

1. Add to `pkg/types/types.go`
2. Add to `internal/detector/detector.go`
3. Create `internal/language/<lang>.go`
4. Register in `internal/language/registry.go`
5. Add test fixture
6. Update docs

Time: ~15 minutes

## Documentation

| Doc | Purpose |
|-----|---------|
| [README.md](README.md) | Main docs |
| [QUICKSTART.md](docs/QUICKSTART.md) | Get started |
| [DEMO.md](DEMO.md) | Visual demo |
| [EXAMPLES.md](docs/EXAMPLES.md) | Examples |
| [LANGUAGES.md](docs/LANGUAGES.md) | Languages |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribute |
| [GETTING_STARTED.md](GETTING_STARTED.md) | Dev setup |

## Tips & Tricks

### Create an alias
```bash
alias pg='polyglot run'
```

### Use with find
```bash
find . -name "*.py" -exec polyglot run {} \;
```

### In CI/CD
```yaml
- run: polyglot run tests/test.py
```

### With watch
```bash
watch -n 1 polyglot run script.py
```

## Troubleshooting

### Command not found
```bash
# Check installation
which polyglot

# Add to PATH
export PATH="/usr/local/bin:$PATH"
```

### Toolchain not found
```bash
# Check what's installed
polyglot check

# Install missing toolchain
# Ubuntu: sudo apt install python3
# macOS: brew install python3
```

### Permission denied
```bash
# Make file readable
chmod +r file.py
```

## Links

- **Repo**: https://github.com/velo4705/polyglot
- **Issues**: https://github.com/velo4705/polyglot/issues
- **Docs**: [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md)

## Version

**Current**: v1.3.0  
**Milestone**: Production Ready with Auto-Update  
**Status**: Production Ready ✅

---

**Print this page for quick reference!**
