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

func NewNullLiteral(null Attrib) (Attrib, error) {
	lit, err := ast.NewNullLiteral()
	if err != nil {
		return nil, err
	}

	return InitLiteral(lit, null)
}

func NewJSONLiteral(json, content Attrib) (Attrib, error) {
	lit, err := ast.NewJSONLiteral()
	if err != nil {
		return nil, err
	}

	value, err := InitLiteral(lit, content)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(value, json)
}
