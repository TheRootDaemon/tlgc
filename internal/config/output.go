package config

type OptionStyle string

const (
	OptionStyleShort    OptionStyle = "short"
	OptionStyleLong     OptionStyle = "long"
	OptionStyleCombined OptionStyle = "both"
)

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
