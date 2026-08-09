package installer

import (
	"os/exec"
	"runtime"
)

// PackageManager represents a system package manager
type PackageManager struct {
	Name       string
	Command    string
	InstallCmd []string
	UpdateCmd  []string
	SearchCmd  []string
	Available  bool
}

// DetectPackageManager detects the available package manager on the system
func DetectPackageManager() *PackageManager {
	managers := []*PackageManager{
		{
			Name:       "dnf",
			Command:    "dnf",
			InstallCmd: []string{"sudo", "dnf", "install", "-y"},
			UpdateCmd:  []string{"sudo", "dnf", "update"},
			SearchCmd:  []string{"dnf", "search"},
		},
		{
			Name:       "apt",
			Command:    "apt-get",
			InstallCmd: []string{"sudo", "apt-get", "install", "-y"},
			UpdateCmd:  []string{"sudo", "apt-get", "update"},
			SearchCmd:  []string{"apt-cache", "search"},
		},
		{
			Name:       "brew",
			Command:    "brew",
			InstallCmd: []string{"brew", "install"},
			UpdateCmd:  []string{"brew", "update"},
			SearchCmd:  []string{"brew", "search"},
		},
		{
			Name:       "pacman",
			Command:    "pacman",
			InstallCmd: []string{"sudo", "pacman", "-S", "--noconfirm"},
			UpdateCmd:  []string{"sudo", "pacman", "-Sy"},
			SearchCmd:  []string{"pacman", "-Ss"},
		},
		{
			Name:       "zypper",
			Command:    "zypper",
			InstallCmd: []string{"sudo", "zypper", "install", "-y"},
			UpdateCmd:  []string{"sudo", "zypper", "refresh"},
			SearchCmd:  []string{"zypper", "search"},
		},
		{
			Name:       "apk",
			Command:    "apk",
			InstallCmd: []string{"sudo", "apk", "add"},
			UpdateCmd:  []string{"sudo", "apk", "update"},
			SearchCmd:  []string{"apk", "search"},
		},
	}

	// Check each package manager
	for _, mgr := range managers {
		if _, err := exec.LookPath(mgr.Command); err == nil {
			mgr.Available = true
			return mgr
		}
	}

	return nil
}

