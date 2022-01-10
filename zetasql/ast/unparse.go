package ast

import (
	"strings"
)

type unparser struct {
	out   strings.Builder
	depth int

	Operation
}

func Unparse(n NodeHandler) string {
	u := &unparser{strings.Builder{}, 0, Operation{}}
	u.Operation.visitor = u
	n.Accept(u, nil)
	u.out.WriteString("\n")

	return u.out.String()
}

func (u *unparser) VisitSelect(n *Select, d interface{}) {
	u.out.WriteString("SELECT")
	u.depth++
	n.SelectList.Accept(u, d)
	u.depth--

	if n.FromClause != nil {
		u.indentedNewline()
		u.out.WriteString("FROM")
		u.depth++
		u.indentedNewline()
		n.FromClause.Accept(u, d)
		u.depth--
	}

	if n.WhereClause != nil {
		u.indentedNewline()
		u.out.WriteString("WHERE")
		u.depth++
		u.indentedNewline()
		n.WhereClause.Accept(u, d)
		u.depth--
	}
}

func (u *unparser) VisitSelectList(n *SelectList, d interface{}) {
	for i, col := range n.Columns {
		if i > 0 {
			u.out.WriteString(",")
		}

		u.indentedNewline()
		col.Accept(u, d)
	}
}

func (u *unparser) VisitSelectColumn(n *SelectColumn, d interface{}) {
	n.Expression.Accept(u, d)

	if n.Alias != nil {
		u.out.WriteString(" AS ")
		u.out.WriteString(ToIdentifierLiteral(n.Alias.Identifier.IDString))
	}
}

func (u *unparser) VisitIntLiteral(n *IntLiteral, d interface{}) {
	u.out.WriteString(n.image)
}

func (u *unparser) VisitStringLiteral(n *StringLiteral, d interface{}) {
	u.out.WriteString(n.image)
}

func (u *unparser) VisitBinaryExpression(n *BinaryExpression, d interface{}) {
	if n.IsParenthesized() {
		u.out.WriteString("(")
	}

	n.LHS.Accept(u, d)
	u.out.WriteRune(' ')

	if n.IsNot {
		switch n.Op {
		case BinaryIs:
			u.out.WriteString("IS NOT ")
		case BinaryLike:
			u.out.WriteString("NOT LIKE ")
		}
	} else {
		u.out.WriteString(n.Op.String())
		u.out.WriteRune(' ')
	}

	n.RHS.Accept(u, d)

	if n.IsParenthesized() {
		u.out.WriteString(")")
	}
}

func (u *unparser) VisitUnaryExpression(n *UnaryExpression, d interface{}) {
	if n.IsParenthesized() {
		u.out.WriteString("(")
	}

	u.out.WriteString(n.Op.String())

	if n.Op == UnaryNot {
		u.out.WriteRune(' ')
	}

	switch n.Operand.(type) {
	case *UnaryExpression:
		// Add space to avoid generating something like "--1"
		u.out.WriteString(" ")
		n.Operand.Accept(u, d)
	default:
		n.Operand.Accept(u, d)
	}

	if n.IsParenthesized() {
		u.out.WriteString(")")
	}
}

func (u *unparser) VisitBetweenExpression(
	n *BetweenExpression, d interface{}) {
	n.LHS.Accept(u, d)
	u.out.WriteString(" BETWEEN ")
	n.Low.Accept(u, d)
	u.out.WriteString(" AND ")
	n.High.Accept(u, d)
}

func (u *unparser) VisitJoin(n *Join, d interface{}) {
	n.LHS.Accept(u, d)

	switch n.JoinType {
	case DefaultJoin:
		u.indentedNewline()
		u.out.WriteString("JOIN")
	case CommaJoin:
		u.out.WriteString(",")
	case CrossJoin:
		u.indentedNewline()
		u.out.WriteString("CROSS JOIN")
	case FullJoin:
		u.indentedNewline()
		u.out.WriteString("FULL OUTER JOIN")
	case InnerJoin:
		u.indentedNewline()
		u.out.WriteString("INNER JOIN")
	case LeftJoin:
		u.indentedNewline()
		u.out.WriteString("LEFT JOIN")
	case RightJoin:
		u.indentedNewline()
		u.out.WriteString("RIGHT JOIN")
	}

	u.indentedNewline()
	n.RHS.Accept(u, d)
}

func (u *unparser) VisitTablePathExpression(
	n *TablePathExpression, d interface{},
) {
	if n.PathExpr != nil {
		n.PathExpr.Accept(u, d)
	} else {
		n.UnnestExpr.Accept(u, d)
	}

	if n.Alias != nil {
		u.out.WriteString(" AS " + n.Alias.Identifier.IDString)
	}
}

func (u *unparser) VisitTableSubquery(n *TableSubquery, d interface{}) {
	u.out.WriteString("(")
	u.depth++
	u.indentedNewline()
	n.Subquery.Accept(u, d)
	u.depth--
	u.indentedNewline()
	u.out.WriteString(")")

	if n.Alias != nil {
		u.out.WriteString(" AS " + n.Alias.Identifier.IDString)
	}
}

