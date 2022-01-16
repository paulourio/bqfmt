package ast

import (
	"strings"
)

func Unparse(n NodeHandler) string {
	u := &unparser{
		fmt: &unparseFormatter{
			buffer:      strings.Builder{},
			unparsed:    strings.Builder{},
			depth:       0,
			numColLimit: 100,
		},
		Operation: Operation{},
	}
	u.Operation.visitor = u
	n.Accept(u, nil)
	u.fmt.FlushLine()

	return u.fmt.unparsed.String()
}

type unparser struct {
	fmt *unparseFormatter

	Operation
}

func (u *unparser) print(s string) {
	u.fmt.Format(s)
}

func (u *unparser) println(s string) {
	u.fmt.FormatLine(s)
}

func (u *unparser) incDepth() {
	u.fmt.depth++
}

func (u *unparser) decDepth() {
	u.fmt.depth--
}

func (u *unparser) printOpenParenIfNeeded(n NodeHandler) {
	if !(n.IsExpression() || n.IsQueryExpression()) {
		panic("parenthesization is not allowed for " + n.Kind().String())
	}

	if expr, ok := n.(ExpressionHandler); ok && expr.IsParenthesized() {
		u.print("(")
		return
	}

	if expr, ok := n.(QueryExpressionHandler); ok && expr.IsParenthesized() {
		u.print("(")
		return
	}
}

func (u *unparser) printCloseParenIfNeeded(n NodeHandler) {
	if !(n.IsExpression() || n.IsQueryExpression()) {
		panic("parenthesization is not allowed for " + n.Kind().String())
	}

	if expr, ok := n.(ExpressionHandler); ok && expr.IsParenthesized() {
		u.print(")")
		return
	}

	if expr, ok := n.(QueryExpressionHandler); ok && expr.IsParenthesized() {
		u.print(")")
		return
	}
}

