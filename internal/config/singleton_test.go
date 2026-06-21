package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetCurrentConfig clears the singleton so tests start from a clean state.
func resetCurrentConfig() {
	currentConfig.Store(nil)
}

func TestC_Defaults(t *testing.T) {
	resetCurrentConfig()
	defer resetCurrentConfig()

	cfg := C()
	require.NotNil(t, cfg)
	assert.Equal(t, Default(), *cfg)
}

func TestInitialize_and_C(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte("[output]\nshow_title = false\n"), 0o644)
	require.NoError(t, err)

	resetCurrentConfig()
	defer resetCurrentConfig()

	t.Setenv("TLGC_CONFIG", path)
	err = Initialize()
	require.NoError(t, err)

	cfg := C()
	assert.False(t, cfg.Output.ShowTitle)
	assert.Equal(t, DefaultCacheConfig(), cfg.Cache)
}

func TestInitialize_Error(t *testing.T) {
	resetCurrentConfig()
	defer resetCurrentConfig()

	t.Setenv("TLGC_CONFIG", "/nonexistent/path/config.toml")
	err := Initialize()
	require.Error(t, err)
}

func TestCache_Accessor(t *testing.T) {
	resetCurrentConfig()
	defer resetCurrentConfig()

	cc := Cache()
	assert.Equal(t, DefaultCacheConfig(), cc)
}

func TestStyle_Accessor(t *testing.T) {
	resetCurrentConfig()
	defer resetCurrentConfig()

	sc := Style()
	assert.Equal(t, DefaultStyleConfig(), sc)
}

func TestIndent_Accessor(t *testing.T) {
	resetCurrentConfig()
	defer resetCurrentConfig()

	ic := Indent()
	assert.Equal(t, DefaultIndentConfig(), ic)
}

func TestOutput_Accessor(t *testing.T) {
	resetCurrentConfig()
	defer resetCurrentConfig()

	oc := Output()
	assert.Equal(t, DefaultOutputConfig(), oc)
}

func TestInitialize_then_Accessors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	input := `[cache]
dir = "/custom/cache"
languages = ["en", "de"]

[output]
show_title = false
`
	err := os.WriteFile(path, []byte(input), 0o644)
	require.NoError(t, err)

	resetCurrentConfig()
	defer resetCurrentConfig()

	t.Setenv("TLGC_CONFIG", path)
	err = Initialize()
	require.NoError(t, err)

	// Accessors should reflect the loaded config.
	cc := Cache()
	assert.Equal(t, "/custom/cache", cc.Dir)
	assert.Equal(t, []string{"en", "de"}, cc.Languages)

	oc := Output()
	assert.False(t, oc.ShowTitle)

	// Unchanged sections should still have defaults.
	sc := Style()
	assert.Equal(t, DefaultStyleConfig(), sc)
}