func (u *unparser) VisitIdentifier(n *Identifier, d interface{}) {
	u.out.WriteString(ToIdentifierLiteral(n.IDString))
}

func (u *unparser) VisitDotIdentifier(n *DotIdentifier, d interface{}) {
	n.Expr.Accept(u, d)
	u.out.WriteString(".")
	n.Name.Accept(u, d)
}

func (u *unparser) VisitPathExpression(n *PathExpression, d interface{}) {
	if n.IsParenthesized() {
		u.out.WriteRune('(')
	}

	for i, p := range n.Names {
		if i > 0 {
			u.out.WriteRune('.')
		}

		p.Accept(u, d)
	}

	if n.IsParenthesized() {
		u.out.WriteRune(')')
	}
}

func (u *unparser) VisitCastExpression(n *CastExpression, d interface{}) {
	if n.IsSafeCast {
		u.out.WriteString("SAFE_CAST(")
	} else {
		u.out.WriteString("CAST(")
	}

	n.Expr.Accept(u, d)
	u.out.WriteString(" AS ")
	n.Type.Accept(u, d)

	if n.Format != nil {
		if n.Format.Format != nil {
			u.out.WriteString(" FORMAT ")
			n.Format.Format.Accept(u, d)
		}

		if n.Format.TimeZoneExpr != nil {
			u.out.WriteString(" AT TIME ZONE ")
			n.Format.TimeZoneExpr.Accept(u, d)
		}
	}

	u.out.WriteRune(')')
}

func (u *unparser) VisitFunctionCall(n *FunctionCall, d interface{}) {
	n.Function.Accept(u, n)
	u.out.WriteRune('(')

	if n.Distinct {
		u.out.WriteString("DISTINCT")

		if len(n.Arguments) > 0 {
			u.out.WriteRune(' ')
		}
	}

	for i, arg := range n.Arguments {
		if i > 0 {
			u.out.WriteString(", ")
		}

		arg.Accept(u, d)
	}

	u.out.WriteRune(')')
}

func (u *unparser) VisitAndExpr(n *AndExpr, d interface{}) {
	if n.IsParenthesized() {
		u.out.WriteRune('(')
	}

	for i, c := range n.Conjuncts {
		if i > 0 {
			u.out.WriteString(" AND ")
		}

		c.Accept(u, d)
	}

	if n.IsParenthesized() {
		u.out.WriteRune(')')
	}
}

func (u *unparser) VisitOrExpr(n *OrExpr, d interface{}) {
	if n.IsParenthesized() {
		u.out.WriteRune('(')
	}

	for i, c := range n.Disjuncts {
		if i > 0 {
			u.out.WriteString(" OR ")
		}

		c.Accept(u, d)
	}

	if n.IsParenthesized() {
		u.out.WriteRune(')')
	}
}

func (u *unparser) VisitStar(n *Star, d interface{}) {
	u.out.WriteString(n.Image())
}

func (u *unparser) VisitBooleanLiteral(n *BooleanLiteral, d interface{}) {
	u.out.WriteString(n.Image())
}

func (u *unparser) VisitNamedType(n *NamedType, d interface{}) {
	n.Name.Accept(u, d)

	if t := n.TypeParameters(); t != nil {
		u.out.WriteRune('(')

		for i, p := range t.Parameters {
			if i > 0 {
				u.out.WriteString(", ")
			}

			p.Accept(u, d)
		}

		u.out.WriteRune(')')
	}
}

func (u *unparser) VisitArrayType(n *ArrayType, d interface{}) {
	u.out.WriteString("ARRAY< ")
	n.ElementType.Accept(u, d)
	u.out.WriteString(" >")
}

func (u *unparser) VisitStructType(n *StructType, d interface{}) {
	u.out.WriteString("STRUCT< ")

	for i, f := range n.StructFields {
		if i > 0 {
			u.out.WriteString(", ")
		}

		f.Accept(u, d)
	}

	u.out.WriteString(" >")
}

func (u *unparser) VisitStructField(n *StructField, d interface{}) {
	if n.Name != nil {
		n.Name.Accept(u, d)
		u.out.WriteRune(' ')
	}

	n.Type.Accept(u, d)
}

func (u *unparser) VisitArrayConstructor(n *ArrayConstructor, d interface{}) {
	u.out.WriteString("ARRAY[")

	for i, e := range n.Elements {
		if i > 0 {
			u.out.WriteString(", ")
		}

		e.Accept(u, d)
	}

	u.out.WriteRune(']')
}

func (u *unparser) VisitArrayElement(n *ArrayElement, d interface{}) {
	n.Array.Accept(u, d)
	u.out.WriteRune('[')
	n.Position.Accept(u, d)
	u.out.WriteRune(']')
}

func (u *unparser) VisitNullLiteral(n *NullLiteral, d interface{}) {
	u.out.WriteString("NULL")
}

func (u *unparser) indent() {
	u.out.WriteString(strings.Repeat(" ", u.depth*2))
}

func (u *unparser) indentedNewline() {
	u.out.WriteString("\n")
	u.indent()
}
