package config

import (
	"os/user"
	"path/filepath"
)

// Mirror is the default URL
// for downloading tldr-pages archives.
const Mirror string = "https://github.com/tldr-pages/tldr/releases/latest/download"

// CacheConfig configures the page cache.
type CacheConfig struct {
	Dir             string   `toml:"dir"`
	Mirror          string   `toml:"mirror"`
	AutoUpdate      bool     `toml:"auto_update"`
	DeferAutoUpdate bool     `toml:"defer_auto_update"`
	MaxAge          uint64   `toml:"max_age"`
	Languages       []string `toml:"languages"`
}

// DefaultCacheConfig returns the default cache settings.
//
// The cache directory defaults to ~/.cache/tlgc,
// auto-update is enabled with a 2-week max age,
// and no specific languages are set (auto-detected from environment).
func DefaultCacheConfig() CacheConfig {
	user, _ := user.Current()
	defaultDir := filepath.Join(user.HomeDir, ".cache", "tlgc")

	return CacheConfig{
		Dir:             defaultDir,
		Mirror:          Mirror,
		AutoUpdate:      true,
		DeferAutoUpdate: false,
		MaxAge:          336,
		Languages:       []string{},
	}
}
