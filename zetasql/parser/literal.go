package parser

import "github.com/paulourio/bqfmt/zetasql/ast"

func NewBigNumericLiteral(prefix, str Attrib) (Attrib, error) {
	lit, err := ast.NewBigNumericLiteral()
	if err != nil {
		return nil, err
	}

	value, err := InitLiteral(lit, str)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(value, prefix)
}

func NewNumericLiteral(prefix, str Attrib) (Attrib, error) {
	lit, err := ast.NewNumericLiteral()
	if err != nil {
		return nil, err
	}

	value, err := InitLiteral(lit, str)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(value, prefix)
}
