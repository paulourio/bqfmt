package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/paulourio/bqfmt/zetasql/literal"
)

// Sprint returns the formatted SQL of the AST tree.  The printing may
// may be configured with the optional opts param.
func Sprint(node NodeHandler, opts *PrintOptions) string {
	if node == nil {
		return ""
	}

	if opts == nil {
		opts = &PrintOptions{
			SoftMaxColumns:          70,
			NewlineBeforeClause:     true,
			AlignLogicalWithClauses: true,
			Indentation:             2,
			FunctionNameStyle:       UpperCase,
			IdentifierStyle:         LowerCase,
			KeywordStyle:            UpperCase,
			TypeStyle:               UpperCase,
			StringStyle:             PreferSingleQuote,
			MultilineStringStyle:    AlwaysSingleQuote,
		}
	}

	p := &printer{
		fmt: &formatter{
			opts: opts,
		},
	}

	p.Operation = Operation{visitor: p}

	node.Accept(p, nil)
	p.fmt.FlushLine()

	return strings.ReplaceAll(p.unnest(), "\v", "")
}

type PrintOptions struct {
	// MaxCol is a soft limit of the maximum number of characters to be
	// formatted into a single line. This limit may ignore whitespaces
	// at the beginning of each line.
	SoftMaxColumns int
	// NewlineBeforeClause sets whether new lines should be generated
	// before clauses at the same level, like FROM, WHERE, etc.
	NewlineBeforeClause bool
	// AlignLogicalWithClauses specifies if AND and OR expressions
	// should be aligned with clauses like WHERE.
	// When false you get
	//  WHERE     a = 1
	//        AND b = 2
	//         OR c = 3
	// and when true you get
	//  WHERE a = 1
	//    AND b = 3
	//     OR c = 3
	AlignLogicalWithClauses bool
	// Indentation sets the minimum amount of indentation when certain
	// expressions need to be split across lines.
	Indentation int
	// FunctionName sets how to style the name of function calls with
	// unquoted names.
	FunctionNameStyle PrintCase
	// IdentifierStyle sets how identifiers, such as column names and
	// aliases should be printed.
	IdentifierStyle PrintCase
	// KeywordStyle sets how keywords should be printed.
	KeywordStyle PrintCase
	// TypeStyle sets how type names should be printed.
	TypeStyle PrintCase
	// StringStyle sets how single-line strings should be printed.
	StringStyle StringStyle
	// MultilineStringStyle sets how multi-line strings should be
	// printed. In some cases the string will always be considered
	// multi-line, such as in definition of functions with body in
	// another language.
	MultilineStringStyle StringStyle
}

type PrintCase int

const (
	AsIsCase PrintCase = iota
	LowerCase
	UpperCase
)

type StringStyle int

const (
	// AsIsStringStyle prints the string as is from the input source.
	AsIsStringStyle StringStyle = iota
	// PreferSingleQuote prints prefers single quotes but allows double
	// quotes when the string contains as single quote.
	// It prefers 'ab"bc' but allows "ab'cd"
	PreferSingleQuote
	// PreferDoubleQuote prints prefers double quotes but allows double
	// quotes when the string contains as double quote.
	// It prefers "ab'bc" but allows 'ab"cd'
	PreferDoubleQuote
	// AlwaysSingleQuote forces all strings to use single quotes.
	AlwaysSingleQuote
	// AlwaysDoubleQuote forces all strings to use double quotes.
	AlwaysDoubleQuote
)

// onelinePrintConfig is used to estimate the length of the SQL
// representation of some node.
var onelinePrintConfig = &PrintOptions{
	SoftMaxColumns:      10000,
	NewlineBeforeClause: false,
}

type printer struct {
	fmt *formatter

	Operation
}

type formatter struct {
	opts                   *PrintOptions
	buffer                 strings.Builder
	formatted              strings.Builder
	maxLength              int
	depth                  int
	last                   rune
	lastWasSingleCharUnary bool
	// noFlushInNextFormat disables line flushing in the next call to Format()
	noFlushInNextFormat bool
}

func (p *printer) String() string {
	p.fmt.FlushLine()
	return strings.Trim(p.fmt.formatted.String(), "\n")
}

func (p *printer) print(s string) {
	p.fmt.Format(s)
}

func (p *printer) println(s string) {
	p.fmt.FormatLine(s)
}

func (p *printer) incDepth() {
	p.fmt.depth++
}

func (p *printer) decDepth() {
	p.fmt.depth--
}

// nest returns a new printer with the same options to perform printing
// on a nested section of the tree.
func (p *printer) nest() *printer {
	currSize := len(p.fmt.buffer.String())
	capacity := p.fmt.opts.SoftMaxColumns - currSize

	if capacity < 40 {
		capacity = 40
	}

	n := &printer{
		fmt: &formatter{
			opts:      p.fmt.opts,
			maxLength: capacity,
		},
	}

	n.Operation = Operation{visitor: n}

	return n
}

// unnest flushes the buffer and returns the strings with alignment
// symbols at the beginning of each line.
func (p *printer) unnest() string {
	trimmed := p.String()
	aligned := alignNested(trimmed)
	aligned = "\v" + aligned
	aligned = strings.ReplaceAll(aligned, "\n", "\n\v")

	return aligned
}

