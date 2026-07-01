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
	// Dir is the path to the local directory where downloaded pages are cached.
	Dir string `toml:"dir"`

	// Mirror is the URL used to download tldr-pages archives.
	Mirror string `toml:"mirror"`

	// AutoUpdate controls whether the cache is refreshed automatically on startup.
	AutoUpdate bool `toml:"auto_update"`

	// DeferAutoUpdate controls whether auto-updates run in the background
	// instead of blocking startup.
	DeferAutoUpdate bool `toml:"defer_auto_update"`

	// MaxAge is the maximum age of cached pages, in hours, before they are
	// re-downloaded.
	MaxAge uint64 `toml:"max_age"`

	// Languages is the list of preferred page languages.
	// When empty the language is auto-detected from the environment.
	Languages []string `toml:"languages"`
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
