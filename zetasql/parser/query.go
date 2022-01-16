package parser

import "github.com/paulourio/bqfmt/zetasql/ast"

func SetQueryParenthesized(open, query, close Attrib) (Attrib, error) {
	q := query.(*ast.Query)

	q.SetParenthesized(true)

	return WrapWithLoc(q, open, close)
}

func NewTableSubquery(
	open, query, close, pivotUnpivotAlias, sample Attrib) (Attrib, error) {
	p := pivotUnpivotAlias.(*pivotOrUnpivotAndAlias)

	t, err := ast.NewTableSubquery(
		query, p.Alias, p.PivotClause, p.UnpivotClause, sample)
	if err != nil {
		return nil, err
	}

	t.Subquery.IsNested = true

	return UpdateLoc(t, open, close)
}

func PivotUnpivotAliasWithAlias(as, alias Attrib) (Attrib, error) {
	a, err := UpdateLoc(alias, as)
	if err != nil {
		return nil, err
	}

	return &pivotOrUnpivotAndAlias{Alias: a}, nil
}

// pivotOrUnpivotAndAlias is a temporary holder of information when
// constructing a table path expression.
type pivotOrUnpivotAndAlias struct {
	Alias         Attrib
	PivotClause   Attrib
	UnpivotClause Attrib
}

func NewLimitOffset(limitTok, limit, offset Attrib) (Attrib, error) {
	c, err := ast.NewLimitOffset(limit, offset)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(c, limitTok)
}

func NewOrderBy(order, expr Attrib) (Attrib, error) {
	c, err := ast.NewOrderBy(expr)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(c, order)
}

func NewNullOrder(nulls, order Attrib, nullsFirst bool) (Attrib, error) {
	e, err := ast.NewNullOrder(nullsFirst)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, nulls, order)
}

func NewWithOffset(with, offset, alias Attrib) (Attrib, error) {
	e, err := ast.NewWithOffset(alias)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, with, offset)
}

func NewCastIntLiteralOrParam(
	cast, expr, typ, format, close Attrib) (Attrib, error) {
	const isSafeCast = false

	e, err := ast.NewCastExpression(expr, typ, format, isSafeCast)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, cast, close)
}
