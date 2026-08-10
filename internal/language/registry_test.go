package language

import (
	"testing"

	"github.com/velo4705/polyglot/pkg/types"
)

func TestAllHandlersRegistered(t *testing.T) {
	languages := types.AllLanguages()
	if len(languages) != 51 {
		t.Fatalf("expected 51 languages, got %d", len(languages))
	}

	for _, lang := range languages {
		h := GetHandler(lang)
		if h == nil {
			t.Errorf("language %q (%d): handler is nil", lang.Name(), lang)
		}
	}
}

func TestHandlerNames(t *testing.T) {
	cases := []struct {
		lang types.Language
		name string
	}{
		{types.Python, "Python"},
		{types.Go, "Go"},
		{types.JavaScript, "JavaScript"},
		{types.Cpp, "C++"},
		{types.FSharp, "F#"},
		{types.CommonLisp, "Common Lisp"},
		{types.Nasm, "NASM"},
	}

	for _, tc := range cases {
		h := GetHandler(tc.lang)
		if h == nil {
			t.Errorf("handler for %s is nil", tc.name)
			continue
		}
		if h.Name() != tc.name {
			t.Errorf("%s: handler.Name() = %q, want %q", tc.name, h.Name(), tc.name)
		}
	}
}

var expectedTypes = map[types.Language]types.LanguageType{
	types.Python:        types.Interpreted,
	types.Go:            types.Compiled,
	types.JavaScript:    types.Interpreted,
	types.Java:          types.Compiled,
	types.C:             types.Compiled,
	types.Cpp:           types.Compiled,
	types.Rust:          types.Compiled,
	types.Ruby:          types.Interpreted,
	types.PHP:           types.Interpreted,
	types.Perl:          types.Interpreted,
	types.Lua:           types.Interpreted,
	types.Zig:           types.Compiled,
	types.Nim:           types.Compiled,
	types.Crystal:       types.Compiled,
	types.D:             types.Compiled,
	types.Shell:         types.Interpreted,
	types.Haskell:       types.Compiled,
	types.OCaml:         types.Compiled,
	types.Elixir:        types.Interpreted,
	types.Erlang:        types.Compiled,
	types.FSharp:        types.Interpreted,
	types.Kotlin:        types.Compiled,
	types.Scala:         types.Compiled,
	types.Groovy:        types.Interpreted,
	types.Swift:         types.Compiled,
	types.Dart:          types.Compiled,
	types.TypeScript:    types.Compiled,
	types.R:             types.Interpreted,
	types.Julia:         types.Interpreted,
	types.Brainfuck:     types.Interpreted,
	types.Fortran:       types.Compiled,
	types.Pascal:        types.Compiled,
	types.Ada:           types.Compiled,
	types.COBOL:         types.Compiled,
	types.Scheme:        types.Interpreted,
	types.CommonLisp:    types.Interpreted,
	types.Forth:         types.Interpreted,
	types.Prolog:        types.Interpreted,
	types.Tcl:           types.Interpreted,
	types.Clojure:       types.Interpreted,
	types.Gleam:         types.Compiled,
	types.Elm:           types.Compiled,
	types.PureScript:    types.Compiled,
	types.Roc:           types.Compiled,
	types.V:             types.Compiled,
	types.Odin:          types.Compiled,
	types.Nasm:          types.Compiled,
	types.Gas:           types.Compiled,
	types.ArmAssembly:   types.Compiled,
	types.MipsAssembly:  types.Compiled,
	types.RiscVAssembly: types.Compiled,
}

func TestHandlerTypes(t *testing.T) {
	for lang, wantType := range expectedTypes {
		h := GetHandler(lang)
		if h == nil {
			t.Errorf("handler for %s is nil", lang.Name())
			continue
		}
		if gotType := h.Type(); gotType != wantType {
			t.Errorf("%s: Type() = %v, want %v", lang.Name(), gotType, wantType)
		}
	}
}

func TestHandlerExtensions(t *testing.T) {
	for _, lang := range types.AllLanguages() {
		h := GetHandler(lang)
		if h == nil {
			t.Errorf("handler for %s is nil", lang.Name())
			continue
		}
		exts := h.Extensions()
		if len(exts) == 0 {
			t.Errorf("%s: no extensions defined", lang.Name())
		}
	}
}

func TestNoHandlerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	h := GetHandler(types.Unknown)
	if h != nil {
		t.Error("expected nil handler for Unknown language")
	}
}
