package config

import (
	"sync/atomic"
)

// currentConfig represents the singleton
// of the current client configuration.
var currentConfig atomic.Pointer[Config]

// Initialize loads a TOML config file and sets it as the global singleton.
// It is safe to call from multiple goroutines.
// Subsequent calls replace the current singleton.
func Initialize() error {
	cfg, err := LoadConfig(ConfigPath())
	if err != nil {
		return err
	}

	currentConfig.Store(cfg)
	return nil
}

// C returns the global Config singleton.
// If Initialize has not been called, it returns a pointer
// to the default configuration.
func C() *Config {
	if cfg := currentConfig.Load(); cfg != nil {
		return cfg
	}

	d := Default()
	return &d
}

// Cache returns the cache configuration from the global singleton.
func Cache() CacheConfig {
	return C().Cache
}

// Style returns the style configuration from the global singleton.
func Style() StyleConfig {
	return C().Style
}

// Indent returns the indent configuration from the global singleton.
func Indent() IndentConfig {
	return C().Indent
}

// Output returns the output configuration from the global singleton.
func Output() OutputConfig {
	return C().Output
}
