# Polyglot - Frequently Asked Questions

## General Questions

### What is Polyglot?

Polyglot is a universal compiler and interpreter wrapper that automatically detects programming languages and executes them with the appropriate toolchain using a single, consistent command.

Instead of remembering different commands for different languages, just use:
```bash
polyglot run <any-file>
```

### How many languages does it support?

Currently **51 languages**:
- **Interpreted:** Python, JavaScript, Ruby, PHP, Perl, Lua, Shell, R, Julia, Elixir, Scheme, Common Lisp, Forth, Prolog, Tcl, Clojure, Brainfuck, F#
- **Compiled:** Go, Java, C, C++, Rust, Zig, Nim, Crystal, D, Haskell, OCaml, Erlang, Kotlin, Scala, Swift, Dart, TypeScript, Elm, Gleam, PureScript, Roc, V, Odin, Fortran, Pascal, Ada, COBOL
- **Assembly:** NASM, GAS, ARM Assembly, MIPS Assembly, RISC-V Assembly

> **Note:** Some "compiled" languages (e.g., Go, Zig, Crystal, Kotlin, Scala, Swift, Dart, TypeScript, V, Nim) support run-in-one-step execution via their toolchain's built-in runner (e.g., `go run`, `zig run`), so no separate compile step is needed. See `polyglot list` for full details including required toolchains.

### Is it free?

Yes! Polyglot is open source under the MIT License.

---

## Installation & Setup

### Do I need to install all 51 language toolchains?

**Nope,** You only need to install the languages you want to use.

Polyglot is a **wrapper tool** - it uses your system's installed compilers and interpreters. It doesn't bundle them.

```bash
# Check what you have installed
polyglot check

# Install only what you need
sudo dnf install python3 gcc nodejs  # Example
```

### Why doesn't Polyglot bundle the compilers?

**Current Design (Wrapper Model):**
- ✅ Small binary (~9MB)
- ✅ Uses system-installed, up-to-date compilers
- ✅ Users control versions
- ✅ No licensing issues
- ✅ Respects system package management

**Yes!** Polyglot can auto-install missing toolchains using your system package manager.

When you try to run code in a language that's not installed, Polyglot will offer to install it for you.


### Can Polyglot auto-install missing languages?

**Yes!** Polyglot will prompt you to install any missing toolchain when you try to run a file. You can also install manually with `polyglot install <language>`.

```bash
$ polyglot run hello.py
⚠️  Python not found in your system

Polyglot can install Python for you:
  • Installation method: System package manager (dnf)
  • Package: python3
  • Will be installed system-wide (available to all programs)
  • Requires: sudo password

Install Python now? [Y/n]
```

Polyglot detects your system's package manager (apt, dnf, brew, pacman, zypper, apk) and installs the correct package automatically. Alternatively, install manually:

```bash
polyglot install python
```

### How do I install language toolchains?

**Fedora/RHEL:**
```bash
sudo dnf install python3 nodejs gcc g++ rustc ruby php perl lua
```

**Ubuntu/Debian:**
```bash
sudo apt install python3 nodejs build-essential rustc ruby php perl lua
```

**macOS:**
```bash
brew install python3 node gcc rust ruby php perl lua
```

Then verify:
```bash
polyglot check
```

---

## Usage Questions

### How do I run a program?

```bash
polyglot run <file>
```

Examples:
```bash
polyglot run hello.py
polyglot run main.go
polyglot run app.js
polyglot run hello.c
```

### How do I pass arguments to my program?

```bash
polyglot run script.py --args hello --args world
polyglot run script.py --args "hello world"
```

### How do I suppress Polyglot's messages?

Use quiet mode:
```bash
polyglot run -q script.py
```

Perfect for piping:
```bash
polyglot run -q data.py | jq .
```

### How do I see detailed execution info?

Use verbose mode:
```bash
polyglot run -v hello.py
```

### What if my language isn't supported?