// GetPackageName returns the package name for a language on a specific package manager
func GetPackageName(language, pkgManager string) string {
	packageMap := map[string]map[string]string{
		"python": {
			"dnf":    "python3",
			"apt":    "python3",
			"brew":   "python3",
			"pacman": "python",
			"zypper": "python3",
			"apk":    "python3",
		},
		"node": {
			"dnf":    "nodejs",
			"apt":    "nodejs",
			"brew":   "node",
			"pacman": "nodejs",
			"zypper": "nodejs",
			"apk":    "nodejs",
		},
		"gcc": {
			"dnf":    "gcc",
			"apt":    "build-essential",
			"brew":   "gcc",
			"pacman": "gcc",
			"zypper": "gcc",
			"apk":    "gcc",
		},
		"g++": {
			"dnf":    "gcc-c++",
			"apt":    "build-essential",
			"brew":   "gcc",
			"pacman": "gcc",
			"zypper": "gcc-c++",
			"apk":    "g++",
		},
		"rustc": {
			"dnf":    "rust",
			"apt":    "rustc",
			"brew":   "rust",
			"pacman": "rust",
			"zypper": "rust",
			"apk":    "rust",
		},
		"javac": {
			"dnf":    "java-latest-openjdk-devel",
			"apt":    "default-jdk",
			"brew":   "openjdk",
			"pacman": "jdk-openjdk",
			"zypper": "java-openjdk-devel",
			"apk":    "openjdk11",
		},
		"ruby": {
			"dnf":    "ruby",
			"apt":    "ruby",
			"brew":   "ruby",
			"pacman": "ruby",
			"zypper": "ruby",
			"apk":    "ruby",
		},
		"php": {
			"dnf":    "php",
			"apt":    "php",
			"brew":   "php",
			"pacman": "php",
			"zypper": "php",
			"apk":    "php",
		},
		"perl": {
			"dnf":    "perl",
			"apt":    "perl",
			"brew":   "perl",
			"pacman": "perl",
			"zypper": "perl",
			"apk":    "perl",
		},
		"lua": {
			"dnf":    "lua",
			"apt":    "lua5.4",
			"brew":   "lua",
			"pacman": "lua",
			"zypper": "lua",
			"apk":    "lua",
		},
		"go": {
			"dnf":    "golang",
			"apt":    "golang",
			"brew":   "go",
			"pacman": "go",
			"zypper": "go",
			"apk":    "go",
		},
		"ghc": {
			"dnf":    "ghc",
			"apt":    "ghc",
			"brew":   "ghc",
			"pacman": "ghc",
			"zypper": "ghc",
			"apk":    "ghc",
		},
		"elixir": {
			"dnf":    "elixir",
			"apt":    "elixir",
			"brew":   "elixir",
			"pacman": "elixir",
			"zypper": "elixir",
			"apk":    "elixir",
		},
		"julia": {
			"dnf":    "julia",
			"apt":    "julia",
			"brew":   "julia",
			"pacman": "julia",
			"zypper": "julia",
			"apk":    "julia",
		},
		"Rscript": {
			"dnf":    "R",
			"apt":    "r-base",
			"brew":   "r",
			"pacman": "r",
			"zypper": "R-base",
			"apk":    "R",
		},
		"dart": {
			"dnf":    "dart",
			"apt":    "dart",
			"brew":   "dart",
			"pacman": "dart",
			"zypper": "dart",
			"apk":    "dart",
		},
		"kotlin": {
			"dnf":    "kotlin",
			"apt":    "kotlin",
			"brew":   "kotlin",
			"pacman": "kotlin",
			"zypper": "kotlin",
			"apk":    "kotlin",
		},
		"scala": {
			"dnf":    "scala",
			"apt":    "scala",
			"brew":   "scala",
			"pacman": "scala",
			"zypper": "scala",
			"apk":    "scala",
		},
		"groovy": {
			"dnf":    "groovy",
			"apt":    "groovy",
			"brew":   "groovy",
			"pacman": "groovy",
			"zypper": "groovy",
			"apk":    "groovy",
		},
		"swift": {
			"dnf":    "swift-lang",
			"apt":    "swift",
			"brew":   "swift",
			"pacman": "swift",
			"zypper": "swift",
			"apk":    "swift",
		},
		"ts-node": {
			"dnf":    "nodejs",
			"apt":    "nodejs",
			"brew":   "node",
			"pacman": "nodejs",
			"zypper": "nodejs",
			"apk":    "nodejs",
		},
		"gfortran": {
			"dnf":    "gcc-gfortran",
			"apt":    "gfortran",
			"brew":   "gcc",
			"pacman": "gcc-fortran",
			"zypper": "gcc-fortran",
			"apk":    "gfortran",
		},
		"fpc": {
			"dnf":    "fpc",
			"apt":    "fp-compiler",
			"brew":   "fpc",
			"pacman": "fpc",
			"zypper": "fpc",
			"apk":    "fpc",
		},
		"gnatmake": {
			"dnf":    "gcc-gnat",
			"apt":    "gnat",
			"brew":   "gcc",
			"pacman": "gcc-ada",
			"zypper": "gcc13-ada",
			"apk":    "gnat",
		},
		"cobc": {
			"dnf":    "open-cobol",
			"apt":    "open-cobol",
			"brew":   "gnu-cobol",
			"pacman": "open-cobol",
			"zypper": "open-cobol",
			"apk":    "open-cobol",
		},
		"guile": {
			"dnf":    "guile-scheme",
			"apt":    "guile-3.0-libs",
			"brew":   "guile",
			"pacman": "guile",
			"zypper": "guile",
			"apk":    "guile",
		},
		"sbcl": {
			"dnf":    "sbcl",
			"apt":    "sbcl",
			"brew":   "sbcl",
			"pacman": "sbcl",
			"zypper": "sbcl",
			"apk":    "sbcl",
		},
		"gforth": {
			"dnf":    "gforth",
			"apt":    "gforth",
			"brew":   "gforth",
			"pacman": "gforth",
			"zypper": "gforth",
			"apk":    "gforth",
		},
		"swipl": {
			"dnf":    "pl",
			"apt":    "swi-prolog",
			"brew":   "swi-prolog",
			"pacman": "swi-prolog",
			"zypper": "pl",
			"apk":    "swi-prolog",
		},
		"tclsh": {
			"dnf":    "tcl",
			"apt":    "tcl",
			"brew":   "tcl-tk",
			"pacman": "tcl",
			"zypper": "tcl",
			"apk":    "tcl",
		},
		"clojure": {
			"dnf":    "clojure",
			"apt":    "clojure",
			"brew":   "clojure",
			"pacman": "clojure",
			"zypper": "clojure",
			"apk":    "clojure",
		},
		"gleam": {
			"dnf":    "gleam",
			"apt":    "gleam",
			"brew":   "gleam",
			"pacman": "gleam",
			"zypper": "gleam",
			"apk":    "gleam",
		},
		"elm": {
			"dnf":    "nodejs-npm",
			"apt":    "npm",
			"brew":   "elm",
			"pacman": "elm",
			"zypper": "nodejs-npm",
			"apk":    "elm",
		},
		"spago": {
			"dnf":    "purescript",
			"apt":    "purescript",
			"brew":   "purescript",
			"pacman": "purescript",
			"zypper": "purescript",
			"apk":    "purescript",
		},
		"roc": {
			"dnf":    "roc",
			"apt":    "roc",
			"brew":   "roc",
			"pacman": "roc",
			"zypper": "roc",
			"apk":    "roc",
		},
		"odin": {
			"dnf":    "odin",
			"apt":    "odin",
			"brew":   "odin",
			"pacman": "odin",
			"zypper": "odin",
			"apk":    "odin",
		},
		"nasm": {
			"dnf":    "nasm",
			"apt":    "nasm",
			"brew":   "nasm",
			"pacman": "nasm",
			"zypper": "nasm",
			"apk":    "nasm",
		},
		"arm-linux-gnueabi-gcc": {
			"dnf":    "gcc-arm-linux-gnu",
			"apt":    "gcc-arm-linux-gnueabihf",
			"brew":   "arm-linux-gnueabihf-gcc",
			"pacman": "mingw-w64-gcc",
			"zypper": "cross-arm-linux-gnu-gcc",
			"apk":    "arm-none-eabi-gcc",
		},
		"mips-linux-gnu-gcc": {
			"dnf":    "gcc-mips-linux-gnu",
			"apt":    "gcc-mips-linux-gnu",
			"brew":   "mips-linux-gnu-gcc",
			"pacman": "gcc",
			"zypper": "cross-mips-linux-gnu-gcc",
			"apk":    "gcc",
		},
		"riscv64-linux-gnu-gcc": {
			"dnf":    "gcc-riscv64-linux-gnu",
			"apt":    "gcc-riscv64-linux-gnu",
			"brew":   "riscv64-linux-gnu-gcc",
			"pacman": "gcc",
			"zypper": "cross-riscv64-linux-gnu-gcc",
			"apk":    "gcc",
		},
	}

	if langMap, ok := packageMap[language]; ok {
		if pkgName, ok := langMap[pkgManager]; ok {
			return pkgName
		}
	}

	// Default: return the language name
	return language
}

// GetOS returns the operating system name
func GetOS() string {
	return runtime.GOOS
}
