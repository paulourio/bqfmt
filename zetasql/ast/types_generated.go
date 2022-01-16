package ast

// types_generated.go is generated from type_generated.go.j2 by
// gen_types.py.

import "fmt"

// defaultCapacity is the default capacity for new slices.
const defaultCapacity = 4

// QueryStatement represents a single query statement.
type QueryStatement struct {
	Query *Query

	Statement
}

type Query struct {
	// WithClause is the WITH clause wrapping this query.
	WithClause *WithClause
	// QueryExpr can be a single Select, or a more complex structure
	// composed out of nodes like SetOperation and Query.
	QueryExpr QueryExpressionHandler
	// OrderBy applies, if present, to the result of QueryExpr.
	OrderBy *OrderBy
	// LimitOffset applies, if present, after the result of QueryExpr
	// and OrderBy.
	LimitOffset *LimitOffset
	IsNested    bool

	QueryExpression
}

type Select struct {
	Distinct     bool
	SelectAs     *SelectAs
	SelectList   *SelectList
	FromClause   *FromClause
	WhereClause  *WhereClause
	GroupBy      *GroupBy
	Having       *Having
	Qualify      *Qualify
	WindowClause *WindowClause

	QueryExpression
}

type SelectList struct {
	Columns []*SelectColumn

	Node
}

type SelectColumn struct {
	Expression ExpressionHandler
	Alias      *Alias

	Node
}

type IntLiteral struct {
	Leaf
}

type Identifier struct {
	IDString string

	Expression
}

type Alias struct {
	Identifier *Identifier

	Node
}

// PathExpression is used for dotted identifier paths only, not
// dotting into arbitrary expressions (see DotIdentifier).
type PathExpression struct {
	Names []*Identifier

	Expression
}

// TablePathExpression is a table expression than introduce a single
// scan, referenced by a path expression or UNNEST, can optionally have
// aliases.  Exactly one of PathExpr or UnnestExpr must be non nil.
type TablePathExpression struct {
	PathExpr      *PathExpression
	UnnestExpr    *UnnestExpression
	Alias         *Alias
	WithOffset    *WithOffset
	PivotClause   interface{}
	UnpivotClause interface{}
	SampleClause  *SampleClause

	TableExpression
}

type FromClause struct {
	// TableExpression has exactly one table expression child.  If
	// the FROM clause has commas, they will expressed as a tree of Join
	// nodes with JoinType=Comma.
	TableExpression TableExpressionHandler

	Node
}

type WhereClause struct {
	Expression ExpressionHandler

	Node
}

type BooleanLiteral struct {
	Value bool

	Leaf
}

type AndExpr struct {
	Conjuncts []ExpressionHandler

	Expression
}

type BinaryExpression struct {
	Op  BinaryOp
	LHS ExpressionHandler
	RHS ExpressionHandler
	// IsNot indicates whether the binary operator has a preceding
	// NOT to it.  For NOT LIKE and IS NOT.
	IsNot bool

	Expression
}

type StringLiteral struct {
	StringValue string

	Leaf
}

type Star struct {
	Leaf
}

type OrExpr struct {
	Disjuncts []ExpressionHandler

	Expression
}

// GroupingItem represents a grouping item, which is either an
// expression (a regular group by key) or a rollup list. Exactly one of
// Expression and Rollup will be non-nil.
type GroupingItem struct {
	Expression ExpressionHandler
	Rollup     *Rollup

	Node
}

type GroupBy struct {
	GroupingItems []*GroupingItem

	Node
}

type OrderingExpression struct {
	Expression   ExpressionHandler
	NullOrder    *NullOrder
	OrderingSpec OrderingSpec

	Node
}

type OrderBy struct {
	OrderingExpression []*OrderingExpression

	Node
}

type LimitOffset struct {
	// Limit is the LIMIT value, never nil.
	Limit ExpressionHandler
	// Offset is the optional OFFSET value, or nil.
	Offset ExpressionHandler

	Node
}

type FloatLiteral struct {
	Leaf
}

type NullLiteral struct {
	Leaf
}

type JSONLiteral struct {
	Leaf
}

type OnClause struct {
	Expression ExpressionHandler

	Node
}

type WithClauseEntry struct {
	Alias *Identifier
	Query *Query

	Node
}

// Join can introduce multiple scans and cannot have aliases. It can
// also represent a JOIN with a list of consecutive ON/USING clauses.
type Join struct {
	LHS         TableExpressionHandler
	RHS         TableExpressionHandler
	JoinType    JoinType
	OnClause    *OnClause
	UsingClause *UsingClause
	// When consecutive ON/USING clauses are encountered, they are
	// saved as ClauseList, and both OnClause and UsingClause will be
	// nil.
	ClauseList        *OnOrUsingClauseList
	ContainsCommaJoin bool
	// Indicates whether this node needs to be transformed. This is
	// true if contains if ClauseList is non nil, or if there is a JOIN
	// with ON/USING clause list on the lhs of the tree path.
	TransformationNeeded bool

	TableExpression
}

type UsingClause struct {
	Keys []*Identifier

	Node
}

type WithClause struct {
	With []*WithClauseEntry

	Node
}

type Having struct {
	Expr ExpressionHandler

	Node
}

type NamedType struct {
	Name *PathExpression

	Type
}

type ArrayType struct {
	ElementType TypeHandler

	Type
}

type StructType struct {
	StructFields      []*StructField
	TypeParameterList *TypeParameterList

	Type
}

type StructField struct {
	// Name will be nil for anonymous fields like in STRUCT<int,
	// string>.
	Name *Identifier
	Type TypeHandler

	Node
}

type CastExpression struct {
	Expr       ExpressionHandler
	Type       TypeHandler
	Format     *FormatClause
	IsSafeCast bool

	Expression
}

// SelectAs represents a SELECT with AS clause giving it an output
// type. Exactly one of SELECT AS STRUCT, SELECT AS VALUE, SELECT AS
// <TypeName> is present.
type SelectAs struct {
	TypeName *PathExpression
	AsMode   AsMode

	Node
}

type Rollup struct {
	Expressions []ExpressionHandler

	Node
}

type FunctionCall struct {
	Function  ExpressionHandler
	Arguments []ExpressionHandler
	// OrderBy is set when the function is called with FUNC(args
	// ORDER BY cols).
	OrderBy *OrderBy
	// LimitOffset is set when the function is called with FUNC(args
	// LIMIT n).
	LimitOffset *LimitOffset
	// NullHandlingModifier is set when the function is called with
	// FUNC(args {IGNORE|RESPECT} NULLS).
	NullHandlingModifier NullHandlingModifier
	// Distinct is true when the function is called with
	// FUNC(DISTINCT args).
	Distinct bool

	Expression
}

type ArrayConstructor struct {
	// Type may be nil, depending on whether the array is constructed
	// through ARRAY<type>[...] syntax or ARRAY[...] or [...].
	Type     *ArrayType
	Elements []ExpressionHandler

	Expression
}

type StructConstructorArg struct {
	Expression ExpressionHandler
	Alias      *Alias

	Node
}

// StructConstructorWithParens is resulted from structs constructed
// with (expr, expr, ...) with at least two expressions.
type StructConstructorWithParens struct {
	FieldExpressions []ExpressionHandler

	Node
}

// StructConstructorWithKeyword is resulted from structs constructed
// with STRUCT(expr [AS alias], ...) or STRUCT<...>(expr [AS alias],
// ...). Both forms support empty field lists.  The StructType is non-
// nil when the type is explicitly defined.
type StructConstructorWithKeyword struct {
	StructType *StructType
	Fields     []*StructConstructorArg

	Expression
}

// InExpression is resulted from expr IN (expr, expr, ...), expr IN
// UNNEST(...), and expr IN (query). Exactly one of InList, Query, or
// UnnestExpr is present.
type InExpression struct {
	LHS        ExpressionHandler
	InList     *InList
	Query      *Query
	UnnestExpr *UnnestExpression
	// IsNot signifies whether the IN operator as a preceding NOT to
	// it.
	IsNot bool

	Expression
}

// InList is shared with the IN operator and LIKE ANY/SOME/ALL.
type InList struct {
	// List contains the expressions present in the InList node.
	List []ExpressionHandler

	Node
}

// BetweenExpression is resulted through <LHS> BETWEEN <Low> AND
// <High>.
type BetweenExpression struct {
	LHS  ExpressionHandler
	Low  ExpressionHandler
	High ExpressionHandler
	// IsNot signifies whether the BETWEEN operator has a preceding
	// NOT to it.
	IsNot bool

	Expression
}

type NumericLiteral struct {
	Leaf
}

type BigNumericLiteral struct {
	Leaf
}

type BytesLiteral struct {
	BytesValue []byte

	Leaf
}

type DateOrTimeLiteral struct {
	StringLiteral *StringLiteral
	TypeKind      TypeKind

	Expression
}

type CaseValueExpression struct {
	Arguments []ExpressionHandler

	Expression
}

type CaseNoValueExpression struct {
	Arguments []ExpressionHandler

	Expression
}

type ArrayElement struct {
	Array    ExpressionHandler
	Position ExpressionHandler

	Expression
}

type BitwiseShiftExpression struct {
	LHS ExpressionHandler
	RHS ExpressionHandler
	// IsLeftShift signifies whether the bitwise shift is of left
	// shift type "<<" or right shift type ">>".
	IsLeftShift bool

	Expression
}

// DotGeneralizedField is a generalized form of extracting a field
// from an expression. It uses a parenthesized PathExpression instead of
// a single identifier ot select a field.
type DotGeneralizedField struct {
	Expr ExpressionHandler
	Path *PathExpression

	Expression
}

// DotIdentifier is used for using dot to extract a field from an
// arbitrary expression. Is cases where we know the left side is always
// an identifier path, we use PathExpression instead.
type DotIdentifier struct {
	Expr ExpressionHandler
	Name *Identifier

	Expression
}

type DotStar struct {
	Expr ExpressionHandler

	Expression
}

// DotStarWithModifiers is an expression constructed through SELECT
// x.* EXCEPT (...) REPLACE (...).
type DotStarWithModifiers struct {
	Expr      ExpressionHandler
	Modifiers *StarModifiers

	Expression
}

// ExpressionSubquery is a subquery in an expression. (Not in the
// FROM clause.)
type ExpressionSubquery struct {
	Query *Query
	// Modifier is the syntactic modifier on this expression
	// subquery.
	Modifier SubqueryModifier

	Expression
}

// ExtractExpression is resulted from EXTRACT(<LHS> FROM <RHS>
// <TimeZone>).
type ExtractExpression struct {
	LHS      ExpressionHandler
	RHS      ExpressionHandler
	TimeZone ExpressionHandler

	Expression
}

type IntervalExpr struct {
	IntervalValue  ExpressionHandler
	DatePartName   ExpressionHandler
	DatePartNameTo ExpressionHandler

	Expression
}

type NullOrder struct {
	NullsFirst bool

	Node
}

type OnOrUsingClauseList struct {
	// List is a list of OnClause and UsingClause elements.
	List []NodeHandler

	Node
}

type ParenthesizedJoin struct {
	Join         *Join
	SampleClause *SampleClause

	TableExpression
}

type PartitionBy struct {
	PartitioningExpressions []ExpressionHandler

	Node
}

type SetOperation struct {
	Inputs   []QueryExpressionHandler
	OpType   SetOp
	Distinct bool

	QueryExpression
}

type StarExceptList struct {
	Identifiers []*Identifier

	Node
}

// StarModifiers is resulted from SELECT * EXCEPT (...) REPLACE
// (...).
type StarModifiers struct {
	ExceptList   *StarExceptList
	ReplaceItems []*StarReplaceItem

	Node
}

type StarReplaceItem struct {
	Expression ExpressionHandler
	Alias      *Identifier

	Node
}

// StarModifiers is resulted from SELECT * EXCEPT (...) REPLACE
// (...).
type StarWithModifiers struct {
	Modifiers *StarModifiers

	Expression
}

// TableSubquery contains the table subquery, which can contain
// either a PivotClause or an UnpivotClause.
type TableSubquery struct {
	Subquery      *Query
	Alias         *Alias
	PivotClause   interface{}
	UnpivotClause interface{}
	SampleClause  *SampleClause

	TableExpression
}

type UnaryExpression struct {
	Op      UnaryOp
	Operand ExpressionHandler

	Expression
}

type UnnestExpression struct {
	Expression ExpressionHandler

	Node
}

type WindowClause struct {
	Windows []*WindowDefinition

	Node
}

type WindowDefinition struct {
	Name       *Identifier
	WindowSpec *WindowSpecification

	Node
}

type WindowFrame struct {
	StartExpr *WindowFrameExpr
	EndExpr   *WindowFrameExpr
	FrameUnit FrameUnit

	Node
}

type WindowFrameExpr struct {
	// Expression specifies the boundary as a logical or physical
	// offset to current row. It is present when BoundaryType is
	// OffsetPreceding or OffsetFollowing.
	Expression   ExpressionHandler
	BoundaryType BoundaryType

	Node
}

type LikeExpression struct {
	LHS    ExpressionHandler
	InList *InList
	IsNot  bool

	Expression
}

type WindowSpecification struct {
	BaseWindowName *Identifier
	PartitionBy    *PartitionBy
	OrderBy        *OrderBy
	WindowFrame    *WindowFrame

	Node
}

type WithOffset struct {
	Alias *Alias

	Node
}

type TypeParameterList struct {
	Parameters []LeafHandler

	Node
}

type SampleClause struct {
	SampleMethod *Identifier
	SampleSize   *SampleSize
	SampleSuffix *SampleSuffix

	Node
}

type SampleSize struct {
	Size        ExpressionHandler
	PartitionBy *PartitionBy
	Unit        SampleSizeUnit

	Node
}

type SampleSuffix struct {
	Weight *WithWeight
	Repeat *RepeatableClause

	Node
}

type WithWeight struct {
	Alias *Alias

	Node
}

type RepeatableClause struct {
	Argument ExpressionHandler

	Node
}

type Qualify struct {
	Expression ExpressionHandler

	Node
}

type FormatClause struct {
	Format       ExpressionHandler
	TimeZoneExpr ExpressionHandler

	Node
}

type ParameterExpr struct {
	Name *Identifier

	Expression
}

type AnalyticFunctionCall struct {
	Expr       ExpressionHandler
	WindowSpec *WindowSpecification

	Expression
}

func NewQueryStatement(
	query interface{},
) (*QueryStatement, error) {
	var err error

	nn := &QueryStatement{}
	nn.SetKind(QueryStatementKind)

	err = nn.InitQuery(query)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *QueryStatement) InitQuery(d interface{}) error {
	return n.initQuery(d)
}

