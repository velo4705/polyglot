# Supported Languages

Polyglot supports 30 programming languages across interpreted, compiled, and esoteric categories.

## Overview

| Category | Count | Languages |
|----------|-------|-----------|
| Interpreted | 16 | Python, JavaScript, Ruby, PHP, Perl, Lua, Shell, Elixir, Erlang, Groovy, R, Julia, Haskell, OCaml, TypeScript, Dart |
| Compiled | 13 | Go, Java, C, C++, Rust, Zig, Nim, Crystal, D, Kotlin, Scala, Swift |
| Esoteric | 1 | Brainfuck |
| **Total** | **30** | |

## Interpreted Languages

| Language | Extensions | Command | Status |
|----------|-----------|---------|--------|
| Python | `.py` | `python3` | ✅ |
| JavaScript | `.js` | `node` | ✅ |
| Ruby | `.rb` | `ruby` | ✅ |
| PHP | `.php` | `php` | ✅ |
| Perl | `.pl` | `perl` | ✅ |
| Lua | `.lua` | `lua` | ✅ |
| Shell | `.sh` `.bash` | `bash` | ✅ |
| Elixir | `.ex` `.exs` | `elixir` | ✅ |
| Erlang | `.erl` | `escript` | ✅ |
| Groovy | `.groovy` | `groovy` | ✅ |
| R | `.r` `.R` | `Rscript` | ✅ |
| Julia | `.jl` | `julia` | ✅ |
| Haskell | `.hs` | `runhaskell` | ✅ |
| OCaml | `.ml` `.mli` | `ocaml` | ✅ |
| TypeScript | `.ts` | `ts-node` | ✅ |
| Dart | `.dart` | `dart` | ✅ |

## Compiled Languages

| Language | Extensions | Compiler | Status |
|----------|-----------|----------|--------|
| Go | `.go` | `go run` | ✅ |
| Java | `.java` | `javac` + `java` | ✅ |
| C | `.c` | `gcc` | ✅ |
| C++ | `.cpp` `.cc` `.cxx` `.C` | `g++` | ✅ |
| Rust | `.rs` | `rustc` | ✅ |
| Zig | `.zig` | `zig run` | ✅ |
| Nim | `.nim` | `nim c` | ✅ |
| Crystal | `.cr` | `crystal run` | ✅ |
| D | `.d` | `dmd` | ✅ |
| Kotlin | `.kt` `.kts` | `kotlinc` + `kotlin` | ✅ |
| Scala | `.scala` | `scalac` + `scala` | ✅ |
| Swift | `.swift` | `swift` | ✅ |

## Esoteric Languages

| Language | Extensions | Interpreter | Status |
|----------|-----------|-------------|--------|
| Brainfuck | `.bf` `.b` | `bf` | ✅ |

## Testing Language Support

To check which languages are available on your system:

```bash
polyglot check
```

To see all languages polyglot can detect:

```bash
polyglot list
```

## Adding New Languages

See [CONTRIBUTING.md](../CONTRIBUTING.md) for instructions on adding support for new languages.

## Language-Specific Notes

### Python
- Uses `python3` by default
- Supports all Python 3.x versions

### Go
- Uses `go run` for quick execution
- Use `polyglot compile` for standalone binaries

### JavaScript
- Requires Node.js
- Does not support browser-specific APIs

### Ruby
- Requires Ruby interpreter
- Tested with Ruby 2.7+

### PHP
- Requires PHP CLI
- Tested with PHP 7.4+

## Requesting New Languages

To request support for a new language, please open an issue with:
- Language name
- File extension(s)
- Compiler/interpreter command
- Example hello world program
- Use case / why you need it
