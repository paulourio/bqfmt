package errors

import (
	"errors"
	"fmt"

	"github.com/paulourio/bqfmt/zetasql/ast"
)

var (
	// ErrMalformedParser indicates a bug in the parser logic.
	ErrMalformedParser         = errors.New("malformed parser")
	ErrInvalidDashedIdentifier = errors.New("invalid dashed identifier")
	ErrEmptyIdentifier         = errors.New("invalid empty identifier")
	ErrInvalidIdentifier       = errors.New("invalid identifier")
)

type SyntaxError struct {
	Loc ast.Loc
	Msg string
}

func NewSyntaxError(loc ast.Loc, msg string) *SyntaxError {
	return &SyntaxError{loc, msg}
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error: %s", e.Msg)
}
