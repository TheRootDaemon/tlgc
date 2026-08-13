package app

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/stretchr/testify/require"
)

// validPage is a valid tldr page used by app tests.
const validPage = "# tar\n\n> Archiving utility.\n\n- Create an archive:\n\n`tar cf archive.tar`\n"

// writePage writes content to a file in dir and returns its path.
func writePage(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

// newTestApp creates an App with buffers
// for capturing standard output and standard error.
func newTestApp() (*App, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := &App{
		Stdout: stdout,
		Stderr: stderr,
	}

	return app, stdout, stderr
}

func TestMain(m *testing.M) {
	logger.SetDefault(logger.NewWithWriter(true, 0, io.Discard))
	os.Exit(m.Run())
}
