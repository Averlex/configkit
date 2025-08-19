# configkit

[![Go version](https://img.shields.io/badge/go-1.24.2+-blue.svg)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/Averlex/configkit.svg)](https://pkg.go.dev/github.com/Averlex/configkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/Averlex/configkit)](https://goreportcard.com/report/github.com/Averlex/configkit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Simple configuration bootstrapping for Go services — load from file, override with env, and get CLI flags (--config, --version, --help) out of the box.

**Minimal configuration loader for Go services**

Zero-boilerplate setup for config loading from file, environment variables, and CLI flags.

```go
loader := configkit.NewLoader("myapp", "My App", "", "config.yaml", "MYAPP")
result, err := loader.Load(&cfg, configkit.PlainVersionPrinter("v1.0.0"), os.Stdout)
if err != nil {
    log.Fatal(err)
}
if result == configkit.LoadResultStop {
    return // --help or --version was used
}
// continue with initialized config
```

## 📚 Table of Contents

- [Why configkit?](#-why-configkit)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Version Output](#-version-output)
- [API Reference](#-api-reference)
- [Testing](#-testing)
- [Example](#-example)

## 🚀 Why configkit?

Setting up configuration in Go services often involves repetitive boilerplate:

- parsing `--config flag`,
- loading YAML/JSON,
- overriding with environment variables,
- handling `--version` and `--help`.

`configkit` solves this with a single, reusable `Load()` call — while staying lightweight and predictable.

## 📦 Features

- ✅ File + Env + CLI — load config from file, override values **and config path** with env vars (`PREFIX_KEY=value`), use `--config` flag.
- ✅ Built-in `--version` and `--help` — with customizable output.
- ✅ Stateless & isolated — each `Load()` uses its own `viper` instance. Safe for reuse.
- ✅ No global state pollution — doesn't touch `viper.GetViper()` or `pflag.CommandLine`.
- ✅ Flexible version printing — plain text, JSON, or custom format.
- ✅ Explicit control flow — `LoadResultStop` vs `LoadResultContinue` avoids `os.Exit()` surprises.

## 🛠 Installation

```bash
go get github.com/Averlex/configkit
```

## 🧩 Usage

1. Define your config struct
   ```go
   type Config struct {
       Port int `mapstructure:"port"`
       DB   struct {
           URL string `mapstructure:"url"`
       } `mapstructure:"db"`
   }
   ```
2. Create and use the loader

   ```go
   package main

   import (
       "log"
       "os"

       "github.com/Averlex/configkit"
   )

   func main() {
       var cfg Config

       loader := configkit.NewLoader(
           "myapp",                    // command name
           "My awesome service",       // short description
           "Long description...",      // long description (optional)
           "config.yaml",              // default config path
           "MYAPP",                    // env prefix (e.g. MYAPP_DB_URL)
       )

       result, err := loader.Load(
           &cfg,
           configkit.PlainVersionPrinter("v1.0.0"),
           os.Stdout,
       )
       if err != nil {
           log.Fatal(err)
       }
       if result == configkit.LoadResultStop {
           return // --help or --version was used
       }

       log.Printf("Config loaded: %+v", cfg)
   }
   ```

3. Configuration file (`config.yaml`)

   ```yaml
   port: 8080
   db:
   url: "localhost:5432"
   ```

4. Override via environment

   ```bash
   MYAPP_DB_URL=prod-db:5432 ./myapp --config prod.yaml
   ```

5. Override config path via environment

   The `--config` flag itself can be overridden by an environment variable:

   ```bash
   MYAPP_CONFIG=staging.yaml ./myapp
   ```

## 🖨 Version Output

Use built-in helpers:

```go
configkit.PlainVersionPrinter("v1.0.0")
configkit.JSONVersionPrinter("v1.0.0", "abc123", "2025-04-05")
```

Or define your own:

```go
func(w io.Writer) error {
    return json.NewEncoder(w).Encode(map[string]string{
        "version": "dev",
        "built":   time.Now().Format(time.RFC3339),
    })
}
```

## 📚 API Reference

```go
NewLoader(name, short, long, configPath, envPrefix) *Loader
```

**Creates a new configuration loader**.

- `name`: command name (e.g., `myapp`).
- `short`, `long`: descriptions for `--help`.
- `configPath`: fallback path if `--config` is not provided.
- `envPrefix`: prefix for environment variables (e.g., `APP_CONFIG`, `APP_LOG_LEVEL`).
  **Automatically binds `PREFIX_CONFIG` to the `--config` flag**.

```go
Load(cfg, printVersion, writer) (LoadResult, error)
```

**Loads configuration** into `cfg`.

- `cfg`: must be a pointer to a struct.
- `printVersion`: function to call when `--version` is used.
- `writer`: where version/help output is written (e.g., `os.Stdout`).

**Returns**:

- `LoadResultContinue`: config loaded successfully.
- `LoadResultStop`: `--help` or `--version` was used — stop execution.
- `error`: failed to load config (e.g., file not found).

> ⚠️ **Concurrency note**: While `Loader` is stateless and uses isolated `viper` instances, concurrent calls to `Load()` with CLI flag parsing are not recommended, as underlying libraries (such as `cobra`) are not designed for concurrent use. Use `Load()` sequentially during application initialization.

## 🧪 Testing

`configkit` is designed to be testable:

- Accepts any `io.Writer` (use `*bytes.Buffer` in tests).
- Uses isolated `viper` instances.
- No global state.

See `loader_test.go` for examples of testing config loading, version flag, and env override.

## 🧰 Example

Check `examples/simple` for a complete working example.
