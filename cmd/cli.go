package cmd

// CLI holds all parsed command-line flags and positional arguments.
type CLI struct {
	// Quiet suppresses informational and warning messages.
	Quiet bool

	// GenConfig requests printing the default configuration.
	GenConfig bool

	// ConfigPath requests printing the config file path.
	ConfigPath bool

	// Verbose controls the verbosity level (0–2).
	Verbose uint8
}
