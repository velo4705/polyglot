package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/config"
	yaml "gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Polyglot configuration",
	Long: `View and modify Polyglot configuration settings.

Commands:
  show                  Show current configuration
  init                  Initialize configuration file
  set <key> <value>     Set a configuration value
  get <key>             Get a configuration value
  path                  Show configuration file paths

Examples:
  polyglot config show
  polyglot config init
  polyglot config set auto_install.enabled true
  polyglot config set languages.python.version 3.11
  polyglot config set custom_extensions..myext python
  polyglot config get auto_install.enabled
  polyglot config path`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE:  configInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all configuration",
	RunE:  configShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE:  configSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE:  configGet,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file paths",
	RunE:  configPath,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configPathCmd)
}

func configInit(cmd *cobra.Command, args []string) error {
	configPath := config.GetConfigPath()

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists: %s\n", configPath)
		fmt.Print("Overwrite? [y/N]: ")
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	// Create default config
	cfg := config.DefaultConfig()
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Configuration initialized: %s\n", configPath)
	fmt.Println("\nDefault settings:")
	fmt.Println("  auto_install.enabled: false")
	fmt.Println("  auto_install.method: system")
	fmt.Println("  auto_install.prompt: true")
	fmt.Println("\nEdit the file to customize your settings.")

	return nil
}

func configShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	fmt.Println("Current Configuration:")
	fmt.Println("=====================")
	fmt.Print(string(data))
	fmt.Println("\nConfiguration files:")
	fmt.Printf("  Global: %s\n", config.GetConfigPath())
	fmt.Printf("  Project: %s\n", config.GetProjectConfigPath())

	return nil
}

func configSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Parse the key path and set the value
	parts := strings.Split(key, ".")
	if err := setConfigValue(cfg, parts, value); err != nil {
		return err
	}

	// Save configuration
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Set %s = %s\n", key, value)
	return nil
}

func configGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Parse the key path and get the value
	parts := strings.Split(key, ".")
	value, err := getConfigValue(cfg, parts)
	if err != nil {
		return err
	}

	fmt.Printf("%s = %v\n", key, value)
	return nil
}

func configPath(cmd *cobra.Command, args []string) error {
	fmt.Printf("Global config: %s\n", config.GetConfigPath())
	fmt.Printf("Project config: %s\n", config.GetProjectConfigPath())
	return nil
}

// Helper functions for nested config access
func setConfigValue(cfg *config.Config, parts []string, value string) error {
	if len(parts) == 0 {
		return fmt.Errorf("invalid key path")
	}

	switch parts[0] {
	case "auto_install":
		if len(parts) < 2 {
			return fmt.Errorf("invalid auto_install key")
		}
		switch parts[1] {
		case "enabled":
			cfg.AutoInstall.Enabled = value == "true"
		case "method":
			cfg.AutoInstall.Method = value
		case "prompt":
			cfg.AutoInstall.Prompt = value == "true"
		default:
			return fmt.Errorf("unknown auto_install key: %s", parts[1])
		}
	case "languages":
		if len(parts) < 3 {
			return fmt.Errorf("invalid languages key (format: languages.<lang>.<property>)")
		}
		lang := parts[1]
		prop := parts[2]

		if cfg.Languages == nil {
			cfg.Languages = make(map[string]config.LanguageConfig)
		}

		langCfg := cfg.Languages[lang]
		switch prop {
		case "enabled":
			langCfg.Enabled = value == "true"
		case "command":
			langCfg.Command = value
		case "version":
			langCfg.Version = value
		case "install_method":
			langCfg.InstallMethod = value
		default:
			return fmt.Errorf("unknown language property: %s", prop)
		}
		cfg.Languages[lang] = langCfg
	case "custom_extensions":
		if len(parts) < 2 {
			return fmt.Errorf("invalid custom_extensions key")
		}
		ext := parts[1]
		if cfg.CustomExtensions == nil {
			cfg.CustomExtensions = make(map[string]string)
		}
		cfg.CustomExtensions[ext] = value
	case "environment":
		if len(parts) < 2 {
			return fmt.Errorf("invalid environment key")
		}
		envVar := parts[1]
		if cfg.Environment == nil {
			cfg.Environment = make(map[string]string)
		}
		cfg.Environment[envVar] = value
	default:
		return fmt.Errorf("unknown config key: %s", parts[0])
	}

	return nil
}

func getConfigValue(cfg *config.Config, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid key path")
	}

	switch parts[0] {
	case "auto_install":
		if len(parts) < 2 {
			return cfg.AutoInstall, nil
		}
		switch parts[1] {
		case "enabled":
			return cfg.AutoInstall.Enabled, nil
		case "method":
			return cfg.AutoInstall.Method, nil
		case "prompt":
			return cfg.AutoInstall.Prompt, nil
		default:
			return nil, fmt.Errorf("unknown auto_install key: %s", parts[1])
		}
	case "languages":
		if len(parts) < 2 {
			return cfg.Languages, nil
		}
		lang := parts[1]
		if len(parts) < 3 {
			return cfg.Languages[lang], nil
		}
		langCfg := cfg.Languages[lang]
		switch parts[2] {
		case "enabled":
			return langCfg.Enabled, nil
		case "command":
			return langCfg.Command, nil
		case "version":
			return langCfg.Version, nil
		case "install_method":
			return langCfg.InstallMethod, nil
		default:
			return nil, fmt.Errorf("unknown language property: %s", parts[2])
		}
	case "custom_extensions":
		if len(parts) < 2 {
			return cfg.CustomExtensions, nil
		}
		ext := parts[1]
		if ext == "" && len(parts) > 2 {
			// Handle double dot case: custom_extensions..pyx
			ext = "." + parts[2]
		}
		return cfg.CustomExtensions[ext], nil
	case "environment":
		if len(parts) < 2 {
			return cfg.Environment, nil
		}
		return cfg.Environment[parts[1]], nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", parts[0])
	}
}
