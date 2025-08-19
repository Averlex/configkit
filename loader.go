package configkit

import (
	"fmt"
	"io"
	"reflect"

	"github.com/spf13/viper"
)

// Loader is a configuration loader with CLI support.
type Loader struct {
	name, short, long string // Root command attributes.
	configPath        string
	envPrefix         string
}

// NewLoader returns a new viper loader.
//
//   - name: short name of the service, as well as the name of the root command.
//   - short, long: short and long descriptions of the service for the root command.
//   - configPath is the path to the configuration file. It will be overrided with a value,
//     received via the --config flag. If the flag is not set, Loader will use the configPath.
//   - envPrefix: prefix for environment variables (e.g., "APP" → APP_LOG_LEVEL).
func NewLoader(name, short, long, configPath, envPrefix string) *Loader {
	return &Loader{
		configPath: configPath,
		envPrefix:  envPrefix,
		name:       name,
		short:      short,
		long:       long,
	}
}

// Load loads configuration from a file and environment variables into cfg.
//
// Arguments:
//   - cfg must be a pointer to a struct.
//   - printVersion is called when --version is used.
//     Package provides to helpers (JSONVersionPrinter and PlainVersionPrinter) for a quick setup.
//   - writer is used for output (can be os.Stdout, buffer, etc. - any writer to handle printVersion).
//
// Note: While Loader is stateless and uses isolated viper instances,
// calling Load() concurrently with CLI flag parsing is not recommended,
// as underlying libraries (such as cobra) are not designed for concurrent use.
// Use Load() sequentially during application startup.
//
// Returns:
//   - LoadResultStop: if --help or --version was used (no error).
//   - LoadResultContinue: if config was loaded successfully.
//   - error: if there was a problem (e.g. config file not found).
func (l *Loader) Load(cfg any, printVersion func(io.Writer) error, writer io.Writer) (LoadResult, error) {
	// Validate the input.
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return LoadResultStop, fmt.Errorf("cfg must be a pointer to a struct - got %s", reflect.ValueOf(cfg).Kind().String())
	}
	if printVersion == nil {
		return LoadResultStop, fmt.Errorf("printVersion must be a function")
	}
	if writer == nil {
		return LoadResultStop, fmt.Errorf("writer must be a non-nil writer")
	}

	v := viper.New()

	cmd, err := l.buildRootCommand(v, cfg, printVersion, writer)
	if err != nil {
		return LoadResultStop, fmt.Errorf("build root command: %w", err)
	}

	if err := cmd.Execute(); err != nil {
		return LoadResultStop, fmt.Errorf("execute root command: %w", err)
	}

	// If --help or --version was triggered, stop gracefully.
	if cmd.Flags().Changed("help") || cmd.Flags().Changed("version") {
		return LoadResultStop, nil
	}

	return LoadResultContinue, nil
}
