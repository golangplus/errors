package errorsp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golangplus/testing/assert"
)

func TestWithStacks_NoError(t *testing.T) {
	assert.Equal(t, "WithStacks(nil)", WithStacks(nil), nil)
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

func TestNewWithStacks(t *testing.T) {
	err := NewWithStacks("%s:%d", "a", 123).(*ErrorWithStacks)
	assert.Equal(t, "err.Err.Error()", err.Err.Error(), "a:123")
}

func TestWithStacksAndMessage_NoError(t *testing.T) {
	assert.Equal(t, "WithStacksAndMessage(nil)", WithStacksAndMessage(nil, ""), nil)
}

func TestWithStacksAndMessage(t *testing.T) {
	origin := fmt.Errorf("origin")
	wrapped, isErrWithStacks := WithStacksAndMessage(origin, "%s:%d", "a", 123).(*ErrorWithStacks)
	assert.True(t, "isErrWithStacks", isErrWithStacks)
	assert.Equal(t, "wrapped.Err", wrapped.Err, origin)
	assert.ValueShould(t, "wrapped.Stacks[0]", wrapped.Stacks[0], strings.HasSuffix(wrapped.Stacks[0], ": a:123"), "does not contain message correctly")

	wrapped2, isErrWithStacks := WithStacksAndMessage(wrapped, "%s:%d", "b", 456).(*ErrorWithStacks)
	assert.True(t, "isErrWithStacks", isErrWithStacks)
	assert.Equal(t, "wrapped.Err", wrapped2.Err, origin)
	assert.ValueShould(t, "wrapped2.Stacks[0]", wrapped2.Stacks[0], strings.HasSuffix(wrapped2.Stacks[0], ": a:123: b:456"), "does not contain message correctly")
}

func TestErrorPrintf(t *testing.T) {
	assert.Equal(t, "errorPrintf", errorPrintf("abc: %s", "a\nb"), `abc: a
  b`)
}
