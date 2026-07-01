package config

// IndentConfig controls indentation
// for each section of a rendered page.
type IndentConfig struct {
	// Title is the number of spaces used to indent the page title.
	Title int `toml:"title"`

	// Description is the number of spaces used to indent description lines.
	Description int `toml:"description"`

	// Bullet is the number of spaces used to indent bullet items (example descriptions).
	Bullet int `toml:"bullet"`

	// Example is the number of spaces used to indent command example blocks.
	Example int `toml:"example"`
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
