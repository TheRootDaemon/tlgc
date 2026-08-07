package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLint(t *testing.T) {
	clean := "# App\n\n> Description.\n\n- List files:\n\n`ls`\n"
	missingPeriod := "# App\n\n> Description\n\n- List files:\n\n`ls`\n"

	tests := []struct {
		name     string
		content  string
		badName  bool
		ignore   []string
		wantCode []string
		line     int
	}{
		{"clean page reports nothing", clean, false, nil, nil, 0},
		{"missing period", missingPeriod, false, nil, []string{"TLDR004"}, 2},
		{"ignored rule skipped", missingPeriod, false, []string{"TLDR004"}, nil, 0},
		{"filename rules run on bad name", clean, true, nil, []string{"TLDR107", "TLDR108", "TLDR109"}, 0},
		{"bad ignore code is harmless", clean, false, []string{"TLDR999"}, nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := "app.md"
			if tt.badName {
				name = "Bad File.txt"
			}

			path := filepath.Join(t.TempDir(), name)
			require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o644))

			f, err := os.Open(path)
			require.NoError(t, err)
			defer func() {
				_ = f.Close()
			}()

			r, err := Lint(f, tt.ignore...)
			require.NoError(t, err)
			require.Equal(t, tt.wantCode, errorCodes(r))
			if tt.line != 0 {
				require.Equal(t, tt.line, r.Errors[0].Line)
			}
		})
	}
}
