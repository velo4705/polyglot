package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/velo4705/polyglot/internal/installer"
	"github.com/velo4705/polyglot/internal/ui"
)

// ClassifyError checks if the error message is due to a missing library/dependency.
// If it is, it returns true and a one-line distribution-specific instruction to install it.
// Otherwise, it returns false, "".
// ClassifyError checks if the error message is due to a missing library/dependency.
// If it is, it returns true and a one-line distribution-specific instruction to install it.
// Otherwise, it returns false, "".
func ClassifyError(errStr string, langName string) (bool, string) {
	errLower := strings.ToLower(errStr)

	switch langName {
	case "Python":
		re := regexp.MustCompile(`(?:modulenotfounderror|importerror):\s*no\s+module\s+named\s+['"]([^'"]+)['"]`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: pip install %s", pkg)
		}
		// Also catch generic "no module named"
		if strings.Contains(errLower, "no module named") {
			re2 := regexp.MustCompile(`no\s+module\s+named\s+['"]?([^'"]+)['"]?`)
			if matches := re2.FindStringSubmatch(errLower); len(matches) > 1 {
				pkg := matches[1]
				return true, fmt.Sprintf("Missing library detected. To install, run: pip install %s", pkg)
			}
		}

	case "JavaScript", "TypeScript":
		re := regexp.MustCompile(`cannot\s+find\s+module\s+['"]([^'"]+)['"]`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: npm install %s", pkg)
		}

	case "Ruby":
		re := regexp.MustCompile(`cannot\s+load\s+such\s+file\s+--\s+([^\s']+)`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: gem install %s", pkg)
		}

	case "Go":
		if strings.Contains(errLower, "no required module provides package") {
			re := regexp.MustCompile(`no\s+required\s+module\s+provides\s+package\s+([^\s:]+)`)
			if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
				pkg := matches[1]
				return true, fmt.Sprintf("Missing library detected. To install, run: go get %s", pkg)
			}
		}
		if strings.Contains(errLower, "cannot find package") {
			re := regexp.MustCompile(`cannot\s+find\s+package\s+["']([^"']+)["']`)
			if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
				pkg := matches[1]
				return true, fmt.Sprintf("Missing library detected. To install, run: go get %s", pkg)
			}
		}

	case "Rust":
		if strings.Contains(errLower, "unresolved import") || strings.Contains(errLower, "use of undeclared crate or module") {
			re := regexp.MustCompile("(?:unresolved\\s+import|undeclared\\s+crate\\s+or\\s+module)\\s+[`'\"]?([^`'\"]+)[`'\"]?")
			if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
				pkg := matches[1]
				// Get root crate name if it has subpaths like a::b
				parts := strings.Split(pkg, "::")
				return true, fmt.Sprintf("Missing library detected. To install, run: cargo add %s", parts[0])
			}
		}

	case "PHP":
		re := regexp.MustCompile(`class\s+['"]?([^'"\s]+)['"]?\s+not\s+found`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: composer require %s", strings.ToLower(pkg))
		}

	case "Java":
		re := regexp.MustCompile(`error:\s+package\s+([^\s]+)\s+does\s+not\s+exist`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please add library '%s' to your Maven/Gradle dependencies.", pkg)
		}

	case "Kotlin":
		re := regexp.MustCompile(`error:\s+unresolved\s+reference:\s+([^\s]+)`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please verify package/dependency '%s' is present in your build.gradle.", pkg)
		}

	case "Scala":
		re := regexp.MustCompile(`error:\s+object\s+([^\s]+)\s+is\s+not\s+a\s+member\s+of\s+package`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please add '%s' to your sbt build dependencies.", pkg)
		}

	case "Groovy":
		re := regexp.MustCompile(`unable\s+to\s+resolve\s+class\s+([^\s]+)`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please add package '%s' to your classpath.", pkg)
		}

	case "Swift":
		re := regexp.MustCompile(`error:\s+no\s+such\s+module\s+['"]([^'"]+)['"]`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please add '%s' to your Package.swift dependencies.", pkg)
		}

	case "Dart":
		re := regexp.MustCompile(`error:\s+target\s+of\s+uri\s+doesn't\s+exist:\s*['"]package:([^/]+)/`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. To install, run: dart pub add %s", pkg)
		}

	case "Julia":
		re := regexp.MustCompile(`argumenterror:\s+package\s+([^\s]+)\s+not\s+found`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. To install, run: julia -e 'using Pkg; Pkg.add(\"%s\")'", pkg)
		}

	case "R":
		re := regexp.MustCompile(`there\s+is\s+no\s+package\s+called\s+['"]([^'"]+)['"]`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. To install, run: Rscript -e 'install.packages(\"%s\")'", pkg)
		}

	case "Haskell":
		re := regexp.MustCompile("could\\s+not\\s+find\\s+module\\s+[`'\"]([^`'\"]+)[`'\"]")
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. To install, run: cabal install %s", pkg)
		}

	case "OCaml":
		re := regexp.MustCompile(`error:\s+unbound\s+module\s+([^\s]+)`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. To install, run: opam install %s", pkg)
		}

	case "Elixir":
		re := regexp.MustCompile(`module\s+([^\s]+)\s+is\s+not\s+loaded`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please add '%s' to your deps in mix.exs.", pkg)
		}

	case "Erlang":
		re := regexp.MustCompile(`include_lib\s+failed:\s*"([^"]+)"`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. Please ensure '%s' is present in your Erlang path.", pkg)
		}

	case "F#":
		re := regexp.MustCompile(`error\s+fs0039:\s+the\s+namespace\s+or\s+module\s+['"]([^'"]+)['"]\s+is\s+not\s+defined`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing dependency detected. To install, run: dotnet add package %s", pkg)
		}

	case "Perl":
		re := regexp.MustCompile(`can't\s+locate\s+([^\s]+)\s+in\s+@inc`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			// Convert Foo/Bar.pm to Foo::Bar
			pkg = strings.TrimSuffix(pkg, ".pm")
			pkg = strings.ReplaceAll(pkg, "/", "::")
			return true, fmt.Sprintf("Missing library detected. To install, run: cpan %s", pkg)
		}

	case "Lua":
		re := regexp.MustCompile(`module\s+['"]([^'"]+)['"]\s+not\s+found`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: luarocks install %s", pkg)
		}

	case "Zig":
		re := regexp.MustCompile(`error:\s+import\s+of\s+['"]([^'"]+)['"]\s+failed`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing import detected. Please verify '%s' is configured in your build.zig dependencies.", pkg)
		}

	case "Nim":
		re := regexp.MustCompile(`error:\s+cannot\s+open\s+['"]([^'"]+)['"]`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: nimble install %s", pkg)
		}

	case "Crystal":
		re := regexp.MustCompile(`error:\s+can't\s+find\s+file\s+['"]([^'"]+)['"]`)
		if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. Please add '%s' to shard.yml and run: shards install", pkg)
		}

	case "D":
		re := regexp.MustCompile(`error:\s+module\s+([^\s]+)\s+is\s+missing`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			pkg := matches[1]
			return true, fmt.Sprintf("Missing library detected. To install, run: dub add %s", pkg)
		}

	case "Shell":
		if strings.Contains(errLower, "command not found") {
			re := regexp.MustCompile(`(?:command\s+not\s+found:\s*|:\s*command\s+not\s+found)\s*([^\s]+)`)
			if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
				pkg := matches[1]
				return true, fmt.Sprintf("Missing CLI tool detected. Please install '%s' using your distribution package manager.", pkg)
			}
		}

	case "C", "C++":
		re := regexp.MustCompile(`fatal\s+error:\s*([^\s:]+):\s*no\s+such\s+file\s+or\s+directory`)
		if matches := re.FindStringSubmatch(errLower); len(matches) > 1 {
			header := matches[1]
			libGuess := strings.TrimSuffix(header, filepath.Ext(header))
			libGuess = strings.TrimPrefix(libGuess, "sys/")

			pm := installer.DetectPackageManager()
			if pm != nil {
				switch pm.Name {
				case "apt":
					return true, fmt.Sprintf("Missing header detected. To install, run: sudo apt-get install lib%s-dev", libGuess)
				case "dnf":
					return true, fmt.Sprintf("Missing header detected. To install, run: sudo dnf install %s-devel", libGuess)
				case "brew":
					return true, fmt.Sprintf("Missing header detected. To install, run: brew install %s", libGuess)
				case "pacman":
					return true, fmt.Sprintf("Missing header detected. To install, run: sudo pacman -S %s", libGuess)
				case "apk":
					return true, fmt.Sprintf("Missing header detected. To install, run: sudo apk add %s-dev", libGuess)
				case "zypper":
					return true, fmt.Sprintf("Missing header detected. To install, run: sudo zypper install %s-devel", libGuess)
				}
			}
			return true, fmt.Sprintf("Missing header '%s' detected. Please install the corresponding development package using your distribution's package manager.", header)
		}
	}

	return false, ""
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiCandidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// SelfCorrectFile reads the file, sends it with the error to Gemini API,
// overwrites the file with the corrected content, and returns nil if successful.
func SelfCorrectFile(filename string, errStr string, langName string) error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	// Read the original file
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file for self-correction: %w", err)
	}
	content := string(contentBytes)

	// Formulate prompt
	prompt := fmt.Sprintf(`You are a code fixer for the Polyglot CLI runner.
The user ran a %s file '%s' and got the following execution error:
---
%s
---
Here is the content of the file:
---
%s
---
Please fix the error(s) inside this file.
Return ONLY the corrected, raw code. Do NOT include markdown code blocks, explanations, backticks, or any other wrapper text. Just the fixed code.`, langName, filepath.Base(filename), errStr, content)

	// Construct request body
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshalling gemini request: %w", err)
	}

	ui.Info("Self-correcting error in %s via Gemini API...", ui.File(filename))

	// Send POST request
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("calling gemini api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return fmt.Errorf("decoding gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("received empty response from gemini api")
	}

	fixedCode := geminiResp.Candidates[0].Content.Parts[0].Text

	// Clean up markdown block wrapping if present
	fixedCode = cleanMarkdownWrapping(fixedCode)

	// Overwrite original file
	if err := os.WriteFile(filename, []byte(fixedCode), 0644); err != nil {
		return fmt.Errorf("writing self-corrected file: %w", err)
	}

	ui.Success("File successfully corrected. Re-executing...")
	return nil
}

// cleanMarkdownWrapping strips ```lang ... ``` wrappers if the LLM returned code blocks
func cleanMarkdownWrapping(code string) string {
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "```") {
		// Split by newline
		lines := strings.Split(code, "\n")
		if len(lines) > 2 {
			// Find the matching end tag
			endIdx := -1
			for i := len(lines) - 1; i > 0; i-- {
				if strings.HasPrefix(lines[i], "```") {
					endIdx = i
					break
				}
			}
			if endIdx != -1 {
				// Join middle lines
				return strings.Join(lines[1:endIdx], "\n")
			}
		}
	}
	return code
}