// unnest flushes the buffer and returns the strings with alignment
// symbols at the beginning of each line.
func (p *printer) unnestWithDepth(d int) string {
	trimmed := p.String()
	aligned := alignNested(trimmed)
	aligned = "\v" + aligned
	alignment := strings.Repeat("\v", d)
	aligned = strings.ReplaceAll(aligned, "\n", "\n"+alignment)

	return aligned
}

// printNestedWithSep receives a slice of NodeHandler items and print each
// in a nested printer.
func (p *printer) printNestedWithSep(
	items interface{}, d interface{}, sep string) {
	pp := p.nest()

	list := reflect.ValueOf(items)

	for i := 0; i < list.Len(); i++ {
		if i > 0 {
			pp.print(sep)
		}

		pp.visitNested(list.Index(i).Interface().(NodeHandler), d)
	}

	p.print(pp.unnest())
}

// visitNested visits a node with a nested printer.
func (p *printer) visitNested(n NodeHandler, d interface{}) {
	pp := p.nest()
	n.Accept(pp, d)
	p.print(pp.unnest())
}

func debugContent(s string) string {
	d := strings.ReplaceAll(s, "\v", "|")
	d = strings.ReplaceAll(d, "\b", "%")

	return d
}

// unnest flushes the buffer and returns the strings with alignment
// symbols at the beginning of each line.
func (p *printer) unnestLeft() string {
	aligned := leftAlignNested(p.String())
	return "\v" + strings.ReplaceAll(aligned, "\n", "\n\v")
}

func (p *printer) printOpenParenIfNeeded(n NodeHandler) {
	if !(n.IsExpression() || n.IsQueryExpression()) {
		panic("parenthesization is not allowed for " + n.Kind().String())
	}

	if expr, ok := n.(ExpressionHandler); ok && expr.IsParenthesized() {
		p.print("(")

		return
	}

	if expr, ok := n.(QueryExpressionHandler); ok && expr.IsParenthesized() {
		p.print("(")

		return
	}
}

func (p *printer) printCloseParenIfNeeded(n NodeHandler) {
	if !(n.IsExpression() || n.IsQueryExpression()) {
		panic("parenthesization is not allowed for " + n.Kind().String())
	}

	if expr, ok := n.(ExpressionHandler); ok && expr.IsParenthesized() {
		p.print(")")

		return
	}

	if expr, ok := n.(QueryExpressionHandler); ok && expr.IsParenthesized() {
		p.print(")")

		return
	}
}

func (p *printer) printClause(s string) {
	if p.fmt.opts.NewlineBeforeClause {
		p.println("")
	}

	p.print(s)
}

func (p *printer) identifier(s string) string {
	switch p.fmt.opts.IdentifierStyle {
	case AsIsCase:
		return s
	case UpperCase:
		return strings.ToUpper(s)
	case LowerCase:
		return strings.ToLower(s)
	}

	return ""
}

func (p *printer) function(s string) string {
	switch p.fmt.opts.FunctionNameStyle {
	case AsIsCase:
		return s
	case UpperCase:
		return strings.ToUpper(s)
	case LowerCase:
		return strings.ToLower(s)
	}

	return ""
}

func (p *printer) keyword(s string) string {
	switch p.fmt.opts.KeywordStyle {
	case UpperCase:
		return strings.ToUpper(s)
	case LowerCase:
		return strings.ToLower(s)
	case AsIsCase:
		panic("KeywordStyle must be either UpperCase or LowerCase")
	}

	return ""
}

func (p *printer) typename(s string) string {
	switch p.fmt.opts.TypeStyle {
	case UpperCase:
		return strings.ToUpper(s)
	case LowerCase:
		return strings.ToLower(s)
	case AsIsCase:
		panic("TypeStyle must be either UpperCase or LowerCase")
	}

	return ""
}

func (p *printer) VisitLeafHandler(n LeafHandler, d interface{}) {
	p.print(n.Image())
}

func (p *printer) VisitQuery(n *Query, d interface{}) {
	p.printOpenParenIfNeeded(n)
	p1 := p.nest()

	if n.WithClause != nil {
		n.WithClause.Accept(p1, d)
		// Force an empty line after the WithClause.
		p1.println("\n")
	}

	n.QueryExpr.Accept(p1, d)

	if n.OrderBy != nil {
		p1.printClause(p1.keyword("ORDER"))
		p2 := p1.nest()
		p2.print(p2.keyword("BY"))
		p2.printNestedWithSep(n.OrderBy.OrderingExpressions, d, ",")
		p1.print(p2.unnest())
	}

	if n.LimitOffset != nil {
		p1.printClause(p1.keyword("LIMIT"))

		p2 := p1.nest()

		n.LimitOffset.Limit.Accept(p2, d)

		if n.LimitOffset.Offset != nil {
			p2.print(p2.keyword("OFFSET"))
			n.LimitOffset.Offset.Accept(p2, d)
		}

		p1.print(p2.unnest())
	}

	p1.printCloseParenIfNeeded(n)
	p.print(p1.unnest())
}

func (p *printer) VisitExpressionSubquery(
	n *ExpressionSubquery, d interface{}) {
	pp := p.nest()

	pp.print(n.Modifier.ToSQL())
	pp.print("(")
	pp.visitNested(n.Query, d)
	pp.print(")")
	p.print(pp.unnest())
}

