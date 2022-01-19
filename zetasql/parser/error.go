package parser

import (
	"github.com/paulourio/bqfmt/zetasql/ast"
	zerrors "github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func NewSyntaxError(pos Attrib, msg string) (Attrib, error) {
	tok := pos.(*token.Token)
	start := tok.Offset
	end := start + len(tok.Lit)

	return nil, zerrors.NewSyntaxError(ast.Loc{Start: start, End: end}, msg)
}

func NewInternalError(pos Attrib, msg string) (Attrib, error) {
	barePos, _ := unwrap(pos)

	tok := barePos.(*token.Token)
	start := tok.Offset
	end := start + len(tok.Lit)

	return nil, zerrors.NewInternalError(ast.Loc{Start: start, End: end}, msg)
}
