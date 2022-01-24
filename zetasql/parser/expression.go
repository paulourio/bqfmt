package parser

import (
	"fmt"
	"reflect"

	"github.com/paulourio/bqfmt/zetasql/ast"
	zerrors "github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func SetExpressionParenthesized(open, expr, close Attrib) (Attrib, error) {
	e, loc := getExpressionHandler(expr)

	e.SetParenthesized(true)

	// Do not include the location in the parentheses. Semantic
	// error messages about this expression should point at the
	// start of the expression, not at the opening parentheses.
	return WrapWithLoc(e, open, close, loc)
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

func NewExpressionSubquery(
	modifierTok, open, query, close, modifier Attrib) (Attrib, error) {
	e, err := ast.NewExpressionSubquery(query, modifier)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, modifierTok, open, close)
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

func NewWindowFrameClause(unit, start, end Attrib) (Attrib, error) {
	return ast.NewWindowFrame(start, end, unit)
}

func NewAndExpr(lhs, rhs Attrib) (Attrib, error) {
	if left, ok := lhs.(*ast.AndExpr); ok && !left.IsParenthesized() {
		// Embrace and extend the lhs node to flatten a series of ANDs.
		return WithExtraChild(left, rhs)
	}

	return ast.NewAndExpr(List(lhs, rhs))
}

func NewOrExpr(lhs, rhs Attrib) (Attrib, error) {
	if left, ok := lhs.(*ast.OrExpr); ok && !left.IsParenthesized() {
		// Embrace and extend the lhs node to flatten a series of ORs.
		return WithExtraChild(left, rhs)
	}

	return ast.NewOrExpr(List(lhs, rhs))
}

func NewFunctionCall(
	expr, nulls, orderby, limit, closetok Attrib,
) (Attrib, error) {
	f, _ := getFunctionCall(expr)

	if nulls != nil {
		err := f.InitNullHandlingModifier(nulls)
		if err != nil {
			return nil, err
		}
	}

	if orderby != nil {
		err := f.InitOrderBy(orderby)
		if err != nil {
			return nil, err
		}
	}

	if limit != nil {
		err := f.InitLimitOffset(limit)
		if err != nil {
			return nil, err
		}
	}

	return UpdateLoc(f, closetok)
}

func getFunctionCall(v interface{}) (*ast.FunctionCall, ast.Loc) {
	switch t := v.(type) {
	case *ast.FunctionCall:
		return t, ast.Loc{Start: t.StartLoc(), End: t.EndLoc()}
	case *ast.Wrapped:
		// Return the inner node but with the current location.
		e, _ := getFunctionCall(t.Value)
		return e, t.Loc
	}

	panic(fmt.Errorf("%w: could not get FunctionCall from %v",
		zerrors.ErrMalformedParser, reflect.TypeOf(v)))
}

func NewStructConstructorWithKeyword(kw Attrib) (Attrib, error) {
	if tok, ok := kw.(*token.Token); ok {
		t, err := ast.NewStructConstructorWithKeyword(nil)
		if err != nil {
			return nil, err
		}

		return UpdateLoc(t, tok)
	}

	return ast.NewStructConstructorWithKeyword(kw)
}

func NewStructConstructorWithParens(
	open, field1, field2 Attrib) (Attrib, error) {
	t, err := ast.NewStructConstructorWithParens(List(field1, field2))
	if err != nil {
		return nil, err
	}

	return UpdateLoc(t, open)
}

func NewCaseValueExpression(caseKw, value, when, then Attrib) (Attrib, error) {
	e, err := ast.NewCaseValueExpression(List(value, when, then))
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, caseKw)
}

func NewCaseNoValueExpression(caseKw, when, then Attrib) (Attrib, error) {
	e, err := ast.NewCaseNoValueExpression(List(when, then))
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, caseKw)
}

func CaseExpressionWithElse(expr, elseValue, end Attrib) (Attrib, error) {
	e, err := WithExtraChild(expr, elseValue)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, end)
}

func getExpressionHandler(v interface{}) (ast.ExpressionHandler, ast.Loc) {
	switch t := v.(type) {
	case ast.ExpressionHandler:
		return t, ast.Loc{Start: t.StartLoc(), End: t.EndLoc()}
	case *ast.Wrapped:
		// Return the inner expression but with the current location.
		e, _ := getExpressionHandler(t.Value)
		return e, t.Loc
	}

	panic(fmt.Errorf("%w: could not get ExpressionHandler from %v",
		zerrors.ErrMalformedParser, reflect.TypeOf(v)))
}
