package errorsp

import (
	"fmt"
	"path"
	"runtime"
)

var (
	// The maximum depth in ErrorWithStacks.Stacks.
	// The last line is set to "..." if some call stacks are ignored.
	MaxStackDepth int = 10
)

// ErrorWithStacks is a struct containing the original error and the call stacks.
type ErrorWithStacks struct {
	Err    error
	Stacks []string
}

var _ error = (*ErrorWithStacks)(nil)

// Error implements the error interface.
func (e *ErrorWithStacks) Error() string {
	msg := e.Err.Error()
	for _, s := range e.Stacks {
		msg += "\n  at " + s
	}
	return msg
}

// WithStacks returns a *ErrorWithStacks error with stacks set.
// If err has been a *ErrorWithStacks, it is directly returned.
// If err is nil, a nil is returned.
func WithStacks(err error) error {
	if err == nil {
		// Remain no-error.
		return nil
	}
	if _, ok := err.(*ErrorWithStacks); ok {
		// If err has been a ErrorWithStacks, no need to wrap it.
		return err
	}
	e := &ErrorWithStacks{
		Err: err,
	}
	for i := 0; i <= MaxStackDepth; i++ {
		_, file, line, ok := runtime.Caller(i + 1)
		if !ok {
			break
		}

		if i >= MaxStackDepth {
			e.Stacks = append(e.Stacks, "...")
		} else {
			e.Stacks = append(e.Stacks, fmt.Sprintf("%s:%d", path.Base(file), line))
		}
	}
	return e
}