func (p *printer) VisitSetOperation(n *SetOperation, d interface{}) {
	p.printOpenParenIfNeeded(n)

	for i, query := range n.Inputs {
		if i > 0 {
			p.print(n.OpType.ToSQL())

			if n.Distinct {
				p.print(p.keyword("DISTINCT"))
			} else {
				p.print(p.keyword("ALL"))
			}
		}

		query.Accept(p, d)
	}

	p.printCloseParenIfNeeded(n)
}

func (p *printer) VisitPartitionBy(n *PartitionBy, d interface{}) {
	if n.Parent().Kind() == QueryKind {
		p.printClause(p.keyword("PARTITION"))
	} else {
		p.print(p.keyword("PARTITION"))
	}

	p1 := p.nest()
	p1.print(p1.keyword("BY"))
	p1.printNestedWithSep(n.PartitioningExpressions, d, ",")
	p.print(p1.unnest())
}

func (p *printer) VisitOrderBy(n *OrderBy, d interface{}) {
	if n.Parent().Kind() == QueryKind {
		p.printClause(p.keyword("ORDER"))
	} else {
		p.print(p.keyword("ORDER"))
	}

	p1 := p.nest()
	p1.print(p1.keyword("BY"))
	p1.printNestedWithSep(n.OrderingExpressions, d, ",")
	p.print(p1.unnest())
}

func (p *printer) VisitOrderingExpression(
	n *OrderingExpression, d interface{}) {
	n.Expression.Accept(p, d)

	switch n.OrderingSpec {
	case DescendingOrder:
		p.print(p.keyword("DESC"))
	case AscendingOrder:
		p.print(p.keyword("ASC"))
	case NoOrderingSpec:
		/* nothing */
		break
	}

	if n.NullOrder != nil {
		n.NullOrder.Accept(p, d)
	}
}

func (p *printer) VisitNullOrder(n *NullOrder, d interface{}) {
	if n.NullsFirst {
		p.print(p.keyword("NULLS FIRST"))
	} else {
		p.print(p.keyword("NULLS LAST"))
	}
}

func (p *printer) VisitLimitOffset(n *LimitOffset, d interface{}) {
	p.printClause(p.keyword("LIMIT") + "\v ")
	n.Limit.Accept(p, d)

	if n.Offset != nil {
		p.print(p.keyword("OFFSET"))
		n.Offset.Accept(p, d)
	}
}

func (p *printer) VisitSelect(n *Select, d interface{}) {
	pp := p.nest()

	pp.printOpenParenIfNeeded(n)

	pp2 := pp.nest()
	pp2.printClause(pp2.keyword("SELECT"))

	pp3 := pp2.nest()

	if n.Distinct {
		pp3.println(pp3.keyword("DISTINCT"))
	}

	n.SelectList.Accept(pp3, d)
	pp2.print(pp3.unnest())

	if n.FromClause != nil {
		pp2.printClause(pp2.keyword("FROM"))
		pp2.visitNested(n.FromClause, d)
	}

	if n.WhereClause != nil {
		pp2.printClause(pp2.keyword("WHERE"))

		pp3 := pp2.nest()

		// If the WHERE clause contains AND or OR, we will format them
		// as if they were clauses, right-aligned with the WHERE clause.
		if n.WhereClause.Expression.Kind() == AndExprKind {
			bin := n.WhereClause.Expression.(*AndExpr)

			for i, conjunct := range bin.Conjuncts {
				if i > 0 {
					if p.fmt.opts.AlignLogicalWithClauses {
						// Clear buffer and write AND as a clause.
						pp2.print(pp3.unnest())
						pp2.printClause(pp2.keyword("AND"))
						// Create new nested builder.
						pp3 = pp2.nest()
					} else {
						pp2.printClause(pp2.keyword("AND"))
					}
				}

				pp3.visitNested(conjunct, d)
			}
		} else if n.WhereClause.Expression.Kind() == OrExprKind {
			bin := n.WhereClause.Expression.(*OrExpr)
			for i, disjunct := range bin.Disjuncts {
				if i > 0 {
					if p.fmt.opts.AlignLogicalWithClauses {
						pp2.print(pp3.unnest())
						pp2.printClause(pp2.keyword("OR"))
						pp2 = pp2.nest()
					} else {
						pp2.printClause(pp2.keyword("OR"))
					}
				}

				pp3.visitNested(disjunct, d)
			}
		} else {
			n.WhereClause.Accept(pp3, d)
		}

		pp2.print(pp3.unnest())
	}

	if n.GroupBy != nil {
		pp2.printClause(pp2.keyword("GROUP") + " \v" + pp2.keyword("BY"))
		pp3 := pp2.nest()
		pp3.printNestedWithSep(n.GroupBy.GroupingItems, d, ",")
		pp2.print(pp3.unnest())
	}

	if n.Having != nil {
		pp2.printClause(pp2.keyword("HAVING"))
		pp2.visitNested(n.Having, d)
	}

	if n.Qualify != nil {
		pp2.printClause(pp2.keyword("QUALIFY"))
		pp2.visitNested(n.Qualify, d)
	}

	if n.WindowClause != nil {
		pp2.printClause(pp2.keyword("WINDOW"))
		pp2.visitNested(n.WindowClause, d)
	}

	// If this select is inside a Query node, we want to possibly align
	// SELECT, FROM, WHERE and other clauses with Query's ORDER BY and LIMIT.
	// Thus, we will not unnest is this case.
	if n.Parent() != nil && n.Parent().Kind() == QueryKind {
		pp.print(pp2.String())
		p.print(pp.String())
	} else {
		pp.print(pp2.unnest())
		pp.printCloseParenIfNeeded(n)
		pp.println("")
		p.print(pp.unnest())
	}
}

