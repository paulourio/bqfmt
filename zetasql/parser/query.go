package parser

import "github.com/paulourio/bqfmt/zetasql/ast"

func SetQueryParenthesized(open, query, close Attrib) (Attrib, error) {
	q := query.(*ast.Query)

	q.SetParenthesized(true)

	return WrapWithLoc(q, open, close)
}

func ValidateParenthesizedSubquery(expr Attrib) (Attrib, error) {
	n := expr.(ast.NodeHandler)
	if n.Kind() != ast.ExpressionSubqueryKind {
		return NewSyntaxError(
			expr,
			"Parenthesized expression cannot be parsed as an expression, "+
				"struct constructor, or subquery")
	}

	e := n.(*ast.ExpressionSubquery)
	e.Query.SetParenthesized(true)
	return e, nil
}

func ExtractQueryFromSubquery(expr Attrib) (Attrib, error) {
	subquery := expr.(*ast.ExpressionSubquery)

	if subquery.Query == nil {
		return NewInternalError(expr, "expected query child of subquery")
	}

	return subquery.Query, nil
}

func NewSubqueryOrInList(open, expr, close Attrib) (Attrib, error) {
	node := expr.(ast.NodeHandler)

	if subquery, ok := node.(*ast.ExpressionSubquery); ok {
		if subquery.Modifier == ast.NoSubqueryModifier {
			// To match Bison and JavaCC parser, we prefer interpreting
			// IN ((Query)) as IN (Query) with a parenthesized Query,
			// not a value IN list containing a scalar expression query.
			// Return the contained Query, wrapped in another Query to
			// replace the parentheses.
			query := subquery.Query

			if query == nil {
				return NewInternalError(
					expr, "expected query child of parenthesized subquery")
			}

			query.SetParenthesized(true)

			return ast.NewQuery(nil, query, nil, nil)
		} else {
			// The ExpressionSubquery is an EXISTS or ARRAY subquery,
			// which is a scalar expression and is not interpreted as
			// a Query. Treat this as an InList with a single element.
			// Do not include the paretheses in the location, to match
			// the JavaCC parser.
			return ast.NewInList(subquery)
		}
	} else {
		// Do not include the parentheses in the location, to match the
		// JavaCC parser.
		return ast.NewInList(expr)
	}
}

func NewSetOperation(a, setop, allOrDistinct, b Attrib) (Attrib, error) {
	var (
		op ast.SetOp
		distinct bool
	)

	if kw, ok := setop.(*ast.Wrapped); ok {
		op = kw.Value.(ast.SetOp)
	} else {
		op = setop.(ast.SetOp)
	}

	if kw, ok := allOrDistinct.(*ast.Wrapped); ok {
		distinct = kw.Value == DistinctKeyword
	} else {
		distinct = allOrDistinct.(allOrDistinctKeyword) == DistinctKeyword
	}

	if s, ok := a.(*ast.SetOperation); ok {
		if s.OpType != op || s.Distinct != distinct {
			return NewSyntaxError(
				setop,
				"Different set operations cannot be used in the same query "+
				"without using parentheses for grouping")
		}

		return WithExtraChild(s, b)
	}

	return ast.NewSetOperation(List(a, b), op, distinct)
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

type allOrDistinctKeyword int

const (
	NoAllOrDistinctKeyword allOrDistinctKeyword = iota
	AllKeyword
	DistinctKeyword
)
