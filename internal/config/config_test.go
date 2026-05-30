package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	assert.Equal(t, DefaultCacheConfig(), cfg.Cache)
	assert.Equal(t, DefaultIndentConfig(), cfg.Indent)
	assert.Equal(t, DefaultOutputConfig(), cfg.Output)
	assert.Equal(t, DefaultStyleConfig(), cfg.Style)
}

func TestDefaultConfig(t *testing.T) {
	s, err := DefaultConfig()
	require.NoError(t, err)
	require.NotEmpty(t, s)
}

func TestDefaultConfigRoundTrip(t *testing.T) {
	s, err := DefaultConfig()
	require.NoError(t, err)

	var cfg Config
	_, err = toml.Decode(s, &cfg)
	require.NoError(t, err)

	assert.Equal(t, DefaultCacheConfig(), cfg.Cache)
	assert.Equal(t, DefaultIndentConfig(), cfg.Indent)
	assert.Equal(t, DefaultOutputConfig(), cfg.Output)
	assert.Equal(t, DefaultStyleConfig(), cfg.Style)
}

func TestConfigPath(t *testing.T) {
	t.Run("uses env var when set", func(t *testing.T) {
		t.Setenv("TLGC_CONFIG", "/custom/path/config.toml")
		assert.Equal(t, "/custom/path/config.toml", ConfigPath())
	})

	t.Run("falls back to user config dir", func(t *testing.T) {
		t.Setenv("TLGC_CONFIG", "")
		dir, err := os.UserConfigDir()
		require.NoError(t, err)
		expected := filepath.Join(dir, "tlgc", "config.toml")
		assert.Equal(t, expected, ConfigPath())
	})
}

func TestLoadConfig_fileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.toml")
	require.Error(t, err)
}

func TestLoadConfig_validFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	s, err := DefaultConfig()
	require.NoError(t, err)
	err = os.WriteFile(path, []byte(s), 0o644)
	require.NoError(t, err)

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, DefaultCacheConfig(), cfg.Cache)
	assert.Equal(t, DefaultIndentConfig(), cfg.Indent)
	assert.Equal(t, DefaultOutputConfig(), cfg.Output)
	assert.Equal(t, DefaultStyleConfig(), cfg.Style)
}

func TestLoadConfig_partialOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	input := `[output]
show_title = false
option_style = "short"
`
	err := os.WriteFile(path, []byte(input), 0o644)
	require.NoError(t, err)

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.False(t, cfg.Output.ShowTitle)
	assert.Equal(t, OptionStyleShort, cfg.Output.OptionStyle)

	assert.Equal(t, DefaultCacheConfig(), cfg.Cache)
	assert.Equal(t, DefaultIndentConfig(), cfg.Indent)
	assert.Equal(t, DefaultStyleConfig(), cfg.Style)
}

func TestDefaultCacheConfig(t *testing.T) {
	cfg := DefaultCacheConfig()

	assert.NotEmpty(t, cfg.Dir)
	assert.Equal(t, Mirror, cfg.Mirror)
	assert.True(t, cfg.AutoUpdate)
	assert.False(t, cfg.DeferAutoUpdate)
	assert.Equal(t, uint64(336), cfg.MaxAge)
	assert.Empty(t, cfg.Languages)
}

func TestDefaultIndentConfig(t *testing.T) {
	cfg := DefaultIndentConfig()

	assert.Equal(t, 2, cfg.Title)
	assert.Equal(t, 2, cfg.Description)
	assert.Equal(t, 2, cfg.Bullet)
	assert.Equal(t, 4, cfg.Example)
}

func TestDefaultOutputConfig(t *testing.T) {
	cfg := DefaultOutputConfig()

	assert.True(t, cfg.ShowTitle)
	assert.False(t, cfg.PlatformTitle)
	assert.False(t, cfg.ShowHyphens)
	assert.False(t, cfg.EditLink)
	assert.Equal(t, "- ", cfg.ExamplePrefix)
	assert.Equal(t, 0, cfg.LineLength)
	assert.False(t, cfg.Compact)
	assert.Equal(t, OptionStyleLong, cfg.OptionStyle)
	assert.False(t, cfg.RawMarkdown)
}

func TestOptionStyleConstants(t *testing.T) {
	assert.Equal(t, OptionStyle("short"), OptionStyleShort)
	assert.Equal(t, OptionStyle("long"), OptionStyleLong)
	assert.Equal(t, OptionStyle("both"), OptionStyleCombined)
}

func TestDefaultStyleConfig(t *testing.T) {
	cfg := DefaultStyleConfig()

	magenta := OutputColor{Kind: ColorKindNamed, Named: ColorMagenta}
	green := OutputColor{Kind: ColorKindNamed, Named: ColorGreen}
	cyan := OutputColor{Kind: ColorKindNamed, Named: ColorCyan}
	red := OutputColor{Kind: ColorKindNamed, Named: ColorRed}
	yellow := OutputColor{Kind: ColorKindNamed, Named: ColorYellow}
	defBg := DefaultColor()

	assert.Equal(t, OutputStyle{Color: magenta, Background: defBg, Bold: true}, cfg.Title)
	assert.Equal(t, OutputStyle{Color: magenta, Background: defBg}, cfg.Description)
	assert.Equal(t, OutputStyle{Color: green, Background: defBg}, cfg.Bullet)
	assert.Equal(t, OutputStyle{Color: cyan, Background: defBg}, cfg.Example)
	assert.Equal(t, OutputStyle{Color: red, Background: defBg, Italic: true}, cfg.URL)
	assert.Equal(t, OutputStyle{Color: yellow, Background: defBg, Italic: true}, cfg.InlineCode)
	assert.Equal(t, OutputStyle{Color: red, Background: defBg, Italic: true}, cfg.Placeholder)
}
