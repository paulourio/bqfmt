package ast

import "errors"

var (
	ErrMissingRequiredField    = errors.New("missing required field")
	ErrInvalidDashedIdentifier = errors.New("invalid dashed identifier")
	ErrFieldAlreadyInitialized = errors.New("field already initialized")
)