func (p *printer) VisitSelectList(n *SelectList, d interface{}) {
	pp := p.nest()

	for i, col := range n.Columns {
		if i > 0 {
			pp.println(",")
		}

		pp2 := pp.nest()
		col.Accept(pp2, d)
		pp.print(pp2.String())
	}

	p.print(pp.unnestLeft())
}

func (p *printer) VisitSelectColumn(n *SelectColumn, d interface{}) {
	pp := p.nest()
	n.Expression.Accept(pp, d)
	p.print(pp.unnest())

	if n.Alias != nil {
		pp = p.nest()
		pp.print(pp.keyword("AS"))
		pp.print(pp.identifier(n.Alias.Identifier.IDString))
		p.print(pp.unnest())
	}
}

func (p *printer) VisitGroupBy(n *GroupBy, d interface{}) {
	p.printClause(p.keyword("GROUP") + "\v " + p.keyword("BY"))

	for i, item := range n.GroupingItems {
		if i > 0 {
			p.print(",")
		}

		item.Expression.Accept(p, d)
	}
}

func (p *printer) VisitJoin(n *Join, d interface{}) {
	pp := p.nest()

	n.LHS.Accept(pp, d)

	switch n.JoinType {
	case DefaultJoin:
		pp.println("")
		pp.println("\v")
		pp.print("JOIN")
	case CommaJoin:
		pp.print(",")
	case CrossJoin:
		pp.println("")
		pp.println("\v")
		pp.print("CROSS JOIN")
	case FullJoin:
		pp.println("")
		pp.println("\v")
		pp.print("FULL JOIN")
	case InnerJoin:
		pp.println("")
		pp.println("\v")
		pp.print("INNER JOIN")
	case LeftJoin:
		pp.println("")
		pp.println("\v")
		pp.print("LEFT JOIN")
	case RightJoin:
		pp.println("")
		pp.println("\v")
		pp.print("RIGHT JOIN")
	}

	pp.println("")

	pp2 := p.nest()
	n.RHS.Accept(pp2, d)
	pp.print(pp2.unnest())

	if n.ClauseList != nil {
		pp.println("")
		pp.visitNested(n.ClauseList, d)
	}

	if n.OnClause != nil {
		pp.println("")
		pp.visitNested(n.OnClause, d)
	}

	if n.UsingClause != nil {
		pp.println("")
		pp.visitNested(n.UsingClause, d)
	}

	p.print(pp.unnestLeft())
}

func (p *printer) VisitOnClause(n *OnClause, d interface{}) {
	p1 := p.nest()
	p1.printClause(p1.keyword("ON"))
	n.Expression.Accept(p1, d)
	p.print(p1.unnestLeft())
}

func (p *printer) VisitUnnestExpression(
	n *UnnestExpression, d interface{}) {
	p.print(p.keyword("UNNEST") + "(")
	n.Expression.Accept(p, d)
	p.print(")")
}

func (p *printer) VisitBetweenExpression(
	n *BetweenExpression, d interface{}) {
	n.LHS.Accept(p, d)
	p.print(p.keyword("BETWEEN"))
	n.Low.Accept(p, d)
	p.print(p.keyword("AND"))
	n.High.Accept(p, d)
}

func (p *printer) VisitBinaryExpression(n *BinaryExpression, d interface{}) {
	var (
		ctx map[string]int
		ok  bool
	)

	p.printOpenParenIfNeeded(n)

	if ctx, ok = d.(map[string]int); !ok {
		ctx = make(map[string]int, 1)
	}

	var align string
	capacity := ctx["alignBinaryOp"]

	if capacity > 0 {
		ctx["alignBinaryOp"]--
		align = " \v"
	}

	n.LHS.Accept(p, ctx)

	if n.IsNot {
		switch n.Op { //nolint:exhaustive
		case BinaryIs:
			p.print(align + p.keyword("IS NOT") + align)
		case BinaryLike:
			p.print(align + p.keyword("NOT LIKE") + align)
		}
	} else {
		// The binary operator may be interpreted as a separator that allows
		// breaking line. In binary expressions, we do not want to allow
		// line breaking right after the operator.
		p.fmt.noFlushInNextFormat = true
		p.print(align + n.Op.String() + align)
	}

	n.RHS.Accept(p, ctx)

	p.printCloseParenIfNeeded(n)
}

