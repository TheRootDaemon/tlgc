package config

// OptionStyle controls how command-line options
// appear in placeholder tokens.
type OptionStyle string

const (
	// OptionStyleShort shows only the short flag: -s.
	OptionStyleShort OptionStyle = "short"
	// OptionStyleLong shows only the long flag: --long.
	OptionStyleLong OptionStyle = "long"
	// OptionStyleCombined shows both flags: [-s|--long].
	OptionStyleCombined OptionStyle = "both"
)

// OutputConfig controls how rendered pages are displayed.
type OutputConfig struct {
	ShowTitle     bool        `toml:"show_title"`
	PlatformTitle bool        `toml:"platform_title"`
	ShowHyphens   bool        `toml:"show_hyphens"`
	EditLink      bool        `toml:"edit_link"`
	ExamplePrefix string      `toml:"example_prefix"`
	LineLength    int         `toml:"line_length"`
	Compact       bool        `toml:"compact"`
	OptionStyle   OptionStyle `toml:"option_style"`
	RawMarkdown   bool        `toml:"raw_markdown"`
}

// DefaultOutputConfig returns
// the default display settings.
//
// By default, the title is shown,
// hyphens and edit links are hidden,
// the example prefix is "- ",
// and options are shown in their long form.
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		ShowTitle:     true,
		PlatformTitle: false,
		ShowHyphens:   false,
		EditLink:      false,
		ExamplePrefix: "- ",
		LineLength:    0,
		Compact:       false,
		OptionStyle:   OptionStyleLong,
		RawMarkdown:   false,
	}
}
