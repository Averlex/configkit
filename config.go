// Package configkit provides a minimal configuration loader for Go services.
// It supports loading from file, overriding with environment variables,
// and built-in CLI flags (--config, --version, --help).
package configkit

import (
	"errors"
)

// Additional errors to control the execution flow.
var (
	// ErrShouldStop is returned when the execution should be stopped. E.g. on -v and -h flags.
	ErrShouldStop = errors.New("execution should be stopped")
)

// ServiceConfig is an interface generalizing service config.
type ServiceConfig interface {
	// GetSubConfig returns the part of the config that corresponds to the key.
	GetSubConfig(key string) (map[string]any, error)
}