func (p *printer) VisitInExpression(n *InExpression, d interface{}) {
	pp := p.nest()

	pp2 := pp.nest()

	n.LHS.Accept(pp2, d)
	pp2.print(p.keyword("IN"))
	pp.print(pp2.unnest())

	if n.InList != nil {
		pp2 = pp.nest()
		n.InList.Accept(pp2, d)
		pp.print(pp2.unnest())
	}

	if n.Query != nil {
		pp2 = pp.nest()
		n.Query.Accept(pp2, d)
		pp.print(pp2.unnest())
	}

	if n.UnnestExpr != nil {
		pp2 = pp.nest()
		n.UnnestExpr.Accept(pp2, d)
		pp.print(pp2.unnest())
	}

	p.print(pp.unnest())
}

func (p *printer) VisitCastExpression(n *CastExpression, d interface{}) {
	pp := p.nest()
	pp.print(p.keyword("CAST") + "(")
	n.Expr.Accept(pp, d)
	pp.print(p.keyword("AS"))
	n.Type.Accept(pp, d)

	if n.Format != nil {
		pp.print(p.keyword("FORMAT"))
		n.Format.Accept(pp, d)
	}

	pp.print(")")
	p.print(pp.unnest())
}

func (p *printer) VisitAnalyticFunctionCall(
	n *AnalyticFunctionCall, d interface{}) {
	pp := p.nest()

	pp2 := p.nest()
	n.Expr.Accept(pp2, d)
	pp.print(strings.Trim(pp2.unnest(), "\v"))

	pp.print(p.keyword("OVER") + " (")

	pp2 = pp.nest()
	n.WindowSpec.Accept(pp2, d)
	pp.print(pp2.unnest())

	pp.print(")")
	p.print(pp.unnest())
}

func (p *printer) VisitFunctionCall(n *FunctionCall, d interface{}) {
	pp := p.nest()
	pp.printOpenParenIfNeeded(n)
	n.Function.Accept(pp, n)

	// Strip off the alignment symbol at the beginning.
	expr := pp.unnest()[1:]

	pp = p.nest()

	pp.print(pp.function(expr))
	pp.print("(")

	if n.Distinct {
		pp.print(pp.keyword("DISTINCT"))
	}

	pp2 := pp.nest()

	for i, arg := range n.Arguments {
		if i > 0 {
			pp2.print(",")
		}

		pp3 := pp2.nest()
		arg.Accept(pp3, d)
		printedArg := strings.Trim(pp3.String(), "\n")
		if strings.Contains(printedArg, "\n") {
			pp2.println("")
		}
		pp2.print(printedArg)
	}

	pp.print(pp2.unnest())

	switch n.NullHandlingModifier {
	case DefaultNullHandling:
	case IgnoreNulls:
		pp.print(pp.keyword("IGNORE NULLS"))
	case RespectNulls:
		pp.print(pp.keyword("RESPECT NULLS"))
	}

	if n.OrderBy != nil {
		n.OrderBy.Accept(pp, d)
	}

	if n.LimitOffset != nil {
		n.LimitOffset.Accept(pp, d)
	}

	pp.print(")")
	pp.printCloseParenIfNeeded(n)

	p.print(pp.unnest())
}

func (p *printer) VisitWindowSpecification(
	n *WindowSpecification, d interface{}) {
	pp := p.nest()

	if n.BaseWindowName != nil {
		n.BaseWindowName.Accept(pp, d)
		pp.print(p.keyword("AS") + " (")
	}

	pp2 := pp.nest()
	forceAcrossLines := n.WindowFrame != nil

	if n.PartitionBy != nil {
		n.PartitionBy.Accept(pp2, d)
	}

	if n.OrderBy != nil {
		if forceAcrossLines && n.PartitionBy != nil {
			pp2.println("")
		}

		n.OrderBy.Accept(pp2, d)
	}

	if n.WindowFrame != nil {
		if forceAcrossLines && (n.PartitionBy != nil || n.OrderBy != nil) {
			pp2.println("")
		}

		n.WindowFrame.Accept(pp2, d)
	}

	pp.print(pp2.unnest())

	if n.BaseWindowName != nil {
		pp.print(")")
	}

	p.print(pp.unnest())
}

func (p *printer) VisitWindowFrame(n *WindowFrame, d interface{}) {
	pp := p.nest()
	pp.print(p.keyword(n.FrameUnit.String()) + " \v")
	pp.print(p.keyword("BETWEEN"))
	n.StartExpr.Accept(pp, d)
	pp.print(p.keyword("AND"))
	n.EndExpr.Accept(pp, d)
	p.print(strings.Trim(pp.String(), "\n"))
}

func (p *printer) VisitWindowFrameExpr(n *WindowFrameExpr, d interface{}) {
	pp := p.nest()

	if n.Expression != nil {
		n.Expression.Accept(pp, d)
	}

	p.print(p.keyword(n.BoundaryType.ToSQL()))
	p.print(pp.unnest())
}

func (p *printer) VisitAndExpr(n *AndExpr, d interface{}) {
	p.printOpenParenIfNeeded(n)

	p1 := p.nest()

	for i, conjunct := range n.Conjuncts {
		if i > 0 {
			if p.fmt.opts.AlignLogicalWithClauses &&
				(isInsideOfWhereClause(n) || isInsideOfOnClause(n)) {
				// Clear buffer and write AND as a clause.
				p.print(p1.unnest())
				p.printClause(p.keyword("AND"))
				// Create new nested builder.
				p1 = p.nest()
			} else {
				p1.printClause(p1.keyword("AND"))
			}
		}

		if conjunct.Kind() == AndExprKind {
			conjunct.Accept(p1, d)
		} else {
			p1.visitNested(conjunct, d)
		}
	}

	p.print(p1.unnestLeft())

	p.printCloseParenIfNeeded(n)
}

