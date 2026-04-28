package config

import (
	"os/user"
	"path/filepath"
)

const Mirror string = "https://github.com/tldr-pages/tldr/releases/latest/download"

type CacheConfig struct {
	Dir             string
	Mirror          string
	AutoUpdate      bool
	DeferAutoUpdate bool
	MaxAge          uint64
	Languages       []string
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
