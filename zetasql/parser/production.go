package parser

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/paulourio/bqfmt/zetasql/ast"
	zerrors "github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/literal"
	"github.com/paulourio/bqfmt/zetasql/token"
)

// OverrideLoc resets the location of a node and use the given tokens
// as reference instead.
func OverrideLoc(node Attrib, tokens ...Attrib) (ast.NodeHandler, error) {
	n := node.(ast.NodeHandler)
	first := true

	for i, t := range tokens {
		if t == nil {
			continue
		}

		var start, end int

		switch elem := t.(type) {
		case *token.Token:
			start = elem.Pos.Offset
			end = elem.Pos.Offset + len(elem.Lit)
		case ast.NodeHandler:
			start = elem.StartLoc()
			end = elem.EndLoc()
		case *ast.Wrapped:
			start = elem.Loc.Start
			end = elem.Loc.End
		case ast.Loc:
			start = elem.Start
			end = elem.End
		default:
			return nil, fmt.Errorf(
				"OverrideLoc: invalid argument %d of type %v",
				i+1, reflect.TypeOf(t))
		}

		if first {
			first = false

			n.SetStartLoc(start)
			n.SetEndLoc(end)
		} else {
			n.ExpandLoc(start, end)
		}
	}

	return n, nil
}

// UpdateLoc expands the localization of a node with a list of tokens
// or locations, from which a location range [min, max) is inferred.
func UpdateLoc(node Attrib, tokens ...Attrib) (ast.NodeHandler, error) {
	n := node.(ast.NodeHandler)

	for _, t := range tokens {
		if t == nil {
			continue
		}

		switch v := t.(type) {
		case *token.Token:
			n.ExpandLoc(v.Pos.Offset, v.Pos.Offset+len(v.Lit))
		case ast.Loc:
			n.ExpandLoc(v.Start, v.End)
		case ast.NodeHandler:
			n.ExpandLoc(v.StartLoc(), v.EndLoc())
		default:
			return nil, fmt.Errorf(
				"%w: cannot UpdateLoc with type %v",
				zerrors.ErrMalformedParser,
				reflect.TypeOf(t))
		}
	}

	return n, nil
}

// WithExtraChild adds a child node to the node.
func WithExtraChild(a Attrib, c Attrib) (Attrib, error) {
	node, locLeft := getNodeHandler(a)
	child, locRight := getNodeHandler(c)

	node.AddChild(child)

	return WrapWithLoc(node, locLeft, locRight)
}

// WithExtraChildren adds a child node to the node.
func WithExtraChildren(a Attrib, children ...Attrib) (Attrib, error) {
	node := a.(ast.NodeHandler)

	for _, c := range children {
		if c != nil {
			n, loc := getNodeHandler(c)
			node.AddChild(n)
			node.ExpandLoc(loc.Start, loc.End)
		}
	}

	return node, nil
}

// WrapWithLoc returns a value wrapped with a location range [min, max)
// inferred from a given list of tokens.
func WrapWithLoc(elem Attrib, tokens ...Attrib) (Attrib, error) {
	start := maxInt
	end := 0

	if elem == nil {
		return elem, nil
	}

	for _, t := range tokens {
		if t == nil {
			continue
		}

		loc := mustGetLoc(t)

		if loc.Start == -1 || loc.End == -1 {
			// Skip invalid location.
			continue
		}

		if loc.Start < start {
			start = loc.Start
		}

		if loc.End > end {
			end = loc.End
		}
	}

	// Do not wrap when the [start, end) is invalid
	if start >= end {
		return elem, nil
	}

	// Check if we really need to wrap the elem.
	if node, ok := elem.(ast.NodeHandler); ok {
		if node.StartLoc() == start && node.EndLoc() == end {
			return node, nil
		}
	}

	return ast.WrapWithLoc(elem, start, end)
}

// List casts a list of parser attributes to a generic list.
func List(a ...Attrib) Attrib {
	r := make([]interface{}, 0, len(a))
	l := make([]Attrib, 0, len(a))

	for _, e := range a {
		node, loc := getNodeHandler(e)
		r = append(r, node)
		l = append(l, loc)
	}

	w, err := WrapWithLoc(r, l...)
	if err != nil {
		panic(fmt.Errorf("%w: failed to created %v wrapped with %v",
			zerrors.ErrMalformedParser, reflect.TypeOf(r), reflect.TypeOf(l)))
	}

	return w
}

func InitLiteral(lit ast.LeafHandler, t Attrib) (Attrib, error) {
	tok := t.(*token.Token)
	lit.SetImage(string(tok.Lit))
	lit.SetStartLoc(tok.Pos.Offset)
	lit.SetEndLoc(tok.Pos.Offset + len(tok.Lit))

	return lit, nil
}