func (p *printer) VisitOrExpr(n *OrExpr, d interface{}) {
	p1 := p.nest()

	p1.printOpenParenIfNeeded(n)

	for i, disjunct := range n.Disjuncts {
		if i > 0 {
			if p1.fmt.opts.AlignLogicalWithClauses &&
				(isInsideOfWhereClause(n) && isInsideOfOnClause(n)) {
				// Clear buffer and write AND as a clause.
				p.print(p1.unnest())
				p.printClause(p.keyword("OR"))
				// Create new nested builder.
				p1 = p.nest()
			} else {
				p1.printClause(p1.keyword("OR"))
			}
		}

		p1.visitNested(disjunct, d)
	}

	p1.printCloseParenIfNeeded(n)
	p.print(p1.unnest())
}

// isInsideOfWhereClause returns true when the current node is is inside
// of a WHERE clause directly. The node can be inside of other AndExpr
// and OrExpr.
func isInsideOfWhereClause(n NodeHandler) bool {
	for p := n.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == WhereClauseKind {
			return true
		}

		if p.Kind() != AndExprKind && p.Kind() != OrExprKind {
			return false
		}
	}

	return false
}

// isInsideOfOnClause returns true when the current node is is inside
// of an ON clause directly. The node can be inside of other AndExpr
// and OrExpr.
func isInsideOfOnClause(n NodeHandler) bool {
	for p := n.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == OnClauseKind {
			return true
		}

		if p.Kind() != AndExprKind && p.Kind() != OrExprKind {
			return false
		}
	}

	return false
}

func (p *printer) VisitCaseNoValueExpression(
	n *CaseNoValueExpression, d interface{}) {
	p1 := p.nest()

	p1.print(p1.keyword("CASE") + " \v")

	args := n.Arguments
	for len(args) >= 2 {
		ctx := map[string]int{
			"alignBinaryOp": 1,
		}

		p1.println("")
		p1.print(p1.keyword("WHEN") + " \v")

		p2 := p1.nest()
		args[0].Accept(p2, ctx)
		p1.print(strings.Trim(p2.String(), "\n\v"))

		count := strings.Count(p1.fmt.buffer.String(), "\v")
		if count == 1 {
			p1.print("\v")
		}

		p1.print(p1.keyword("THEN") + " \v")
		p2 = p1.nest()
		args[1].Accept(p2, d)
		p1.print(p2.unnestWithDepth(4))
		args = args[2:]
	}

	if len(args) == 1 {
		p1.println("")
		p1.print(" ")
		p1.print("\v\v\v")
		p1.print(p1.keyword("ELSE") + " \v")
		p2 := p1.nest()
		args[0].Accept(p2, d)
		p1.print(p2.unnestWithDepth(4))
	}

	p1.println("")
	p1.print(p1.keyword("END"))

	p.print(p1.unnest())
}

func (p *printer) VisitCaseValueExpression(
	n *CaseValueExpression, d interface{}) {
	pp := p.nest()

	pp.print(pp.keyword("CASE"))

	pp2 := pp.nest()
	n.Arguments[0].Accept(pp2, d)
	pp.print(pp2.unnest())

	args := n.Arguments[1:]
	for len(args) >= 2 {
		ctx := map[string]int{
			"alignBinaryOp": 1,
		}

		pp.println("")
		pp.print(pp.keyword("WHEN"))

		pp2 := pp.nest()
		args[0].Accept(pp2, ctx)
		pp.print("\v" + strings.Trim(pp2.String(), "\n"))

		count := strings.Count(pp.fmt.buffer.String(), "\v")
		if count == 1 {
			pp.print("\v")
		}

		pp.print("\v" + pp.keyword("THEN"))

		pp2 = pp.nest()
		args[1].Accept(pp2, d)
		pp.print(pp2.unnestWithDepth(4))

		args = args[2:]
	}

	if len(args) == 1 {
		pp.println("")
		pp.print(" ")
		pp.print("\v\v\v")
		pp.print(pp.keyword("ELSE") + " \v")
		p2 := pp.nest()
		args[0].Accept(p2, d)
		pp.print(p2.unnestWithDepth(4))
	}

	pp.println("")
	pp.print(pp.keyword("END"))

	p.print(pp.unnest())
}

func (p *printer) VisitAlias(n *Alias, d interface{}) {
	p.print(p.keyword("AS"))
	p.print(p.identifier(n.Identifier.IDString))
}

func (p *printer) VisitTablePathExpression(
	n *TablePathExpression, d interface{},
) {
	if n.PathExpr != nil {
		n.PathExpr.Accept(p, d)
	} else {
		n.UnnestExpr.Accept(p, d)
	}

	if n.Alias != nil {
		n.Alias.Accept(p, d)
	}

	if n.WithOffset != nil {
		n.WithOffset.Accept(p, d)
	}

	// TODO PivotClause
	// TODO UnpivotClause
	// TODO SampleClause
}

func (p *printer) VisitWithOffset(n *WithOffset, d interface{}) {
	p.print(p.keyword("WITH OFFSET"))

	if n.Alias != nil {
		n.Alias.Accept(p, d)
	}
}

