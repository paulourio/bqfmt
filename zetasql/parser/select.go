package parser

import (
	"strings"

	"github.com/paulourio/bqfmt/zetasql/ast"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func NewSelect(
	selectTok, allOrDistinct, selectas, selectlist, fromclause, whereclause,
	groupby, having, qualify, windowclause Attrib,
) (Attrib, error) {
	distinct := allOrDistinct == DistinctKeyword

	s, err := ast.NewSelect(distinct, selectas, selectlist, fromclause,
		whereclause, groupby, having, qualify, windowclause)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(s, selectTok)
}

func NewSelectAs(as, what Attrib) (Attrib, error) {
	var (
		selectAs *ast.SelectAs
		err      error
	)

	if tok, ok := what.(*token.Token); ok {
		if tok.Type == token.TokMap.Type("struct") {
			selectAs, err = ast.NewSelectAs(nil, ast.AsStruct)
		}
	}

	if expr, ok := what.(*ast.PathExpression); ok {
		if len(expr.Names) == 1 {
			value := strings.ToUpper(expr.Names[0].IDString)
			if value == "VALUE" {
				selectAs, err = ast.NewSelectAs(nil, ast.AsValue)
			}
		}
	}

	if selectAs == nil {
		selectAs, err = ast.NewSelectAs(what, ast.AsTypeName)
	}

	if err != nil {
		return nil, err
	}

	return UpdateLoc(selectAs, as)
}

func NewSelectStar(star Attrib) (Attrib, error) {
	s, err := ast.NewStar()
	if err != nil {
		return nil, err
	}

	starTok := star.(*token.Token)
	s.SetImage(string(starTok.Lit))

	c, err := UpdateLoc(s, star)
	if err != nil {
		return nil, err
	}

	return ast.NewSelectColumn(c, nil)
}

func NewSelectStarWithModifiers(star, modifiers Attrib) (Attrib, error) {
	s, err := ast.NewStarWithModifiers(modifiers)
	if err != nil {
		return nil, err
	}

	c, err := UpdateLoc(s, star)
	if err != nil {
		return nil, err
	}

	return ast.NewSelectColumn(c, nil)
}

func NewStarExceptList(except, identifier Attrib) (Attrib, error) {
	e, err := ast.NewStarExceptList(identifier)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, except)
}

func NewStarModifiers(except, replace, item Attrib) (Attrib, error) {
	e, err := ast.NewStarModifiers(except, item)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, replace)
}

func NewSelectDotStar(expr, dot Attrib) (Attrib, error) {
	ds, err := ast.NewDotStar(expr)
	if err != nil {
		return nil, err
	}

	d, err := OverrideLoc(ds, dot)
	if err != nil {
		return nil, err
	}

	s, err := ast.NewSelectColumn(d, nil)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(s, expr, dot)
}

func NewSelectDotStarWithModifiers(
	expr, dot, modifiers Attrib) (Attrib, error) {
	e, err := ast.NewDotStarWithModifiers(expr, modifiers)
	if err != nil {
		return nil, err
	}

	return ast.NewSelectColumn(e, nil)
}

func NewSelectExprAsAlias(expr, as, alias Attrib) (Attrib, error) {
	a, err := UpdateLoc(alias, as)
	if err != nil {
		return nil, err
	}

	return ast.NewSelectColumn(expr, a)
}

func NewFromClause(from, contents Attrib) (Attrib, error) {
	node, err := TransformJoinExpression(contents)
	if err != nil {
		return nil, err
	}

	f, err := ast.NewFromClause(node)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(f, from)
}

func NewGroupBy(group, item Attrib) (Attrib, error) {
	g, err := ast.NewGroupBy(item)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(g, group)
}

func NewRollup(rollup, expr Attrib) (Attrib, error) {
	e, err := ast.NewRollup(expr)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, rollup)
}

func NewHavingClause(having, expr Attrib) (Attrib, error) {
	e, err := ast.NewHaving(expr)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, having)
}

func NewQualify(qualify, expr Attrib) (Attrib, error) {
	q, err := ast.NewQualify(expr)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(q, qualify)
}

func NewWindowClause(window, definition Attrib) (Attrib, error) {
	w, err := ast.NewWindowClause(definition)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(w, window)
}

func NewWithClause(with, entry Attrib) (Attrib, error) {
	w, err := ast.NewWithClause(entry)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(w, with)
}

func NewWithClauseEntry(alias, query, parens Attrib) (Attrib, error) {
	e, err := ast.NewWithClauseEntry(alias, query)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(e, parens)
}

type allOrDistinctKeyword int

const (
	NoAllOrDistinctKeyword allOrDistinctKeyword = iota
	AllKeyword
	DistinctKeyword
)
