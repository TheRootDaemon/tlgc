package browser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid_command",
			command: "echo",
			args:    []string{"hello"},
			wantErr: false,
		},
		{
			name:    "invalid_command",
			command: "nonexistent-command-xyz",
			args:    nil,
			wantErr: true,
		},
		{
			name:    "command_with_empty_args",
			command: "echo",
			args:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := browse(tt.command, tt.args...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
