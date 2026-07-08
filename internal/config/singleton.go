package config

import (
	"os"
	"sync/atomic"
)

// currentConfig represents the singleton
// of the current client configuration.
var currentConfig atomic.Pointer[Config]

// Initialize loads the TOML config file and sets the global singleton.
// If the config file does not exist, the singleton is set to defaults
// and no error is returned. It is safe to call from multiple goroutines.
func Initialize() error {
	path := ConfigPath()
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			d := Default()
			currentConfig.Store(&d)
			return nil
		}
		return err
	}

	cfg, err := LoadConfig(path)
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

// ResetForTesting clears the global config singleton.
// This function is only intended for use by tests.
func ResetForTesting() {
	currentConfig.Store(nil)
}
