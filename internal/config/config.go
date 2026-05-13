package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Cache  CacheConfig  `toml:"cache"`
	Indent IndentConfig `toml:"indent"`
	Output OutputConfig `toml:"output"`
	Style  StyleConfig  `toml:"style"`
}

func Default() Config {
	return Config{
		Cache:  DefaultCacheConfig(),
		Indent: DefaultIndentConfig(),
		Output: DefaultOutputConfig(),
		Style:  DefaultStyleConfig(),
	}
}

func DefaultConfig() (string, error) {
	var buf strings.Builder
	err := toml.NewEncoder(&buf).Encode(Default())
	return buf.String(), err
}

func ConfigPath() string {
	if p := os.Getenv("TLGC_CONFIG"); p != "" {
		return p
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "tlgc", "config.toml")
}

func LoadConfig(path string) (*Config, error) {
	cfg := Default()
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
