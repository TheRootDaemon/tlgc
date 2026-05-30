package config

// IndentConfig controls indentation
// for each section of a rendered page.
type IndentConfig struct {
	Title       int `toml:"title"`
	Description int `toml:"description"`
	Bullet      int `toml:"bullet"`
	Example     int `toml:"example"`
}

// DefaultIndentConfig returns
// the default indentation settings.
//
// Titles and descriptions are indented 2 spaces,
// examples are indented 4 spaces.
func DefaultIndentConfig() IndentConfig {
	return IndentConfig{
		Title:       2,
		Description: 2,
		Bullet:      2,
		Example:     4,
	}
}
