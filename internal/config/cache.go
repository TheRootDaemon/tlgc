package config

import (
	"os/user"
	"path/filepath"
)

const Mirror string = "https://github.com/tldr-pages/tldr/releases/latest/download"

type CacheConfig struct {
	Dir             string   `toml:"dir"`
	Mirror          string   `toml:"mirror"`
	AutoUpdate      bool     `toml:"auto_update"`
	DeferAutoUpdate bool     `toml:"defer_auto_update"`
	MaxAge          uint64   `toml:"max_age"`
	Languages       []string `toml:"languages"`
}

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
