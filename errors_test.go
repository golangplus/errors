package errorsp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golangplus/testing/assert"
)

func TestWithStacks_NoError(t *testing.T) {
	assert.Equal(t, "WithStacks", WithStacks(nil), nil)
}

func TestWithStacks(t *testing.T) {
	origin := fmt.Errorf("myerror")
	wrapped, isErrWithStacks := WithStacks(origin).(*ErrorWithStacks)
	assert.True(t, "isErrWithStacks", isErrWithStacks)
	assert.Equal(t, "wrapped.Err", wrapped.Err, origin)

	wrapped2, isErrWithStacks := WithStacks(wrapped).(*ErrorWithStacks)
	assert.True(t, "isErrWithStacks", isErrWithStacks)
	assert.Equal(t, "wrapped2", wrapped2, wrapped)
}

func TestMaxStackDepth(t *testing.T) {
	MaxStackDepth = 1
	err := WithStacks(fmt.Errorf("myerror")).(*ErrorWithStacks)
	assert.Equal(t, "len(err.Stacks)", len(err.Stacks), 2)
	assert.Equal(t, "err.Stacks[1]", err.Stacks[1], "...")
	assert.Equal(t, "lines of err.Error()", len(strings.Split(err.Error(), "\n")), 3)
}

func TestCause(t *testing.T) {
	cause := fmt.Errorf("myerror")
	assert.Equal(t, "Cause(cause)", Cause(cause), cause)
	assert.Equal(t, "Cause(WithStacks(cause))", Cause(WithStacks(cause)), cause)
}
