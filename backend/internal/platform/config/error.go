package config

import (
	"fmt"

	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror/stack"
)

type ErrorKind int

const (
	ErrorKindUnknown ErrorKind = iota
	ErrorKindProcess
)

type Error struct {
	error
	Kind  ErrorKind
	stack []string
}

func newError(kind ErrorKind, err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		error: err,
		Kind:  kind,
		stack: []string{stack.New(2)},
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("ConfigError(Kind: %d): %v", e.Kind, e.error)
}

func (e Error) Stack() []string {
	return e.stack
}
