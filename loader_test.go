package configkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

var errSomeError = fmt.Errorf("some error")

type LoaderSuite struct {
	suite.Suite
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) SetupTest() {
	viper.Reset()
	os.Args = []string{"testapp"}
}

func (s *LoaderSuite) TestLoad_ArgsValidation() {
	testCases := []struct {
		name          string
		cfg           any
		printVersion  func(io.Writer) error
		writer        io.Writer
		args          []string
		expected      LoadResult
		expectedError error
	}{
		{
			name:          "nonpointer config",
			cfg:           struct{}{},
			printVersion:  PlainVersionPrinter("v1.0.0"),
			writer:        &bytes.Buffer{},
			args:          []string{},
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
		{
			name:          "nil printVersion",
			cfg:           &struct{}{},
			printVersion:  nil,
			writer:        &bytes.Buffer{},
			args:          []string{},
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
		{
			name:          "nil writer",
			cfg:           &struct{}{},
			printVersion:  PlainVersionPrinter("v1.0.0"),
			writer:        nil,
			args:          []string{},
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.name, func() {
			s.T().Parallel()
			loader := NewLoader("testapp", "Test App", "", "config.yaml", "TESTAPP")
			os.Args = append([]string{"testapp"}, tC.args...)
			result, err := loader.Load(tC.cfg, tC.printVersion, tC.writer)

			if tC.expectedError != nil {
				s.Require().Error(err, "expected error, got nil")
			} else {
				s.Require().NoError(err, "expected nil, got error")
			}
			s.Require().Equal(tC.expected, result, "unexpected load result: got %v, want %v", result, tC.expected)
		})
	}
}

func (s *LoaderSuite) TestLoad_FlagsProcessing() {
	testCases := []struct {
		name          string
		cfg           any
		printVersion  func(io.Writer) error
		writer        io.Writer
		args          []string
		configPath    string
		expected      LoadResult
		expectedError error
	}{
		{
			name:         "version",
			cfg:          &struct{}{},
			printVersion: PlainVersionPrinter("v1.0.0"),
			writer:       &bytes.Buffer{},
			args:         []string{"--version"},
			configPath:   "config.yaml",
			expected:     LoadResultStop,
		},
		{
			name:         "help",
			cfg:          &struct{}{},
			printVersion: PlainVersionPrinter("v1.0.0"),
			writer:       &bytes.Buffer{},
			args:         []string{"--help"},
			configPath:   "config.yaml",
			expected:     LoadResultStop,
		},
		{
			name:          "config flag override",
			cfg:           &struct{}{},
			printVersion:  PlainVersionPrinter("v1.0.0"),
			writer:        &bytes.Buffer{},
			args:          []string{"--config", "nonexistent.yaml"},
			configPath:    "config.yaml",
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
		{
			name:          "config flag unset",
			cfg:           &struct{}{},
			printVersion:  PlainVersionPrinter("v1.0.0"),
			writer:        &bytes.Buffer{},
			args:          []string{},
			configPath:    "nonexistent.yaml",
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
		{
			name:         "version and config flags",
			cfg:          &struct{}{},
			printVersion: PlainVersionPrinter("v1.0.0"),
			writer:       &bytes.Buffer{},
			args:         []string{"--version", "--config", "nonexistent.yaml"},
			configPath:   "config.yaml",
			expected:     LoadResultStop,
		},
		{
			name:         "help and config flags",
			cfg:          &struct{}{},
			printVersion: PlainVersionPrinter("v1.0.0"),
			writer:       &bytes.Buffer{},
			args:         []string{"--help", "--config", "nonexistent.yaml"},
			configPath:   "config.yaml",
			expected:     LoadResultStop,
		},
		{
			name:         "version and help flags",
			cfg:          &struct{}{},
			printVersion: PlainVersionPrinter("v1.0.0"),
			writer:       &bytes.Buffer{},
			args:         []string{"--version", "--help"},
			configPath:   "config.yaml",
			expected:     LoadResultStop,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.name, func() {
			s.T().Parallel()
			loader := NewLoader("testapp", "Test App", "", tC.configPath, "TESTAPP")
			os.Args = append([]string{"testapp"}, tC.args...)
			result, err := loader.Load(tC.cfg, tC.printVersion, tC.writer)

			if tC.expectedError != nil {
				s.Require().Error(err, "expected error, got nil")
			} else {
				s.Require().NoError(err, "expected nil, got error")
			}
			s.Require().Equal(tC.expected, result, "unexpected load result: got %v, want %v", result, tC.expected)
		})
	}
}

func (s *LoaderSuite) TestLoad_IncorrectConfig() {
	// Create a temporary file with invalid YAML content.
	invalidConfigPath := filepath.Join(os.TempDir(), "invalid_config.yaml")
	err := os.WriteFile(invalidConfigPath, []byte("invalid: yaml: content"), 0o600)
	s.Require().NoError(err, "failed to create invalid config file")
	defer os.Remove(invalidConfigPath)

	testCases := []struct {
		name          string
		cfg           any
		printVersion  func(io.Writer) error
		writer        io.Writer
		args          []string
		configPath    string
		expected      LoadResult
		expectedError error
	}{
		{
			name:          "nonexistent path",
			cfg:           &struct{}{},
			printVersion:  PlainVersionPrinter("v1.0.0"),
			writer:        &bytes.Buffer{},
			args:          []string{},
			configPath:    "nonexistent.yaml",
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
		{
			name:          "invalid config",
			cfg:           &struct{}{},
			printVersion:  PlainVersionPrinter("v1.0.0"),
			writer:        &bytes.Buffer{},
			args:          []string{},
			configPath:    invalidConfigPath,
			expected:      LoadResultStop,
			expectedError: errSomeError,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.name, func() {
			s.T().Parallel()
			loader := NewLoader("testapp", "Test App", "", tC.configPath, "TESTAPP")
			os.Args = append([]string{"testapp"}, tC.args...)
			result, err := loader.Load(tC.cfg, tC.printVersion, tC.writer)

			if tC.expectedError != nil {
				s.Require().Error(err, "expected error, got nil")
			} else {
				s.Require().NoError(err, "expected nil, got error")
			}
			s.Require().Equal(tC.expected, result, "unexpected load result: got %v, want %v", result, tC.expected)
		})
	}
}

func (s *LoaderSuite) TestLoad_ValidCases() {
	type testConfig struct {
		LogLevel string `mapstructure:"log_level"`
		Port     int    `mapstructure:"port"`
	}

	testCases := []struct {
		name           string
		setupFunc      func() (string, func())
		envVars        map[string]string
		args           []string
		configPath     string
		expectedConfig testConfig
		expected       LoadResult
		expectedError  error
	}{
		{
			name: "unset config flag",
			setupFunc: func() (string, func()) {
				configPath := filepath.Join(os.TempDir(), "default_config.yaml")
				content := map[string]any{"log_level": "info", "port": 8080}
				data, err := yaml.Marshal(content)
				s.Require().NoError(err, "marshal yaml")
				err = os.WriteFile(configPath, data, 0o600)
				s.Require().NoError(err, "write config file")
				return configPath, func() { os.Remove(configPath) }
			},
			args:           []string{},
			expectedConfig: testConfig{LogLevel: "info", Port: 8080},
			expected:       LoadResultContinue,
		},
		{
			name: "flag overrides default path",
			setupFunc: func() (string, func()) {
				configPath := filepath.Join(os.TempDir(), "override_config.yaml")
				content := map[string]any{"log_level": "warn", "port": 9090}
				data, err := yaml.Marshal(content)
				s.Require().NoError(err, "marshal yaml")
				err = os.WriteFile(configPath, data, 0o600)
				s.Require().NoError(err, "write config file")
				return configPath, func() { os.Remove(configPath) }
			},
			args:           []string{"--config", ""}, // Will be set in test.
			configPath:     "default.yaml",           // Default, but overridden.
			expectedConfig: testConfig{LogLevel: "warn", Port: 9090},
			expected:       LoadResultContinue,
		},
		{
			name: "config path overridden by env",
			setupFunc: func() (string, func()) {
				configPath := filepath.Join(os.TempDir(), "env_path_config.yaml")
				content := map[string]any{"log_level": "error", "port": 6060}
				data, err := yaml.Marshal(content)
				s.Require().NoError(err, "marshal yaml")
				err = os.WriteFile(configPath, data, 0o600)
				s.Require().NoError(err, "write config file")
				return configPath, func() { os.Remove(configPath) }
			},
			envVars:        map[string]string{"TESTAPP_CONFIG": ""}, // Will be set in test.
			args:           []string{},
			configPath:     "default.yaml", // Default, but overridden by env.
			expectedConfig: testConfig{LogLevel: "error", Port: 6060},
			expected:       LoadResultContinue,
		},
		{
			name: "field overridden by env",
			setupFunc: func() (string, func()) {
				configPath := filepath.Join(os.TempDir(), "env_config.yaml")
				content := map[string]any{"log_level": "info", "port": 8080}
				data, err := yaml.Marshal(content)
				s.Require().NoError(err, "marshal yaml")
				err = os.WriteFile(configPath, data, 0o600)
				s.Require().NoError(err, "write config file")
				return configPath, func() { os.Remove(configPath) }
			},
			envVars:        map[string]string{"TESTAPP_LOG_LEVEL": "debug"},
			args:           []string{},
			expectedConfig: testConfig{LogLevel: "debug", Port: 8080},
			expected:       LoadResultContinue,
		},
		{
			name: "config and env combined",
			setupFunc: func() (string, func()) {
				configPath := filepath.Join(os.TempDir(), "combined_config.yaml")
				content := map[string]any{"log_level": "info", "port": 8080}
				data, err := yaml.Marshal(content)
				s.Require().NoError(err, "marshal yaml")
				err = os.WriteFile(configPath, data, 0o600)
				s.Require().NoError(err, "write config file")
				return configPath, func() { os.Remove(configPath) }
			},
			envVars:        map[string]string{"TESTAPP_PORT": "7070"},
			args:           []string{},
			expectedConfig: testConfig{LogLevel: "info", Port: 7070},
			expected:       LoadResultContinue,
		},
		{
			name: "json config",
			setupFunc: func() (string, func()) {
				configPath := filepath.Join(os.TempDir(), "json_config.json")
				content := map[string]any{"log_level": "error", "port": 6060}
				data, err := json.Marshal(content)
				s.Require().NoError(err, "marshal json")
				err = os.WriteFile(configPath, data, 0o600)
				s.Require().NoError(err, "write config file")
				return configPath, func() { os.Remove(configPath) }
			},
			args:           []string{},
			expectedConfig: testConfig{LogLevel: "error", Port: 6060},
			expected:       LoadResultContinue,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.name, func() {
			s.T().Parallel()
			var cleanup func()
			var actualConfigPath string
			if tC.setupFunc != nil {
				actualConfigPath, cleanup = tC.setupFunc()
				defer cleanup()
			}
			loader := NewLoader("testapp", "Test App", "", actualConfigPath, "TESTAPP")

			for k, v := range tC.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			if len(tC.args) > 0 && tC.args[0] == "--config" {
				tC.args[1] = actualConfigPath
			}

			os.Args = append([]string{"testapp"}, tC.args...)
			cfg := &testConfig{}
			result, err := loader.Load(cfg, PlainVersionPrinter("v1.0.0"), &bytes.Buffer{})

			if tC.expectedError != nil {
				s.Require().Error(err, "expected error, got nil")
			} else {
				s.Require().NoError(err, "expected nil, got error")
			}
			s.Require().Equal(tC.expected, result, "unexpected load result: got %v, want %v", result, tC.expected)
			s.Require().Equal(tC.expectedConfig, *cfg, "unexpected config: got %+v, want %+v", *cfg, tC.expectedConfig)
		})
	}
}

func (s *LoaderSuite) TestLoad_SequentialLoads() {
	type testConfig struct {
		LogLevel string `mapstructure:"log_level"`
		Port     int    `mapstructure:"port"`
	}

	s.Run("sequential loads", func() {
		// Create the first config.
		configPath1 := filepath.Join(os.TempDir(), "seq_config1.yaml")
		content1 := map[string]any{"log_level": "info", "port": 8080}
		data1, err := yaml.Marshal(content1)
		s.Require().NoError(err, "marshal yaml")
		err = os.WriteFile(configPath1, data1, 0o600)
		s.Require().NoError(err, "write config file")
		defer os.Remove(configPath1)

		loader1 := NewLoader("testapp", "Test App", "", configPath1, "TESTAPP")
		os.Args = []string{"testapp"}
		cfg1 := &testConfig{}
		result1, err1 := loader1.Load(cfg1, PlainVersionPrinter("v1.0.0"), &bytes.Buffer{})
		s.Require().NoError(err1, "expected nil, got error")
		s.Require().Equal(LoadResultContinue, result1, "unexpected load result: got %v, want %v", result1, LoadResultContinue)
		s.Require().Equal(
			testConfig{LogLevel: "info", Port: 8080},
			*cfg1,
			"unexpected config: got %+v, want %+v", *cfg1, testConfig{LogLevel: "info", Port: 8080},
		)

		// Create the second config.
		configPath2 := filepath.Join(os.TempDir(), "seq_config2.yaml")
		content2 := map[string]any{"log_level": "debug", "port": 9090}
		data2, err := yaml.Marshal(content2)
		s.Require().NoError(err, "marshal yaml")
		err = os.WriteFile(configPath2, data2, 0o600)
		s.Require().NoError(err, "write config file")
		defer os.Remove(configPath2)

		loader2 := NewLoader("testapp", "Test App", "", configPath2, "TESTAPP")
		os.Args = []string{"testapp"}
		cfg2 := &testConfig{}
		result2, err2 := loader2.Load(cfg2, PlainVersionPrinter("v1.0.0"), &bytes.Buffer{})
		s.Require().NoError(err2, "expected nil, got error")
		s.Require().Equal(LoadResultContinue, result2, "unexpected load result: got %v, want %v", result2, LoadResultContinue)
		s.Require().Equal(
			testConfig{LogLevel: "debug", Port: 9090},
			*cfg2,
			"unexpected config: got %+v, want %+v", *cfg2, testConfig{LogLevel: "debug", Port: 9090},
		)
	})
}
