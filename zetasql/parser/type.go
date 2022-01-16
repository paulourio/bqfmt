package parser

import (
	"fmt"
	"reflect"

	"github.com/paulourio/bqfmt/zetasql/ast"
	"github.com/paulourio/bqfmt/zetasql/errors"
)

func NewNamedType(typ, params Attrib) (Attrib, error) {
	switch t := params.(type) {
	case nil:
		return typ, nil
	case *ast.TypeParameterList:
		n := typ.(ast.TypeHandler)

		n.SetTypeParameters(t)
		n.AddChild(t)

		return n, nil
	}

	return nil, fmt.Errorf("%w: cannot create NamedType from %v",
		errors.ErrMalformedParser, reflect.TypeOf(params))
}

func NewArrayType(array, elemType, close Attrib) (Attrib, error) {
	t, err := ast.NewArrayType(elemType)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(t, array, close)
}

func NewStructType(structWord, field, close Attrib) (Attrib, error) {
	t, err := ast.NewStructType(field)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(t, structWord, close)
}