func (n *QueryStatement) initQuery(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("QueryStatement.Query: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initQuery(t.Value)
	case NodeHandler:
		n.Query = d.(*Query)
		n.Statement.AddChild(t)
	default:
		n.Query = d.(*Query)
		n.Statement.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewQuery(
	withclause interface{},
	queryexpr interface{},
	orderby interface{},
	limitoffset interface{},
) (*Query, error) {
	var err error

	nn := &Query{}
	nn.SetKind(QueryKind)

	err = nn.InitWithClause(withclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitQueryExpr(queryexpr)
	if err != nil {
		return nil, err
	}

	err = nn.InitOrderBy(orderby)
	if err != nil {
		return nil, err
	}

	err = nn.InitLimitOffset(limitoffset)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Query) InitWithClause(d interface{}) error {
	return n.initWithClause(d)
}

func (n *Query) initWithClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.WithClause = d.(*WithClause)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWithClause(t.Value)
	default:
		n.WithClause = d.(*WithClause)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Query) InitQueryExpr(d interface{}) error {
	return n.initQueryExpr(d)
}

func (n *Query) initQueryExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Query.QueryExpr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initQueryExpr(t.Value)
	case NodeHandler:
		n.QueryExpr = d.(QueryExpressionHandler)
		n.QueryExpression.AddChild(t)
	default:
		n.QueryExpr = d.(QueryExpressionHandler)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Query) InitOrderBy(d interface{}) error {
	return n.initOrderBy(d)
}

func (n *Query) initOrderBy(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.OrderBy = d.(*OrderBy)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOrderBy(t.Value)
	default:
		n.OrderBy = d.(*OrderBy)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Query) InitLimitOffset(d interface{}) error {
	return n.initLimitOffset(d)
}

func (n *Query) initLimitOffset(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.LimitOffset = d.(*LimitOffset)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLimitOffset(t.Value)
	default:
		n.LimitOffset = d.(*LimitOffset)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Query) InitIsNested(d interface{}) error {
	return n.initIsNested(d)
}

func (n *Query) initIsNested(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsNested(t.Value)
	default:
		n.IsNested = d.(bool)
	}

	return nil
}