1. Check if it's planned: [ROADMAP.md](ROADMAP.md)
2. Request it: Open an issue
3. Add it yourself: See [CONTRIBUTING.md](CONTRIBUTING.md) - takes ~15 minutes!

---

## Compilation Questions

### Does Polyglot compile my code?

**For compiled languages (C, C++, Rust, Java):** Yes!

Polyglot automatically:
1. Compiles your code
2. Runs the compiled binary
3. Cleans up the binary after execution

```bash
$ polyglot run hello.c
ℹ Detected: C
→ Compiling... ✓
→ Executing: ./hello
Hello from C!
# Binary automatically cleaned up
```

### Where do compiled binaries go?

They're created in the same directory as your source file and automatically deleted after execution.

If you want to keep the binary, use:
```bash
polyglot compile hello.c
# Creates 'hello' binary without running it
```

### Can I see the compilation command?

Use verbose mode:
```bash
polyglot run -v hello.c
```

---

## Technical Questions

### How does language detection work?

Polyglot looks at the file extension:
- `.py` → Python
- `.js` → JavaScript
- `.c` → C
- `.rs` → Rust
- etc.

See the full list:
```bash
polyglot list
```

### What if I have multiple versions of a language?

Polyglot uses whatever is in your PATH:
```bash
which python3  # Shows which Python Polyglot will use
```

### Does Polyglot work on Windows?

**Yes!** Polyglot supports Windows, macOS, and Linux.

### How big is the Polyglot binary?

About 13MB. It's a single, self-contained binary with no dependencies (except the language toolchains you want to use).

### Is Polyglot fast?

Yes! Polyglot adds minimal overhead:
- Detection: < 1ms
- Execution: Direct process spawn (no overhead)
- Compilation: Same speed as running the compiler directly

---

## Self‑Correct API

### What is Self‑Correct?

The `--self-correct` flag enables Polyglot to automatically fix syntax or compilation errors by sending the source code to a language model. You must provide your own API key for the LLM you want to use.

### How do I provide my API key?

Set one of the supported environment variables before running Polyglot:

- **Gemini:** `export GEMINI_API_KEY=your_key`
- **OpenAI:** `export OPENAI_API_KEY=your_key`
- **Groq:** `export GROQ_API_KEY=your_key`
- **Anthropic:** `export ANTHROPIC_API_KEY=your_key`

If more than one key is present, specify the provider with `--provider <gemini|openai|groq|anthropic>`.

### How to use Self‑Correct

```bash
# Basic usage (Gemini)
export GEMINI_API_KEY=your_key
polyglot run broken.py --self-correct

# Specify provider explicitly
polyglot run broken.py --self-correct --provider openai
```

The command runs the program, and if a syntax error is detected, Polyglot sends the file to the selected LLM, receives a corrected version, rewrites the file, and re‑executes it automatically.

### Supported providers

- Gemini (default)
- OpenAI
- Groq

### Limitations

- Only syntax/file‑level errors are auto‑corrected.
- The feature is opt‑in; without `--self-correct` the program runs normally.
- Provider support beyond Gemini is currently a stub and will return “not yet implemented” until fully wired.

### How do I enable the other providers?

Each provider requires a specific environment variable:

| Provider | Env Var | Notes |
|----------|---------|-------|
| Gemini | `GEMINI_API_KEY` | Fully implemented (gemini-2.5-flash) |
| OpenAI | `OPENAI_API_KEY` | Not yet implemented |
| Groq | `GROQ_API_KEY` | Not yet implemented |
| Anthropic | `ANTHROPIC_API_KEY` | Not yet implemented |



## Comparison Questions

### How is this different from Docker?

**Polyglot:**
- ✅ Lightweight wrapper
- ✅ Uses system toolchains
- ✅ Fast startup
- ✅ No containerization overhead

**Docker:**
- ✅ Complete isolation
- ✅ Reproducible environments
- ❌ Slower startup
- ❌ More complex setup

