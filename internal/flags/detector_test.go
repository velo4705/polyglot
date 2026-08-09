package flags

import (
	"regexp"
	"testing"
)

func TestDetectorDefaults(t *testing.T) {
	d := &Detector{
		Defaults: []string{"-Wall", "-Wextra"},
		Rules:    nil,
	}
	flags := d.Detect([]byte("int main() {}"))
	if len(flags) != 2 {
		t.Errorf("expected 2 default flags, got %d: %v", len(flags), flags)
	}
}

func TestDetectorRuleMatch(t *testing.T) {
	d := &Detector{
		Defaults: []string{"-Wall"},
		Rules: []Rule{
			{regexp.MustCompile(`pthread_create`), []string{"-pthread"}},
		},
	}
	flags := d.Detect([]byte("pthread_create(&t, NULL, func, NULL)"))
	found := false
	for _, f := range flags {
		if f == "-pthread" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected -pthread in flags, got %v", flags)
	}
}

func TestDetectorRuleNoMatch(t *testing.T) {
	d := &Detector{
		Defaults: []string{"-Wall"},
		Rules: []Rule{
			{regexp.MustCompile(`pthread_create`), []string{"-pthread"}},
		},
	}
	flags := d.Detect([]byte("int main() { return 0; }"))
	for _, f := range flags {
		if f == "-pthread" {
			t.Errorf("should not have -pthread in flags for code without pthread")
		}
	}
}

func TestDetectorDeduplication(t *testing.T) {
	d := &Detector{
		Defaults: []string{"-Wall"},
		Rules: []Rule{
			{regexp.MustCompile(`std::thread`), []string{"-pthread"}},
			{regexp.MustCompile(`#include\s*<thread>`), []string{"-pthread"}},
		},
	}
	flags := d.Detect([]byte(`
		#include <thread>
		std::thread t;
	`))
	count := 0
	for _, f := range flags {
		if f == "-pthread" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected -pthread once (deduplicated), got %d times in %v", count, flags)
	}
}

func TestDetectorMultipleRules(t *testing.T) {
	d := &Detector{
		Defaults: []string{"-Wall"},
		Rules: []Rule{
			{regexp.MustCompile(`std::thread`), []string{"-pthread"}},
			{regexp.MustCompile(`std::optional`), []string{"-std=c++17"}},
			{regexp.MustCompile(`std::format`), []string{"-std=c++20"}},
		},
	}
	flags := d.Detect([]byte(`
		#include <thread>
		#include <optional>
		std::thread t;
		std::optional<int> v;
		std::format("hi");
	`))
	flagSet := make(map[string]bool)
	for _, f := range flags {
		flagSet[f] = true
	}
	if !flagSet["-Wall"] {
		t.Error("missing -Wall")
	}
	if !flagSet["-pthread"] {
		t.Error("missing -pthread")
	}
	if !flagSet["-std=c++17"] {
		t.Error("missing -std=c++17")
	}
	if !flagSet["-std=c++20"] {
		t.Error("missing -std=c++20")
	}
}

func TestRegistry(t *testing.T) {
	d := &Detector{Defaults: []string{"-test"}}
	Register("TestLang", d)

	got := Get("TestLang")
	if got == nil {
		t.Fatal("Get returned nil for registered language")
	}
	if len(got.Defaults) != 1 || got.Defaults[0] != "-test" {
		t.Errorf("unexpected defaults: %v", got.Defaults)
	}
}

func TestRegistryMiss(t *testing.T) {
	got := Get("NonExistentLang")
	if got != nil {
		t.Error("expected nil for unregistered language")
	}
}

func TestCppDetector(t *testing.T) {
	d := Get("C++")
	if d == nil {
		t.Fatal("C++ detector not registered")
	}

	tests := []struct {
		name     string
		source   string
		mustHave []string
	}{
		{
			name:     "threading",
			source:   `std::thread t;`,
			mustHave: []string{"-pthread"},
		},
		{
			name:     "C++17 optional",
			source:   `std::optional<int> v;`,
			mustHave: []string{"-std=c++17"},
		},
		{
			name:     "C++20 format",
			source:   `std::format("hi");`,
			mustHave: []string{"-std=c++20"},
		},
		{
			name:     "C++20 concepts",
			source:   `template<typename T> concept Sortable = requires(T a) { a < a; };`,
			mustHave: []string{"-std=c++20"},
		},
		{
			name:     "C++20 coroutines",
			source:   `co_await something;`,
			mustHave: []string{"-std=c++20"},
		},
		{
			name:     "OpenMP",
			source:   `#pragma omp parallel for`,
			mustHave: []string{"-fopenmp"},
		},
		{
			name:     "SIMD",
			source:   `#include <immintrin.h>`,
			mustHave: []string{"-msse4.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := d.Detect([]byte(tt.source))
			flagSet := make(map[string]bool)
			for _, f := range flags {
				flagSet[f] = true
			}
			for _, want := range tt.mustHave {
				if !flagSet[want] {
					t.Errorf("source %q: expected %q in flags %v", tt.source, want, flags)
				}
			}
		})
	}
}

func TestCDetector(t *testing.T) {
	d := Get("C")
	if d == nil {
		t.Fatal("C detector not registered")
	}

	tests := []struct {
		name     string
		source   string
		mustHave []string
	}{
		{
			name:     "pthread",
			source:   `#include <pthread.h>`,
			mustHave: []string{"-pthread"},
		},
		{
			name:     "math",
			source:   `#include <math.h>`,
			mustHave: []string{"-lm"},
		},
		{
			name:     "atomic",
			source:   `#include <stdatomic.h>`,
			mustHave: []string{"-std=c11"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := d.Detect([]byte(tt.source))
			flagSet := make(map[string]bool)
			for _, f := range flags {
				flagSet[f] = true
			}
			for _, want := range tt.mustHave {
				if !flagSet[want] {
					t.Errorf("source %q: expected %q in flags %v", tt.source, want, flags)
				}
			}
		})
	}
}

func TestRustDetector(t *testing.T) {
	d := Get("Rust")
	if d == nil {
		t.Fatal("Rust detector not registered")
	}

	// Default should be edition 2021
	flags := d.Detect([]byte(`fn main() {}`))
	flagSet := make(map[string]bool)
	for _, f := range flags {
		flagSet[f] = true
	}
	if !flagSet["--edition"] {
		t.Error("missing --edition in defaults")
	}
	if !flagSet["2021"] {
		t.Error("missing 2021 edition in defaults")
	}
}

func TestEmptySource(t *testing.T) {
	d := Get("C++")
	if d == nil {
		t.Fatal("C++ detector not registered")
	}
	flags := d.Detect([]byte{})
	// Should still have defaults
	if len(flags) != 2 {
		t.Errorf("expected 2 default flags for empty source, got %d: %v", len(flags), flags)
	}
}
