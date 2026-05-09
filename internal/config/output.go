package config

type OptionStyle int

const (
	OptionStyleShort OptionStyle = iota
	OptionStyleLong
	OptionStyleCombined
)

type OutputConfig struct {
	ShowTitle     bool
	PlatformTitle bool
	ShowHyphens   bool
	EditLink      bool
	ExamplePrefix string
	LineLength    int
	Compact       bool
	OptionStyle   OptionStyle
	RawMarkdown   bool
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
