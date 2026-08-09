# Polyglot

A universal compiler and interpreter wrapper. Detects language by file extension and runs it with the right toolchain.

Supports **51 languages** — Python, JavaScript, Go, Java, C, C++, Rust, Ruby, PHP, and more.
Implements Auto-Correction of Syntax errors and Auto-Applying Compilation Flags depending on Code.

## Install

### One-Liner Installation
```bash
curl -sSL https://raw.githubusercontent.com/velo4705/polyglot/main/scripts/install.sh | bash
```

### Source Installation
```bash
git clone https://github.com/velo4705/polyglot.git && cd polyglot
make build && sudo make install
```

## Usage

```bash
polyglot run hello.py          # run any file
polyglot watch app.js          # auto-rerun on changes
polyglot compile hello.cpp     # compile without running
polyglot list                  # show supported languages
polyglot check                 # check installed toolchains
polyglot update                # self-update
```

### Compiler flags

Polyglot auto-detects required compiler flags from your source code:

```bash
# Flags detected automatically (e.g. -std=c++17, -pthread, -lm)
polyglot run hello.cpp

# Override or add flags manually
polyglot run --compile-flags "-O3 -march=native" hello.cpp
```

### Dry run

```bash
polyglot run --dry-run hello.c   # see what would execute without running
```

### Stdin

```bash
echo 'print("hello")' | polyglot run --lang Python
```

### Arguments

```bash
polyglot run --args "arg1 arg2" hello.py
```

## Supported Languages

**Interpreted:** Python, JavaScript, Ruby, PHP, Perl, Lua, Shell, TypeScript, Dart, Elixir, Erlang, Groovy, R, Julia, Scheme, Common Lisp, Forth, Prolog, Tcl, Clojure, Gleam, PureScript, Roc

**Compiled:** Go, Java, C, C++, Rust, Zig, Nim, Crystal, D, Swift, Kotlin, Scala, Haskell, OCaml, F#, Fortran, Pascal, Ada, COBOL, V, Odin, Elm

**Assembly:** NASM, GAS, ARM Assembly, MIPS Assembly, RISC-V Assembly

**Esoteric:** Brainfuck

**Total: 51 languages** with 55+ file extensions

## Configuration

Create `.polyglot.yaml` in your project:

```yaml
auto_install:
  enabled: true
  prompt: false

languages:
  python:
    command: python3.11

custom_extensions:
  .pyx: Python
  .jsx: JavaScript

sandbox:
  enabled: true
  timeout: 30
```

Or use the CLI:

```bash
polyglot config init
polyglot config set auto_install.enabled true
polyglot config show
```

## Docker

```bash
docker build -t polyglot .
docker run -v $(pwd):/workspace polyglot run hello.py
```

## Testing

```bash
go test -race ./pkg/... ./internal/... -v
```

## License

MIT
