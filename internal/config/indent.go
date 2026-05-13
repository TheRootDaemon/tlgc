package config

type IndentConfig struct {
	Title       int `toml:"title"`
	Description int `toml:"description"`
	Bullet      int `toml:"bullet"`
	Example     int `toml:"example"`
}

func DefaultIndentConfig() IndentConfig {
	return IndentConfig{
		Title:       2,
		Description: 2,
		Bullet:      2,
		Example:     4,
	}
}
