// Package configkit provides a minimal configuration loader for Go services.
// It supports:
//   - Loading config from file (--config flag)
//   - Overriding values with environment variables (prefix-based)
//   - Built-in --version and --help flags
//   - Custom version output formatting
//
// Example:
//
//	loader := configkit.NewLoader("myapp", "My App", "", "config.yaml", "MYAPP")
//	result, err := loader.Load(&cfg, configkit.PlainVersionPrinter("v1.0.0"), os.Stdout)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if result == configkit.LoadResultStop {
//	    return
//	}
//	// continue with initialized config
package configkit

// LoadResult indicates the outcome of the Load operation.
type LoadResult int

const (
	// LoadResultContinue means the program should continue after config load.
	LoadResultContinue LoadResult = iota

	// LoadResultStop means the program should stop (e.g. after --help or --version).
	LoadResultStop
)
