// Simple example of using configkit.
package main

import (
	"log"
	"os"

	"github.com/Averlex/configkit"
)

type Config struct {
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
	DB       struct {
		URL      string `mapstructure:"url"`
		PoolSize int    `mapstructure:"pool_size"`
	} `mapstructure:"db"`
}

func main() {
	var cfg Config

	loader := configkit.NewLoader(
		"simple-example",
		"A simple configkit demo",
		"This service demonstrates how to use configkit for config loading.",
		"config.yaml",
		"EXAMPLE",
	)

	result, err := loader.Load(
		&cfg,
		configkit.PlainVersionPrinter("v0.1.0"),
		os.Stdout,
	)
	if err != nil {
		log.Fatal(err)
	}
	if result == configkit.LoadResultStop {
		return
	}

	log.Printf("✅ Config loaded:")
	log.Printf("  Port: %d", cfg.Port)
	log.Printf("  Log Level: %s", cfg.LogLevel)
	log.Printf("  DB URL: %s", cfg.DB.URL)
	log.Printf("  DB Pool Size: %d", cfg.DB.PoolSize)
}
