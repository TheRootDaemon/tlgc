package exitcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitCodes(t *testing.T) {
	assert.Equal(t, 0, Success)
	assert.Equal(t, 1, InvalidArguments)
	assert.Equal(t, 2, Default)
	assert.Equal(t, 3, CacheUpdateFailure)
	assert.Equal(t, 4, LintFailure)
	assert.Equal(t, 5, ConfigParseFailure)
	assert.Equal(t, 6, PageNotFound)
}
