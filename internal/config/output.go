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
	// ShowTitle controls whether the page title is displayed.
	ShowTitle bool `toml:"show_title"`

	// PlatformTitle controls whether the platform name is shown as a title prefix.
	PlatformTitle bool `toml:"platform_title"`

	// ShowHyphens controls whether a hyphen is displayed before each example description.
	ShowHyphens bool `toml:"show_hyphens"`

	// EditLink controls whether a GitHub edit link is shown at the top of the page.
	EditLink bool `toml:"edit_link"`

	// ExamplePrefix is the string prepended to each example description
	// when ShowHyphens is true.
	ExamplePrefix string `toml:"example_prefix"`

	// LineLength is the maximum line length for text wrapping.
	// A value of zero or less disables wrapping.
	LineLength int `toml:"line_length"`

	// Compact controls whether blank separator lines between sections are omitted.
	Compact bool `toml:"compact"`

	// OptionStyle controls how command-line option placeholders are displayed:
	// short form, long form, or both combined.
	OptionStyle OptionStyle `toml:"option_style"`

	// RawMarkdown controls whether the raw markdown source is written
	// instead of the formatted output.
	RawMarkdown bool `toml:"raw_markdown"`
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
