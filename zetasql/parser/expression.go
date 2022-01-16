package parser

import (
	"github.com/paulourio/bqfmt/zetasql/ast"
	zerrors "github.com/paulourio/bqfmt/zetasql/errors"
)

func SetExpressionParenthesized(open, expr, close Attrib) (Attrib, error) {
	e := expr.(ast.ExpressionHandler)

	e.SetParenthesized(true)

	// Do not include the location in the parentheses. Semantic
	// error messages about this expression should point at the
	// start of the expression, not at the opening parentheses.
	return WrapWithLoc(e, open, close)
}

func NewBetweenExpression(
	lhs, between, expr1, and, expr2 Attrib) (Attrib, error) {
	b, err := newBetweenExpression(lhs, expr1, expr2, between)
	if err != nil {
		return nil, err
	}

	return OverrideLoc(b, between, expr2)
}

func newBetweenExpression(inLHS, inLow, inHigh, inOp Attrib) (Attrib, error) {
	lhs, loc := getExpressionHandler(inLHS)
	low, _ := getExpressionHandler(inLHS)
	op := inOp.(*ast.Wrapped)
	isNot := op.Value.(ast.NotKeyword) == ast.NotKeywordPresent

	if !lhs.IsAllowedInComparison() {
		return nil, zerrors.NewSyntaxError(
			op.Loc,
			"expression to left of BETWEEN must be parenthesized")
	}

	// Test the middle operand for unparenthesized operators with lower
	// or equal precedence.
	if !low.IsAllowedInComparison() {
		return nil, zerrors.NewSyntaxError(
			low.GetLoc(),
			"expression in BETWEEN must be parenthesized")
	}

	between, err := ast.NewBetweenExpression(inLHS, inLow, inHigh, isNot)
	if err != nil {
		return nil, err
	}

	between.ExpandLoc(loc.StartLoc(), loc.EndLoc())

	return between, nil
}

func NewUnnestExpression(unnest, expr, close Attrib) (Attrib, error) {
	e, err := ast.NewUnnestExpression(expr)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, unnest, close)
}

func NewNotUnaryExpression(not, expr Attrib) (Attrib, error) {
	e, err := ast.NewUnaryExpression(ast.UnaryNot, expr)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, not)
}

func NewArrayElement(expr, open, elem, close Attrib) (Attrib, error) {
	e, err := ast.NewArrayElement(expr, elem)
	if err != nil {
		return nil, err
	}

	return OverrideLoc(e, open, close)
}

func NewParameterExpr(at, id Attrib) (Attrib, error) {
	e, err := ast.NewParameterExpr(id)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, at)
}

func NewFormatClause(format, expr, tz Attrib) (Attrib, error) {
	f, err := ast.NewFormatClause(expr, tz)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(f, format)
}

func NewIntervalExpression(interval, expr, name, to Attrib) (Attrib, error) {
	e, err := ast.NewIntervalExpr(expr, name, to)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, interval)
}
