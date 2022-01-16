package parser

import (
	"fmt"
	"reflect"

	"github.com/paulourio/bqfmt/zetasql/ast"
	"github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func NewJoin(lhs, rhs, clauses, typ, joinKw Attrib) (Attrib, error) {
	var (
		clauseList *ast.OnOrUsingClauseList
		clauseLoc  ast.Loc
	)

	lhsNode, _ := getNodeHandler(lhs)
	rhsNode, rhsLoc := getNodeHandler(rhs)
	joinType, joinLoc, joinHasLoc := getJoinType(typ)

	// Raise an error if we have a RIGHT or FULL JOIN following a comma
	// join since our left-to-right binding would violate the standard.
	if (joinType == ast.FullJoin || joinType == ast.RightJoin) &&
		lhsNode.Kind() == ast.JoinKind {
		joinInput := lhsNode.(*ast.Join)

		for {
			if joinInput.JoinType == ast.CommaJoin {
				return nil, errors.NewSyntaxError(
					joinLoc,
					fmt.Sprintf(
						"%v JOIN must be parenthesized when following a "+
							"comma join.  Also, if the preceding comma join "+
							"is a correlated CROSS JOIN that unnests an "+
							"array, then CROSS JOIN syntax must be used in "+
							"place of the comma join",
						joinType))
			}

			if joinInput.LHS.Kind() == ast.JoinKind {
				// Look deeper only if the left input is an unparenthesized
				// join.
				joinInput = joinInput.LHS.(*ast.Join)
			} else {
				break
			}
		}
	}

	join, err := ast.NewJoin(lhsNode, rhsNode, joinType)
	if err != nil {
		return nil, err
	}

	if clauses != nil {
		clauseList, clauseLoc = getOnOrUsingClauseList(clauses)

		n, err := UpdateLoc(join, clauseLoc)
		if err != nil {
			return nil, err
		}

		join = n.(*ast.Join)
	}

	if clauseList == nil { //nolint:gocritic
		// Nothing.
	} else if len(clauseList.List) == 1 {
		switch c := clauseList.List[0].(type) {
		case *ast.OnClause:
			err := join.InitOnClause(c)
			if err != nil {
				return nil, err
			}
		case *ast.UsingClause:
			err := join.InitUsingClause(c)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf(
				"%w: cannot create join with ClauseList of type %#v",
				errors.ErrMalformedParser,
				reflect.TypeOf(clauseList.List[0]))
		}
	} else {
		err := join.InitClauseList(clauseList)
		if err != nil {
			return nil, err
		}

		join.TransformationNeeded = true
	}

	node, err := OverrideLoc(join, joinKw, rhsLoc, clauses)
	if err != nil {
		return nil, err
	}

	if joinHasLoc {
		node, err = UpdateLoc(node, joinLoc)
		if err != nil {
			return nil, err
		}
	}

	return node.(*ast.Join), nil
}

func NewCommaJoin(lhs, rhs, comma Attrib) (Attrib, error) {
	if IsTransformationNeeded(lhs) {
		return nil, errors.NewSyntaxError(
			mustGetLoc(comma),
			"Comma join is not allowed after consecutives ON/USING clauses")
	}

	c, err := ast.NewJoin(lhs, rhs, ast.CommaJoin)
	if err != nil {
		return nil, err
	}

	c.ContainsCommaJoin = true

	// Not sure why, but reference implementation sets start
	// location at the comma token.
	c.SetStartLoc(comma.(*token.Token).Offset)

	return c, nil
}

func NewParenthesizedJoin(lhs, sample, open, close Attrib) (Attrib, error) {
	node, err := TransformJoinExpression(lhs)
	if err != nil {
		return nil, err
	}

	join, err := ast.NewParenthesizedJoin(node, sample)
	if err != nil {
		return nil, err
	}

	return UpdateLoc(join, open, close)
}

func IsTransformationNeeded(expr Attrib) bool {
	if join, ok := expr.(*ast.Join); ok {
		return join.TransformationNeeded
	}

	return false
}

func TransformJoinExpression(expr Attrib) (Attrib, error) {
	return expr, nil
}

func getOnOrUsingClauseList(
	v interface{},
) (*ast.OnOrUsingClauseList, ast.Loc) {
	n, loc := getNodeHandler(v)
	return n.(*ast.OnOrUsingClauseList), loc
}

// getJoinType returns the join type, location, and a boolean indicating
// the value has a location associated.
func getJoinType(v interface{}) (ast.JoinType, ast.Loc, bool) {
	switch t := v.(type) {
	case ast.JoinType:
		return t, ast.Loc{}, false
	case *ast.Wrapped:
		typ, _, _ := getJoinType(t.Value)
		return typ, t.Loc, true
	}

	panic(fmt.Errorf("%w: could not get JoinType from %v",
		errors.ErrMalformedParser, reflect.TypeOf(v)))
}
