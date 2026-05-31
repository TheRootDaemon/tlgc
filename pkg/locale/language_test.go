package locale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLanguages(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		language string
		want     []string
	}{
		{
			name:     "simple LANG only",
			lang:     "cz",
			language: "",
			want:     []string{"cz"},
		},
		{
			name:     "LANGUAGE and LANG",
			lang:     "cz",
			language: "it:cz:de",
			want:     []string{"it", "cz", "de", "cz"},
		},
		{
			name:     "different order",
			lang:     "cz",
			language: "it:de:fr",
			want:     []string{"it", "de", "fr", "cz"},
		},
		{
			name:     "no LANG set",
			lang:     "",
			language: "it:cz",
			want:     nil,
		},
		{
			name:     "no LANG and no LANGUAGE",
			lang:     "",
			language: "",
			want:     nil,
		},
		{
			name:     "locale format with UTF-8 suffix",
			lang:     "en_US.UTF-8",
			language: "de_DE.UTF-8:pl:en",
			want:     []string{"de_DE", "de", "pl", "en", "en_US", "en"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LANG", tt.lang)
			t.Setenv("LANGUAGE", tt.language)

			var out []string
			GetLanguages(&out)

			if tt.want == nil {
				assert.Empty(t, out)
			} else {
				assert.Equal(t, tt.want, out)
			}
		})
	}
}
