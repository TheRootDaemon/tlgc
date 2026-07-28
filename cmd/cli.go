package cmd

// CLI holds all parsed command-line flags and positional arguments.
type CLI struct {
	// operations

	// Page holds the positional arguments (page names to look up).
	Page []string

	// Update requests a cache update.
	Update bool

	// List requests listing pages for the current platform.
	List bool

	// ListAll requests listing all pages across all platforms.
	ListAll bool

	// Search requests a keyword search across pages.
	Search string

	// Browse requests opening the page in the default web browser.
	Browse bool

	// ListPlatforms requests listing available platforms.
	ListPlatforms bool

	// ListLanguages requests listing installed languages.
	ListLanguages bool

	// Info requests showing cache information.
	Info bool

	// Render requests rendering a local markdown file.
	Render string

	// CleanCache requests interactively cleaning the cache.
	CleanCache bool

	// GenConfig requests printing the default configuration.
	GenConfig bool

	// ConfigPath requests printing the config file path.
	ConfigPath bool

	// ShowVersion requests printing the version string.
	ShowVersion bool

	// ShowHelp requests printing the help text.
	ShowHelp bool

	// options

	// Platform overrides the platform used for page lookup.
	Platform string

	// Languages overrides the language list.
	Languages []string

	// ShortOptions requests displaying short option forms.
	ShortOptions bool

	// LongOptions requests displaying long option forms.
	LongOptions bool

	// Offline suppresses automatic cache updates.
	Offline bool

	// Compact strips empty lines from output.
	Compact bool

	// NoCompact overrides --compact, preserving empty lines in output.
	NoCompact bool

	// Raw prints pages in raw markdown.
	Raw bool

	// NoRaw overrides --raw, rendering pages instead of printing raw content.
	NoRaw bool

	// Quiet suppresses informational and warning messages.
	Quiet bool

	// Verbose controls the verbosity level (0–2).
	Verbose uint8

	// Color controls when to enable color output: auto, always, never.
	Color string

	// Config specifies an alternative config file path.
	Config string

	// page actions

	// Edit requests displaying a GitHub edit link.
	Edit bool

	// HasArgs indicates whether any command-line arguments were provided.
	HasArgs bool
}
