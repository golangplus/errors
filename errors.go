package errorsp

import (
	"fmt"
	"math"
	"path"
	"runtime"
)

var (
	// The maximum depth in ErrorWithStacks.Stacks.
	// The last line is set to "..." if some call stacks are ignored.
	// Only set this if some very deep callstack can happen, e.g. deep recursive calling.
	MaxStackDepth int = math.MaxInt32
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

// Returns the root cause of the error. Will not return a *ErrorWithStacks.
func Cause(err error) error {
	for {
		ews, ok := err.(*ErrorWithStacks)
		if !ok {
			return err
		}
		err = ews.Err
	}
}

func stacks(skip int) []string {
	var stacks []string
	for i := 0; i <= MaxStackDepth; i++ {
		pc, file, line, ok := runtime.Caller(i + skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)

		if i >= MaxStackDepth {
			stacks = append(stacks, "...")
		} else {
			stacks = append(stacks, fmt.Sprintf("%s(%s:%d)", fn.Name(), path.Base(file), line))
		}
	}
	return stacks
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
		// If err has been an ErrorWithStacks, no need to wrap it.
		return err
	}
	return &ErrorWithStacks{
		Err:    err,
		Stacks: stacks(2),
	}
}

func WithStacksAndMessage(err error, format string, args ...interface{}) error {
	if err == nil {
		// Remain no-error.
		return nil
	}
	s := stacks(2)
	ews, ok := err.(*ErrorWithStacks)
	if !ok || len(ews.Stacks) < len(s) {
		s[0] += ": " + fmt.Sprintf(format, args...)
		return &ErrorWithStacks{
			Err:    err,
			Stacks: s,
		}
	}
	ews.Stacks[len(ews.Stacks)-len(s)] = ": " + fmt.Sprintf(format, args...)
	return ews
}
