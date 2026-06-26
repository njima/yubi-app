package stack

import (
	"fmt"
	"runtime"
)

// New returns the call stack entry at the given skip depth as a "pkg.FuncName:line" string.
// skip is passed directly to runtime.Caller:
//   - skip=0: the New function itself
//   - skip=1: the caller of New
//   - skip=2: the caller's caller, and so on
//
// Returns an empty string if the stack information cannot be retrieved.
func New(skip int) string {
	if pc, _, line, ok := runtime.Caller(skip); ok {
		return fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
	}

	return ""
}

// Stacker is an interface for types that can return a stack trace as a slice of strings.
// Implement this on error types to expose the call stack at the point of error creation.
type Stacker interface {
	Stack() []string
}
