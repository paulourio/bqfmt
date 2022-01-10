package parser

import (
	"github.com/paulourio/bqfmt/zetasql/ast"
	"github.com/paulourio/bqfmt/zetasql/errors"
)

// ParseIdentifier parses a ZetaSQL identifier.  The string s may be
// quoted with backticks.  If so, unquote and unescape it.
func ParseIdentifier(s string, allowReservedKw bool) (string, error) {
	var quoted bool

	if len(s) > 0 && s[0] == '`' {
		quoted = true
	}

	if quoted {
		const quotesLength = 1 // Starts after the opening quote '`'.`

		if len(s) == quotesLength*2 {
			return "", errors.ErrEmptyIdentifier
		}

		unquoted := s[1 : len(s)-1]

		// TODO: unescape
		return unquoted, nil
	} else {
		if !ast.IsValidUnquotedIdentifier(s) && !allowReservedKw {
			return "", errors.ErrInvalidIdentifier
		}

		return s, nil
	}
}
