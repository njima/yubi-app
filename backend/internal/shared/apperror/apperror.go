package apperror

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// Detail holds a single error code and a developer-facing message.
type Detail struct {
	code                Code
	messageForDeveloper string
}

// NewMessage creates a Detail with the given Code and a formatted developer message.
func NewMessage(c Code, format string, val ...any) Detail {
	return Detail{
		code:                c,
		messageForDeveloper: fmt.Sprintf(format, val...),
	}
}

// Error is the application error type that carries one or more Detail entries
// along with a stack trace and an optional wrapped cause.
type Error struct {
	details []Detail
	err     error
	frames  *runtime.Frames
}

// Error implements the error interface, returning a concatenation of all developer messages in the chain.
func (e *Error) Error() string {
	var buf bytes.Buffer
	var originalErr error = e
	for {
		appErr := ConvertAppError(originalErr)
		if appErr == nil {
			io.WriteString(&buf, originalErr.Error())
			break
		} else {
			io.WriteString(&buf, appErr.messageForDeveloper())
		}

		originalErr = unWrap(originalErr)
		if originalErr == nil {
			break
		}
	}
	return buf.String()
}

func (e *Error) messageForDeveloper() string {
	var buf bytes.Buffer
	for _, m := range e.details {
		if m.messageForDeveloper == "" {
			continue
		}
		io.WriteString(&buf, m.messageForDeveloper)
		io.WriteString(&buf, "/")
	}
	return buf.String()
}

// NewError creates a new Error with the given Detail entries and no wrapped cause.
func NewError(d ...Detail) error {
	return newError(nil, d...)
}

func newError(err error, d ...Detail) error {
	return &Error{
		details: d,
		err:     err,
		frames:  newFrames(),
	}
}

// Wrap wraps an existing error as an apperror.Error with no additional detail.
func Wrap(err error) error {
	return newError(err, NewMessage(CodeEmpty, ""))
}

// WrapWithDevMessage wraps an existing error and attaches a formatted developer message with no Code.
func WrapWithDevMessage(err error, format string, val ...any) error {
	return newError(err, NewMessage(CodeEmpty, format, val...))
}

// WrapWithMessage wraps an existing error and attaches one or more Detail entries.
func WrapWithMessage(err error, m ...Detail) error {
	return newError(err, m...)
}

// GetCodes walks the error chain and returns all Code values from contiguous apperror.Error layers
// that share the same Kind. Stops at the first Kind boundary or non-apperror error.
func GetCodes(err error) []Code {
	codes := []Code{}
	currentErr := err
	prevErr := NewError(NewMessage(CodeEmpty, "")).(*Error)
	for {
		appErr, yes := currentErr.(*Error)
		if !yes {
			break
		}

		if appErr.details[0].code.Kind == KindEmpty {
			currentErr = unWrap(appErr)
			prevErr = appErr
			continue
		}

		if prevErr.details[0].code.Kind != KindEmpty {
			if appErr.details[0].code.Kind != prevErr.details[0].code.Kind {
				break
			}
		}

		codes = append(codes, appErr.getCodes()...)

		currentErr = unWrap(appErr)
		prevErr = appErr
	}

	return codes
}

func unWrap(err error) error {
	appErr := ConvertAppError(err)
	if appErr == nil {
		return nil
	}
	return appErr.err
}

// ConvertAppError type-asserts err to *Error. Returns nil if err is not an *Error.
func ConvertAppError(err error) *Error {
	appErr, yes := err.(*Error)
	if !yes {
		return nil
	}
	return appErr
}

func (e *Error) getCodes() []Code {
	codes := make([]Code, 0)
	for _, m := range e.details {
		if m.code.Kind == KindEmpty {
			continue
		}
		codes = append(codes, m.code)
	}
	return codes
}

// StackTrace returns the stack trace of the innermost apperror.Error in the chain.
// Returns "not apperror.Error" if err is not an *Error.
func StackTrace(err error) string {
	appErr := ConvertAppError(err)
	if appErr == nil {
		return "not apperror.Error"
	}
	return appErr.StackTrace()
}

// StackTrace returns the formatted stack trace captured at the creation of the innermost apperror.Error.
func (e *Error) StackTrace() string {
	originalErr := e

	for {
		err := unWrap(originalErr)
		if err == nil {
			break
		}
		appErr, yes := err.(*Error)
		if !yes {
			break
		}
		originalErr = appErr
	}

	var buf bytes.Buffer
	for {
		frame, more := originalErr.frames.Next()

		io.WriteString(&buf, frame.Function)
		io.WriteString(&buf, "\n\t")
		io.WriteString(&buf, frame.File)
		io.WriteString(&buf, ":")
		io.WriteString(&buf, strconv.Itoa(frame.Line))

		if !more {
			break
		}
	}

	return buf.String()
}

func newFrames() *runtime.Frames {
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)
	pc = pc[:n]
	return runtime.CallersFrames(pc)
}

// SameKind reports whether the outermost meaningful Kind in err matches k.
func SameKind(err error, k Kind) bool {
	appErr := ConvertAppError(err)
	if appErr == nil {
		return false
	}

	for _, m := range appErr.details {
		if m.code.Kind == KindEmpty {
			return SameKind(appErr.err, k)
		}
		if m.code.Kind == k {
			return true
		}
	}
	return false
}

// ErrorResponse is the transport-neutral error response shape used by HTTP adapters.
type ErrorResponse struct {
	Code    int
	Message string
}

// NewErrorResponse converts a slice of Code values into an error response.
// Uses the first Code's Kind to determine the HTTP status and message.
// Falls back to 500 Internal Server Error when cs is empty.
func NewErrorResponse(cs []Code) ErrorResponse {
	var status int
	var message string

	if len(cs) > 0 {
		if cs[0].MessageForEndUser != "" {
			message = cs[0].MessageForEndUser
		} else {
			message = cs[0].Kind.ResponseWrapErrorMessage()
		}
		status = cs[0].Kind.HTTPStatus()
	} else {
		status = KindInternal.HTTPStatus()
		message = KindInternal.ResponseWrapErrorMessage()
	}

	return ErrorResponse{
		Code:    status,
		Message: message,
	}
}
