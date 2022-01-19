package errors

import (
	"fmt"

	"github.com/paulourio/bqfmt/zetasql/ast"
)

type InternalError struct {
	Loc ast.Loc
	Msg string
}

func NewInternalError(loc ast.Loc, msg string) *InternalError {
	return &InternalError{loc, msg}
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s", e.Msg)
}
