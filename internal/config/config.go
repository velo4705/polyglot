package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents Polyglot configuration
type Config struct {
	AutoInstall      AutoInstallConfig         `yaml:"auto_install"`
	Languages        map[string]LanguageConfig `yaml:"languages,omitempty"`
	CustomExtensions map[string]string         `yaml:"custom_extensions,omitempty"`
	Environment      map[string]string         `yaml:"environment,omitempty"`
	Sandbox          SandboxConfig             `yaml:"sandbox,omitempty"`
}

// AutoInstallConfig represents auto-install settings
type AutoInstallConfig struct {
	Enabled bool   `yaml:"enabled"`
	Method  string `yaml:"method"` // system, script, managed
	Prompt  bool   `yaml:"prompt"`
}

// SandboxConfig represents sandboxing settings
type SandboxConfig struct {
	Enabled       bool  `yaml:"enabled"`
	MemoryLimitMB int64 `yaml:"memory_limit_mb"` // Memory limit in MB
	CPULimit      int   `yaml:"cpu_limit"`       // CPU limit in seconds
	Timeout       int   `yaml:"timeout"`         // Timeout limit in seconds
}

// LanguageConfig represents per-language configuration
type LanguageConfig struct {
	Enabled       bool              `yaml:"enabled"`
	Command       string            `yaml:"command,omitempty"`
	CompileCmd    string            `yaml:"compile_cmd,omitempty"`
	RunCmd        string            `yaml:"run_cmd,omitempty"`
	Flags         []string          `yaml:"flags,omitempty"`
	Version       string            `yaml:"version,omitempty"`
	InstallMethod string            `yaml:"install_method,omitempty"` // system, script, binary, manual
	InstallScript string            `yaml:"install_script,omitempty"`
	BinaryURL     string            `yaml:"binary_url,omitempty"`
	Environment   map[string]string `yaml:"environment,omitempty"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		AutoInstall: AutoInstallConfig{
			Enabled: false, // Disabled by default for safety
			Method:  "system",
			Prompt:  true,
		},
		Sandbox: SandboxConfig{
			Enabled:       false,
			MemoryLimitMB: 256, // 256MB default
			CPULimit:      10,  // 10s CPU time limit
			Timeout:       15,  // 15s wall clock time limit
		},
		Languages:        make(map[string]LanguageConfig),
		CustomExtensions: make(map[string]string),
		Environment:      make(map[string]string),
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".polyglot", "config.yaml")
}

// GetProjectConfigPath returns the path to the project config file
func GetProjectConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(cwd, ".polyglot.yaml")
}

// Load loads configuration from file (merges global and project configs)
func Load() (*Config, error) {
	// Start with default config
	config := DefaultConfig()

	// Load global config
	globalPath := GetConfigPath()
	if _, err := os.Stat(globalPath); err == nil {
		globalConfig, err := loadFromFile(globalPath)
		if err != nil {
			return nil, err
		}
		config = mergeConfigs(config, globalConfig)
	}

	// Load project config (overrides global)
	projectPath := GetProjectConfigPath()
	if _, err := os.Stat(projectPath); err == nil {
		projectConfig, err := loadFromFile(projectPath)
		if err != nil {
			return nil, err
		}
		config = mergeConfigs(config, projectConfig)
	}

	return config, nil
}

// loadFromFile loads configuration from a specific file
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// mergeConfigs merges two configs (second overrides first)
func mergeConfigs(base, override *Config) *Config {
	result := *base

	// Merge auto-install settings
	if override.AutoInstall.Enabled {
		result.AutoInstall = override.AutoInstall
	}

	// Merge languages
	if result.Languages == nil {
		result.Languages = make(map[string]LanguageConfig)
	}
	for k, v := range override.Languages {
		result.Languages[k] = v
	}

	// Merge custom extensions
	if result.CustomExtensions == nil {
		result.CustomExtensions = make(map[string]string)
	}
	for k, v := range override.CustomExtensions {
		result.CustomExtensions[k] = v
	}

	// Merge environment
	if result.Environment == nil {
		result.Environment = make(map[string]string)
	}
	for k, v := range override.Environment {
		result.Environment[k] = v
	}

	// Merge sandbox settings
	if override.Sandbox.Enabled {
		result.Sandbox.Enabled = override.Sandbox.Enabled
	}
	if override.Sandbox.MemoryLimitMB != 0 {
		result.Sandbox.MemoryLimitMB = override.Sandbox.MemoryLimitMB
	}
	if override.Sandbox.CPULimit != 0 {
		result.Sandbox.CPULimit = override.Sandbox.CPULimit
	}
	if override.Sandbox.Timeout != 0 {
		result.Sandbox.Timeout = override.Sandbox.Timeout
	}

	return &result
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := GetConfigPath()

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// GetLanguageConfig returns configuration for a specific language
func (c *Config) GetLanguageConfig(language string) *LanguageConfig {
	if cfg, ok := c.Languages[language]; ok {
		return &cfg
	}
	return nil
}

// GetCustomExtension returns the language for a custom extension
func (c *Config) GetCustomExtension(ext string) string {
	if lang, ok := c.CustomExtensions[ext]; ok {
		return lang
	}
	return ""
}

// GetEnvironment returns all environment variables (global + language-specific)
func (c *Config) GetEnvironment(language string) map[string]string {
	env := make(map[string]string)

	// Add global environment
	for k, v := range c.Environment {
		env[k] = v
	}

	// Add language-specific environment
	if langCfg := c.GetLanguageConfig(language); langCfg != nil {
		for k, v := range langCfg.Environment {
			env[k] = v
		}
	}

	return env
}

// IsLanguageEnabled checks if a language is enabled
func (c *Config) IsLanguageEnabled(language string) bool {
	if langCfg := c.GetLanguageConfig(language); langCfg != nil {
		return langCfg.Enabled
	}
	return true // Default to enabled if not configured
}
