package config

type IndentConfig struct {
	Title       int
	Description int
	Bullet      int
	Example     int
}

func DefaultIndentConfig() IndentConfig {
	return IndentConfig{
		Title:       2,
		Description: 2,
		Bullet:      2,
		Example:     4,
	}
}