### How is this different from asdf/mise?

**Polyglot:**
- ✅ Automatic language detection
- ✅ Single command for all languages
- ✅ No configuration needed
- ❌ No version management (yet)

**asdf/mise:**
- ✅ Version management
- ✅ Per-project versions
- ❌ Requires configuration
- ❌ Separate commands per language

**Future:** We're planning to combine the best of both!

### How is this different from Nix?

**Polyglot:**
- ✅ Simple wrapper
- ✅ Easy to understand
- ✅ Works with system packages
- ❌ No reproducibility guarantees

**Nix:**
- ✅ Reproducible builds
- ✅ Declarative configuration
- ❌ Steep learning curve
- ❌ Complex setup

---

## Contributing Questions

### How do I add a new language?

It takes about 15 minutes! See [GETTING_STARTED.md](docs/development/GETTING_STARTED.md) for a step-by-step tutorial.

Basic steps:
1. Add language constant to `pkg/types/types.go` (in the `Language` enum and `LanguageInfo` registry)
2. Create handler in `internal/language/<lang>.go`
3. Register in `internal/language/registry.go`
4. Add compile-flag detector in `internal/flags/<lang>.go` (if applicable)
5. Add test for extension detection
6. Update docs

### Can I contribute?

Yes! We welcome contributions. See [CONTRIBUTING.md](CONTRIBUTING.md).

Good first issues:
- Add a new language
- Improve error messages
- Add more examples
- Improve documentation

### How is the project organized?

See [PROJECT_STRUCTURE.md](docs/development/PROJECT_STRUCTURE.md) for a complete overview.

---

## Troubleshooting

### "command not found: polyglot"

Polyglot isn't in your PATH. Either:
```bash
# Use full path
./bin/polyglot run hello.py

# Or install it
make install
```

### "unsupported file extension"

The language isn't supported yet. Check:
```bash
polyglot list
```

### "executable file not found in $PATH"

The language toolchain isn't installed. Check:
```bash
polyglot check
```

Then install the missing toolchain.

### Compilation fails

Check if the compiler is installed:
```bash
which gcc    # For C
which rustc  # For Rust
which javac  # For Java
```

Use verbose mode to see the exact error:
```bash
polyglot run -v hello.c
```

---

## Future Plans

### What's coming next?

See [ROADMAP.md](ROADMAP.md) for future plans.

**Upcoming features:**
- Multi-file project support
- Build system integration
- Dependency management
- Container-based execution
- Performance optimizations

### Will Polyglot support my favorite language?

Probably! Check the [ROADMAP.md](ROADMAP.md) or open an issue to request it.

---

## Getting Help

### Where can I get help?

- 📖 Read the [Quick Start Guide](docs/QUICKSTART.md)
- 💬 Open an [Issue](https://github.com/velo4705/polyglot/issues)
- 📚 Check [Examples](docs/EXAMPLES.md)
- 🎯 See the [Demo](docs/development/DEMO.md)

### How do I report a bug?

Open an issue with:
- Polyglot version (`polyglot version`)
- Operating system
- Language/toolchain versions
- Steps to reproduce
- Expected vs actual behavior

### How do I request a feature?

Open an issue describing:
- What you want to do
- Why it's useful
- How you envision it working

---

## Philosophy

### Why build Polyglot?

**Problem:** Developers working with multiple languages face friction:
- Different commands for different languages
- Need to remember language-specific flags
- Inconsistent interfaces

**Solution:** One command to run them all!
```bash
polyglot run <any-file>
```

### What's the vision?

Make it effortless to run programs in any programming language with a single, consistent command.

**Current:** Wrapper for system toolchains  
**Future:** Smart toolchain manager with auto-download  
**Ultimate:** Universal development environment

---

## More Questions?

Open an issue or check the documentation:
- [README.md](README.md)
- [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md)
- [All Docs](docs/)