func NewStringLiteral(a Attrib) (Attrib, error) {
	lit, err := ast.NewStringLiteral()
	if err != nil {
		return nil, err
	}

	tok := a.(*token.Token)

	value, err := literal.ParseString(string(tok.Lit))
	if err != nil {
		var litErr *literal.UnescapeError

		// If the string contains invalid escape sequences, we need to
		// adjust the string offset error to raise with offset computed
		// for the global input.
		if errors.As(err, &litErr) {
			return nil, &literal.UnescapeError{
				Msg:    litErr.Msg,
				Offset: tok.Offset + litErr.Offset,
			}
		}

		return nil, err
	}

	lit.StringValue = value

	return InitLiteral(lit, a)
}

func NewBytesLiteral(a Attrib) (Attrib, error) {
	lit, err := ast.NewBytesLiteral()
	if err != nil {
		return nil, err
	}

	tok := a.(*token.Token)

	value, err := literal.ParseBytes(string(tok.Lit))
	if err != nil {
		var unescapeErr *literal.UnescapeError

		if errors.As(err, &unescapeErr) {
			return nil, &literal.UnescapeError{
				Msg:    unescapeErr.Msg,
				Offset: tok.Offset + unescapeErr.Offset,
			}
		}

		return nil, err
	}

	lit.BytesValue = value

	return InitLiteral(lit, a)
}

func NewDashedIdentifier(lhs Attrib, rhs Attrib) (*ast.Identifier, error) {
	a := lhs.(*token.Token)
	b := rhs.(*token.Token)
	actual := b.Offset - a.Offset + 1
	expected := len(a.Lit) + 1

	if actual != expected {
		return nil, zerrors.ErrInvalidDashedIdentifier
	}

	if a.Lit[0] == '`' || b.Lit[0] == '`' {
		return nil, zerrors.ErrInvalidDashedIdentifier
	}

	values := a.Lit
	values = append(values, '-')
	values = append(values, b.Lit...)

	return ast.NewIdentifier(string(values))
}

func NewIdentifier(a Attrib, allowReservedKw bool) (Attrib, error) {
	tok := a.(*token.Token)
	raw := string(tok.Lit)

	id, err := ParseIdentifier(raw, allowReservedKw)
	if err != nil {
		return nil, err
	}

	wrapped, err := WrapWithLoc(id, tok)
	if err != nil {
		return nil, err
	}

	return ast.NewIdentifier(wrapped)
}

func NewIntLiteral(a Attrib) (Attrib, error) {
	lit, err := ast.NewIntLiteral()
	if err != nil {
		return nil, err
	}

	return InitLiteral(lit, a)
}

func NewFloatLiteral(a Attrib) (Attrib, error) {
	lit, err := ast.NewFloatLiteral()
	if err != nil {
		return nil, err
	}

	return InitLiteral(lit, a)
}

func NewBooleanLiteral(a Attrib) (Attrib, error) {
	lit, err := ast.NewBooleanLiteral()
	if err != nil {
		return nil, err
	}

	return InitLiteral(lit, a)
}

// ExpandPathExpressionOrNewDotIdentifier tries to build path
// expressions as long as identifiers are added. As soon as a dotted
// path contains anything else, we use generalized DotIdentifier.
func ExpandPathExpressionOrNewDotIdentifier(
	expr, dot, name Attrib,
) (Attrib, error) {
	e, _ := getExpressionHandler(expr)
	if e.Kind() == ast.PathExpressionKind && !e.IsParenthesized() {
		return WithExtraChild(e, name)
	} else {
		var (
			d   ast.NodeHandler
			err error
		)

		d, err = ast.NewDotIdentifier(expr, name)
		if err != nil {
			return nil, err
		}

		d, err = OverrideLoc(d, dot, name)
		if err != nil {
			return nil, err
		}

		return WrapWithLoc(d, expr, name)
	}
}

func NewLikeBinaryExpression(inOp, inLHS, inRHS Attrib) (Attrib, error) {
	lhs, loc := getExpressionHandler(inLHS)
	op := inOp.(*ast.Wrapped)
	isNot := op.Value.(ast.NotKeyword) == ast.NotKeywordPresent

	if !lhs.IsAllowedInComparison() {
		return nil, zerrors.NewSyntaxError(
			op.Loc,
			"expression to left of LIKE must be parenthesized")
	}

	bin, err := ast.NewBinaryExpression(ast.BinaryLike, lhs, inRHS, isNot)
	if err != nil {
		return nil, err
	}

	bin.ExpandLoc(loc.StartLoc(), loc.EndLoc())

	return bin, nil
}

