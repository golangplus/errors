package errors

import (
	"fmt"
	"path"
	"runtime"
)

type ErrorWithStacks struct {
	Err    error
	Stacks []string
}

var _ error = (*ErrorWithStacks)(nil)

func (e *ErrorWithStacks) Error() string {
	msg := e.Err.Error()
	for _, s := range e.Stacks {
		msg += "\n  at " + s
	}
	return msg
}

func WithStacks(err error) error {
	if err == nil {
		return nil
	}
	e := &ErrorWithStacks{
		Err: err,
	}
	for i := 0; i < 10; i++ {
		_, file, line, ok := runtime.Caller(i + 1)
		if !ok {
			break
		}

		e.Stacks = append(e.Stacks, fmt.Sprintf("%s:%d", path.Base(file), line))
	}
	return e
}
