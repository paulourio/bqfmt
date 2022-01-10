package ast

import "errors"

var (
	ErrMissingRequiredField    = errors.New("missing required field")
	ErrFieldAlreadyInitialized = errors.New("field already initialized")
)