func NewInExpression(inOp, inLHS, inRHS Attrib) (Attrib, error) {
	lhs, locStart := getExpressionHandler(inLHS)
	rhs, locEnd := unwrap(inRHS)
	op := inOp.(*ast.Wrapped)
	isNot := op.Value.(ast.NotKeyword) == ast.NotKeywordPresent

	if !lhs.IsAllowedInComparison() {
		return nil, zerrors.NewSyntaxError(
			op.Loc,
			"expression to left of IN must be parenthesized")
	}

	var (
		expr *ast.InExpression
		err  error
	)

	switch r := rhs.(type) {
	case *ast.InList:
		expr, err = ast.NewInExpression(lhs, r, nil, nil, isNot)
	case *ast.Query:
		expr, err = ast.NewInExpression(lhs, nil, r, nil, isNot)
	case *ast.UnnestExpression:
		expr, err = ast.NewInExpression(lhs, nil, nil, r, isNot)
	default:
		return NewInternalError(
			inRHS,
			fmt.Sprintf("unexpected type %v", reflect.TypeOf(rhs)))
	}

	if err != nil {
		return nil, err
	}

	return WrapWithLoc(expr, locStart, locEnd)
}

func NewInList(open, a, b Attrib) (Attrib, error) {
	node, err := ast.NewInList(List(a, b))
	if err != nil {
		return nil, err
	}

	return WrapWithLoc(node, open)
}

func NewIsBinaryExpression(inOp, inLHS, inRHS Attrib) (Attrib, error) {
	lhs, loc := getExpressionHandler(inLHS)
	op := inOp.(*ast.Wrapped)
	isNot := op.Value.(ast.NotKeyword) == ast.NotKeywordPresent

	if !lhs.IsAllowedInComparison() {
		return nil, zerrors.NewSyntaxError(
			op.Loc,
			"expression to left of IS must be parenthesized")
	}

	bin, err := ast.NewBinaryExpression(ast.BinaryIs, lhs, inRHS, isNot)
	if err != nil {
		return nil, err
	}

	bin.ExpandLoc(loc.StartLoc(), loc.EndLoc())

	return bin, nil
}

func NewTablePathExpression(
	base, pivotUnpivotAlias, offset, sample Attrib,
) (*ast.TablePathExpression, error) {
	var (
		path   ast.NodeHandler
		unnest ast.NodeHandler
	)

	p := pivotUnpivotAlias.(*pivotOrUnpivotAndAlias)

	switch v := base.(type) {
	case *ast.PathExpression:
		path = v
	case *ast.UnnestExpression:
		unnest = v
	default:
		return nil, fmt.Errorf(
			"%w: NewTablePathExperession path/unnest with type %#v",
			zerrors.ErrMalformedParser,
			reflect.TypeOf(base))
	}

	if offset != nil {
		if p.PivotClause != nil {
			return nil, zerrors.NewSyntaxError(
				mustGetLoc(offset),
				"PIVOT and WITH OFFSET cannot be combined")
		}

		if p.UnpivotClause != nil {
			return nil, zerrors.NewSyntaxError(
				mustGetLoc(offset),
				"UNPIVOT and WITH OFFSET cannot be combined")
		}
	}

	return ast.NewTablePathExpression(
		path, unnest, p.Alias, offset, p.PivotClause, p.UnpivotClause, sample)
}

func IsUnparenthesizedNotExpression(a Attrib) (r bool) {
	if n, ok := a.(*ast.UnaryExpression); ok {
		r = !n.IsParenthesized() && n.Op != ast.UnaryNot
	}

	return
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

func getNodeHandler(v interface{}) (ast.NodeHandler, ast.Loc) {
	switch t := v.(type) {
	case ast.NodeHandler:
		return t, ast.Loc{Start: t.StartLoc(), End: t.EndLoc()}
	case *ast.Wrapped:
		// Return the inner node but with the current location.
		e, _ := getNodeHandler(t.Value)
		return e, t.Loc
	}

	panic(fmt.Errorf("%w: could not get NodeHandler from %v",
		zerrors.ErrMalformedParser, reflect.TypeOf(v)))
}

// Try get a loc from attribute, but panic if fail.
func mustGetLoc(a Attrib) ast.Loc {
	switch t := a.(type) {
	case ast.Loc:
		return t
	case ast.NodeHandler:
		return t.GetLoc()
	case *ast.Wrapped:
		return t.Loc
	case *token.Token:
		start := t.Offset
		end := t.Offset + len(t.Lit)

		return ast.Loc{Start: start, End: end}
	}

	panic(fmt.Errorf("%w: could not extract Loc from %v",
		zerrors.ErrMalformedParser, reflect.TypeOf(a)))
}

func unwrap(elem Attrib) (Attrib, ast.Loc) {
	if wrap, ok := elem.(*ast.Wrapped); ok {
		return wrap.Value, wrap.Loc
	}

	return elem, ast.Loc{Start: -1, End: -1}
}

var maxInt = int(^uint(0) >> 1) // largest int
