package errs

import (
	"testing"

	"github.com/jf-tech/omniparser/errs"
	"github.com/stretchr/testify/assert"
)

func TestIsNonErrorRecordSkipped(t *testing.T) {
	assert.True(t, IsNonErrorRecordSkipped(NonErrorRecordSkipped("test")))
	assert.Equal(t, "test", NonErrorRecordSkipped("test").Error())
	assert.False(t, IsNonErrorRecordSkipped(errs.ErrTransformFailed("test")))
}