func NewSelect(
	distinct interface{},
	selectas interface{},
	selectlist interface{},
	fromclause interface{},
	whereclause interface{},
	groupby interface{},
	having interface{},
	qualify interface{},
	windowclause interface{},
) (*Select, error) {
	var err error

	nn := &Select{}
	nn.SetKind(SelectKind)

	err = nn.InitDistinct(distinct)
	if err != nil {
		return nil, err
	}

	err = nn.InitSelectAs(selectas)
	if err != nil {
		return nil, err
	}

	err = nn.InitSelectList(selectlist)
	if err != nil {
		return nil, err
	}

	err = nn.InitFromClause(fromclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitWhereClause(whereclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitGroupBy(groupby)
	if err != nil {
		return nil, err
	}

	err = nn.InitHaving(having)
	if err != nil {
		return nil, err
	}

	err = nn.InitQualify(qualify)
	if err != nil {
		return nil, err
	}

	err = nn.InitWindowClause(windowclause)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Select) InitDistinct(d interface{}) error {
	return n.initDistinct(d)
}

func (n *Select) initDistinct(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initDistinct(t.Value)
	default:
		n.Distinct = d.(bool)
	}

	return nil
}

func (n *Select) InitSelectAs(d interface{}) error {
	return n.initSelectAs(d)
}

func (n *Select) initSelectAs(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.SelectAs = d.(*SelectAs)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSelectAs(t.Value)
	default:
		n.SelectAs = d.(*SelectAs)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitSelectList(d interface{}) error {
	return n.initSelectList(d)
}

func (n *Select) initSelectList(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Select.SelectList: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSelectList(t.Value)
	case NodeHandler:
		n.SelectList = d.(*SelectList)
		n.QueryExpression.AddChild(t)
	default:
		n.SelectList = d.(*SelectList)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitFromClause(d interface{}) error {
	return n.initFromClause(d)
}

func (n *Select) initFromClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.FromClause = d.(*FromClause)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFromClause(t.Value)
	default:
		n.FromClause = d.(*FromClause)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitWhereClause(d interface{}) error {
	return n.initWhereClause(d)
}

func (n *Select) initWhereClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.WhereClause = d.(*WhereClause)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWhereClause(t.Value)
	default:
		n.WhereClause = d.(*WhereClause)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitGroupBy(d interface{}) error {
	return n.initGroupBy(d)
}

func (n *Select) initGroupBy(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.GroupBy = d.(*GroupBy)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initGroupBy(t.Value)
	default:
		n.GroupBy = d.(*GroupBy)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitHaving(d interface{}) error {
	return n.initHaving(d)
}

func (n *Select) initHaving(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Having = d.(*Having)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initHaving(t.Value)
	default:
		n.Having = d.(*Having)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitQualify(d interface{}) error {
	return n.initQualify(d)
}

func (n *Select) initQualify(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Qualify = d.(*Qualify)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initQualify(t.Value)
	default:
		n.Qualify = d.(*Qualify)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Select) InitWindowClause(d interface{}) error {
	return n.initWindowClause(d)
}

func (n *Select) initWindowClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.WindowClause = d.(*WindowClause)
		n.QueryExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWindowClause(t.Value)
	default:
		n.WindowClause = d.(*WindowClause)
		n.QueryExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewSelectList(
	columns interface{},
) (*SelectList, error) {
	var err error

	nn := &SelectList{}
	nn.SetKind(SelectListKind)

	err = nn.InitColumns(columns)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SelectList) InitColumns(d interface{}) error {
	if n.Columns != nil {
		return fmt.Errorf("SelectList.Columns: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Columns = make([]*SelectColumn, 0, defaultCapacity)

	return n.initColumns(d)
}

func (n *SelectList) initColumns(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initColumns(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*SelectColumn)
		n.Columns = append(n.Columns, newElem)
	}

	return nil
}

func (n *SelectList) AddChild(c NodeHandler) {
	n.Columns = append(n.Columns, c.(*SelectColumn))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *SelectList) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Columns = append(n.Columns, c.(*SelectColumn))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewSelectColumn(
	expression interface{},
	alias interface{},
) (*SelectColumn, error) {
	var err error

	nn := &SelectColumn{}
	nn.SetKind(SelectColumnKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SelectColumn) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *SelectColumn) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("SelectColumn.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SelectColumn) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *SelectColumn) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Alias = d.(*Alias)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	default:
		n.Alias = d.(*Alias)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewIntLiteral() (*IntLiteral, error) {
	var err error

	nn := &IntLiteral{}
	nn.SetKind(IntLiteralKind)

	return nn, err
}

func NewIdentifier(
	idstring interface{},
) (*Identifier, error) {
	var err error

	nn := &Identifier{}
	nn.SetKind(IdentifierKind)

	err = nn.InitIDString(idstring)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Identifier) InitIDString(d interface{}) error {
	return n.initIDString(d)
}

func (n *Identifier) initIDString(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIDString(t.Value)
	default:
		n.IDString = d.(string)
	}

	return nil
}

func NewAlias(
	identifier interface{},
) (*Alias, error) {
	var err error

	nn := &Alias{}
	nn.SetKind(AliasKind)

	err = nn.InitIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Alias) InitIdentifier(d interface{}) error {
	return n.initIdentifier(d)
}

func (n *Alias) initIdentifier(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Alias.Identifier: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIdentifier(t.Value)
	case NodeHandler:
		n.Identifier = d.(*Identifier)
		n.Node.AddChild(t)
	default:
		n.Identifier = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewPathExpression(
	names interface{},
) (*PathExpression, error) {
	var err error

	nn := &PathExpression{}
	nn.SetKind(PathExpressionKind)

	err = nn.InitNames(names)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *PathExpression) InitNames(d interface{}) error {
	if n.Names != nil {
		return fmt.Errorf("PathExpression.Names: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Names = make([]*Identifier, 0, defaultCapacity)

	return n.initNames(d)
}

func (n *PathExpression) initNames(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initNames(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*Identifier)
		n.Names = append(n.Names, newElem)
	}

	return nil
}

func (n *PathExpression) AddChild(c NodeHandler) {
	n.Names = append(n.Names, c.(*Identifier))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *PathExpression) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Names = append(n.Names, c.(*Identifier))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewTablePathExpression(
	pathexpr interface{},
	unnestexpr interface{},
	alias interface{},
	withoffset interface{},
	pivotclause interface{},
	unpivotclause interface{},
	sampleclause interface{},
) (*TablePathExpression, error) {
	var err error

	nn := &TablePathExpression{}
	nn.SetKind(TablePathExpressionKind)

	err = nn.InitPathExpr(pathexpr)
	if err != nil {
		return nil, err
	}

	err = nn.InitUnnestExpr(unnestexpr)
	if err != nil {
		return nil, err
	}

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	err = nn.InitWithOffset(withoffset)
	if err != nil {
		return nil, err
	}

	err = nn.InitPivotClause(pivotclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitUnpivotClause(unpivotclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitSampleClause(sampleclause)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *TablePathExpression) InitPathExpr(d interface{}) error {
	return n.initPathExpr(d)
}

func (n *TablePathExpression) initPathExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.PathExpr = d.(*PathExpression)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPathExpr(t.Value)
	default:
		n.PathExpr = d.(*PathExpression)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TablePathExpression) InitUnnestExpr(d interface{}) error {
	return n.initUnnestExpr(d)
}

func (n *TablePathExpression) initUnnestExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.UnnestExpr = d.(*UnnestExpression)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initUnnestExpr(t.Value)
	default:
		n.UnnestExpr = d.(*UnnestExpression)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TablePathExpression) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *TablePathExpression) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Alias = d.(*Alias)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	default:
		n.Alias = d.(*Alias)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TablePathExpression) InitWithOffset(d interface{}) error {
	return n.initWithOffset(d)
}

func (n *TablePathExpression) initWithOffset(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.WithOffset = d.(*WithOffset)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWithOffset(t.Value)
	default:
		n.WithOffset = d.(*WithOffset)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TablePathExpression) InitPivotClause(d interface{}) error {
	return n.initPivotClause(d)
}

func (n *TablePathExpression) initPivotClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.PivotClause = d
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPivotClause(t.Value)
	default:
		n.PivotClause = d
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TablePathExpression) InitUnpivotClause(d interface{}) error {
	return n.initUnpivotClause(d)
}

func (n *TablePathExpression) initUnpivotClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.UnpivotClause = d
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initUnpivotClause(t.Value)
	default:
		n.UnpivotClause = d
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TablePathExpression) InitSampleClause(d interface{}) error {
	return n.initSampleClause(d)
}

func (n *TablePathExpression) initSampleClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.SampleClause = d.(*SampleClause)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSampleClause(t.Value)
	default:
		n.SampleClause = d.(*SampleClause)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewFromClause(
	tableexpression interface{},
) (*FromClause, error) {
	var err error

	nn := &FromClause{}
	nn.SetKind(FromClauseKind)

	err = nn.InitTableExpression(tableexpression)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *FromClause) InitTableExpression(d interface{}) error {
	return n.initTableExpression(d)
}

func (n *FromClause) initTableExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("FromClause.TableExpression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTableExpression(t.Value)
	case NodeHandler:
		n.TableExpression = d.(TableExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.TableExpression = d.(TableExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWhereClause(
	expression interface{},
) (*WhereClause, error) {
	var err error

	nn := &WhereClause{}
	nn.SetKind(WhereClauseKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WhereClause) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *WhereClause) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WhereClause.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewBooleanLiteral() (*BooleanLiteral, error) {
	var err error

	nn := &BooleanLiteral{}
	nn.SetKind(BooleanLiteralKind)

	return nn, err
}

func (n *BooleanLiteral) InitValue(d interface{}) error {
	return n.initValue(d)
}

func (n *BooleanLiteral) initValue(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initValue(t.Value)
	default:
		n.Value = d.(bool)
	}

	return nil
}

func NewAndExpr(
	conjuncts interface{},
) (*AndExpr, error) {
	var err error

	nn := &AndExpr{}
	nn.SetKind(AndExprKind)

	err = nn.InitConjuncts(conjuncts)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *AndExpr) InitConjuncts(d interface{}) error {
	if n.Conjuncts != nil {
		return fmt.Errorf("AndExpr.Conjuncts: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Conjuncts = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initConjuncts(d)
}

func (n *AndExpr) initConjuncts(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initConjuncts(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Conjuncts = append(n.Conjuncts, newElem)
	}

	return nil
}

func (n *AndExpr) AddChild(c NodeHandler) {
	n.Conjuncts = append(n.Conjuncts, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *AndExpr) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Conjuncts = append(n.Conjuncts, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewBinaryExpression(
	op interface{},
	lhs interface{},
	rhs interface{},
	isnot interface{},
) (*BinaryExpression, error) {
	var err error

	nn := &BinaryExpression{}
	nn.SetKind(BinaryExpressionKind)

	err = nn.InitOp(op)
	if err != nil {
		return nil, err
	}

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitRHS(rhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitIsNot(isnot)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *BinaryExpression) InitOp(d interface{}) error {
	return n.initOp(d)
}

func (n *BinaryExpression) initOp(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOp(t.Value)
	default:
		n.Op = d.(BinaryOp)
	}

	return nil
}

func (n *BinaryExpression) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *BinaryExpression) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BinaryExpression.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BinaryExpression) InitRHS(d interface{}) error {
	return n.initRHS(d)
}

func (n *BinaryExpression) initRHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BinaryExpression.RHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initRHS(t.Value)
	case NodeHandler:
		n.RHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.RHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BinaryExpression) InitIsNot(d interface{}) error {
	return n.initIsNot(d)
}

func (n *BinaryExpression) initIsNot(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsNot(t.Value)
	default:
		n.IsNot = d.(bool)
	}

	return nil
}

func NewStringLiteral() (*StringLiteral, error) {
	var err error

	nn := &StringLiteral{}
	nn.SetKind(StringLiteralKind)

	return nn, err
}

func (n *StringLiteral) InitStringValue(d interface{}) error {
	return n.initStringValue(d)
}

func (n *StringLiteral) initStringValue(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.StringValue = d.(string)
		n.Leaf.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initStringValue(t.Value)
	default:
		n.StringValue = d.(string)
		n.Leaf.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewStar() (*Star, error) {
	var err error

	nn := &Star{}
	nn.SetKind(StarKind)

	return nn, err
}

func NewOrExpr(
	disjuncts interface{},
) (*OrExpr, error) {
	var err error

	nn := &OrExpr{}
	nn.SetKind(OrExprKind)

	err = nn.InitDisjuncts(disjuncts)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *OrExpr) InitDisjuncts(d interface{}) error {
	if n.Disjuncts != nil {
		return fmt.Errorf("OrExpr.Disjuncts: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Disjuncts = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initDisjuncts(d)
}

func (n *OrExpr) initDisjuncts(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initDisjuncts(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Disjuncts = append(n.Disjuncts, newElem)
	}

	return nil
}

func (n *OrExpr) AddChild(c NodeHandler) {
	n.Disjuncts = append(n.Disjuncts, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *OrExpr) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Disjuncts = append(n.Disjuncts, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewGroupingItem(
	expression interface{},
	rollup interface{},
) (*GroupingItem, error) {
	var err error

	nn := &GroupingItem{}
	nn.SetKind(GroupingItemKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	err = nn.InitRollup(rollup)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *GroupingItem) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *GroupingItem) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *GroupingItem) InitRollup(d interface{}) error {
	return n.initRollup(d)
}

func (n *GroupingItem) initRollup(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Rollup = d.(*Rollup)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initRollup(t.Value)
	default:
		n.Rollup = d.(*Rollup)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewGroupBy(
	groupingitems interface{},
) (*GroupBy, error) {
	var err error

	nn := &GroupBy{}
	nn.SetKind(GroupByKind)

	err = nn.InitGroupingItems(groupingitems)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *GroupBy) InitGroupingItems(d interface{}) error {
	if n.GroupingItems != nil {
		return fmt.Errorf("GroupBy.GroupingItems: %w",
			ErrFieldAlreadyInitialized)
	}

	n.GroupingItems = make([]*GroupingItem, 0, defaultCapacity)

	return n.initGroupingItems(d)
}

func (n *GroupBy) initGroupingItems(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initGroupingItems(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*GroupingItem)
		n.GroupingItems = append(n.GroupingItems, newElem)
	}

	return nil
}

func (n *GroupBy) AddChild(c NodeHandler) {
	n.GroupingItems = append(n.GroupingItems, c.(*GroupingItem))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *GroupBy) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.GroupingItems = append(n.GroupingItems, c.(*GroupingItem))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewOrderingExpression(
	expression interface{},
	nullorder interface{},
	orderingspec interface{},
) (*OrderingExpression, error) {
	var err error

	nn := &OrderingExpression{}
	nn.SetKind(OrderingExpressionKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	err = nn.InitNullOrder(nullorder)
	if err != nil {
		return nil, err
	}

	err = nn.InitOrderingSpec(orderingspec)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *OrderingExpression) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *OrderingExpression) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("OrderingExpression.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *OrderingExpression) InitNullOrder(d interface{}) error {
	return n.initNullOrder(d)
}

func (n *OrderingExpression) initNullOrder(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.NullOrder = d.(*NullOrder)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initNullOrder(t.Value)
	default:
		n.NullOrder = d.(*NullOrder)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *OrderingExpression) InitOrderingSpec(d interface{}) error {
	return n.initOrderingSpec(d)
}

func (n *OrderingExpression) initOrderingSpec(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOrderingSpec(t.Value)
	default:
		n.OrderingSpec = d.(OrderingSpec)
	}

	return nil
}

func NewOrderBy(
	orderingexpression interface{},
) (*OrderBy, error) {
	var err error

	nn := &OrderBy{}
	nn.SetKind(OrderByKind)

	err = nn.InitOrderingExpression(orderingexpression)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *OrderBy) InitOrderingExpression(d interface{}) error {
	if n.OrderingExpression != nil {
		return fmt.Errorf("OrderBy.OrderingExpression: %w",
			ErrFieldAlreadyInitialized)
	}

	n.OrderingExpression = make([]*OrderingExpression, 0, defaultCapacity)

	return n.initOrderingExpression(d)
}

func (n *OrderBy) initOrderingExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOrderingExpression(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*OrderingExpression)
		n.OrderingExpression = append(n.OrderingExpression, newElem)
	}

	return nil
}

func (n *OrderBy) AddChild(c NodeHandler) {
	n.OrderingExpression = append(n.OrderingExpression, c.(*OrderingExpression))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *OrderBy) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.OrderingExpression = append(n.OrderingExpression, c.(*OrderingExpression))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewLimitOffset(
	limit interface{},
	offset interface{},
) (*LimitOffset, error) {
	var err error

	nn := &LimitOffset{}
	nn.SetKind(LimitOffsetKind)

	err = nn.InitLimit(limit)
	if err != nil {
		return nil, err
	}

	err = nn.InitOffset(offset)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *LimitOffset) InitLimit(d interface{}) error {
	return n.initLimit(d)
}

func (n *LimitOffset) initLimit(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("LimitOffset.Limit: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLimit(t.Value)
	case NodeHandler:
		n.Limit = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Limit = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *LimitOffset) InitOffset(d interface{}) error {
	return n.initOffset(d)
}

func (n *LimitOffset) initOffset(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Offset = d.(ExpressionHandler)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOffset(t.Value)
	default:
		n.Offset = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewFloatLiteral() (*FloatLiteral, error) {
	var err error

	nn := &FloatLiteral{}
	nn.SetKind(FloatLiteralKind)

	return nn, err
}

func NewNullLiteral() (*NullLiteral, error) {
	var err error

	nn := &NullLiteral{}
	nn.SetKind(NullLiteralKind)

	return nn, err
}

func NewJSONLiteral() (*JSONLiteral, error) {
	var err error

	nn := &JSONLiteral{}
	nn.SetKind(JSONLiteralKind)

	return nn, err
}

func NewOnClause(
	expression interface{},
) (*OnClause, error) {
	var err error

	nn := &OnClause{}
	nn.SetKind(OnClauseKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *OnClause) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *OnClause) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("OnClause.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWithClauseEntry(
	alias interface{},
	query interface{},
) (*WithClauseEntry, error) {
	var err error

	nn := &WithClauseEntry{}
	nn.SetKind(WithClauseEntryKind)

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	err = nn.InitQuery(query)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WithClauseEntry) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *WithClauseEntry) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WithClauseEntry.Alias: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	case NodeHandler:
		n.Alias = d.(*Identifier)
		n.Node.AddChild(t)
	default:
		n.Alias = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WithClauseEntry) InitQuery(d interface{}) error {
	return n.initQuery(d)
}

func (n *WithClauseEntry) initQuery(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WithClauseEntry.Query: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initQuery(t.Value)
	case NodeHandler:
		n.Query = d.(*Query)
		n.Node.AddChild(t)
	default:
		n.Query = d.(*Query)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewJoin(
	lhs interface{},
	rhs interface{},
	jointype interface{},
) (*Join, error) {
	var err error

	nn := &Join{}
	nn.SetKind(JoinKind)

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitRHS(rhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitJoinType(jointype)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Join) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *Join) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Join.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(TableExpressionHandler)
		n.TableExpression.AddChild(t)
	default:
		n.LHS = d.(TableExpressionHandler)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Join) InitRHS(d interface{}) error {
	return n.initRHS(d)
}

func (n *Join) initRHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Join.RHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initRHS(t.Value)
	case NodeHandler:
		n.RHS = d.(TableExpressionHandler)
		n.TableExpression.AddChild(t)
	default:
		n.RHS = d.(TableExpressionHandler)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Join) InitJoinType(d interface{}) error {
	return n.initJoinType(d)
}

func (n *Join) initJoinType(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initJoinType(t.Value)
	default:
		n.JoinType = d.(JoinType)
	}

	return nil
}

func (n *Join) InitOnClause(d interface{}) error {
	return n.initOnClause(d)
}

func (n *Join) initOnClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.OnClause = d.(*OnClause)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOnClause(t.Value)
	default:
		n.OnClause = d.(*OnClause)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Join) InitUsingClause(d interface{}) error {
	return n.initUsingClause(d)
}

func (n *Join) initUsingClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.UsingClause = d.(*UsingClause)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initUsingClause(t.Value)
	default:
		n.UsingClause = d.(*UsingClause)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Join) InitClauseList(d interface{}) error {
	return n.initClauseList(d)
}

func (n *Join) initClauseList(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.ClauseList = d.(*OnOrUsingClauseList)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initClauseList(t.Value)
	default:
		n.ClauseList = d.(*OnOrUsingClauseList)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Join) InitContainsCommaJoin(d interface{}) error {
	return n.initContainsCommaJoin(d)
}

func (n *Join) initContainsCommaJoin(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.ContainsCommaJoin = d.(bool)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initContainsCommaJoin(t.Value)
	default:
		n.ContainsCommaJoin = d.(bool)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *Join) InitTransformationNeeded(d interface{}) error {
	return n.initTransformationNeeded(d)
}

func (n *Join) initTransformationNeeded(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.TransformationNeeded = d.(bool)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTransformationNeeded(t.Value)
	default:
		n.TransformationNeeded = d.(bool)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewUsingClause(
	keys interface{},
) (*UsingClause, error) {
	var err error

	nn := &UsingClause{}
	nn.SetKind(UsingClauseKind)

	err = nn.InitKeys(keys)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *UsingClause) InitKeys(d interface{}) error {
	if n.Keys != nil {
		return fmt.Errorf("UsingClause.Keys: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Keys = make([]*Identifier, 0, defaultCapacity)

	return n.initKeys(d)
}

func (n *UsingClause) initKeys(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initKeys(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*Identifier)
		n.Keys = append(n.Keys, newElem)
	}

	return nil
}

func (n *UsingClause) AddChild(c NodeHandler) {
	n.Keys = append(n.Keys, c.(*Identifier))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *UsingClause) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Keys = append(n.Keys, c.(*Identifier))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewWithClause(
	with interface{},
) (*WithClause, error) {
	var err error

	nn := &WithClause{}
	nn.SetKind(WithClauseKind)

	err = nn.InitWith(with)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WithClause) InitWith(d interface{}) error {
	if n.With != nil {
		return fmt.Errorf("WithClause.With: %w",
			ErrFieldAlreadyInitialized)
	}

	n.With = make([]*WithClauseEntry, 0, defaultCapacity)

	return n.initWith(d)
}

func (n *WithClause) initWith(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWith(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*WithClauseEntry)
		n.With = append(n.With, newElem)
	}

	return nil
}

func (n *WithClause) AddChild(c NodeHandler) {
	n.With = append(n.With, c.(*WithClauseEntry))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *WithClause) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.With = append(n.With, c.(*WithClauseEntry))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewHaving(
	expr interface{},
) (*Having, error) {
	var err error

	nn := &Having{}
	nn.SetKind(HavingKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Having) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *Having) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Having.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewNamedType(
	name interface{},
) (*NamedType, error) {
	var err error

	nn := &NamedType{}
	nn.SetKind(NamedTypeKind)

	err = nn.InitName(name)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *NamedType) InitName(d interface{}) error {
	return n.initName(d)
}

func (n *NamedType) initName(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("NamedType.Name: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initName(t.Value)
	case NodeHandler:
		n.Name = d.(*PathExpression)
		n.Type.AddChild(t)
	default:
		n.Name = d.(*PathExpression)
		n.Type.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewArrayType(
	elementtype interface{},
) (*ArrayType, error) {
	var err error

	nn := &ArrayType{}
	nn.SetKind(ArrayTypeKind)

	err = nn.InitElementType(elementtype)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ArrayType) InitElementType(d interface{}) error {
	return n.initElementType(d)
}

func (n *ArrayType) initElementType(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ArrayType.ElementType: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initElementType(t.Value)
	case NodeHandler:
		n.ElementType = d.(TypeHandler)
		n.Type.AddChild(t)
	default:
		n.ElementType = d.(TypeHandler)
		n.Type.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewStructType(
	structfields interface{},
) (*StructType, error) {
	var err error

	nn := &StructType{}
	nn.SetKind(StructTypeKind)

	err = nn.InitStructFields(structfields)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StructType) InitStructFields(d interface{}) error {
	if n.StructFields != nil {
		return fmt.Errorf("StructType.StructFields: %w",
			ErrFieldAlreadyInitialized)
	}

	n.StructFields = make([]*StructField, 0, defaultCapacity)

	return n.initStructFields(d)
}

func (n *StructType) initStructFields(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initStructFields(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*StructField)
		n.StructFields = append(n.StructFields, newElem)
	}

	return nil
}

func (n *StructType) InitTypeParameterList(d interface{}) error {
	return n.initTypeParameterList(d)
}

func (n *StructType) initTypeParameterList(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.TypeParameterList = d.(*TypeParameterList)
		n.Type.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTypeParameterList(t.Value)
	default:
		n.TypeParameterList = d.(*TypeParameterList)
		n.Type.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *StructType) AddChild(c NodeHandler) {
	n.StructFields = append(n.StructFields, c.(*StructField))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *StructType) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.StructFields = append(n.StructFields, c.(*StructField))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewStructField(
	name interface{},
	typ interface{},
) (*StructField, error) {
	var err error

	nn := &StructField{}
	nn.SetKind(StructFieldKind)

	err = nn.InitName(name)
	if err != nil {
		return nil, err
	}

	err = nn.InitType(typ)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StructField) InitName(d interface{}) error {
	return n.initName(d)
}

func (n *StructField) initName(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Name = d.(*Identifier)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initName(t.Value)
	default:
		n.Name = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *StructField) InitType(d interface{}) error {
	return n.initType(d)
}

func (n *StructField) initType(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("StructField.Type: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initType(t.Value)
	case NodeHandler:
		n.Type = d.(TypeHandler)
		n.Node.AddChild(t)
	default:
		n.Type = d.(TypeHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewCastExpression(
	expr interface{},
	typ interface{},
	format interface{},
	issafecast interface{},
) (*CastExpression, error) {
	var err error

	nn := &CastExpression{}
	nn.SetKind(CastExpressionKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	err = nn.InitType(typ)
	if err != nil {
		return nil, err
	}

	err = nn.InitFormat(format)
	if err != nil {
		return nil, err
	}

	err = nn.InitIsSafeCast(issafecast)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *CastExpression) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *CastExpression) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("CastExpression.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *CastExpression) InitType(d interface{}) error {
	return n.initType(d)
}

func (n *CastExpression) initType(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("CastExpression.Type: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initType(t.Value)
	case NodeHandler:
		n.Type = d.(TypeHandler)
		n.Expression.AddChild(t)
	default:
		n.Type = d.(TypeHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *CastExpression) InitFormat(d interface{}) error {
	return n.initFormat(d)
}

func (n *CastExpression) initFormat(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Format = d.(*FormatClause)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFormat(t.Value)
	default:
		n.Format = d.(*FormatClause)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *CastExpression) InitIsSafeCast(d interface{}) error {
	return n.initIsSafeCast(d)
}

func (n *CastExpression) initIsSafeCast(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsSafeCast(t.Value)
	default:
		n.IsSafeCast = d.(bool)
	}

	return nil
}

func NewSelectAs(
	typename interface{},
	asmode interface{},
) (*SelectAs, error) {
	var err error

	nn := &SelectAs{}
	nn.SetKind(SelectAsKind)

	err = nn.InitTypeName(typename)
	if err != nil {
		return nil, err
	}

	err = nn.InitAsMode(asmode)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SelectAs) InitTypeName(d interface{}) error {
	return n.initTypeName(d)
}

func (n *SelectAs) initTypeName(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.TypeName = d.(*PathExpression)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTypeName(t.Value)
	default:
		n.TypeName = d.(*PathExpression)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SelectAs) InitAsMode(d interface{}) error {
	return n.initAsMode(d)
}

func (n *SelectAs) initAsMode(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAsMode(t.Value)
	default:
		n.AsMode = d.(AsMode)
	}

	return nil
}

func NewRollup(
	expressions interface{},
) (*Rollup, error) {
	var err error

	nn := &Rollup{}
	nn.SetKind(RollupKind)

	err = nn.InitExpressions(expressions)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Rollup) InitExpressions(d interface{}) error {
	if n.Expressions != nil {
		return fmt.Errorf("Rollup.Expressions: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Expressions = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initExpressions(d)
}

func (n *Rollup) initExpressions(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpressions(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Expressions = append(n.Expressions, newElem)
	}

	return nil
}

func (n *Rollup) AddChild(c NodeHandler) {
	n.Expressions = append(n.Expressions, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *Rollup) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Expressions = append(n.Expressions, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewFunctionCall(
	function interface{},
	distinct interface{},
) (*FunctionCall, error) {
	var err error

	nn := &FunctionCall{}
	nn.SetKind(FunctionCallKind)

	err = nn.InitFunction(function)
	if err != nil {
		return nil, err
	}

	err = nn.InitDistinct(distinct)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *FunctionCall) InitFunction(d interface{}) error {
	return n.initFunction(d)
}

func (n *FunctionCall) initFunction(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("FunctionCall.Function: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFunction(t.Value)
	case NodeHandler:
		n.Function = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Function = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *FunctionCall) InitArguments(d interface{}) error {
	if n.Arguments != nil {
		return fmt.Errorf("FunctionCall.Arguments: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Arguments = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initArguments(d)
}

func (n *FunctionCall) initArguments(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initArguments(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Arguments = append(n.Arguments, newElem)
	}

	return nil
}

func (n *FunctionCall) InitOrderBy(d interface{}) error {
	return n.initOrderBy(d)
}

func (n *FunctionCall) initOrderBy(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.OrderBy = d.(*OrderBy)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOrderBy(t.Value)
	default:
		n.OrderBy = d.(*OrderBy)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *FunctionCall) InitLimitOffset(d interface{}) error {
	return n.initLimitOffset(d)
}

func (n *FunctionCall) initLimitOffset(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.LimitOffset = d.(*LimitOffset)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLimitOffset(t.Value)
	default:
		n.LimitOffset = d.(*LimitOffset)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *FunctionCall) InitNullHandlingModifier(d interface{}) error {
	return n.initNullHandlingModifier(d)
}

func (n *FunctionCall) initNullHandlingModifier(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initNullHandlingModifier(t.Value)
	default:
		n.NullHandlingModifier = d.(NullHandlingModifier)
	}

	return nil
}

func (n *FunctionCall) InitDistinct(d interface{}) error {
	return n.initDistinct(d)
}

func (n *FunctionCall) initDistinct(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initDistinct(t.Value)
	default:
		n.Distinct = d.(bool)
	}

	return nil
}

func (n *FunctionCall) AddChild(c NodeHandler) {
	n.Arguments = append(n.Arguments, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *FunctionCall) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Arguments = append(n.Arguments, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewArrayConstructor(
	typ interface{},
) (*ArrayConstructor, error) {
	var err error

	nn := &ArrayConstructor{}
	nn.SetKind(ArrayConstructorKind)

	err = nn.InitType(typ)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ArrayConstructor) InitType(d interface{}) error {
	return n.initType(d)
}

func (n *ArrayConstructor) initType(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Type = d.(*ArrayType)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initType(t.Value)
	default:
		n.Type = d.(*ArrayType)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *ArrayConstructor) InitElements(d interface{}) error {
	if n.Elements != nil {
		return fmt.Errorf("ArrayConstructor.Elements: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Elements = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initElements(d)
}

func (n *ArrayConstructor) initElements(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initElements(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Elements = append(n.Elements, newElem)
	}

	return nil
}

func (n *ArrayConstructor) AddChild(c NodeHandler) {
	n.Elements = append(n.Elements, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *ArrayConstructor) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Elements = append(n.Elements, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewStructConstructorArg(
	expression interface{},
	alias interface{},
) (*StructConstructorArg, error) {
	var err error

	nn := &StructConstructorArg{}
	nn.SetKind(StructConstructorArgKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StructConstructorArg) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *StructConstructorArg) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("StructConstructorArg.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *StructConstructorArg) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *StructConstructorArg) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("StructConstructorArg.Alias: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	case NodeHandler:
		n.Alias = d.(*Alias)
		n.Node.AddChild(t)
	default:
		n.Alias = d.(*Alias)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewStructConstructorWithParens(
	fieldexpressions interface{},
) (*StructConstructorWithParens, error) {
	var err error

	nn := &StructConstructorWithParens{}
	nn.SetKind(StructConstructorWithParensKind)

	err = nn.InitFieldExpressions(fieldexpressions)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StructConstructorWithParens) InitFieldExpressions(d interface{}) error {
	if n.FieldExpressions != nil {
		return fmt.Errorf("StructConstructorWithParens.FieldExpressions: %w",
			ErrFieldAlreadyInitialized)
	}

	n.FieldExpressions = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initFieldExpressions(d)
}

func (n *StructConstructorWithParens) initFieldExpressions(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFieldExpressions(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.FieldExpressions = append(n.FieldExpressions, newElem)
	}

	return nil
}

func (n *StructConstructorWithParens) AddChild(c NodeHandler) {
	n.FieldExpressions = append(n.FieldExpressions, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *StructConstructorWithParens) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.FieldExpressions = append(n.FieldExpressions, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewStructConstructorWithKeyword(
	structtype interface{},
) (*StructConstructorWithKeyword, error) {
	var err error

	nn := &StructConstructorWithKeyword{}
	nn.SetKind(StructConstructorWithKeywordKind)

	err = nn.InitStructType(structtype)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StructConstructorWithKeyword) InitStructType(d interface{}) error {
	return n.initStructType(d)
}

func (n *StructConstructorWithKeyword) initStructType(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.StructType = d.(*StructType)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initStructType(t.Value)
	default:
		n.StructType = d.(*StructType)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *StructConstructorWithKeyword) InitFields(d interface{}) error {
	if n.Fields != nil {
		return fmt.Errorf("StructConstructorWithKeyword.Fields: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Fields = make([]*StructConstructorArg, 0, defaultCapacity)

	return n.initFields(d)
}

func (n *StructConstructorWithKeyword) initFields(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFields(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*StructConstructorArg)
		n.Fields = append(n.Fields, newElem)
	}

	return nil
}

func (n *StructConstructorWithKeyword) AddChild(c NodeHandler) {
	n.Fields = append(n.Fields, c.(*StructConstructorArg))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *StructConstructorWithKeyword) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Fields = append(n.Fields, c.(*StructConstructorArg))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewInExpression(
	lhs interface{},
	inlist interface{},
	query interface{},
	unnestexpr interface{},
	isnot interface{},
) (*InExpression, error) {
	var err error

	nn := &InExpression{}
	nn.SetKind(InExpressionKind)

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitInList(inlist)
	if err != nil {
		return nil, err
	}

	err = nn.InitQuery(query)
	if err != nil {
		return nil, err
	}

	err = nn.InitUnnestExpr(unnestexpr)
	if err != nil {
		return nil, err
	}

	err = nn.InitIsNot(isnot)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *InExpression) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *InExpression) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("InExpression.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *InExpression) InitInList(d interface{}) error {
	return n.initInList(d)
}

func (n *InExpression) initInList(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.InList = d.(*InList)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initInList(t.Value)
	default:
		n.InList = d.(*InList)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *InExpression) InitQuery(d interface{}) error {
	return n.initQuery(d)
}

func (n *InExpression) initQuery(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Query = d.(*Query)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initQuery(t.Value)
	default:
		n.Query = d.(*Query)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *InExpression) InitUnnestExpr(d interface{}) error {
	return n.initUnnestExpr(d)
}

func (n *InExpression) initUnnestExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.UnnestExpr = d.(*UnnestExpression)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initUnnestExpr(t.Value)
	default:
		n.UnnestExpr = d.(*UnnestExpression)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *InExpression) InitIsNot(d interface{}) error {
	return n.initIsNot(d)
}

func (n *InExpression) initIsNot(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsNot(t.Value)
	default:
		n.IsNot = d.(bool)
	}

	return nil
}

func NewInList(
	list interface{},
) (*InList, error) {
	var err error

	nn := &InList{}
	nn.SetKind(InListKind)

	err = nn.InitList(list)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *InList) InitList(d interface{}) error {
	if n.List != nil {
		return fmt.Errorf("InList.List: %w",
			ErrFieldAlreadyInitialized)
	}

	n.List = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initList(d)
}

func (n *InList) initList(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initList(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.List = append(n.List, newElem)
	}

	return nil
}

func (n *InList) AddChild(c NodeHandler) {
	n.List = append(n.List, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *InList) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.List = append(n.List, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewBetweenExpression(
	lhs interface{},
	low interface{},
	high interface{},
	isnot interface{},
) (*BetweenExpression, error) {
	var err error

	nn := &BetweenExpression{}
	nn.SetKind(BetweenExpressionKind)

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitLow(low)
	if err != nil {
		return nil, err
	}

	err = nn.InitHigh(high)
	if err != nil {
		return nil, err
	}

	err = nn.InitIsNot(isnot)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *BetweenExpression) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *BetweenExpression) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BetweenExpression.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BetweenExpression) InitLow(d interface{}) error {
	return n.initLow(d)
}

func (n *BetweenExpression) initLow(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BetweenExpression.Low: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLow(t.Value)
	case NodeHandler:
		n.Low = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Low = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BetweenExpression) InitHigh(d interface{}) error {
	return n.initHigh(d)
}

func (n *BetweenExpression) initHigh(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BetweenExpression.High: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initHigh(t.Value)
	case NodeHandler:
		n.High = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.High = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BetweenExpression) InitIsNot(d interface{}) error {
	return n.initIsNot(d)
}

func (n *BetweenExpression) initIsNot(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsNot(t.Value)
	default:
		n.IsNot = d.(bool)
	}

	return nil
}

func NewNumericLiteral() (*NumericLiteral, error) {
	var err error

	nn := &NumericLiteral{}
	nn.SetKind(NumericLiteralKind)

	return nn, err
}

func NewBigNumericLiteral() (*BigNumericLiteral, error) {
	var err error

	nn := &BigNumericLiteral{}
	nn.SetKind(BigNumericLiteralKind)

	return nn, err
}

func NewBytesLiteral() (*BytesLiteral, error) {
	var err error

	nn := &BytesLiteral{}
	nn.SetKind(BytesLiteralKind)

	return nn, err
}

func (n *BytesLiteral) InitBytesValue(d interface{}) error {
	return n.initBytesValue(d)
}

func (n *BytesLiteral) initBytesValue(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.BytesValue = d.([]byte)
		n.Leaf.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initBytesValue(t.Value)
	default:
		n.BytesValue = d.([]byte)
		n.Leaf.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewDateOrTimeLiteral(
	stringliteral interface{},
	typekind interface{},
) (*DateOrTimeLiteral, error) {
	var err error

	nn := &DateOrTimeLiteral{}
	nn.SetKind(DateOrTimeLiteralKind)

	err = nn.InitStringLiteral(stringliteral)
	if err != nil {
		return nil, err
	}

	err = nn.InitTypeKind(typekind)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *DateOrTimeLiteral) InitStringLiteral(d interface{}) error {
	return n.initStringLiteral(d)
}

func (n *DateOrTimeLiteral) initStringLiteral(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DateOrTimeLiteral.StringLiteral: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initStringLiteral(t.Value)
	case NodeHandler:
		n.StringLiteral = d.(*StringLiteral)
		n.Expression.AddChild(t)
	default:
		n.StringLiteral = d.(*StringLiteral)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *DateOrTimeLiteral) InitTypeKind(d interface{}) error {
	return n.initTypeKind(d)
}

func (n *DateOrTimeLiteral) initTypeKind(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.TypeKind = d.(TypeKind)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTypeKind(t.Value)
	default:
		n.TypeKind = d.(TypeKind)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewCaseValueExpression(
	arguments interface{},
) (*CaseValueExpression, error) {
	var err error

	nn := &CaseValueExpression{}
	nn.SetKind(CaseValueExpressionKind)

	err = nn.InitArguments(arguments)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *CaseValueExpression) InitArguments(d interface{}) error {
	if n.Arguments != nil {
		return fmt.Errorf("CaseValueExpression.Arguments: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Arguments = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initArguments(d)
}

func (n *CaseValueExpression) initArguments(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initArguments(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Arguments = append(n.Arguments, newElem)
	}

	return nil
}

func (n *CaseValueExpression) AddChild(c NodeHandler) {
	n.Arguments = append(n.Arguments, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *CaseValueExpression) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Arguments = append(n.Arguments, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewCaseNoValueExpression(
	arguments interface{},
) (*CaseNoValueExpression, error) {
	var err error

	nn := &CaseNoValueExpression{}
	nn.SetKind(CaseNoValueExpressionKind)

	err = nn.InitArguments(arguments)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *CaseNoValueExpression) InitArguments(d interface{}) error {
	if n.Arguments != nil {
		return fmt.Errorf("CaseNoValueExpression.Arguments: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Arguments = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initArguments(d)
}

func (n *CaseNoValueExpression) initArguments(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initArguments(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.Arguments = append(n.Arguments, newElem)
	}

	return nil
}

func (n *CaseNoValueExpression) AddChild(c NodeHandler) {
	n.Arguments = append(n.Arguments, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *CaseNoValueExpression) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Arguments = append(n.Arguments, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewArrayElement(
	array interface{},
	position interface{},
) (*ArrayElement, error) {
	var err error

	nn := &ArrayElement{}
	nn.SetKind(ArrayElementKind)

	err = nn.InitArray(array)
	if err != nil {
		return nil, err
	}

	err = nn.InitPosition(position)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ArrayElement) InitArray(d interface{}) error {
	return n.initArray(d)
}

func (n *ArrayElement) initArray(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ArrayElement.Array: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initArray(t.Value)
	case NodeHandler:
		n.Array = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Array = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *ArrayElement) InitPosition(d interface{}) error {
	return n.initPosition(d)
}

func (n *ArrayElement) initPosition(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ArrayElement.Position: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPosition(t.Value)
	case NodeHandler:
		n.Position = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Position = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewBitwiseShiftExpression(
	lhs interface{},
	rhs interface{},
	isleftshift interface{},
) (*BitwiseShiftExpression, error) {
	var err error

	nn := &BitwiseShiftExpression{}
	nn.SetKind(BitwiseShiftExpressionKind)

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitRHS(rhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitIsLeftShift(isleftshift)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *BitwiseShiftExpression) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *BitwiseShiftExpression) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BitwiseShiftExpression.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BitwiseShiftExpression) InitRHS(d interface{}) error {
	return n.initRHS(d)
}

func (n *BitwiseShiftExpression) initRHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("BitwiseShiftExpression.RHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initRHS(t.Value)
	case NodeHandler:
		n.RHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.RHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *BitwiseShiftExpression) InitIsLeftShift(d interface{}) error {
	return n.initIsLeftShift(d)
}

func (n *BitwiseShiftExpression) initIsLeftShift(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.IsLeftShift = d.(bool)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsLeftShift(t.Value)
	default:
		n.IsLeftShift = d.(bool)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewDotGeneralizedField(
	expr interface{},
	path interface{},
) (*DotGeneralizedField, error) {
	var err error

	nn := &DotGeneralizedField{}
	nn.SetKind(DotGeneralizedFieldKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	err = nn.InitPath(path)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *DotGeneralizedField) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *DotGeneralizedField) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotGeneralizedField.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *DotGeneralizedField) InitPath(d interface{}) error {
	return n.initPath(d)
}

func (n *DotGeneralizedField) initPath(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotGeneralizedField.Path: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPath(t.Value)
	case NodeHandler:
		n.Path = d.(*PathExpression)
		n.Expression.AddChild(t)
	default:
		n.Path = d.(*PathExpression)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewDotIdentifier(
	expr interface{},
	name interface{},
) (*DotIdentifier, error) {
	var err error

	nn := &DotIdentifier{}
	nn.SetKind(DotIdentifierKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	err = nn.InitName(name)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *DotIdentifier) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *DotIdentifier) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotIdentifier.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *DotIdentifier) InitName(d interface{}) error {
	return n.initName(d)
}

func (n *DotIdentifier) initName(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotIdentifier.Name: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initName(t.Value)
	case NodeHandler:
		n.Name = d.(*Identifier)
		n.Expression.AddChild(t)
	default:
		n.Name = d.(*Identifier)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewDotStar(
	expr interface{},
) (*DotStar, error) {
	var err error

	nn := &DotStar{}
	nn.SetKind(DotStarKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *DotStar) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *DotStar) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotStar.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewDotStarWithModifiers(
	expr interface{},
	modifiers interface{},
) (*DotStarWithModifiers, error) {
	var err error

	nn := &DotStarWithModifiers{}
	nn.SetKind(DotStarWithModifiersKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	err = nn.InitModifiers(modifiers)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *DotStarWithModifiers) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *DotStarWithModifiers) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotStarWithModifiers.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *DotStarWithModifiers) InitModifiers(d interface{}) error {
	return n.initModifiers(d)
}

func (n *DotStarWithModifiers) initModifiers(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("DotStarWithModifiers.Modifiers: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initModifiers(t.Value)
	case NodeHandler:
		n.Modifiers = d.(*StarModifiers)
		n.Expression.AddChild(t)
	default:
		n.Modifiers = d.(*StarModifiers)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewExpressionSubquery(
	query interface{},
	modifier interface{},
) (*ExpressionSubquery, error) {
	var err error

	nn := &ExpressionSubquery{}
	nn.SetKind(ExpressionSubqueryKind)

	err = nn.InitQuery(query)
	if err != nil {
		return nil, err
	}

	err = nn.InitModifier(modifier)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ExpressionSubquery) InitQuery(d interface{}) error {
	return n.initQuery(d)
}

func (n *ExpressionSubquery) initQuery(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ExpressionSubquery.Query: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initQuery(t.Value)
	case NodeHandler:
		n.Query = d.(*Query)
		n.Expression.AddChild(t)
	default:
		n.Query = d.(*Query)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *ExpressionSubquery) InitModifier(d interface{}) error {
	return n.initModifier(d)
}

func (n *ExpressionSubquery) initModifier(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ExpressionSubquery.Modifier: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initModifier(t.Value)
	case NodeHandler:
		n.Modifier = d.(SubqueryModifier)
		n.Expression.AddChild(t)
	default:
		n.Modifier = d.(SubqueryModifier)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewExtractExpression(
	lhs interface{},
	rhs interface{},
	timezone interface{},
) (*ExtractExpression, error) {
	var err error

	nn := &ExtractExpression{}
	nn.SetKind(ExtractExpressionKind)

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitRHS(rhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitTimeZone(timezone)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ExtractExpression) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *ExtractExpression) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ExtractExpression.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *ExtractExpression) InitRHS(d interface{}) error {
	return n.initRHS(d)
}

func (n *ExtractExpression) initRHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ExtractExpression.RHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initRHS(t.Value)
	case NodeHandler:
		n.RHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.RHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *ExtractExpression) InitTimeZone(d interface{}) error {
	return n.initTimeZone(d)
}

func (n *ExtractExpression) initTimeZone(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.TimeZone = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTimeZone(t.Value)
	default:
		n.TimeZone = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewIntervalExpr(
	intervalvalue interface{},
	datepartname interface{},
	datepartnameto interface{},
) (*IntervalExpr, error) {
	var err error

	nn := &IntervalExpr{}
	nn.SetKind(IntervalExprKind)

	err = nn.InitIntervalValue(intervalvalue)
	if err != nil {
		return nil, err
	}

	err = nn.InitDatePartName(datepartname)
	if err != nil {
		return nil, err
	}

	err = nn.InitDatePartNameTo(datepartnameto)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *IntervalExpr) InitIntervalValue(d interface{}) error {
	return n.initIntervalValue(d)
}

func (n *IntervalExpr) initIntervalValue(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("IntervalExpr.IntervalValue: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIntervalValue(t.Value)
	case NodeHandler:
		n.IntervalValue = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.IntervalValue = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *IntervalExpr) InitDatePartName(d interface{}) error {
	return n.initDatePartName(d)
}

func (n *IntervalExpr) initDatePartName(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("IntervalExpr.DatePartName: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initDatePartName(t.Value)
	case NodeHandler:
		n.DatePartName = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.DatePartName = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *IntervalExpr) InitDatePartNameTo(d interface{}) error {
	return n.initDatePartNameTo(d)
}

func (n *IntervalExpr) initDatePartNameTo(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.DatePartNameTo = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initDatePartNameTo(t.Value)
	default:
		n.DatePartNameTo = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewNullOrder(
	nullsfirst interface{},
) (*NullOrder, error) {
	var err error

	nn := &NullOrder{}
	nn.SetKind(NullOrderKind)

	err = nn.InitNullsFirst(nullsfirst)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *NullOrder) InitNullsFirst(d interface{}) error {
	return n.initNullsFirst(d)
}

func (n *NullOrder) initNullsFirst(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initNullsFirst(t.Value)
	default:
		n.NullsFirst = d.(bool)
	}

	return nil
}

func NewOnOrUsingClauseList(
	list interface{},
) (*OnOrUsingClauseList, error) {
	var err error

	nn := &OnOrUsingClauseList{}
	nn.SetKind(OnOrUsingClauseListKind)

	err = nn.InitList(list)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *OnOrUsingClauseList) InitList(d interface{}) error {
	if n.List != nil {
		return fmt.Errorf("OnOrUsingClauseList.List: %w",
			ErrFieldAlreadyInitialized)
	}

	n.List = make([]NodeHandler, 0, defaultCapacity)

	return n.initList(d)
}

func (n *OnOrUsingClauseList) initList(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initList(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(NodeHandler)
		n.List = append(n.List, newElem)
	}

	return nil
}

func (n *OnOrUsingClauseList) AddChild(c NodeHandler) {
	n.List = append(n.List, c)

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *OnOrUsingClauseList) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.List = append(n.List, c)

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewParenthesizedJoin(
	join interface{},
	sampleclause interface{},
) (*ParenthesizedJoin, error) {
	var err error

	nn := &ParenthesizedJoin{}
	nn.SetKind(ParenthesizedJoinKind)

	err = nn.InitJoin(join)
	if err != nil {
		return nil, err
	}

	err = nn.InitSampleClause(sampleclause)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ParenthesizedJoin) InitJoin(d interface{}) error {
	return n.initJoin(d)
}

func (n *ParenthesizedJoin) initJoin(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("ParenthesizedJoin.Join: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initJoin(t.Value)
	case NodeHandler:
		n.Join = d.(*Join)
		n.TableExpression.AddChild(t)
	default:
		n.Join = d.(*Join)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *ParenthesizedJoin) InitSampleClause(d interface{}) error {
	return n.initSampleClause(d)
}

func (n *ParenthesizedJoin) initSampleClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.SampleClause = d.(*SampleClause)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSampleClause(t.Value)
	default:
		n.SampleClause = d.(*SampleClause)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewPartitionBy(
	partitioningexpressions interface{},
) (*PartitionBy, error) {
	var err error

	nn := &PartitionBy{}
	nn.SetKind(PartitionByKind)

	err = nn.InitPartitioningExpressions(partitioningexpressions)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *PartitionBy) InitPartitioningExpressions(d interface{}) error {
	if n.PartitioningExpressions != nil {
		return fmt.Errorf("PartitionBy.PartitioningExpressions: %w",
			ErrFieldAlreadyInitialized)
	}

	n.PartitioningExpressions = make([]ExpressionHandler, 0, defaultCapacity)

	return n.initPartitioningExpressions(d)
}

func (n *PartitionBy) initPartitioningExpressions(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPartitioningExpressions(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(ExpressionHandler)
		n.PartitioningExpressions = append(n.PartitioningExpressions, newElem)
	}

	return nil
}

func (n *PartitionBy) AddChild(c NodeHandler) {
	n.PartitioningExpressions = append(n.PartitioningExpressions, c.(ExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *PartitionBy) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.PartitioningExpressions = append(n.PartitioningExpressions, c.(ExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewSetOperation(
	inputs interface{},
	optype interface{},
	distinct interface{},
) (*SetOperation, error) {
	var err error

	nn := &SetOperation{}
	nn.SetKind(SetOperationKind)

	err = nn.InitInputs(inputs)
	if err != nil {
		return nil, err
	}

	err = nn.InitOpType(optype)
	if err != nil {
		return nil, err
	}

	err = nn.InitDistinct(distinct)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SetOperation) InitInputs(d interface{}) error {
	if n.Inputs != nil {
		return fmt.Errorf("SetOperation.Inputs: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Inputs = make([]QueryExpressionHandler, 0, defaultCapacity)

	return n.initInputs(d)
}

func (n *SetOperation) initInputs(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initInputs(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(QueryExpressionHandler)
		n.Inputs = append(n.Inputs, newElem)
	}

	return nil
}

func (n *SetOperation) InitOpType(d interface{}) error {
	return n.initOpType(d)
}

func (n *SetOperation) initOpType(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOpType(t.Value)
	default:
		n.OpType = d.(SetOp)
	}

	return nil
}

func (n *SetOperation) InitDistinct(d interface{}) error {
	return n.initDistinct(d)
}

func (n *SetOperation) initDistinct(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initDistinct(t.Value)
	default:
		n.Distinct = d.(bool)
	}

	return nil
}

func (n *SetOperation) AddChild(c NodeHandler) {
	n.Inputs = append(n.Inputs, c.(QueryExpressionHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *SetOperation) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Inputs = append(n.Inputs, c.(QueryExpressionHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewStarExceptList(
	identifiers interface{},
) (*StarExceptList, error) {
	var err error

	nn := &StarExceptList{}
	nn.SetKind(StarExceptListKind)

	err = nn.InitIdentifiers(identifiers)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StarExceptList) InitIdentifiers(d interface{}) error {
	if n.Identifiers != nil {
		return fmt.Errorf("StarExceptList.Identifiers: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Identifiers = make([]*Identifier, 0, defaultCapacity)

	return n.initIdentifiers(d)
}

func (n *StarExceptList) initIdentifiers(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIdentifiers(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*Identifier)
		n.Identifiers = append(n.Identifiers, newElem)
	}

	return nil
}

func (n *StarExceptList) AddChild(c NodeHandler) {
	n.Identifiers = append(n.Identifiers, c.(*Identifier))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *StarExceptList) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Identifiers = append(n.Identifiers, c.(*Identifier))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewStarModifiers(
	exceptlist interface{},
	replaceitems interface{},
) (*StarModifiers, error) {
	var err error

	nn := &StarModifiers{}
	nn.SetKind(StarModifiersKind)

	err = nn.InitExceptList(exceptlist)
	if err != nil {
		return nil, err
	}

	err = nn.InitReplaceItems(replaceitems)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StarModifiers) InitExceptList(d interface{}) error {
	return n.initExceptList(d)
}

func (n *StarModifiers) initExceptList(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.ExceptList = d.(*StarExceptList)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExceptList(t.Value)
	default:
		n.ExceptList = d.(*StarExceptList)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *StarModifiers) InitReplaceItems(d interface{}) error {
	if n.ReplaceItems != nil {
		return fmt.Errorf("StarModifiers.ReplaceItems: %w",
			ErrFieldAlreadyInitialized)
	}

	n.ReplaceItems = make([]*StarReplaceItem, 0, defaultCapacity)

	return n.initReplaceItems(d)
}

func (n *StarModifiers) initReplaceItems(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initReplaceItems(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*StarReplaceItem)
		n.ReplaceItems = append(n.ReplaceItems, newElem)
	}

	return nil
}

func (n *StarModifiers) AddChild(c NodeHandler) {
	n.ReplaceItems = append(n.ReplaceItems, c.(*StarReplaceItem))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *StarModifiers) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.ReplaceItems = append(n.ReplaceItems, c.(*StarReplaceItem))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewStarReplaceItem(
	expression interface{},
	alias interface{},
) (*StarReplaceItem, error) {
	var err error

	nn := &StarReplaceItem{}
	nn.SetKind(StarReplaceItemKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StarReplaceItem) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *StarReplaceItem) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("StarReplaceItem.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *StarReplaceItem) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *StarReplaceItem) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("StarReplaceItem.Alias: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	case NodeHandler:
		n.Alias = d.(*Identifier)
		n.Node.AddChild(t)
	default:
		n.Alias = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewStarWithModifiers(
	modifiers interface{},
) (*StarWithModifiers, error) {
	var err error

	nn := &StarWithModifiers{}
	nn.SetKind(StarWithModifiersKind)

	err = nn.InitModifiers(modifiers)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *StarWithModifiers) InitModifiers(d interface{}) error {
	return n.initModifiers(d)
}

func (n *StarWithModifiers) initModifiers(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("StarWithModifiers.Modifiers: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initModifiers(t.Value)
	case NodeHandler:
		n.Modifiers = d.(*StarModifiers)
		n.Expression.AddChild(t)
	default:
		n.Modifiers = d.(*StarModifiers)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewTableSubquery(
	subquery interface{},
	alias interface{},
	pivotclause interface{},
	unpivotclause interface{},
	sampleclause interface{},
) (*TableSubquery, error) {
	var err error

	nn := &TableSubquery{}
	nn.SetKind(TableSubqueryKind)

	err = nn.InitSubquery(subquery)
	if err != nil {
		return nil, err
	}

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	err = nn.InitPivotClause(pivotclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitUnpivotClause(unpivotclause)
	if err != nil {
		return nil, err
	}

	err = nn.InitSampleClause(sampleclause)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *TableSubquery) InitSubquery(d interface{}) error {
	return n.initSubquery(d)
}

func (n *TableSubquery) initSubquery(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("TableSubquery.Subquery: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSubquery(t.Value)
	case NodeHandler:
		n.Subquery = d.(*Query)
		n.TableExpression.AddChild(t)
	default:
		n.Subquery = d.(*Query)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TableSubquery) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *TableSubquery) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Alias = d.(*Alias)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	default:
		n.Alias = d.(*Alias)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TableSubquery) InitPivotClause(d interface{}) error {
	return n.initPivotClause(d)
}

func (n *TableSubquery) initPivotClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.PivotClause = d
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPivotClause(t.Value)
	default:
		n.PivotClause = d
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TableSubquery) InitUnpivotClause(d interface{}) error {
	return n.initUnpivotClause(d)
}

func (n *TableSubquery) initUnpivotClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.UnpivotClause = d
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initUnpivotClause(t.Value)
	default:
		n.UnpivotClause = d
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *TableSubquery) InitSampleClause(d interface{}) error {
	return n.initSampleClause(d)
}

func (n *TableSubquery) initSampleClause(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.SampleClause = d.(*SampleClause)
		n.TableExpression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSampleClause(t.Value)
	default:
		n.SampleClause = d.(*SampleClause)
		n.TableExpression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewUnaryExpression(
	op interface{},
	operand interface{},
) (*UnaryExpression, error) {
	var err error

	nn := &UnaryExpression{}
	nn.SetKind(UnaryExpressionKind)

	err = nn.InitOp(op)
	if err != nil {
		return nil, err
	}

	err = nn.InitOperand(operand)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *UnaryExpression) InitOp(d interface{}) error {
	return n.initOp(d)
}

func (n *UnaryExpression) initOp(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOp(t.Value)
	default:
		n.Op = d.(UnaryOp)
	}

	return nil
}

func (n *UnaryExpression) InitOperand(d interface{}) error {
	return n.initOperand(d)
}

func (n *UnaryExpression) initOperand(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("UnaryExpression.Operand: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOperand(t.Value)
	case NodeHandler:
		n.Operand = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Operand = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewUnnestExpression(
	expression interface{},
) (*UnnestExpression, error) {
	var err error

	nn := &UnnestExpression{}
	nn.SetKind(UnnestExpressionKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *UnnestExpression) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *UnnestExpression) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("UnnestExpression.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWindowClause(
	windows interface{},
) (*WindowClause, error) {
	var err error

	nn := &WindowClause{}
	nn.SetKind(WindowClauseKind)

	err = nn.InitWindows(windows)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WindowClause) InitWindows(d interface{}) error {
	if n.Windows != nil {
		return fmt.Errorf("WindowClause.Windows: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Windows = make([]*WindowDefinition, 0, defaultCapacity)

	return n.initWindows(d)
}

func (n *WindowClause) initWindows(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWindows(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(*WindowDefinition)
		n.Windows = append(n.Windows, newElem)
	}

	return nil
}

func (n *WindowClause) AddChild(c NodeHandler) {
	n.Windows = append(n.Windows, c.(*WindowDefinition))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *WindowClause) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Windows = append(n.Windows, c.(*WindowDefinition))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewWindowDefinition(
	name interface{},
	windowspec interface{},
) (*WindowDefinition, error) {
	var err error

	nn := &WindowDefinition{}
	nn.SetKind(WindowDefinitionKind)

	err = nn.InitName(name)
	if err != nil {
		return nil, err
	}

	err = nn.InitWindowSpec(windowspec)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WindowDefinition) InitName(d interface{}) error {
	return n.initName(d)
}

func (n *WindowDefinition) initName(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WindowDefinition.Name: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initName(t.Value)
	case NodeHandler:
		n.Name = d.(*Identifier)
		n.Node.AddChild(t)
	default:
		n.Name = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowDefinition) InitWindowSpec(d interface{}) error {
	return n.initWindowSpec(d)
}

func (n *WindowDefinition) initWindowSpec(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WindowDefinition.WindowSpec: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWindowSpec(t.Value)
	case NodeHandler:
		n.WindowSpec = d.(*WindowSpecification)
		n.Node.AddChild(t)
	default:
		n.WindowSpec = d.(*WindowSpecification)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWindowFrame(
	startexpr interface{},
	endexpr interface{},
	frameunit interface{},
) (*WindowFrame, error) {
	var err error

	nn := &WindowFrame{}
	nn.SetKind(WindowFrameKind)

	err = nn.InitStartExpr(startexpr)
	if err != nil {
		return nil, err
	}

	err = nn.InitEndExpr(endexpr)
	if err != nil {
		return nil, err
	}

	err = nn.InitFrameUnit(frameunit)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WindowFrame) InitStartExpr(d interface{}) error {
	return n.initStartExpr(d)
}

func (n *WindowFrame) initStartExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WindowFrame.StartExpr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initStartExpr(t.Value)
	case NodeHandler:
		n.StartExpr = d.(*WindowFrameExpr)
		n.Node.AddChild(t)
	default:
		n.StartExpr = d.(*WindowFrameExpr)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowFrame) InitEndExpr(d interface{}) error {
	return n.initEndExpr(d)
}

func (n *WindowFrame) initEndExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WindowFrame.EndExpr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initEndExpr(t.Value)
	case NodeHandler:
		n.EndExpr = d.(*WindowFrameExpr)
		n.Node.AddChild(t)
	default:
		n.EndExpr = d.(*WindowFrameExpr)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowFrame) InitFrameUnit(d interface{}) error {
	return n.initFrameUnit(d)
}

func (n *WindowFrame) initFrameUnit(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WindowFrame.FrameUnit: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFrameUnit(t.Value)
	case NodeHandler:
		n.FrameUnit = d.(FrameUnit)
		n.Node.AddChild(t)
	default:
		n.FrameUnit = d.(FrameUnit)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWindowFrameExpr(
	expression interface{},
	boundarytype interface{},
) (*WindowFrameExpr, error) {
	var err error

	nn := &WindowFrameExpr{}
	nn.SetKind(WindowFrameExprKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	err = nn.InitBoundaryType(boundarytype)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WindowFrameExpr) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *WindowFrameExpr) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("WindowFrameExpr.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowFrameExpr) InitBoundaryType(d interface{}) error {
	return n.initBoundaryType(d)
}

func (n *WindowFrameExpr) initBoundaryType(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.BoundaryType = d.(BoundaryType)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initBoundaryType(t.Value)
	default:
		n.BoundaryType = d.(BoundaryType)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewLikeExpression(
	lhs interface{},
	inlist interface{},
	isnot interface{},
) (*LikeExpression, error) {
	var err error

	nn := &LikeExpression{}
	nn.SetKind(LikeExpressionKind)

	err = nn.InitLHS(lhs)
	if err != nil {
		return nil, err
	}

	err = nn.InitInList(inlist)
	if err != nil {
		return nil, err
	}

	err = nn.InitIsNot(isnot)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *LikeExpression) InitLHS(d interface{}) error {
	return n.initLHS(d)
}

func (n *LikeExpression) initLHS(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("LikeExpression.LHS: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initLHS(t.Value)
	case NodeHandler:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.LHS = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *LikeExpression) InitInList(d interface{}) error {
	return n.initInList(d)
}

func (n *LikeExpression) initInList(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("LikeExpression.InList: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initInList(t.Value)
	case NodeHandler:
		n.InList = d.(*InList)
		n.Expression.AddChild(t)
	default:
		n.InList = d.(*InList)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *LikeExpression) InitIsNot(d interface{}) error {
	return n.initIsNot(d)
}

func (n *LikeExpression) initIsNot(d interface{}) error {
	switch t := d.(type) {
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initIsNot(t.Value)
	default:
		n.IsNot = d.(bool)
	}

	return nil
}

func NewWindowSpecification(
	basewindowname interface{},
	partitionby interface{},
	orderby interface{},
	windowframe interface{},
) (*WindowSpecification, error) {
	var err error

	nn := &WindowSpecification{}
	nn.SetKind(WindowSpecificationKind)

	err = nn.InitBaseWindowName(basewindowname)
	if err != nil {
		return nil, err
	}

	err = nn.InitPartitionBy(partitionby)
	if err != nil {
		return nil, err
	}

	err = nn.InitOrderBy(orderby)
	if err != nil {
		return nil, err
	}

	err = nn.InitWindowFrame(windowframe)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WindowSpecification) InitBaseWindowName(d interface{}) error {
	return n.initBaseWindowName(d)
}

func (n *WindowSpecification) initBaseWindowName(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.BaseWindowName = d.(*Identifier)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initBaseWindowName(t.Value)
	default:
		n.BaseWindowName = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowSpecification) InitPartitionBy(d interface{}) error {
	return n.initPartitionBy(d)
}

func (n *WindowSpecification) initPartitionBy(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.PartitionBy = d.(*PartitionBy)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPartitionBy(t.Value)
	default:
		n.PartitionBy = d.(*PartitionBy)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowSpecification) InitOrderBy(d interface{}) error {
	return n.initOrderBy(d)
}

func (n *WindowSpecification) initOrderBy(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.OrderBy = d.(*OrderBy)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initOrderBy(t.Value)
	default:
		n.OrderBy = d.(*OrderBy)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *WindowSpecification) InitWindowFrame(d interface{}) error {
	return n.initWindowFrame(d)
}

func (n *WindowSpecification) initWindowFrame(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.WindowFrame = d.(*WindowFrame)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWindowFrame(t.Value)
	default:
		n.WindowFrame = d.(*WindowFrame)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWithOffset(
	alias interface{},
) (*WithOffset, error) {
	var err error

	nn := &WithOffset{}
	nn.SetKind(WithOffsetKind)

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WithOffset) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *WithOffset) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Alias = d.(*Alias)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	default:
		n.Alias = d.(*Alias)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewTypeParameterList(
	parameters interface{},
) (*TypeParameterList, error) {
	var err error

	nn := &TypeParameterList{}
	nn.SetKind(TypeParameterListKind)

	err = nn.InitParameters(parameters)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *TypeParameterList) InitParameters(d interface{}) error {
	if n.Parameters != nil {
		return fmt.Errorf("TypeParameterList.Parameters: %w",
			ErrFieldAlreadyInitialized)
	}

	n.Parameters = make([]LeafHandler, 0, defaultCapacity)

	return n.initParameters(d)
}

func (n *TypeParameterList) initParameters(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initParameters(t.Value)
	case NodeHandler:
		n.AddChild(t)
	case []interface{}:
		for _, elem := range t {
			n.AddChild(elem.(NodeHandler))
		}
	default:
		newElem := d.(LeafHandler)
		n.Parameters = append(n.Parameters, newElem)
	}

	return nil
}

func (n *TypeParameterList) AddChild(c NodeHandler) {
	n.Parameters = append(n.Parameters, c.(LeafHandler))

	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *TypeParameterList) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}

		n.Parameters = append(n.Parameters, c.(LeafHandler))

		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewSampleClause(
	samplemethod interface{},
	samplesize interface{},
	samplesuffix interface{},
) (*SampleClause, error) {
	var err error

	nn := &SampleClause{}
	nn.SetKind(SampleClauseKind)

	err = nn.InitSampleMethod(samplemethod)
	if err != nil {
		return nil, err
	}

	err = nn.InitSampleSize(samplesize)
	if err != nil {
		return nil, err
	}

	err = nn.InitSampleSuffix(samplesuffix)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SampleClause) InitSampleMethod(d interface{}) error {
	return n.initSampleMethod(d)
}

func (n *SampleClause) initSampleMethod(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("SampleClause.SampleMethod: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSampleMethod(t.Value)
	case NodeHandler:
		n.SampleMethod = d.(*Identifier)
		n.Node.AddChild(t)
	default:
		n.SampleMethod = d.(*Identifier)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SampleClause) InitSampleSize(d interface{}) error {
	return n.initSampleSize(d)
}

func (n *SampleClause) initSampleSize(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("SampleClause.SampleSize: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSampleSize(t.Value)
	case NodeHandler:
		n.SampleSize = d.(*SampleSize)
		n.Node.AddChild(t)
	default:
		n.SampleSize = d.(*SampleSize)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SampleClause) InitSampleSuffix(d interface{}) error {
	return n.initSampleSuffix(d)
}

func (n *SampleClause) initSampleSuffix(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.SampleSuffix = d.(*SampleSuffix)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSampleSuffix(t.Value)
	default:
		n.SampleSuffix = d.(*SampleSuffix)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewSampleSize(
	size interface{},
	partitionby interface{},
	unit interface{},
) (*SampleSize, error) {
	var err error

	nn := &SampleSize{}
	nn.SetKind(SampleSizeKind)

	err = nn.InitSize(size)
	if err != nil {
		return nil, err
	}

	err = nn.InitPartitionBy(partitionby)
	if err != nil {
		return nil, err
	}

	err = nn.InitUnit(unit)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SampleSize) InitSize(d interface{}) error {
	return n.initSize(d)
}

func (n *SampleSize) initSize(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("SampleSize.Size: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initSize(t.Value)
	case NodeHandler:
		n.Size = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Size = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SampleSize) InitPartitionBy(d interface{}) error {
	return n.initPartitionBy(d)
}

func (n *SampleSize) initPartitionBy(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.PartitionBy = d.(*PartitionBy)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initPartitionBy(t.Value)
	default:
		n.PartitionBy = d.(*PartitionBy)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SampleSize) InitUnit(d interface{}) error {
	return n.initUnit(d)
}

func (n *SampleSize) initUnit(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Unit = d.(SampleSizeUnit)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initUnit(t.Value)
	default:
		n.Unit = d.(SampleSizeUnit)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewSampleSuffix(
	weight interface{},
	repeat interface{},
) (*SampleSuffix, error) {
	var err error

	nn := &SampleSuffix{}
	nn.SetKind(SampleSuffixKind)

	err = nn.InitWeight(weight)
	if err != nil {
		return nil, err
	}

	err = nn.InitRepeat(repeat)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *SampleSuffix) InitWeight(d interface{}) error {
	return n.initWeight(d)
}

func (n *SampleSuffix) initWeight(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Weight = d.(*WithWeight)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWeight(t.Value)
	default:
		n.Weight = d.(*WithWeight)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *SampleSuffix) InitRepeat(d interface{}) error {
	return n.initRepeat(d)
}

func (n *SampleSuffix) initRepeat(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Repeat = d.(*RepeatableClause)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initRepeat(t.Value)
	default:
		n.Repeat = d.(*RepeatableClause)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewWithWeight(
	alias interface{},
) (*WithWeight, error) {
	var err error

	nn := &WithWeight{}
	nn.SetKind(WithWeightKind)

	err = nn.InitAlias(alias)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *WithWeight) InitAlias(d interface{}) error {
	return n.initAlias(d)
}

func (n *WithWeight) initAlias(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Alias = d.(*Alias)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initAlias(t.Value)
	default:
		n.Alias = d.(*Alias)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewRepeatableClause(
	argument interface{},
) (*RepeatableClause, error) {
	var err error

	nn := &RepeatableClause{}
	nn.SetKind(RepeatableClauseKind)

	err = nn.InitArgument(argument)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *RepeatableClause) InitArgument(d interface{}) error {
	return n.initArgument(d)
}

func (n *RepeatableClause) initArgument(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("RepeatableClause.Argument: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initArgument(t.Value)
	case NodeHandler:
		n.Argument = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Argument = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewQualify(
	expression interface{},
) (*Qualify, error) {
	var err error

	nn := &Qualify{}
	nn.SetKind(QualifyKind)

	err = nn.InitExpression(expression)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *Qualify) InitExpression(d interface{}) error {
	return n.initExpression(d)
}

func (n *Qualify) initExpression(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("Qualify.Expression: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpression(t.Value)
	case NodeHandler:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Expression = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewFormatClause(
	format interface{},
	timezoneexpr interface{},
) (*FormatClause, error) {
	var err error

	nn := &FormatClause{}
	nn.SetKind(FormatClauseKind)

	err = nn.InitFormat(format)
	if err != nil {
		return nil, err
	}

	err = nn.InitTimeZoneExpr(timezoneexpr)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *FormatClause) InitFormat(d interface{}) error {
	return n.initFormat(d)
}

func (n *FormatClause) initFormat(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("FormatClause.Format: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initFormat(t.Value)
	case NodeHandler:
		n.Format = d.(ExpressionHandler)
		n.Node.AddChild(t)
	default:
		n.Format = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *FormatClause) InitTimeZoneExpr(d interface{}) error {
	return n.initTimeZoneExpr(d)
}

func (n *FormatClause) initTimeZoneExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.TimeZoneExpr = d.(ExpressionHandler)
		n.Node.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initTimeZoneExpr(t.Value)
	default:
		n.TimeZoneExpr = d.(ExpressionHandler)
		n.Node.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewParameterExpr(
	name interface{},
) (*ParameterExpr, error) {
	var err error

	nn := &ParameterExpr{}
	nn.SetKind(ParameterExprKind)

	err = nn.InitName(name)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *ParameterExpr) InitName(d interface{}) error {
	return n.initName(d)
}

func (n *ParameterExpr) initName(d interface{}) error {
	switch t := d.(type) {
	case nil:
	case NodeHandler:
		n.Name = d.(*Identifier)
		n.Expression.AddChild(t)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initName(t.Value)
	default:
		n.Name = d.(*Identifier)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func NewAnalyticFunctionCall(
	expr interface{},
	windowspec interface{},
) (*AnalyticFunctionCall, error) {
	var err error

	nn := &AnalyticFunctionCall{}
	nn.SetKind(AnalyticFunctionCallKind)

	err = nn.InitExpr(expr)
	if err != nil {
		return nil, err
	}

	err = nn.InitWindowSpec(windowspec)
	if err != nil {
		return nil, err
	}

	return nn, err
}

func (n *AnalyticFunctionCall) InitExpr(d interface{}) error {
	return n.initExpr(d)
}

func (n *AnalyticFunctionCall) initExpr(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("AnalyticFunctionCall.Expr: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initExpr(t.Value)
	case NodeHandler:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(t)
	default:
		n.Expr = d.(ExpressionHandler)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

func (n *AnalyticFunctionCall) InitWindowSpec(d interface{}) error {
	return n.initWindowSpec(d)
}

func (n *AnalyticFunctionCall) initWindowSpec(d interface{}) error {
	switch t := d.(type) {
	case nil:
		return fmt.Errorf("AnalyticFunctionCall.WindowSpec: %w",
			ErrMissingRequiredField)
	case *Wrapped:
		n.ExpandLoc(t.Loc.Start, t.Loc.End)
		return n.initWindowSpec(t.Value)
	case NodeHandler:
		n.WindowSpec = d.(*WindowSpecification)
		n.Expression.AddChild(t)
	default:
		n.WindowSpec = d.(*WindowSpecification)
		n.Expression.AddChild(d.(NodeHandler))
	}

	return nil
}

type Visitor interface {
	// Visit by type composition. These one are called when the visit
	// for the concrete type is not implemented.
	VisitNodeHandler(NodeHandler, interface{})
	VisitTableExpressionHandler(TableExpressionHandler, interface{})
	VisitQueryExpressionHandler(QueryExpressionHandler, interface{})
	VisitExpressionHandler(ExpressionHandler, interface{})
	VisitTypeHandler(TypeHandler, interface{})
	VisitLeafHandler(LeafHandler, interface{})
	VisitStatementHandler(StatementHandler, interface{})

	// Visit by concrete type
	VisitQueryStatement(*QueryStatement, interface{})
	VisitQuery(*Query, interface{})
	VisitSelect(*Select, interface{})
	VisitSelectList(*SelectList, interface{})
	VisitSelectColumn(*SelectColumn, interface{})
	VisitIntLiteral(*IntLiteral, interface{})
	VisitIdentifier(*Identifier, interface{})
	VisitAlias(*Alias, interface{})
	VisitPathExpression(*PathExpression, interface{})
	VisitTablePathExpression(*TablePathExpression, interface{})
	VisitFromClause(*FromClause, interface{})
	VisitWhereClause(*WhereClause, interface{})
	VisitBooleanLiteral(*BooleanLiteral, interface{})
	VisitAndExpr(*AndExpr, interface{})
	VisitBinaryExpression(*BinaryExpression, interface{})
	VisitStringLiteral(*StringLiteral, interface{})
	VisitStar(*Star, interface{})
	VisitOrExpr(*OrExpr, interface{})
	VisitGroupingItem(*GroupingItem, interface{})
	VisitGroupBy(*GroupBy, interface{})
	VisitOrderingExpression(*OrderingExpression, interface{})
	VisitOrderBy(*OrderBy, interface{})
	VisitLimitOffset(*LimitOffset, interface{})
	VisitFloatLiteral(*FloatLiteral, interface{})
	VisitNullLiteral(*NullLiteral, interface{})
	VisitJSONLiteral(*JSONLiteral, interface{})
	VisitOnClause(*OnClause, interface{})
	VisitWithClauseEntry(*WithClauseEntry, interface{})
	VisitJoin(*Join, interface{})
	VisitUsingClause(*UsingClause, interface{})
	VisitWithClause(*WithClause, interface{})
	VisitHaving(*Having, interface{})
	VisitNamedType(*NamedType, interface{})
	VisitArrayType(*ArrayType, interface{})
	VisitStructType(*StructType, interface{})
	VisitStructField(*StructField, interface{})
	VisitCastExpression(*CastExpression, interface{})
	VisitSelectAs(*SelectAs, interface{})
	VisitRollup(*Rollup, interface{})
	VisitFunctionCall(*FunctionCall, interface{})
	VisitArrayConstructor(*ArrayConstructor, interface{})
	VisitStructConstructorArg(*StructConstructorArg, interface{})
	VisitStructConstructorWithParens(*StructConstructorWithParens, interface{})
	VisitStructConstructorWithKeyword(*StructConstructorWithKeyword, interface{})
	VisitInExpression(*InExpression, interface{})
	VisitInList(*InList, interface{})
	VisitBetweenExpression(*BetweenExpression, interface{})
	VisitNumericLiteral(*NumericLiteral, interface{})
	VisitBigNumericLiteral(*BigNumericLiteral, interface{})
	VisitBytesLiteral(*BytesLiteral, interface{})
	VisitDateOrTimeLiteral(*DateOrTimeLiteral, interface{})
	VisitCaseValueExpression(*CaseValueExpression, interface{})
	VisitCaseNoValueExpression(*CaseNoValueExpression, interface{})
	VisitArrayElement(*ArrayElement, interface{})
	VisitBitwiseShiftExpression(*BitwiseShiftExpression, interface{})
	VisitDotGeneralizedField(*DotGeneralizedField, interface{})
	VisitDotIdentifier(*DotIdentifier, interface{})
	VisitDotStar(*DotStar, interface{})
	VisitDotStarWithModifiers(*DotStarWithModifiers, interface{})
	VisitExpressionSubquery(*ExpressionSubquery, interface{})
	VisitExtractExpression(*ExtractExpression, interface{})
	VisitIntervalExpr(*IntervalExpr, interface{})
	VisitNullOrder(*NullOrder, interface{})
	VisitOnOrUsingClauseList(*OnOrUsingClauseList, interface{})
	VisitParenthesizedJoin(*ParenthesizedJoin, interface{})
	VisitPartitionBy(*PartitionBy, interface{})
	VisitSetOperation(*SetOperation, interface{})
	VisitStarExceptList(*StarExceptList, interface{})
	VisitStarModifiers(*StarModifiers, interface{})
	VisitStarReplaceItem(*StarReplaceItem, interface{})
	VisitStarWithModifiers(*StarWithModifiers, interface{})
	VisitTableSubquery(*TableSubquery, interface{})
	VisitUnaryExpression(*UnaryExpression, interface{})
	VisitUnnestExpression(*UnnestExpression, interface{})
	VisitWindowClause(*WindowClause, interface{})
	VisitWindowDefinition(*WindowDefinition, interface{})
	VisitWindowFrame(*WindowFrame, interface{})
	VisitWindowFrameExpr(*WindowFrameExpr, interface{})
	VisitLikeExpression(*LikeExpression, interface{})
	VisitWindowSpecification(*WindowSpecification, interface{})
	VisitWithOffset(*WithOffset, interface{})
	VisitTypeParameterList(*TypeParameterList, interface{})
	VisitSampleClause(*SampleClause, interface{})
	VisitSampleSize(*SampleSize, interface{})
	VisitSampleSuffix(*SampleSuffix, interface{})
	VisitWithWeight(*WithWeight, interface{})
	VisitRepeatableClause(*RepeatableClause, interface{})
	VisitQualify(*Qualify, interface{})
	VisitFormatClause(*FormatClause, interface{})
	VisitParameterExpr(*ParameterExpr, interface{})
	VisitAnalyticFunctionCall(*AnalyticFunctionCall, interface{})
}

func (n *QueryStatement) Accept(v Visitor, d interface{}) {
	v.VisitQueryStatement(n, d)
}

func (n *Query) Accept(v Visitor, d interface{}) {
	v.VisitQuery(n, d)
}

func (n *Select) Accept(v Visitor, d interface{}) {
	v.VisitSelect(n, d)
}

func (n *SelectList) Accept(v Visitor, d interface{}) {
	v.VisitSelectList(n, d)
}

func (n *SelectColumn) Accept(v Visitor, d interface{}) {
	v.VisitSelectColumn(n, d)
}

func (n *IntLiteral) Accept(v Visitor, d interface{}) {
	v.VisitIntLiteral(n, d)
}

func (n *Identifier) Accept(v Visitor, d interface{}) {
	v.VisitIdentifier(n, d)
}

func (n *Alias) Accept(v Visitor, d interface{}) {
	v.VisitAlias(n, d)
}

func (n *PathExpression) Accept(v Visitor, d interface{}) {
	v.VisitPathExpression(n, d)
}

func (n *TablePathExpression) Accept(v Visitor, d interface{}) {
	v.VisitTablePathExpression(n, d)
}

func (n *FromClause) Accept(v Visitor, d interface{}) {
	v.VisitFromClause(n, d)
}

func (n *WhereClause) Accept(v Visitor, d interface{}) {
	v.VisitWhereClause(n, d)
}

func (n *BooleanLiteral) Accept(v Visitor, d interface{}) {
	v.VisitBooleanLiteral(n, d)
}

func (n *AndExpr) Accept(v Visitor, d interface{}) {
	v.VisitAndExpr(n, d)
}

func (n *BinaryExpression) Accept(v Visitor, d interface{}) {
	v.VisitBinaryExpression(n, d)
}

func (n *StringLiteral) Accept(v Visitor, d interface{}) {
	v.VisitStringLiteral(n, d)
}

func (n *Star) Accept(v Visitor, d interface{}) {
	v.VisitStar(n, d)
}

func (n *OrExpr) Accept(v Visitor, d interface{}) {
	v.VisitOrExpr(n, d)
}

func (n *GroupingItem) Accept(v Visitor, d interface{}) {
	v.VisitGroupingItem(n, d)
}

func (n *GroupBy) Accept(v Visitor, d interface{}) {
	v.VisitGroupBy(n, d)
}

func (n *OrderingExpression) Accept(v Visitor, d interface{}) {
	v.VisitOrderingExpression(n, d)
}

func (n *OrderBy) Accept(v Visitor, d interface{}) {
	v.VisitOrderBy(n, d)
}

func (n *LimitOffset) Accept(v Visitor, d interface{}) {
	v.VisitLimitOffset(n, d)
}

func (n *FloatLiteral) Accept(v Visitor, d interface{}) {
	v.VisitFloatLiteral(n, d)
}

func (n *NullLiteral) Accept(v Visitor, d interface{}) {
	v.VisitNullLiteral(n, d)
}

func (n *JSONLiteral) Accept(v Visitor, d interface{}) {
	v.VisitJSONLiteral(n, d)
}

func (n *OnClause) Accept(v Visitor, d interface{}) {
	v.VisitOnClause(n, d)
}

func (n *WithClauseEntry) Accept(v Visitor, d interface{}) {
	v.VisitWithClauseEntry(n, d)
}

func (n *Join) Accept(v Visitor, d interface{}) {
	v.VisitJoin(n, d)
}

func (n *UsingClause) Accept(v Visitor, d interface{}) {
	v.VisitUsingClause(n, d)
}

func (n *WithClause) Accept(v Visitor, d interface{}) {
	v.VisitWithClause(n, d)
}

func (n *Having) Accept(v Visitor, d interface{}) {
	v.VisitHaving(n, d)
}

func (n *NamedType) Accept(v Visitor, d interface{}) {
	v.VisitNamedType(n, d)
}

func (n *ArrayType) Accept(v Visitor, d interface{}) {
	v.VisitArrayType(n, d)
}

func (n *StructType) Accept(v Visitor, d interface{}) {
	v.VisitStructType(n, d)
}

func (n *StructField) Accept(v Visitor, d interface{}) {
	v.VisitStructField(n, d)
}

func (n *CastExpression) Accept(v Visitor, d interface{}) {
	v.VisitCastExpression(n, d)
}

func (n *SelectAs) Accept(v Visitor, d interface{}) {
	v.VisitSelectAs(n, d)
}

func (n *Rollup) Accept(v Visitor, d interface{}) {
	v.VisitRollup(n, d)
}

func (n *FunctionCall) Accept(v Visitor, d interface{}) {
	v.VisitFunctionCall(n, d)
}

func (n *ArrayConstructor) Accept(v Visitor, d interface{}) {
	v.VisitArrayConstructor(n, d)
}

func (n *StructConstructorArg) Accept(v Visitor, d interface{}) {
	v.VisitStructConstructorArg(n, d)
}

func (n *StructConstructorWithParens) Accept(v Visitor, d interface{}) {
	v.VisitStructConstructorWithParens(n, d)
}

func (n *StructConstructorWithKeyword) Accept(v Visitor, d interface{}) {
	v.VisitStructConstructorWithKeyword(n, d)
}

func (n *InExpression) Accept(v Visitor, d interface{}) {
	v.VisitInExpression(n, d)
}

func (n *InList) Accept(v Visitor, d interface{}) {
	v.VisitInList(n, d)
}

func (n *BetweenExpression) Accept(v Visitor, d interface{}) {
	v.VisitBetweenExpression(n, d)
}

func (n *NumericLiteral) Accept(v Visitor, d interface{}) {
	v.VisitNumericLiteral(n, d)
}

func (n *BigNumericLiteral) Accept(v Visitor, d interface{}) {
	v.VisitBigNumericLiteral(n, d)
}

func (n *BytesLiteral) Accept(v Visitor, d interface{}) {
	v.VisitBytesLiteral(n, d)
}

func (n *DateOrTimeLiteral) Accept(v Visitor, d interface{}) {
	v.VisitDateOrTimeLiteral(n, d)
}

func (n *CaseValueExpression) Accept(v Visitor, d interface{}) {
	v.VisitCaseValueExpression(n, d)
}

func (n *CaseNoValueExpression) Accept(v Visitor, d interface{}) {
	v.VisitCaseNoValueExpression(n, d)
}

func (n *ArrayElement) Accept(v Visitor, d interface{}) {
	v.VisitArrayElement(n, d)
}

func (n *BitwiseShiftExpression) Accept(v Visitor, d interface{}) {
	v.VisitBitwiseShiftExpression(n, d)
}

func (n *DotGeneralizedField) Accept(v Visitor, d interface{}) {
	v.VisitDotGeneralizedField(n, d)
}

func (n *DotIdentifier) Accept(v Visitor, d interface{}) {
	v.VisitDotIdentifier(n, d)
}

func (n *DotStar) Accept(v Visitor, d interface{}) {
	v.VisitDotStar(n, d)
}

func (n *DotStarWithModifiers) Accept(v Visitor, d interface{}) {
	v.VisitDotStarWithModifiers(n, d)
}

func (n *ExpressionSubquery) Accept(v Visitor, d interface{}) {
	v.VisitExpressionSubquery(n, d)
}

func (n *ExtractExpression) Accept(v Visitor, d interface{}) {
	v.VisitExtractExpression(n, d)
}

func (n *IntervalExpr) Accept(v Visitor, d interface{}) {
	v.VisitIntervalExpr(n, d)
}

func (n *NullOrder) Accept(v Visitor, d interface{}) {
	v.VisitNullOrder(n, d)
}

func (n *OnOrUsingClauseList) Accept(v Visitor, d interface{}) {
	v.VisitOnOrUsingClauseList(n, d)
}

func (n *ParenthesizedJoin) Accept(v Visitor, d interface{}) {
	v.VisitParenthesizedJoin(n, d)
}

func (n *PartitionBy) Accept(v Visitor, d interface{}) {
	v.VisitPartitionBy(n, d)
}

func (n *SetOperation) Accept(v Visitor, d interface{}) {
	v.VisitSetOperation(n, d)
}

func (n *StarExceptList) Accept(v Visitor, d interface{}) {
	v.VisitStarExceptList(n, d)
}

func (n *StarModifiers) Accept(v Visitor, d interface{}) {
	v.VisitStarModifiers(n, d)
}

func (n *StarReplaceItem) Accept(v Visitor, d interface{}) {
	v.VisitStarReplaceItem(n, d)
}

func (n *StarWithModifiers) Accept(v Visitor, d interface{}) {
	v.VisitStarWithModifiers(n, d)
}

func (n *TableSubquery) Accept(v Visitor, d interface{}) {
	v.VisitTableSubquery(n, d)
}

func (n *UnaryExpression) Accept(v Visitor, d interface{}) {
	v.VisitUnaryExpression(n, d)
}

func (n *UnnestExpression) Accept(v Visitor, d interface{}) {
	v.VisitUnnestExpression(n, d)
}

func (n *WindowClause) Accept(v Visitor, d interface{}) {
	v.VisitWindowClause(n, d)
}

func (n *WindowDefinition) Accept(v Visitor, d interface{}) {
	v.VisitWindowDefinition(n, d)
}

func (n *WindowFrame) Accept(v Visitor, d interface{}) {
	v.VisitWindowFrame(n, d)
}

func (n *WindowFrameExpr) Accept(v Visitor, d interface{}) {
	v.VisitWindowFrameExpr(n, d)
}

func (n *LikeExpression) Accept(v Visitor, d interface{}) {
	v.VisitLikeExpression(n, d)
}

func (n *WindowSpecification) Accept(v Visitor, d interface{}) {
	v.VisitWindowSpecification(n, d)
}

func (n *WithOffset) Accept(v Visitor, d interface{}) {
	v.VisitWithOffset(n, d)
}

func (n *TypeParameterList) Accept(v Visitor, d interface{}) {
	v.VisitTypeParameterList(n, d)
}

func (n *SampleClause) Accept(v Visitor, d interface{}) {
	v.VisitSampleClause(n, d)
}

func (n *SampleSize) Accept(v Visitor, d interface{}) {
	v.VisitSampleSize(n, d)
}

func (n *SampleSuffix) Accept(v Visitor, d interface{}) {
	v.VisitSampleSuffix(n, d)
}

func (n *WithWeight) Accept(v Visitor, d interface{}) {
	v.VisitWithWeight(n, d)
}

func (n *RepeatableClause) Accept(v Visitor, d interface{}) {
	v.VisitRepeatableClause(n, d)
}

func (n *Qualify) Accept(v Visitor, d interface{}) {
	v.VisitQualify(n, d)
}

func (n *FormatClause) Accept(v Visitor, d interface{}) {
	v.VisitFormatClause(n, d)
}

func (n *ParameterExpr) Accept(v Visitor, d interface{}) {
	v.VisitParameterExpr(n, d)
}

func (n *AnalyticFunctionCall) Accept(v Visitor, d interface{}) {
	v.VisitAnalyticFunctionCall(n, d)
}

// Operation is the base for new operations using visitors.
type Operation struct {
	// visitor is the visitor to passed when fallbacking to walk.
	visitor Visitor
}

func (o *Operation) VisitNodeHandler(n NodeHandler, d interface{}) {
	for _, c := range n.Children() {
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitTableExpressionHandler(
	n TableExpressionHandler, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitQueryExpressionHandler(
	n QueryExpressionHandler, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitExpressionHandler(
	n ExpressionHandler, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitTypeHandler(n TypeHandler, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitLeafHandler(n LeafHandler, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitStatementHandler(n StatementHandler, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}

func (o *Operation) VisitQueryStatement(n *QueryStatement, d interface{}) {
	o.visitor.VisitStatementHandler(n, d)
}
func (o *Operation) VisitQuery(n *Query, d interface{}) {
	o.visitor.VisitQueryExpressionHandler(n, d)
}
func (o *Operation) VisitSelect(n *Select, d interface{}) {
	o.visitor.VisitQueryExpressionHandler(n, d)
}
func (o *Operation) VisitSelectList(n *SelectList, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitSelectColumn(n *SelectColumn, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitIntLiteral(n *IntLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitIdentifier(n *Identifier, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitAlias(n *Alias, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitPathExpression(n *PathExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitTablePathExpression(
	n *TablePathExpression, d interface{}) {
	o.visitor.VisitTableExpressionHandler(n, d)
}
func (o *Operation) VisitFromClause(n *FromClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWhereClause(n *WhereClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitBooleanLiteral(n *BooleanLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitAndExpr(n *AndExpr, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitBinaryExpression(n *BinaryExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitStringLiteral(n *StringLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitStar(n *Star, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitOrExpr(n *OrExpr, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitGroupingItem(n *GroupingItem, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitGroupBy(n *GroupBy, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitOrderingExpression(
	n *OrderingExpression, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitOrderBy(n *OrderBy, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitLimitOffset(n *LimitOffset, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitFloatLiteral(n *FloatLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitNullLiteral(n *NullLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitJSONLiteral(n *JSONLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitOnClause(n *OnClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWithClauseEntry(n *WithClauseEntry, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitJoin(n *Join, d interface{}) {
	o.visitor.VisitTableExpressionHandler(n, d)
}
func (o *Operation) VisitUsingClause(n *UsingClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWithClause(n *WithClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitHaving(n *Having, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitNamedType(n *NamedType, d interface{}) {
	o.visitor.VisitTypeHandler(n, d)
}
func (o *Operation) VisitArrayType(n *ArrayType, d interface{}) {
	o.visitor.VisitTypeHandler(n, d)
}
func (o *Operation) VisitStructType(n *StructType, d interface{}) {
	o.visitor.VisitTypeHandler(n, d)
}
func (o *Operation) VisitStructField(n *StructField, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitCastExpression(n *CastExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitSelectAs(n *SelectAs, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitRollup(n *Rollup, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitFunctionCall(n *FunctionCall, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitArrayConstructor(n *ArrayConstructor, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitStructConstructorArg(
	n *StructConstructorArg, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitStructConstructorWithParens(
	n *StructConstructorWithParens, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitStructConstructorWithKeyword(
	n *StructConstructorWithKeyword, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitInExpression(n *InExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitInList(n *InList, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitBetweenExpression(
	n *BetweenExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitNumericLiteral(n *NumericLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitBigNumericLiteral(
	n *BigNumericLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitBytesLiteral(n *BytesLiteral, d interface{}) {
	o.visitor.VisitLeafHandler(n, d)
}
func (o *Operation) VisitDateOrTimeLiteral(
	n *DateOrTimeLiteral, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitCaseValueExpression(
	n *CaseValueExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitCaseNoValueExpression(
	n *CaseNoValueExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitArrayElement(n *ArrayElement, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitBitwiseShiftExpression(
	n *BitwiseShiftExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitDotGeneralizedField(
	n *DotGeneralizedField, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitDotIdentifier(n *DotIdentifier, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitDotStar(n *DotStar, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitDotStarWithModifiers(
	n *DotStarWithModifiers, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitExpressionSubquery(
	n *ExpressionSubquery, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitExtractExpression(
	n *ExtractExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitIntervalExpr(n *IntervalExpr, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitNullOrder(n *NullOrder, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitOnOrUsingClauseList(
	n *OnOrUsingClauseList, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitParenthesizedJoin(
	n *ParenthesizedJoin, d interface{}) {
	o.visitor.VisitTableExpressionHandler(n, d)
}
func (o *Operation) VisitPartitionBy(n *PartitionBy, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitSetOperation(n *SetOperation, d interface{}) {
	o.visitor.VisitQueryExpressionHandler(n, d)
}
func (o *Operation) VisitStarExceptList(n *StarExceptList, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitStarModifiers(n *StarModifiers, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitStarReplaceItem(n *StarReplaceItem, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitStarWithModifiers(
	n *StarWithModifiers, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitTableSubquery(n *TableSubquery, d interface{}) {
	o.visitor.VisitTableExpressionHandler(n, d)
}
func (o *Operation) VisitUnaryExpression(n *UnaryExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitUnnestExpression(n *UnnestExpression, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWindowClause(n *WindowClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWindowDefinition(n *WindowDefinition, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWindowFrame(n *WindowFrame, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWindowFrameExpr(n *WindowFrameExpr, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitLikeExpression(n *LikeExpression, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitWindowSpecification(
	n *WindowSpecification, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWithOffset(n *WithOffset, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitTypeParameterList(
	n *TypeParameterList, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitSampleClause(n *SampleClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitSampleSize(n *SampleSize, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitSampleSuffix(n *SampleSuffix, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitWithWeight(n *WithWeight, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitRepeatableClause(n *RepeatableClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitQualify(n *Qualify, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitFormatClause(n *FormatClause, d interface{}) {
	o.visitor.VisitNodeHandler(n, d)
}
func (o *Operation) VisitParameterExpr(n *ParameterExpr, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}
func (o *Operation) VisitAnalyticFunctionCall(
	n *AnalyticFunctionCall, d interface{}) {
	o.visitor.VisitExpressionHandler(n, d)
}

type NodeKind int

const (
	QueryStatementKind NodeKind = iota
	QueryKind
	SelectKind
	SelectListKind
	SelectColumnKind
	IntLiteralKind
	IdentifierKind
	AliasKind
	PathExpressionKind
	TablePathExpressionKind
	FromClauseKind
	WhereClauseKind
	BooleanLiteralKind
	AndExprKind
	BinaryExpressionKind
	StringLiteralKind
	StarKind
	OrExprKind
	GroupingItemKind
	GroupByKind
	OrderingExpressionKind
	OrderByKind
	LimitOffsetKind
	FloatLiteralKind
	NullLiteralKind
	JSONLiteralKind
	OnClauseKind
	WithClauseEntryKind
	JoinKind
	UsingClauseKind
	WithClauseKind
	HavingKind
	NamedTypeKind
	ArrayTypeKind
	StructTypeKind
	StructFieldKind
	CastExpressionKind
	SelectAsKind
	RollupKind
	FunctionCallKind
	ArrayConstructorKind
	StructConstructorArgKind
	StructConstructorWithParensKind
	StructConstructorWithKeywordKind
	InExpressionKind
	InListKind
	BetweenExpressionKind
	NumericLiteralKind
	BigNumericLiteralKind
	BytesLiteralKind
	DateOrTimeLiteralKind
	CaseValueExpressionKind
	CaseNoValueExpressionKind
	ArrayElementKind
	BitwiseShiftExpressionKind
	DotGeneralizedFieldKind
	DotIdentifierKind
	DotStarKind
	DotStarWithModifiersKind
	ExpressionSubqueryKind
	ExtractExpressionKind
	IntervalExprKind
	NullOrderKind
	OnOrUsingClauseListKind
	ParenthesizedJoinKind
	PartitionByKind
	SetOperationKind
	StarExceptListKind
	StarModifiersKind
	StarReplaceItemKind
	StarWithModifiersKind
	TableSubqueryKind
	UnaryExpressionKind
	UnnestExpressionKind
	WindowClauseKind
	WindowDefinitionKind
	WindowFrameKind
	WindowFrameExprKind
	LikeExpressionKind
	WindowSpecificationKind
	WithOffsetKind
	TypeParameterListKind
	SampleClauseKind
	SampleSizeKind
	SampleSuffixKind
	WithWeightKind
	RepeatableClauseKind
	QualifyKind
	FormatClauseKind
	ParameterExprKind
	AnalyticFunctionCallKind
)

func (k NodeKind) String() string { //nolint:cyclop
	switch k {
	case QueryStatementKind:
		return "QueryStatement"
	case QueryKind:
		return "Query"
	case SelectKind:
		return "Select"
	case SelectListKind:
		return "SelectList"
	case SelectColumnKind:
		return "SelectColumn"
	case IntLiteralKind:
		return "IntLiteral"
	case IdentifierKind:
		return "Identifier"
	case AliasKind:
		return "Alias"
	case PathExpressionKind:
		return "PathExpression"
	case TablePathExpressionKind:
		return "TablePathExpression"
	case FromClauseKind:
		return "FromClause"
	case WhereClauseKind:
		return "WhereClause"
	case BooleanLiteralKind:
		return "BooleanLiteral"
	case AndExprKind:
		return "AndExpr"
	case BinaryExpressionKind:
		return "BinaryExpression"
	case StringLiteralKind:
		return "StringLiteral"
	case StarKind:
		return "Star"
	case OrExprKind:
		return "OrExpr"
	case GroupingItemKind:
		return "GroupingItem"
	case GroupByKind:
		return "GroupBy"
	case OrderingExpressionKind:
		return "OrderingExpression"
	case OrderByKind:
		return "OrderBy"
	case LimitOffsetKind:
		return "LimitOffset"
	case FloatLiteralKind:
		return "FloatLiteral"
	case NullLiteralKind:
		return "NullLiteral"
	case JSONLiteralKind:
		return "JSONLiteral"
	case OnClauseKind:
		return "OnClause"
	case WithClauseEntryKind:
		return "WithClauseEntry"
	case JoinKind:
		return "Join"
	case UsingClauseKind:
		return "UsingClause"
	case WithClauseKind:
		return "WithClause"
	case HavingKind:
		return "Having"
	case NamedTypeKind:
		return "NamedType"
	case ArrayTypeKind:
		return "ArrayType"
	case StructTypeKind:
		return "StructType"
	case StructFieldKind:
		return "StructField"
	case CastExpressionKind:
		return "CastExpression"
	case SelectAsKind:
		return "SelectAs"
	case RollupKind:
		return "Rollup"
	case FunctionCallKind:
		return "FunctionCall"
	case ArrayConstructorKind:
		return "ArrayConstructor"
	case StructConstructorArgKind:
		return "StructConstructorArg"
	case StructConstructorWithParensKind:
		return "StructConstructorWithParens"
	case StructConstructorWithKeywordKind:
		return "StructConstructorWithKeyword"
	case InExpressionKind:
		return "InExpression"
	case InListKind:
		return "InList"
	case BetweenExpressionKind:
		return "BetweenExpression"
	case NumericLiteralKind:
		return "NumericLiteral"
	case BigNumericLiteralKind:
		return "BigNumericLiteral"
	case BytesLiteralKind:
		return "BytesLiteral"
	case DateOrTimeLiteralKind:
		return "DateOrTimeLiteral"
	case CaseValueExpressionKind:
		return "CaseValueExpression"
	case CaseNoValueExpressionKind:
		return "CaseNoValueExpression"
	case ArrayElementKind:
		return "ArrayElement"
	case BitwiseShiftExpressionKind:
		return "BitwiseShiftExpression"
	case DotGeneralizedFieldKind:
		return "DotGeneralizedField"
	case DotIdentifierKind:
		return "DotIdentifier"
	case DotStarKind:
		return "DotStar"
	case DotStarWithModifiersKind:
		return "DotStarWithModifiers"
	case ExpressionSubqueryKind:
		return "ExpressionSubquery"
	case ExtractExpressionKind:
		return "ExtractExpression"
	case IntervalExprKind:
		return "IntervalExpr"
	case NullOrderKind:
		return "NullOrder"
	case OnOrUsingClauseListKind:
		return "OnOrUsingClauseList"
	case ParenthesizedJoinKind:
		return "ParenthesizedJoin"
	case PartitionByKind:
		return "PartitionBy"
	case SetOperationKind:
		return "SetOperation"
	case StarExceptListKind:
		return "StarExceptList"
	case StarModifiersKind:
		return "StarModifiers"
	case StarReplaceItemKind:
		return "StarReplaceItem"
	case StarWithModifiersKind:
		return "StarWithModifiers"
	case TableSubqueryKind:
		return "TableSubquery"
	case UnaryExpressionKind:
		return "UnaryExpression"
	case UnnestExpressionKind:
		return "UnnestExpression"
	case WindowClauseKind:
		return "WindowClause"
	case WindowDefinitionKind:
		return "WindowDefinition"
	case WindowFrameKind:
		return "WindowFrame"
	case WindowFrameExprKind:
		return "WindowFrameExpr"
	case LikeExpressionKind:
		return "LikeExpression"
	case WindowSpecificationKind:
		return "WindowSpecification"
	case WithOffsetKind:
		return "WithOffset"
	case TypeParameterListKind:
		return "TypeParameterList"
	case SampleClauseKind:
		return "SampleClause"
	case SampleSizeKind:
		return "SampleSize"
	case SampleSuffixKind:
		return "SampleSuffix"
	case WithWeightKind:
		return "WithWeight"
	case RepeatableClauseKind:
		return "RepeatableClause"
	case QualifyKind:
		return "Qualify"
	case FormatClauseKind:
		return "FormatClause"
	case ParameterExprKind:
		return "ParameterExpr"
	case AnalyticFunctionCallKind:
		return "AnalyticFunctionCall"
	}

	panic("unexpected kind")
}
