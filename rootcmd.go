package configkit

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buildRootCommand builds root command.
//
// Method declares flags and binds them to actions. It also enables env variables.
// If any of the env variables is set, it will overrride the config file.
//
// The result of the execution - is a built up command, ready to execute.
// All reading/setting is perfromed on the pre-run stage.
func (l *Loader) buildRootCommand(
	v *viper.Viper,
	cfg any,
	printVersion func(io.Writer) error,
	writer io.Writer,
) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   l.name,
		Short: l.short,
		Long:  l.long,
		Run: func(_ *cobra.Command, _ []string) {
			// Service logic is expected to be handled elsewhere.
		},
	}

	// Define flags.
	rootCmd.Flags().StringP("config", "c", "", "Path to configuration file")
	rootCmd.Flags().BoolP("version", "v", false, "Show version info")

	// Setup viper.
	v.SetEnvPrefix(l.envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Binding flags to viper.
	if err := v.BindPFlag("config", rootCmd.Flags().Lookup("config")); err != nil {
		return nil, fmt.Errorf("bind config flag: %w", err)
	}
	if err := v.BindPFlag("version", rootCmd.Flags().Lookup("version")); err != nil {
		return nil, fmt.Errorf("bind version flag: %w", err)
	}

	// Pre-run hook: load config or show version.
	rootCmd.PreRunE = func(_ *cobra.Command, _ []string) error {
		// Processing -v flag preemptively.
		if versionFlag := v.GetBool("version"); versionFlag {
			if err := printVersion(writer); err != nil {
				return fmt.Errorf("print version: %w", err)
			}
			return nil
		}

		// Setting the config. If none are set, use the value, passed directly to the loader.
		configPath := v.GetString("config")
		if configPath == "" {
			configPath = l.configPath
		}

		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if errors.As(err, &notFound) {
				return fmt.Errorf("config file not found at %q", configPath)
			}
			return fmt.Errorf("read main config at %q: %w", configPath, err)
		}

		if err := v.Unmarshal(cfg); err != nil {
			return fmt.Errorf("unmarshal main config: %w", err)
		}

		return nil
	}

	return rootCmd, nil
}
