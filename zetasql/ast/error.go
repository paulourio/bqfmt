package ast

import "errors"

var (
	// ErrMalformedParser indicates a bug in the parser logic.
	ErrMalformedParser         = errors.New("malformed parser")
	ErrMissingRequiredField    = errors.New("missing required field")
	ErrInvalidDashedIdentifier = errors.New("invalid dashed identifier")
	ErrFieldAlreadyInitialized = errors.New("field already initialized")
)
