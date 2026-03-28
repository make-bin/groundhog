// @AI_GENERATED
package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewConfigCommand creates the "openclaw config" command group.
func NewConfigCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage OpenClaw configuration",
		Long:  "Read and write configuration values in config.yaml using dot-notation keys.",
	}

	cmd.PersistentFlags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")

	cmd.AddCommand(newConfigGetCommand(&configPath))
	cmd.AddCommand(newConfigSetCommand(&configPath))
	cmd.AddCommand(newConfigListCommand(&configPath))

	return cmd
}

// newConfigGetCommand creates "openclaw config get <key>".
func newConfigGetCommand(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value by dot-notation key",
		Example: `  openclaw config get database.host
  openclaw config get server.port`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			v, err := loadViper(*configPath)
			if err != nil {
				return err
			}

			if !v.IsSet(key) {
				return fmt.Errorf("key %q not found in config", key)
			}

			val := v.Get(key)
			fmt.Printf("%v\n", val)
			return nil
		},
	}
}

// newConfigSetCommand creates "openclaw config set <key> <value>".
func newConfigSetCommand(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value by dot-notation key",
		Example: `  openclaw config set database.host localhost
  openclaw config set server.port 9090`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			v, err := loadViper(*configPath)
			if err != nil {
				return err
			}

			v.Set(key, value)

			if err := v.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("%s✓%s  Set %s = %s\n", colorGreen, colorReset, key, value)
			return nil
		},
	}
}

// newConfigListCommand creates "openclaw config list".
func newConfigListCommand(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration keys and values",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := loadViper(*configPath)
			if err != nil {
				return err
			}

			settings := v.AllSettings()
			keys := flattenKeys(settings, "")
			sort.Strings(keys)

			fmt.Printf("\n%s%sConfiguration: %s%s\n\n", colorBold, colorGreen, *configPath, colorReset)

			const keyWidth = 40
			fmt.Printf("  %-*s  %s\n", keyWidth, "Key", "Value")
			fmt.Printf("  %s  %s\n", repeatChar('─', keyWidth), repeatChar('─', 40))

			for _, k := range keys {
				val := v.Get(k)
				valStr := fmt.Sprintf("%v", val)
				// Mask sensitive keys
				if isSensitiveKey(k) {
					valStr = maskValue(valStr)
				}
				fmt.Printf("  %-*s  %s\n", keyWidth, k, valStr)
			}

			fmt.Println()
			return nil
		},
	}
}

// loadViper creates a viper instance loaded from the given config file path.
func loadViper(configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	return v, nil
}

// flattenKeys recursively flattens a nested map into dot-notation keys.
func flattenKeys(m map[string]interface{}, prefix string) []string {
	var keys []string
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		if nested, ok := v.(map[string]interface{}); ok {
			keys = append(keys, flattenKeys(nested, fullKey)...)
		} else {
			keys = append(keys, fullKey)
		}
	}
	return keys
}

// isSensitiveKey returns true for keys that should be masked in output.
func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	sensitivePatterns := []string{"password", "secret", "api_key", "apikey", "token"}
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// maskValue replaces all but the first 2 characters with asterisks.
func maskValue(val string) string {
	if len(val) <= 2 {
		return strings.Repeat("*", len(val))
	}
	return val[:2] + strings.Repeat("*", len(val)-2)
}

// @AI_GENERATED: end