func (p *printer) VisitPathExpression(n *PathExpression, d interface{}) {
	p.printOpenParenIfNeeded(n)

	for i, name := range n.Names {
		if i > 0 {
			p.print(".")
		}

		name.Accept(p, d)
	}

	p.printCloseParenIfNeeded(n)
}

func (p *printer) VisitIdentifier(n *Identifier, d interface{}) {
	id := ToIdentifierLiteral(n.IDString)
	if id[0] == '`' {
		// Cannot format quoted identifiers.
		p.print(id)
		return
	}

	switch p.fmt.opts.IdentifierStyle {
	case AsIsCase:
		p.print(id)
	case LowerCase:
		p.print(strings.ToLower(id))
	case UpperCase:
		p.print(strings.ToUpper(id))
	}
}

func (p *printer) VisitNamedType(n *NamedType, d interface{}) {
	pp := p.nest()
	n.Name.Accept(pp, d)
	typename := strings.Trim(pp.String(), "\n")
	p.print(p.typename(typename))
}

func (p *printer) VisitBooleanLiteral(n *BooleanLiteral, d interface{}) {
	if n.Value {
		p.print(p.keyword("TRUE"))
	} else {
		p.print(p.keyword("FALSE"))
	}
}

func (p *printer) VisitNullLiteral(n *NullLiteral, d interface{}) {
	p.print(p.keyword("NULL"))
}

func (p *printer) VisitIntervalExpr(n *IntervalExpr, d interface{}) {
	p.print(p.keyword("INTERVAL"))
	n.IntervalValue.Accept(p, d)

	pp := p.nest()
	n.DatePartName.Accept(pp, d)
	p.print(p.keyword(pp.unnest()))
}

func (p *printer) VisitUsingClause(n *UsingClause, d interface{}) {
	p.printClause(p.keyword("USING") + " (")
	p.printNestedWithSep(n.Keys, d, ",")
	p.print(")")
}

func (p *printer) VisitWithClause(n *WithClause, d interface{}) {
	pp := p.nest()
	pp.println(p.keyword("WITH"))

	for i, entry := range n.With {
		if i > 0 {
			// Force empty line.
			pp.print("\n")
		}

		entry.Alias.Accept(pp, d)
		pp.println(p.keyword("AS") + " (")
		entry.Query.Accept(pp, d)
		pp.print(")")

		if i+1 < len(n.With) {
			pp.print(",")
		}

		pp.println("")
	}

	p.print(pp.unnest())
}

func (p *printer) VisitDateOrTimeLiteral(n *DateOrTimeLiteral, d interface{}) {
	p.print(n.TypeKind.ToSQL())
	n.StringLiteral.Accept(p, d)
}

func (p *printer) VisitStringLiteral(n *StringLiteral, d interface{}) {
	var style literal.StringStyle

	switch p.fmt.opts.StringStyle {
	case AsIsStringStyle:
		p.print(n.Image())

		return
	case PreferSingleQuote:
		style = literal.PreferSingleQuote
	case PreferDoubleQuote:
		style = literal.PreferDoubleQuote
	case AlwaysSingleQuote:
		style = literal.AlwaysSingleQuote
	case AlwaysDoubleQuote:
		style = literal.AlwaysDoubleQuote
	}

	p.print(literal.Escape(n.StringValue, style))
}

func (p *printer) VisitArrayConstructor(n *ArrayConstructor, d interface{}) {
	pp := p.nest()

	if n.Type != nil {
		n.Type.Accept(pp, d)
	}

	pp.print("[")
	pp.printNestedWithSep(n.Elements, d, ",")
	pp.print("]")
	p.print(pp.unnest())
}

func (p *printer) VisitArrayType(n *ArrayType, d interface{}) {
	pp := p.nest()

	pp2 := pp.nest()
	n.ElementType.Accept(pp2, d)
	typeSpec := strings.Trim(pp2.String(), "\n")

	pp.print(p.keyword("ARRAY") + "<" + typeSpec + ">")
	p.print(pp.unnest())
}

func (p *printer) VisitArrayElement(n *ArrayElement, d interface{}) {
	pp := p.nest()

	pp2 := pp.nest()
	n.Array.Accept(pp2, d)
	pp.print(pp2.unnest())

	pp.print("[")
	pp.visitNested(n.Position, d)
	pp.print("]")

	p.print(pp.unnest())
}

func (p *printer) VisitStructConstructorWithKeyword(
	n *StructConstructorWithKeyword, d interface{}) {
	pp := p.nest()

	if n.StructType != nil {
		n.StructType.Accept(pp, d)
	} else {
		pp.print(pp.keyword("STRUCT"))
	}

	pp.print("(")
	pp.printNestedWithSep(n.Fields, d, ",")
	pp.print(")")

	p.print(pp.unnest())
}

func (p *printer) VisitStructConstructorWithParens(
	n *StructConstructorWithParens, d interface{}) {
	pp := p.nest()

	pp.print("(")
	pp.printNestedWithSep(n.FieldExpressions, d, ",")
	pp.print(")")
	p.print(pp.unnest())
}

