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
		u.out.WriteString(n.Alias.Identifier.IDString)
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
	u.out.WriteString(" " + n.Op.String() + " ")
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
	u.out.WriteString(toIdentifierLiteral(n.IDString))
}

func (u *unparser) VisitDotIdentifier(n *DotIdentifier, d interface{}) {
	n.Expr.Accept(u, d)
	u.out.WriteString(".")
	n.Name.Accept(u, d)
}

func (u *unparser) VisitPathExpression(n *PathExpression, d interface{}) {
	for i, p := range n.Names {
		if i > 0 {
			u.out.WriteRune('.')
		}
		p.Accept(u, d)
	}
}

func (u *unparser) VisitFunctionCall(n *FunctionCall, d interface{}) {
	n.Function.Accept(u, n)
	u.out.WriteString("(")
	if n.Distinct {
		u.out.WriteString("DISTINCT ")
	}
	for i, arg := range n.Arguments {
		if i > 0 {
			u.out.WriteString(", ")
		}
		arg.Accept(u, d)
	}
	u.out.WriteString(")")
}

func (u *unparser) indent() {
	u.out.WriteString(strings.Repeat(" ", u.depth*2))
}

func (u *unparser) indentedNewline() {
	u.out.WriteString("\n")
	u.indent()
}
