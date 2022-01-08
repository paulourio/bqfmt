package ast

import (
	"fmt"
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
	fmt.Println("VisitSelect called!")
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

func (u *unparser) VisitIdentifier(n *Identifier, d interface{}) {
	u.out.WriteString(n.IDString)
}

func (u *unparser) indent() {
	u.out.WriteString(strings.Repeat(" ", u.depth*2))
}

func (u *unparser) indentedNewline() {
	u.out.WriteString("\n")
	u.indent()
}
