package detector

import (
	"testing"

	"github.com/velo4705/polyglot/pkg/types"
)

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		file string
		want types.Language
	}{
		{"main.py", types.Python},
		{"main.go", types.Go},
		{"app.js", types.JavaScript},
		{"App.java", types.Java},
		{"main.c", types.C},
		{"main.cpp", types.Cpp},
		{"main.rs", types.Rust},
		{"main.rb", types.Ruby},
		{"main.php", types.PHP},
		{"main.pl", types.Perl},
		{"main.lua", types.Lua},
		{"main.zig", types.Zig},
		{"main.nim", types.Nim},
		{"main.cr", types.Crystal},
		{"main.d", types.D},
		{"main.sh", types.Shell},
		{"main.hs", types.Haskell},
		{"main.ml", types.OCaml},
		{"main.ex", types.Elixir},
		{"main.erl", types.Erlang},
		{"main.kt", types.Kotlin},
		{"main.scala", types.Scala},
		{"main.groovy", types.Groovy},
		{"main.swift", types.Swift},
		{"main.dart", types.Dart},
		{"main.ts", types.TypeScript},
		{"main.r", types.R},
		{"main.jq.jl", types.Julia},
		{"main.bf", types.Brainfuck},
		{"main.f90", types.Fortran},
		{"main.pas", types.Pascal},
		{"main.adb", types.Ada},
		{"main.cob", types.COBOL},
		{"main.scm", types.Scheme},
		{"main.lisp", types.CommonLisp},
		{"main.fth", types.Forth},
		{"main.pro", types.Prolog},
		{"main.tcl", types.Tcl},
		{"main.clj", types.Clojure},
		{"main.gleam", types.Gleam},
		{"main.elm", types.Elm},
		{"main.purs", types.PureScript},
		{"main.roc", types.Roc},
		{"main.v", types.V},
		{"main.odin", types.Odin},
		{"main.asm", types.Nasm},
		{"asm.s", types.Gas},
		{"code.arm", types.ArmAssembly},
		{"code.mips", types.MipsAssembly},
		{"code.riscv", types.RiscVAssembly},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			got, err := DetectLanguage(tc.file)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("file %q: got %s, want %s", tc.file, got.Name(), tc.want.Name())
			}
		})
	}
}

func TestDetectUnknownExtension(t *testing.T) {
	got, err := DetectLanguage("file.unknown_ext")
	if err == nil {
		t.Error("expected error for unknown extension")
	}
	if got != types.Unknown {
		t.Errorf("got %s, want Unknown", got.Name())
	}
}