func (u *unparser) VisitQuery(n *Query, d interface{}) {
	u.printOpenParenIfNeeded(n)

	if n.IsNested {
		u.println("")
		u.print("(")
		u.incDepth()
		n.QueryExpr.Accept(u, d)

		if n.OrderBy != nil {
			n.OrderBy.Accept(u, d)
		}

		if n.LimitOffset != nil {
			n.LimitOffset.Accept(u, d)
		}

		u.decDepth()
		u.println("")
		u.print(")")
	} else {
		n.QueryExpr.Accept(u, d)

		if n.OrderBy != nil {
			n.OrderBy.Accept(u, d)
		}

		if n.LimitOffset != nil {
			n.LimitOffset.Accept(u, d)
		}
	}

	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitOrderBy(n *OrderBy, d interface{}) {
	u.println("")
	u.print("ORDER BY")
	u.incDepth()

	for i, item := range n.OrderingExpression {
		if i > 0 {
			u.print(",")
		}

		item.Accept(u, d)
	}

	u.decDepth()
}

func (u *unparser) VisitOrderingExpression(
	n *OrderingExpression, d interface{}) {
	n.Expression.Accept(u, d)

	switch n.OrderingSpec {
	case DescendingOrder:
		u.print("DESC")
	case AscendingOrder:
		u.print("ASC")
	case NoOrderingSpec:
		/* nothing */
		break
	}

	if n.NullOrder != nil {
		n.NullOrder.Accept(u, d)
	}
}

func (u *unparser) VisitNullOrder(n *NullOrder, d interface{}) {
	if n.NullsFirst {
		u.print("NULLS FIRST")
	} else {
		u.print("NULLS LAST")
	}
}

func (u *unparser) VisitLimitOffset(n *LimitOffset, d interface{}) {
	u.println("")
	u.print("LIMIT")
	n.Limit.Accept(u, d)

	if n.Offset != nil {
		u.print("OFFSET")
		n.Offset.Accept(u, d)
	}
}

func (u *unparser) VisitTableSubquery(n *TableSubquery, d interface{}) {
	n.Subquery.Accept(u, d)

	if n.Alias != nil {
		u.print("AS")
		u.print(ToIdentifierLiteral(n.Alias.Identifier.IDString))
	}
}

func (u *unparser) VisitSelect(n *Select, d interface{}) {
	u.printOpenParenIfNeeded(n)
	u.println("")
	u.print("SELECT")
	n.SelectList.Accept(u, d)

	if n.FromClause != nil {
		u.println("")
		u.println("FROM")
		u.incDepth()
		n.FromClause.Accept(u, d)
		u.decDepth()
	}

	if n.WhereClause != nil {
		u.println("")
		u.println("WHERE")
		u.incDepth()
		n.WhereClause.Accept(u, d)
		u.decDepth()
	}

	if n.GroupBy != nil {
		n.GroupBy.Accept(u, d)
	}

	if n.Having != nil {
		u.println("")
		u.println("HAVING")
		u.incDepth()
		n.Having.Accept(u, d)
		u.decDepth()
	}

	if n.Qualify != nil {
		u.println("")
		u.println("QUALIFY")
		u.incDepth()
		n.Qualify.Accept(u, d)
		u.decDepth()
	}

	if n.WindowClause != nil {
		u.println("")
		u.println("WINDOW")
		u.incDepth()
		n.WindowClause.Accept(u, d)
		u.decDepth()
	}

	u.println("")
	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitGroupBy(n *GroupBy, d interface{}) {
	u.println("")
	u.print("GROUP BY")
	u.incDepth()

	for i, item := range n.GroupingItems {
		if i > 0 {
			u.print(",")
		}

		item.Expression.Accept(u, d)
	}

	u.decDepth()
}

func (u *unparser) VisitSelectList(n *SelectList, d interface{}) {
	u.println("")
	u.incDepth()

	for i, col := range n.Columns {
		if i > 0 {
			u.println(",")
		}

		col.Accept(u, d)
	}

	u.decDepth()
}

func (u *unparser) VisitSelectColumn(n *SelectColumn, d interface{}) {
	n.Expression.Accept(u, d)

	if n.Alias != nil {
		u.print("AS")
		u.print(ToIdentifierLiteral(n.Alias.Identifier.IDString))
	}
}

func (u *unparser) VisitLeafHandler(n LeafHandler, d interface{}) {
	u.print(n.Image())
}

func (u *unparser) VisitBinaryExpression(n *BinaryExpression, d interface{}) {
	u.printOpenParenIfNeeded(n)

	n.LHS.Accept(u, d)

	if n.IsNot {
		switch n.Op { //nolint:exhaustive
		case BinaryIs:
			u.print("IS NOT ")
		case BinaryLike:
			u.print("NOT LIKE ")
		}
	} else {
		u.print(n.Op.String())
	}

	n.RHS.Accept(u, d)

	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitUnaryExpression(n *UnaryExpression, d interface{}) {
	u.printOpenParenIfNeeded(n)
	u.fmt.addUnary(n.Op.String())
	n.Operand.Accept(u, d)
	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitBetweenExpression(
	n *BetweenExpression, d interface{}) {
	n.LHS.Accept(u, d)
	u.print("BETWEEN")
	n.Low.Accept(u, d)
	u.print("AND")
	n.High.Accept(u, d)
}

func (u *unparser) VisitJoin(n *Join, d interface{}) {
	n.LHS.Accept(u, d)

	switch n.JoinType {
	case DefaultJoin:
		u.println("")
		u.print("JOIN")
	case CommaJoin:
		u.print(",")
	case CrossJoin:
		u.println("")
		u.print("CROSS JOIN")
	case FullJoin:
		u.println("")
		u.print("FULL JOIN")
	case InnerJoin:
		u.println("")
		u.print("INNER JOIN")
	case LeftJoin:
		u.println("")
		u.print("LEFT JOIN")
	case RightJoin:
		u.println("")
		u.print("RIGHT JOIN")
	}

	u.println("")
	n.RHS.Accept(u, d)

	if n.ClauseList != nil {
		u.println("")
		n.ClauseList.Accept(u, d)
	}

	if n.OnClause != nil {
		u.println("")
		n.OnClause.Accept(u, d)
	}

	if n.UsingClause != nil {
		u.println("")
		n.UsingClause.Accept(u, d)
	}
}

func (u *unparser) VisitParenthesizedJoin(
	n *ParenthesizedJoin, d interface{}) {
	u.println("")
	u.println("(")
	u.incDepth()
	n.Join.Accept(u, d)
	u.decDepth()
	u.println("")
	u.print(")")

	if n.SampleClause != nil {
		n.SampleClause.Accept(u, d)
	}
}

func (u *unparser) VisitUnnestExpression(
	n *UnnestExpression, d interface{}) {
	u.print("UNNEST(")
	n.Expression.Accept(u, d)
	u.print(")")
}

func (u *unparser) VisitSampleClause(n *SampleClause, d interface{}) {
	u.print("TABLESAMPLE")
	n.SampleMethod.Accept(u, d)
	u.print("(")
	n.SampleSize.Accept(u, d)
	u.print(")")

	if n.SampleSuffix != nil {
		n.SampleSuffix.Accept(u, d)
	}
}

func (u *unparser) VisitSampleSize(n *SampleSize, d interface{}) {
	n.Size.Accept(u, d)
	u.print(n.Unit.String())

	if n.PartitionBy != nil {
		n.PartitionBy.Accept(u, d)
	}
}

func (u *unparser) VisitSampleSuffix(n *SampleSuffix, d interface{}) {
	if n.Weight != nil {
		u.print("WITH WEIGHT")
		n.Weight.Alias.Accept(u, d)
	}

	if n.Repeat != nil {
		u.print("REPEATABLE (")
		n.Repeat.Argument.Accept(u, d)
		u.print(")")
	}
}

func (u *unparser) VisitOnClause(n *OnClause, d interface{}) {
	u.println("")
	u.print("ON")
	u.incDepth()
	n.Expression.Accept(u, d)
	u.decDepth()
}

func (u *unparser) VisitUsingClause(n *UsingClause, d interface{}) {
	u.println("")
	u.print("USING(")
	u.incDepth()

	for i, key := range n.Keys {
		if i > 0 {
			u.print(",")
		}

		key.Accept(u, d)
	}

	u.decDepth()
	u.print(")")
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
		n.Alias.Accept(u, d)
	}
}

func (u *unparser) VisitAlias(n *Alias, d interface{}) {
	u.print("AS")
	u.print(ToIdentifierLiteral(n.Identifier.IDString))
}

func (u *unparser) VisitIdentifier(n *Identifier, d interface{}) {
	u.print(ToIdentifierLiteral(n.IDString))
}

func (u *unparser) VisitDotIdentifier(n *DotIdentifier, d interface{}) {
	u.printOpenParenIfNeeded(n)
	n.Expr.Accept(u, d)
	u.print(".")
	n.Name.Accept(u, d)
	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitPathExpression(n *PathExpression, d interface{}) {
	u.printOpenParenIfNeeded(n)

	for i, p := range n.Names {
		if i > 0 {
			u.print(".")
		}

		p.Accept(u, d)
	}

	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitCastExpression(n *CastExpression, d interface{}) {
	if n.IsSafeCast {
		u.print("SAFE_CAST(")
	} else {
		u.print("CAST(")
	}

	n.Expr.Accept(u, d)
	u.print("AS")
	n.Type.Accept(u, d)

	if n.Format != nil {
		if n.Format.Format != nil {
			u.print("FORMAT")
			n.Format.Format.Accept(u, d)
		}

		if n.Format.TimeZoneExpr != nil {
			u.print("AT TIME ZONE")
			n.Format.TimeZoneExpr.Accept(u, d)
		}
	}

	u.print(")")
}

func (u *unparser) VisitFunctionCall(n *FunctionCall, d interface{}) {
	u.printOpenParenIfNeeded(n)
	n.Function.Accept(u, n)
	u.print("(")
	u.incDepth()

	if n.Distinct {
		u.print("DISTINCT")
	}

	for i, arg := range n.Arguments {
		if i > 0 {
			u.print(",")
		}

		arg.Accept(u, d)
	}

	switch n.NullHandlingModifier {
	case DefaultNullHandling:
	case IgnoreNulls:
		u.print("IGNORE NULLS")
	case RespectNulls:
		u.print("RESPECT NULLS")
	}

	if n.OrderBy != nil {
		n.OrderBy.Accept(u, d)
	}

	if n.LimitOffset != nil {
		n.LimitOffset.Accept(u, d)
	}

	u.decDepth()
	u.print(")")
	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitAndExpr(n *AndExpr, d interface{}) {
	u.printOpenParenIfNeeded(n)

	for i, c := range n.Conjuncts {
		if i > 0 {
			u.print("AND")
		}

		c.Accept(u, d)
	}

	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitOrExpr(n *OrExpr, d interface{}) {
	u.printOpenParenIfNeeded(n)

	for i, c := range n.Disjuncts {
		if i > 0 {
			u.print("OR")
		}

		c.Accept(u, d)
	}

	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitStar(n *Star, d interface{}) {
	u.print(n.Image())
}

func (u *unparser) VisitBooleanLiteral(n *BooleanLiteral, d interface{}) {
	u.print(n.Image())
}

func (u *unparser) VisitNumericLiteral(n *NumericLiteral, d interface{}) {
	u.print("NUMERIC " + n.Image())
}

func (u *unparser) VisitBigNumericLiteral(n *BigNumericLiteral, d interface{}) {
	u.print("BIGNUMERIC " + n.Image())
}

func (u *unparser) VisitJSONLiteral(n *JSONLiteral, d interface{}) {
	u.print("JSON " + n.Image())
}

func (u *unparser) VisitNamedType(n *NamedType, d interface{}) {
	n.Name.Accept(u, d)

	if n.TypeParameters() != nil {
		n.TypeParameters().Accept(u, d)
	}
}

func (u *unparser) VisitArrayType(n *ArrayType, d interface{}) {
	u.print("ARRAY<")
	n.ElementType.Accept(u, d)
	u.print(">")

	if n.TypeParameters() != nil {
		n.TypeParameters().Accept(u, d)
	}
}

func (u *unparser) VisitStructType(n *StructType, d interface{}) {
	u.print("STRUCT<")

	for i, f := range n.StructFields {
		if i > 0 {
			u.print(",")
		}

		f.Accept(u, d)
	}

	u.print(">")

	if n.TypeParameters() != nil {
		n.TypeParameters().Accept(u, d)
	}
}

func (u *unparser) VisitTypeParameterList(
	n *TypeParameterList, d interface{}) {
	u.print("(")
	u.incDepth()

	for i, p := range n.Parameters {
		if i > 0 {
			u.print(",")
		}

		p.Accept(u, d)
	}

	u.decDepth()
	u.print(")")
}

func (u *unparser) VisitStructField(n *StructField, d interface{}) {
	if n.Name != nil {
		n.Name.Accept(u, d)
	}

	n.Type.Accept(u, d)
}

func (u *unparser) VisitArrayConstructor(n *ArrayConstructor, d interface{}) {
	if n.Type != nil {
		n.Type.Accept(u, d)
	} else {
		u.print("ARRAY")
	}

	u.print("[")

	for i, e := range n.Elements {
		if i > 0 {
			u.print(",")
		}

		e.Accept(u, d)
	}

	u.print("]")
}

func (u *unparser) VisitArrayElement(n *ArrayElement, d interface{}) {
	u.printOpenParenIfNeeded(n)
	n.Array.Accept(u, d)
	u.print("[")
	n.Position.Accept(u, d)
	u.print("]")
	u.printCloseParenIfNeeded(n)
}

func (u *unparser) VisitNullLiteral(n *NullLiteral, d interface{}) {
	u.print("NULL")
}

func (u *unparser) VisitCaseValueExpression(
	n *CaseValueExpression, d interface{}) {
	u.print("CASE")
	n.Arguments[0].Accept(u, d)
	u.incDepth()

	args := n.Arguments[1:]
	for len(args) >= 2 {
		u.println("")
		u.print("WHEN")
		args[0].Accept(u, d)
		u.print("THEN")
		args[1].Accept(u, d)
		args = args[2:]
	}

	if len(args) == 1 {
		u.println("")
		u.print("ELSE")
		args[0].Accept(u, d)
	}

	u.decDepth()
	u.println("")
	u.print("END")
}

func (u *unparser) VisitCaseNoValueExpression(
	n *CaseNoValueExpression, d interface{}) {
	u.print("CASE")
	u.incDepth()

	args := n.Arguments
	for len(args) >= 2 {
		u.println("")
		u.print("WHEN")
		args[0].Accept(u, d)
		u.print("THEN")
		args[1].Accept(u, d)
		args = args[2:]
	}

	if len(args) == 1 {
		u.println("")
		u.print("ELSE")
		args[0].Accept(u, d)
	}

	u.decDepth()
	u.println("")
	u.print("END")
}

func (u *unparser) VisitStarWithModifiers(
	n *StarWithModifiers, d interface{}) {
	u.print("*")
	n.Modifiers.Accept(u, d)
}

func (u *unparser) VisitStarReplaceItem(n *StarReplaceItem, d interface{}) {
	n.Expression.Accept(u, d)
	u.print("AS")
	n.Alias.Accept(u, d)
}

func (u *unparser) VisitDotStar(n *DotStar, d interface{}) {
	n.Expr.Accept(u, d)
	u.print(".*")
}

func (u *unparser) VisitDotStarWithModifiers(
	n *DotStarWithModifiers, d interface{}) {
	n.Expr.Accept(u, d)
	u.print(".*")
}

func (u *unparser) VisitStarModifiers(n *StarModifiers, d interface{}) {
	if n.ExceptList != nil {
		u.print("EXCEPT (")

		for i, id := range n.ExceptList.Identifiers {
			if i > 0 {
				u.print(",")
			}

			u.print(ToIdentifierLiteral(id.IDString))
		}

		u.print(")")
	}

	if len(n.ReplaceItems) > 0 {
		u.print("REPLACE (")

		for i, item := range n.ReplaceItems {
			if i > 0 {
				u.print(",")
			}

			item.Accept(u, d)
		}

		u.print(")")
	}
}

// This is a fixed formatter used for unparsing an AST.
type unparseFormatter struct {
	// buf is the buffer of formatted output.
	buffer   strings.Builder
	unparsed strings.Builder
	// depth is indentation level that will be prepend to a new line.
	depth       int
	numColLimit int
	// last is the last char written to the buf.
	last                   rune
	lastWasSingleCharUnary bool
}

// Format formats the string automatically according to the context.
// 1. Inserts necessary space between tokens.
// 2. Calls FlushLine() when a line reachs column limit and it is at
//    some point appropriate to break.
// Param s should not contain any leading or trailing whitespace, such
// as ' ' and '\n'.
func (u *unparseFormatter) Format(s string) {
	if len(s) == 0 {
		return
	}

	// At the end we check whether the buffer should be flushed to
	// the unparsed buffer.
	defer func() {
		if u.buffer.Len() >= u.depth*2+u.numColLimit && u.lastIsSeparator() {
			u.FlushLine()
		}

		u.lastWasSingleCharUnary = false
	}()

	data := []rune(s)

	if u.buffer.Len() == 0 {
		u.writeIndent()
		u.writeRunes(data)

		return
	}

	switch u.last {
	case '\n':
		u.writeRunes([]rune{'\n'})
	case '(', '[', '@', '.', '~', ' ':
		u.writeRunes(data)
	default:
		if u.lastWasSingleCharUnary {
			u.writeRunes(data)
			return
		}

		curr := data[0]
		if curr == '(' {
			// Inserts a space if last token is a separator, otherwise
			// regards it as a function call.
			if u.lastIsSeparator() {
				u.writeRunes(append([]rune{' '}, data...))
			} else {
				u.writeRunes(data)
			}

			return
		}

		if curr == ')' ||
			curr == '[' ||
			curr == ']' ||
			// To avoid case like "SELECT 1e10,.1e10"
			(curr == '.' && u.last != ',') ||
			curr == ',' {
			u.writeRunes(data)
			return
		}

		u.writeRunes(append([]rune{' '}, data...))
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
func (u *unparseFormatter) FormatLine(s string) {
	u.Format(s)
	u.FlushLine()
}

// FlushLine flushes buffer to unparsed, with a line break at the end.
// It will do nothing if it is a new line and buffer is empty, to avoid
// empty lines.
// Remember to call FlushLine once after the whole process is over in
// case some content remains in buffer.
func (u *unparseFormatter) FlushLine() {
	unparsed := u.unparsed.String()
	if (len(unparsed) == 0 || unparsed[len(unparsed)-1] == '\n') &&
		u.buffer.Len() == 0 {
		return
	}

	u.unparsed.WriteString(u.buffer.String())
	u.unparsed.WriteByte('\n')
	u.buffer.Reset()
}

func (u *unparseFormatter) lastIsSeparator() bool {
	if u.buffer.Len() == 0 {
		return false
	}

	if !isAlphanum(byte(u.last)) {
		return nonWordSeparators[u.last]
	}

	buf := u.buffer.String()

	i := len(buf) - 1
	for i >= 0 && isAlphanum(buf[i]) {
		i--
	}

	lastTok := buf[i+1:]

	return wordSeparators[lastTok]
}

func (u *unparseFormatter) addUnary(s string) {
	if u.lastWasSingleCharUnary && u.last == '-' && s == "-" {
		u.lastWasSingleCharUnary = false
	}

	u.Format(s)
	u.lastWasSingleCharUnary = len(s) == 1
}

func (u *unparseFormatter) writeIndent() {
	u.buffer.WriteString(strings.Repeat(" ", u.depth*2))
}

func (u *unparseFormatter) writeRunes(d []rune) {
	u.buffer.WriteString(string(d))
	u.last = d[len(d)-1]
}

var wordSeparators = map[string]bool{
	"AND": true,
	"OR":  true,
	"ON":  true,
	"IN":  true,
}

var nonWordSeparators = map[rune]bool{
	',': true,
	'<': true,
	'>': true,
	'-': true,
	'+': true,
	'=': true,
	'*': true,
	'/': true,
	'%': true,
}
