package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var polyglotBinary = "../bin/polyglot"

// TestBuildPolyglot ensures the binary is built before running tests
func TestBuildPolyglot(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", polyglotBinary, "../cmd/polyglot")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build polyglot: %v\n%s", err, output)
	}
}

// TestVersion tests the version command
func TestVersion(t *testing.T) {
	cmd := exec.Command(polyglotBinary, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Polyglot") {
		t.Error("Version output should contain 'Polyglot'")
	}
	if !strings.Contains(outputStr, "1.1.0") {
		t.Error("Version output should contain '1.1.0'")
	}
}

// TestList tests the list command
func TestList(t *testing.T) {
	cmd := exec.Command(polyglotBinary, "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("list command failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	languages := []string{"Python", "JavaScript", "Go", "Java", "C", "C++", "Rust"}
	for _, lang := range languages {
		if !strings.Contains(outputStr, lang) {
			t.Errorf("List output should contain '%s'", lang)
		}
	}
}

// TestRunPython tests running a Python file
func TestRunPython(t *testing.T) {
	// Check if python3 is available
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available, skipping test")
	}

	cmd := exec.Command(polyglotBinary, "run", "fixtures/hello.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run command failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Hello") {
		t.Error("Python output should contain 'Hello'")
	}
}

// TestRunJavaScript tests running a JavaScript file
func TestRunJavaScript(t *testing.T) {
	// Check if node is available
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not available, skipping test")
	}

	cmd := exec.Command(polyglotBinary, "run", "fixtures/hello.js")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run command failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Hello") {
		t.Error("JavaScript output should contain 'Hello'")
	}
}

// TestRunGo tests running a Go file
func TestRunGo(t *testing.T) {
	// Check if go is available
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not available, skipping test")
	}

	cmd := exec.Command(polyglotBinary, "run", "fixtures/hello.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run command failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Hello") {
		t.Error("Go output should contain 'Hello'")
	}
}

// TestDryRun tests the dry-run mode
func TestDryRun(t *testing.T) {
	cmd := exec.Command(polyglotBinary, "run", "--dry-run", "fixtures/hello.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("dry-run command failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Dry Run") {
		t.Error("Dry-run output should contain 'Dry Run'")
	}
	if !strings.Contains(outputStr, "Would execute") {
		t.Error("Dry-run output should contain 'Would execute'")
	}
}

// TestFileNotFound tests error handling for missing files
func TestFileNotFound(t *testing.T) {
	cmd := exec.Command(polyglotBinary, "run", "nonexistent.py")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "not found") && !strings.Contains(outputStr, "File not found") {
		t.Error("Error output should mention file not found")
	}
}

// TestUnsupportedExtension tests error handling for unsupported extensions
func TestUnsupportedExtension(t *testing.T) {
	// Create temporary file with unsupported extension
	tmpFile, err := os.CreateTemp("", "test*.xyz")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cmd := exec.Command(polyglotBinary, "run", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for unsupported extension")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unsupported") {
		t.Error("Error output should mention unsupported extension")
	}
}

// TestConfigCommands tests configuration commands
func TestConfigCommands(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"config show", []string{"config", "show"}, false},
		{"config path", []string{"config", "path"}, false},
		{"config get", []string{"config", "get", "auto_install.enabled"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(polyglotBinary, tt.args...)
			output, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Errorf("Command failed: %v\n%s", err, output)
			}
		})
	}
}

// TestArgumentPassing tests passing arguments to programs
func TestArgumentPassing(t *testing.T) {
	// Check if python3 is available
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available, skipping test")
	}

	cmd := exec.Command(polyglotBinary, "run", "fixtures/args.py", "--args", "arg1,arg2,arg3")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run command with args failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "arg1") || !strings.Contains(outputStr, "arg2") {
		t.Error("Output should contain passed arguments")
	}
}

// TestVerboseMode tests verbose output
func TestVerboseMode(t *testing.T) {
	// Check if python3 is available
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available, skipping test")
	}

	cmd := exec.Command(polyglotBinary, "run", "-v", "fixtures/hello.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("verbose run failed: %v\n%s", err, output)
	}

	// Verbose mode should produce output
	if len(output) == 0 {
		t.Error("Verbose mode should produce output")
	}
}

// TestQuietMode tests quiet output
func TestQuietMode(t *testing.T) {
	// Check if python3 is available
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available, skipping test")
	}

	cmd := exec.Command(polyglotBinary, "run", "-q", "fixtures/hello.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("quiet run failed: %v\n%s", err, output)
	}

	outputStr := string(output)
	// Quiet mode should not show polyglot messages
	if strings.Contains(outputStr, "Detected:") || strings.Contains(outputStr, "Executing:") {
		t.Error("Quiet mode should suppress polyglot messages")
	}
}

// TestAllLanguages tests all supported languages (if toolchains are available)
func TestAllLanguages(t *testing.T) {
	languages := map[string]string{
		"python3": "fixtures/hello.py",
		"node":    "fixtures/hello.js",
		"go":      "fixtures/hello.go",
		"ruby":    "fixtures/hello.rb",
		"php":     "fixtures/hello.php",
		"perl":    "fixtures/hello.pl",
		"lua":     "fixtures/hello.lua",
		"bash":    "fixtures/hello.sh",
	}

	for toolchain, file := range languages {
		t.Run(toolchain, func(t *testing.T) {
			// Check if toolchain is available
			if _, err := exec.LookPath(toolchain); err != nil {
				t.Skipf("%s not available, skipping test", toolchain)
			}

			// Check if test file exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Skipf("Test file %s not found, skipping test", file)
			}

			cmd := exec.Command(polyglotBinary, "run", file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("run command failed: %v\n%s", err, output)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "Hello") {
				t.Errorf("Output should contain 'Hello', got: %s", outputStr)
			}
		})
	}
}

// TestCompiledLanguages tests compiled languages (if toolchains are available)
func TestCompiledLanguages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping compiled languages test in short mode")
	}

	languages := map[string]string{
		"gcc":   "fixtures/hello.c",
		"g++":   "fixtures/hello.cpp",
		"rustc": "fixtures/hello.rs",
	}

	for toolchain, file := range languages {
		t.Run(toolchain, func(t *testing.T) {
			// Check if toolchain is available
			if _, err := exec.LookPath(toolchain); err != nil {
				t.Skipf("%s not available, skipping test", toolchain)
			}

			// Check if test file exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Skipf("Test file %s not found, skipping test", file)
			}

			cmd := exec.Command(polyglotBinary, "run", file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("run command failed: %v\n%s", err, output)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "Hello") {
				t.Errorf("Output should contain 'Hello', got: %s", outputStr)
			}

			// Check that compiled binary was cleaned up
			ext := filepath.Ext(file)
			binaryName := strings.TrimSuffix(file, ext)
			if _, err := os.Stat(binaryName); err == nil {
				t.Errorf("Compiled binary %s should have been cleaned up", binaryName)
			}
		})
	}
}