func (p *printer) VisitStructType(n *StructType, d interface{}) {
	pp := p.nest()
	pp2 := pp.nest()

	pp2.printNestedWithSep(n.StructFields, d, ",")

	pp.print(pp.keyword("STRUCT") + "<" + pp2.unnest() + ">")
	p.print(pp.unnest())
}

func (p *printer) VisitStructField(n *StructField, d interface{}) {
	root := p.nest()

	n.Name.Accept(root, d)
	root.visitNested(n.Type, d)
	p.print(root.unnest())
}

func (p *printer) VisitStructConstructorArg(
	n *StructConstructorArg, d interface{}) {
	pp := p.nest()

	n.Expression.Accept(pp, d)

	if n.Alias != nil {
		n.Alias.Accept(pp, d)
	}

	p.print(pp.unnest())
}

func (p *printer) VisitInList(n *InList, d interface{}) {
	pp := p.nest()

	pp.print("(")
	pp.printNestedWithSep(n.List, d, ",")
	pp.print(")")
	p.print(pp.unnest())
}

func estimateSQLSize(n NodeHandler) int {
	return len(Sprint(n, onelinePrintConfig))
}

func alignNested(s string) string {
	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 0, ' ', tabwriter.AlignRight)

	fmt.Fprint(w, s)
	w.Flush()

	return strings.Trim(buf.String(), "\n")
}

func leftAlignNested(s string) string {
	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 0, ' ', 0)

	fmt.Fprint(w, s)
	w.Flush()

	return strings.Trim(buf.String(), "\n")
}

// Format formats the string automatically according to the context.
// 1. Inserts necessary space between tokens.
// 2. Calls FlushLine() when a line reachs column limit and it is at
//    some point appropriate to break.
// Param s should not contain any leading or trailing whitespace, such
// as ' ' and '\n'.
func (p *formatter) Format(s string) {
	if len(s) == 0 {
		return
	}

	// At the end we check whether the buffer should be flushed to
	// the formatted buffer.
	defer func() {
		if p.buffer.Len() >= p.maxLength &&
			p.lastIsSeparator() &&
			!p.noFlushInNextFormat {
			p.FlushLine()
		}

		p.lastWasSingleCharUnary = false
		p.noFlushInNextFormat = false
	}()

	data := []rune(s)

	if p.buffer.Len() == 0 {
		p.writeIndent()
		p.writeRunes(data)

		return
	}

	switch p.last {
	case '\n':
		p.writeRunes(append([]rune{'\n'}, data...))
	case '(', '[', '@', '.', '~', ' ', '\v', '\b':
		p.writeRunes(data)
	default:
		if p.lastWasSingleCharUnary {
			p.writeRunes(data)
			return
		}

		curr := data[0]
		if curr == '(' {
			// Inserts a space if last token is a separator, otherwise
			// regards it as a function call.
			if p.lastIsSeparator() {
				p.writeRunes(append([]rune{' '}, data...))
			} else {
				p.writeRunes(data)
			}

			return
		}

		if curr == ')' ||
			curr == '[' ||
			curr == ']' ||
			// To avoid case like "SELECT 1e10,.1e10"
			(curr == '.' && p.last != ',') ||
			curr == ',' {
			p.writeRunes(data)
			return
		}

		if p.last == ' ' && data[0] == ' ' {
			p.writeRunes(data)
		} else {
			p.writeRunes(append([]rune{' '}, data...))
		}
	}
}

// FormatLine is like Format, except always calls FlushLine.
// Use this if you explicitly wants to break the line after this string.
// For example:
// 1. To put a newline after SELECT:
// 		FormatLine("SELECT")
// 2. To put close parenthesis on a separate line:
//		FormatLine("")
//		FormatLine(")")
func (p *formatter) FormatLine(s string) {
	p.Format(s)
	p.FlushLine()
}

// FlushLine flushes buffer to formatted, with a line break at the end.
// It will do nothing if it is a new line and buffer is empty, to avoid
// empty lines.
// Remember to call FlushLine once after the whole process is over in
// case some content remains in buffer.
func (p *formatter) FlushLine() {
	fmt := p.formatted.String()
	sz := len(fmt)

	if (sz == 0 || fmt[sz-1] == '\n') && p.buffer.Len() == 0 {
		return
	}

	p.formatted.WriteString(p.buffer.String())
	p.formatted.WriteByte('\n')
	p.buffer.Reset()
}

func (p *formatter) lastIsSeparator() bool {
	if p.buffer.Len() == 0 {
		return false
	}

	if !isAlphanum(byte(p.last)) {
		return nonWordSeparators[p.last]
	}

	buf := p.buffer.String()

	i := len(buf) - 1
	for i >= 0 && isAlphanum(buf[i]) {
		i--
	}

	lastTok := buf[i+1:]

	return wordSeparators[lastTok]
}

func (p *formatter) addUnary(s string) {
	if p.lastWasSingleCharUnary && p.last == '-' && s == "-" {
		p.lastWasSingleCharUnary = false
	}

	p.Format(s)
	p.lastWasSingleCharUnary = len(s) == 1
}

func (p *formatter) writeIndent() {
	p.buffer.WriteString(strings.Repeat(" ", p.depth*2))
}

func (p *formatter) writeRunes(d []rune) {
	p.buffer.WriteString(string(d))
	p.last = d[len(d)-1]
}
