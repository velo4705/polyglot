package types

import "fmt"

// LanguageType represents whether a language is compiled or interpreted
type LanguageType int

const (
	Interpreted LanguageType = iota
	Compiled
)

// Language represents a programming language
type Language int

const (
	Unknown Language = iota
	Python
	Go
	JavaScript
	Java
	C
	Cpp
	Rust
	Ruby
	PHP
	Perl
	Lua
	Zig
	Nim
	Crystal
	D
	Shell
	Haskell
	OCaml
	Elixir
	Erlang
	Kotlin
	Scala
	Groovy
	Swift
	Dart
	TypeScript
	R
	Julia
	Brainfuck
	Fortran
	Pascal
	Ada
	COBOL
	Scheme
	CommonLisp
	Forth
	Prolog
	Tcl
	Clojure
	Gleam
	Elm
	PureScript
	Roc
	V
	Odin
	Nasm
	Gas
	ArmAssembly
	MipsAssembly
	RiscVAssembly
)

type LanguageInfo struct {
	Lang       Language
	Name       string
	Extensions []string
	Type       LanguageType
}

var languageRegistry = []LanguageInfo{
	{Python, "Python", []string{".py"}, Interpreted},
	{Go, "Go", []string{".go"}, Compiled},
	{JavaScript, "JavaScript", []string{".js"}, Interpreted},
	{Java, "Java", []string{".java"}, Compiled},
	{C, "C", []string{".c", ".h"}, Compiled},
	{Cpp, "C++", []string{".cpp", ".cc", ".cxx"}, Compiled},
	{Rust, "Rust", []string{".rs"}, Compiled},
	{Ruby, "Ruby", []string{".rb"}, Interpreted},
	{PHP, "PHP", []string{".php"}, Interpreted},
	{Perl, "Perl", []string{".pl", ".pm"}, Interpreted},
	{Lua, "Lua", []string{".lua"}, Interpreted},
	{Zig, "Zig", []string{".zig"}, Compiled},
	{Nim, "Nim", []string{".nim"}, Compiled},
	{Crystal, "Crystal", []string{".cr"}, Compiled},
	{D, "D", []string{".d"}, Compiled},
	{Shell, "Shell", []string{".sh", ".bash"}, Interpreted},
	{Haskell, "Haskell", []string{".hs", ".lhs"}, Compiled},
	{OCaml, "OCaml", []string{".ml", ".mli"}, Compiled},
	{Elixir, "Elixir", []string{".ex", ".exs"}, Interpreted},
	{Erlang, "Erlang", []string{".erl", ".hrl"}, Compiled},
	{Kotlin, "Kotlin", []string{".kt", ".kts"}, Compiled},
	{Scala, "Scala", []string{".scala", ".sc"}, Compiled},
	{Groovy, "Groovy", []string{".groovy", ".gvy"}, Interpreted},
	{Swift, "Swift", []string{".swift"}, Compiled},
	{Dart, "Dart", []string{".dart"}, Compiled},
	{TypeScript, "TypeScript", []string{".ts", ".tsx"}, Compiled},
	{R, "R", []string{".r", ".R"}, Interpreted},
	{Julia, "Julia", []string{".jl"}, Compiled},
	{Brainfuck, "Brainfuck", []string{".bf", ".b"}, Interpreted},
	{Fortran, "Fortran", []string{".f", ".f90", ".f95", ".f03", ".f08", ".for", ".fpp"}, Compiled},
	{Pascal, "Pascal", []string{".pas", ".pp", ".lpr", ".dpr"}, Compiled},
	{Ada, "Ada", []string{".adb", ".ads"}, Compiled},
	{COBOL, "COBOL", []string{".cob", ".cbl"}, Compiled},
	{Scheme, "Scheme", []string{".scm", ".ss"}, Interpreted},
	{CommonLisp, "Common Lisp", []string{".lisp", ".lsp", ".cl", ".l"}, Interpreted},
	{Forth, "Forth", []string{".fth", ".forth"}, Interpreted},
	{Prolog, "Prolog", []string{".pro", ".prolog"}, Interpreted},
	{Tcl, "Tcl", []string{".tcl", ".tk"}, Interpreted},
	{Clojure, "Clojure", []string{".clj", ".cljs", ".cljc", ".edn"}, Interpreted},
	{Gleam, "Gleam", []string{".gleam"}, Compiled},
	{Elm, "Elm", []string{".elm"}, Compiled},
	{PureScript, "PureScript", []string{".purs"}, Compiled},
	{Roc, "Roc", []string{".roc"}, Compiled},
	{V, "V", []string{".v", ".vsh"}, Compiled},
	{Odin, "Odin", []string{".odin"}, Compiled},
	{Nasm, "NASM", []string{".asm", ".nasm"}, Compiled},
	{Gas, "GAS", []string{".s", ".S"}, Compiled},
	{ArmAssembly, "ARM Assembly", []string{".arm"}, Compiled},
	{MipsAssembly, "MIPS Assembly", []string{".mips"}, Compiled},
	{RiscVAssembly, "RISC-V Assembly", []string{".riscv"}, Compiled},
}

var nameToLanguage map[string]Language
var extToLanguage map[string]Language

func init() {
	nameToLanguage = make(map[string]Language)
	extToLanguage = make(map[string]Language)
	for _, info := range languageRegistry {
		nameToLanguage[info.Name] = info.Lang
		for _, ext := range info.Extensions {
			extToLanguage[ext] = info.Lang
		}
	}
}

func (l Language) String() string {
	return l.Name()
}

func (l Language) Name() string {
	for _, info := range languageRegistry {
		if info.Lang == l {
			return info.Name
		}
	}
	return "Unknown"
}

func (l Language) Extensions() []string {
	for _, info := range languageRegistry {
		if info.Lang == l {
			return info.Extensions
		}
	}
	return nil
}

func (l Language) Type() LanguageType {
	for _, info := range languageRegistry {
		if info.Lang == l {
			return info.Type
		}
	}
	return Interpreted
}

func (l Language) NeedsCompilation() bool {
	return l.Type() == Compiled
}

// String implements fmt.Stringer
func (t LanguageType) String() string {
	if t == Compiled {
		return "Compiled"
	}
	return "Interpreted"
}

// LanguageHandler defines the interface for language-specific operations
type LanguageHandler interface {
	Name() string
	Extensions() []string
	Type() LanguageType
	Compile(source string, output string, extraFlags []string) error
	Run(file string, args []string) ([]byte, error)
	NeedsCompilation() bool
}

// FromString converts a language name string to Language type
func FromString(name string) Language {
	if lang, ok := nameToLanguage[name]; ok {
		return lang
	}
	return Unknown
}

// FromExtension converts a file extension to Language type
func FromExtension(ext string) Language {
	if lang, ok := extToLanguage[ext]; ok {
		return lang
	}
	return Unknown
}

// AllLanguages returns all supported languages
func AllLanguages() []Language {
	langs := make([]Language, 0, len(languageRegistry))
	for _, info := range languageRegistry {
		langs = append(langs, info.Lang)
	}
	return langs
}

// AllLanguageNames returns the display names of all supported languages
func AllLanguageNames() []string {
	names := make([]string, 0, len(languageRegistry))
	for _, info := range languageRegistry {
		names = append(names, info.Name)
	}
	return names
}

// AllExtensions returns all known file extensions
func AllExtensions() []string {
	var exts []string
	for _, info := range languageRegistry {
		exts = append(exts, info.Extensions...)
	}
	return exts
}

// String returns the name of the language
func (l Language) GoString() string {
	return fmt.Sprintf("%s(%d)", l.Name(), l)
}
