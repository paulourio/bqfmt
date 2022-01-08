package ast

import (
	"fmt"
)

// types_generated.go is generated from type_generated.go.j2 by
// gen_types.py.

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
	
	QueryExpression
}

type Select struct {
	Distinct bool
	SelectAs *SelectAs
	SelectList *SelectList
	FromClause *FromClause
	WhereClause *WhereClause
	GroupBy *GroupBy
	Having *Having
	Qualify *Qualify
	WindowClause *WindowClause
	
	QueryExpression
}

type SelectList struct {
	Columns []*SelectColumn
	
	Node
}

type SelectColumn struct {
	Expression ExpressionHandler
	Alias *Alias
	
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
	PathExpr *PathExpression
	UnnestExpr *UnnestExpression
	Alias *Alias
	
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
	Op BinaryOp
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
	Rollup *Rollup
	
	Node
}

type GroupBy struct {
	GroupingItems []*GroupingItem
	
	Node
}

type OrderingExpression struct {
	Expression ExpressionHandler
	NullOrder *NullOrder
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
	LHS TableExpressionHandler
	RHS TableExpressionHandler
	ClauseList *OnOrUsingClauseList
	JoinType JoinType
	
	TableExpression
}

type WithClause struct {
	With []*WithClauseEntry
	
	Node
}

type Having struct {
	
	Node
}

type NamedType struct {
	TypeName *PathExpression
	TypeParameters *TypeParameterList
	
	Type
}

type ArrayType struct {
	ElementType TypeHandler
	TypeParameters *TypeParameterList
	
	Type
}

type StructField struct {
	// Name will be nil for anonymous fields like in STRUCT<int,
    // string>.
	Name *Identifier
	
	Node
}

type StructType struct {
	StructFields *StructField
	TypeParameterList *TypeParameterList
	
	Type
}

type CastExpression struct {
	Expr ExpressionHandler
	Type TypeHandler
	Format *FormatClause
	IsSafeCast bool
	
	Expression
}

// SelectAs represents a SELECT with AS clause giving it an output
// type. Exactly one of SELECT AS STRUCT, SELECT AS VALUE, SELECT AS
// <TypeName> is present.
type SelectAs struct {
	TypeName *PathExpression
	AsMode AsMode
	
	Node
}

type Rollup struct {
	Expressions []ExpressionHandler
	
	Node
}

type FunctionCall struct {
	Function *PathExpression
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
	Type *ArrayType
	Elements []ExpressionHandler
	
	Expression
}

type StructConstructorArg struct {
	Expression ExpressionHandler
	Alias *Alias
	
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
	Fields []*StructConstructorArg
	
	Expression
}

// InExpression is resulted from expr IN (expr, expr, ...), expr IN
// UNNEST(...), and expr IN (query). Exactly one of InList, Query, or
// UnnestExpr is present.
type InExpression struct {
	LHS ExpressionHandler
	InList *InList
	Query *Query
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
	LHS ExpressionHandler
	Low ExpressionHandler
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
	
	Leaf
}

type DateOrTimeLiteral struct {
	StringLiteral *StringLiteral
	TypeKind TypeKind
	
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
	Array ExpressionHandler
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
	Path *PathExpression
	
	Expression
}

type DotStar struct {
	Expr ExpressionHandler
	
	Expression
}

// DotStarWithModifiers is an expression constructed through SELECT
// x.* EXCEPT (...) REPLACE (...).
type DotStarWithModifiers struct {
	Expr ExpressionHandler
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
	LHS ExpressionHandler
	RHS ExpressionHandler
	TimeZone ExpressionHandler
	
	Expression
}

type IntervalExpression struct {
	IntervalValue ExpressionHandler
	DatePartName ExpressionHandler
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
	Join *Join
	SampleClause SampleClause
	
	TableExpression
}

type PartitionBy struct {
	PartitioningExpressions []ExpressionHandler
	
	Node
}

type SetOperation struct {
	Inputs []QueryExpressionHandler
	OpType SetOp
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
	ExceptList *StarExceptList
	ReplaceItems []*StarReplaceItem
	
	Node
}

type StarReplaceItem struct {
	Expression ExpressionHandler
	Alias *Identifier
	
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
	Subquery *Query
	Alias *Alias
	SampleClause *SampleClause
	
	TableExpression
}

type UnaryExpression struct {
	Operand ExpressionHandler
	Op UnaryOp
	
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
	Name *Identifier
	WindowSpec *WindowSpecification
	
	Node
}

type WindowFrame struct {
	StartExpr *WindowFrameExpression
	EndExpr *WindowFrameExpression
	FrameUnit FrameUnit
	
	Node
}

type WindowFrameExpression struct {
	// Expression specifies the boundary as a logical or physical
    // offset to current row. It is present when BoundaryType is
    // OffsetPreceding or OffsetFollowing.
	Expression ExpressionHandler
	BoundaryType BoundaryType
	
	Node
}

type LikeExpression struct {
	LHS ExpressionHandler
	InList *InList
	IsNot bool
	
	Expression
}

type WindowSpecification struct {
	BaseWindowName *Identifier
	PartitionBy *PartitionBy
	OrderBy *OrderBy
	WindowFrame *WindowFrame
	
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
	SampleSize *SampleSize
	SampleSuffix *SampleSuffix
	
	Node
}

type SampleSize struct {
	Size ExpressionHandler
	PartitionBy *PartitionBy
	Unit SampleSizeUnit
	
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
	Format ExpressionHandler
	TimeZoneExpr ExpressionHandler
	
	Node
}

func NewQueryStatement(
	query interface{},
	) (*QueryStatement, error) {
	fmt.Printf("NewQueryStatement(%v)\n",query,
		)
	nn := &QueryStatement{}
	nn.SetKind(QueryStatementKind)
	
	if query == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := query.(NodeHandler); ok {
		nn.Query = query.(*Query)
		nn.AddChild(n)
	} else if w, ok := query.(*Wrapped); ok {
		nn.Query = w.Value.(*Query)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Query = query.(*Query)
	}
	return nn, nil
}



func NewQuery(
	withclause interface{},
	queryexpr interface{},
	orderby interface{},
	limitoffset interface{},
	) (*Query, error) {
	fmt.Printf("NewQuery(%v, %v, %v, %v)\n",withclause,
		queryexpr,
		orderby,
		limitoffset,
		)
	nn := &Query{}
	nn.SetKind(QueryKind)
	
	if withclause != nil {
		if n, ok := withclause.(NodeHandler); ok {
			nn.WithClause = withclause.(*WithClause)
			nn.AddChild(n)
		} else if w, ok := withclause.(*Wrapped); ok {
			nn.WithClause = w.Value.(*WithClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.WithClause = withclause.(*WithClause)
		}
	}
	if queryexpr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := queryexpr.(NodeHandler); ok {
		nn.QueryExpr = queryexpr.(QueryExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := queryexpr.(*Wrapped); ok {
		nn.QueryExpr = w.Value.(QueryExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.QueryExpr = queryexpr.(QueryExpressionHandler)
	}
	if orderby != nil {
		if n, ok := orderby.(NodeHandler); ok {
			nn.OrderBy = orderby.(*OrderBy)
			nn.AddChild(n)
		} else if w, ok := orderby.(*Wrapped); ok {
			nn.OrderBy = w.Value.(*OrderBy)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.OrderBy = orderby.(*OrderBy)
		}
	}
	if limitoffset != nil {
		if n, ok := limitoffset.(NodeHandler); ok {
			nn.LimitOffset = limitoffset.(*LimitOffset)
			nn.AddChild(n)
		} else if w, ok := limitoffset.(*Wrapped); ok {
			nn.LimitOffset = w.Value.(*LimitOffset)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.LimitOffset = limitoffset.(*LimitOffset)
		}
	}
	return nn, nil
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
	fmt.Printf("NewSelect(%v, %v, %v, %v, %v, %v, %v, %v, %v)\n",distinct,
		selectas,
		selectlist,
		fromclause,
		whereclause,
		groupby,
		having,
		qualify,
		windowclause,
		)
	nn := &Select{}
	nn.SetKind(SelectKind)
	
	if distinct != nil {
		if n, ok := distinct.(NodeHandler); ok {
			nn.Distinct = distinct.(bool)
			nn.AddChild(n)
		} else if w, ok := distinct.(*Wrapped); ok {
			nn.Distinct = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Distinct = distinct.(bool)
		}
	}
	if selectas != nil {
		if n, ok := selectas.(NodeHandler); ok {
			nn.SelectAs = selectas.(*SelectAs)
			nn.AddChild(n)
		} else if w, ok := selectas.(*Wrapped); ok {
			nn.SelectAs = w.Value.(*SelectAs)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.SelectAs = selectas.(*SelectAs)
		}
	}
	if selectlist == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := selectlist.(NodeHandler); ok {
		nn.SelectList = selectlist.(*SelectList)
		nn.AddChild(n)
	} else if w, ok := selectlist.(*Wrapped); ok {
		nn.SelectList = w.Value.(*SelectList)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.SelectList = selectlist.(*SelectList)
	}
	if fromclause != nil {
		if n, ok := fromclause.(NodeHandler); ok {
			nn.FromClause = fromclause.(*FromClause)
			nn.AddChild(n)
		} else if w, ok := fromclause.(*Wrapped); ok {
			nn.FromClause = w.Value.(*FromClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.FromClause = fromclause.(*FromClause)
		}
	}
	if whereclause != nil {
		if n, ok := whereclause.(NodeHandler); ok {
			nn.WhereClause = whereclause.(*WhereClause)
			nn.AddChild(n)
		} else if w, ok := whereclause.(*Wrapped); ok {
			nn.WhereClause = w.Value.(*WhereClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.WhereClause = whereclause.(*WhereClause)
		}
	}
	if groupby != nil {
		if n, ok := groupby.(NodeHandler); ok {
			nn.GroupBy = groupby.(*GroupBy)
			nn.AddChild(n)
		} else if w, ok := groupby.(*Wrapped); ok {
			nn.GroupBy = w.Value.(*GroupBy)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.GroupBy = groupby.(*GroupBy)
		}
	}
	if having != nil {
		if n, ok := having.(NodeHandler); ok {
			nn.Having = having.(*Having)
			nn.AddChild(n)
		} else if w, ok := having.(*Wrapped); ok {
			nn.Having = w.Value.(*Having)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Having = having.(*Having)
		}
	}
	if qualify != nil {
		if n, ok := qualify.(NodeHandler); ok {
			nn.Qualify = qualify.(*Qualify)
			nn.AddChild(n)
		} else if w, ok := qualify.(*Wrapped); ok {
			nn.Qualify = w.Value.(*Qualify)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Qualify = qualify.(*Qualify)
		}
	}
	if windowclause != nil {
		if n, ok := windowclause.(NodeHandler); ok {
			nn.WindowClause = windowclause.(*WindowClause)
			nn.AddChild(n)
		} else if w, ok := windowclause.(*Wrapped); ok {
			nn.WindowClause = w.Value.(*WindowClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.WindowClause = windowclause.(*WindowClause)
		}
	}
	return nn, nil
}



func NewSelectList(
	columns interface{},
	) (*SelectList, error) {
	fmt.Printf("NewSelectList(%v)\n",columns,
		)
	nn := &SelectList{}
	nn.SetKind(SelectListKind)
	
	nn.Columns = make([]*SelectColumn, 0, defaultCapacity)
	if n, ok := columns.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := columns.(*Wrapped); ok {
		newElem := w.Value.(*SelectColumn)
		nn.Columns = append(nn.Columns, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := columns.(*SelectColumn)
		nn.Columns = append(nn.Columns, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewSelectColumn(%v, %v)\n",expression,
		alias,
		)
	nn := &SelectColumn{}
	nn.SetKind(SelectColumnKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	if alias != nil {
		if n, ok := alias.(NodeHandler); ok {
			nn.Alias = alias.(*Alias)
			nn.AddChild(n)
		} else if w, ok := alias.(*Wrapped); ok {
			nn.Alias = w.Value.(*Alias)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Alias = alias.(*Alias)
		}
	}
	return nn, nil
}



func NewIntLiteral(
	) (*IntLiteral, error) {
	fmt.Printf("NewIntLiteral()\n",)
	nn := &IntLiteral{}
	nn.SetKind(IntLiteralKind)
	
	return nn, nil
}



func NewIdentifier(
	idstring interface{},
	) (*Identifier, error) {
	fmt.Printf("NewIdentifier(%v)\n",idstring,
		)
	nn := &Identifier{}
	nn.SetKind(IdentifierKind)
	
	if idstring != nil {
		if n, ok := idstring.(NodeHandler); ok {
			nn.IDString = idstring.(string)
			nn.AddChild(n)
		} else if w, ok := idstring.(*Wrapped); ok {
			nn.IDString = w.Value.(string)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IDString = idstring.(string)
		}
	}
	return nn, nil
}



func NewAlias(
	identifier interface{},
	) (*Alias, error) {
	fmt.Printf("NewAlias(%v)\n",identifier,
		)
	nn := &Alias{}
	nn.SetKind(AliasKind)
	
	if identifier == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := identifier.(NodeHandler); ok {
		nn.Identifier = identifier.(*Identifier)
		nn.AddChild(n)
	} else if w, ok := identifier.(*Wrapped); ok {
		nn.Identifier = w.Value.(*Identifier)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Identifier = identifier.(*Identifier)
	}
	return nn, nil
}



func NewPathExpression(
	names interface{},
	) (*PathExpression, error) {
	fmt.Printf("NewPathExpression(%v)\n",names,
		)
	nn := &PathExpression{}
	nn.SetKind(PathExpressionKind)
	
	nn.Names = make([]*Identifier, 0, defaultCapacity)
	if n, ok := names.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := names.(*Wrapped); ok {
		newElem := w.Value.(*Identifier)
		nn.Names = append(nn.Names, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := names.(*Identifier)
		nn.Names = append(nn.Names, newElem)
	}
	return nn, nil
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
	) (*TablePathExpression, error) {
	fmt.Printf("NewTablePathExpression(%v, %v, %v)\n",pathexpr,
		unnestexpr,
		alias,
		)
	nn := &TablePathExpression{}
	nn.SetKind(TablePathExpressionKind)
	
	if pathexpr != nil {
		if n, ok := pathexpr.(NodeHandler); ok {
			nn.PathExpr = pathexpr.(*PathExpression)
			nn.AddChild(n)
		} else if w, ok := pathexpr.(*Wrapped); ok {
			nn.PathExpr = w.Value.(*PathExpression)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.PathExpr = pathexpr.(*PathExpression)
		}
	}
	if unnestexpr != nil {
		if n, ok := unnestexpr.(NodeHandler); ok {
			nn.UnnestExpr = unnestexpr.(*UnnestExpression)
			nn.AddChild(n)
		} else if w, ok := unnestexpr.(*Wrapped); ok {
			nn.UnnestExpr = w.Value.(*UnnestExpression)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.UnnestExpr = unnestexpr.(*UnnestExpression)
		}
	}
	if alias != nil {
		if n, ok := alias.(NodeHandler); ok {
			nn.Alias = alias.(*Alias)
			nn.AddChild(n)
		} else if w, ok := alias.(*Wrapped); ok {
			nn.Alias = w.Value.(*Alias)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Alias = alias.(*Alias)
		}
	}
	return nn, nil
}



func NewFromClause(
	tableexpression interface{},
	) (*FromClause, error) {
	fmt.Printf("NewFromClause(%v)\n",tableexpression,
		)
	nn := &FromClause{}
	nn.SetKind(FromClauseKind)
	
	if tableexpression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := tableexpression.(NodeHandler); ok {
		nn.TableExpression = tableexpression.(TableExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := tableexpression.(*Wrapped); ok {
		nn.TableExpression = w.Value.(TableExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.TableExpression = tableexpression.(TableExpressionHandler)
	}
	return nn, nil
}



func NewWhereClause(
	expression interface{},
	) (*WhereClause, error) {
	fmt.Printf("NewWhereClause(%v)\n",expression,
		)
	nn := &WhereClause{}
	nn.SetKind(WhereClauseKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	return nn, nil
}



func NewBooleanLiteral(
	value interface{},
	) (*BooleanLiteral, error) {
	fmt.Printf("NewBooleanLiteral(%v)\n",value,
		)
	nn := &BooleanLiteral{}
	nn.SetKind(BooleanLiteralKind)
	
	if value != nil {
		if n, ok := value.(NodeHandler); ok {
			nn.Value = value.(bool)
			nn.AddChild(n)
		} else if w, ok := value.(*Wrapped); ok {
			nn.Value = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Value = value.(bool)
		}
	}
	return nn, nil
}



func NewAndExpr(
	conjuncts interface{},
	) (*AndExpr, error) {
	fmt.Printf("NewAndExpr(%v)\n",conjuncts,
		)
	nn := &AndExpr{}
	nn.SetKind(AndExprKind)
	
	nn.Conjuncts = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := conjuncts.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := conjuncts.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Conjuncts = append(nn.Conjuncts, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := conjuncts.(ExpressionHandler)
		nn.Conjuncts = append(nn.Conjuncts, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewBinaryExpression(%v, %v, %v, %v)\n",op,
		lhs,
		rhs,
		isnot,
		)
	nn := &BinaryExpression{}
	nn.SetKind(BinaryExpressionKind)
	
	if op != nil {
		if n, ok := op.(NodeHandler); ok {
			nn.Op = op.(BinaryOp)
			nn.AddChild(n)
		} else if w, ok := op.(*Wrapped); ok {
			nn.Op = w.Value.(BinaryOp)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Op = op.(BinaryOp)
		}
	}
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(ExpressionHandler)
	}
	if rhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := rhs.(NodeHandler); ok {
		nn.RHS = rhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := rhs.(*Wrapped); ok {
		nn.RHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.RHS = rhs.(ExpressionHandler)
	}
	if isnot != nil {
		if n, ok := isnot.(NodeHandler); ok {
			nn.IsNot = isnot.(bool)
			nn.AddChild(n)
		} else if w, ok := isnot.(*Wrapped); ok {
			nn.IsNot = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IsNot = isnot.(bool)
		}
	}
	return nn, nil
}



func NewStringLiteral(
	) (*StringLiteral, error) {
	fmt.Printf("NewStringLiteral()\n",)
	nn := &StringLiteral{}
	nn.SetKind(StringLiteralKind)
	
	return nn, nil
}



func NewStar(
	) (*Star, error) {
	fmt.Printf("NewStar()\n",)
	nn := &Star{}
	nn.SetKind(StarKind)
	
	return nn, nil
}



func NewOrExpr(
	disjuncts interface{},
	) (*OrExpr, error) {
	fmt.Printf("NewOrExpr(%v)\n",disjuncts,
		)
	nn := &OrExpr{}
	nn.SetKind(OrExprKind)
	
	nn.Disjuncts = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := disjuncts.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := disjuncts.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Disjuncts = append(nn.Disjuncts, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := disjuncts.(ExpressionHandler)
		nn.Disjuncts = append(nn.Disjuncts, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewGroupingItem(%v, %v)\n",expression,
		rollup,
		)
	nn := &GroupingItem{}
	nn.SetKind(GroupingItemKind)
	
	if expression != nil {
		if n, ok := expression.(NodeHandler); ok {
			nn.Expression = expression.(ExpressionHandler)
			nn.AddChild(n)
		} else if w, ok := expression.(*Wrapped); ok {
			nn.Expression = w.Value.(ExpressionHandler)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Expression = expression.(ExpressionHandler)
		}
	}
	if rollup != nil {
		if n, ok := rollup.(NodeHandler); ok {
			nn.Rollup = rollup.(*Rollup)
			nn.AddChild(n)
		} else if w, ok := rollup.(*Wrapped); ok {
			nn.Rollup = w.Value.(*Rollup)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Rollup = rollup.(*Rollup)
		}
	}
	return nn, nil
}



func NewGroupBy(
	groupingitems interface{},
	) (*GroupBy, error) {
	fmt.Printf("NewGroupBy(%v)\n",groupingitems,
		)
	nn := &GroupBy{}
	nn.SetKind(GroupByKind)
	
	nn.GroupingItems = make([]*GroupingItem, 0, defaultCapacity)
	if n, ok := groupingitems.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := groupingitems.(*Wrapped); ok {
		newElem := w.Value.(*GroupingItem)
		nn.GroupingItems = append(nn.GroupingItems, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := groupingitems.(*GroupingItem)
		nn.GroupingItems = append(nn.GroupingItems, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewOrderingExpression(%v, %v, %v)\n",expression,
		nullorder,
		orderingspec,
		)
	nn := &OrderingExpression{}
	nn.SetKind(OrderingExpressionKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	if nullorder != nil {
		if n, ok := nullorder.(NodeHandler); ok {
			nn.NullOrder = nullorder.(*NullOrder)
			nn.AddChild(n)
		} else if w, ok := nullorder.(*Wrapped); ok {
			nn.NullOrder = w.Value.(*NullOrder)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.NullOrder = nullorder.(*NullOrder)
		}
	}
	if orderingspec != nil {
		if n, ok := orderingspec.(NodeHandler); ok {
			nn.OrderingSpec = orderingspec.(OrderingSpec)
			nn.AddChild(n)
		} else if w, ok := orderingspec.(*Wrapped); ok {
			nn.OrderingSpec = w.Value.(OrderingSpec)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.OrderingSpec = orderingspec.(OrderingSpec)
		}
	}
	return nn, nil
}



func NewOrderBy(
	orderingexpression interface{},
	) (*OrderBy, error) {
	fmt.Printf("NewOrderBy(%v)\n",orderingexpression,
		)
	nn := &OrderBy{}
	nn.SetKind(OrderByKind)
	
	nn.OrderingExpression = make([]*OrderingExpression, 0, defaultCapacity)
	if n, ok := orderingexpression.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := orderingexpression.(*Wrapped); ok {
		newElem := w.Value.(*OrderingExpression)
		nn.OrderingExpression = append(nn.OrderingExpression, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := orderingexpression.(*OrderingExpression)
		nn.OrderingExpression = append(nn.OrderingExpression, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewLimitOffset(%v, %v)\n",limit,
		offset,
		)
	nn := &LimitOffset{}
	nn.SetKind(LimitOffsetKind)
	
	if limit == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := limit.(NodeHandler); ok {
		nn.Limit = limit.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := limit.(*Wrapped); ok {
		nn.Limit = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Limit = limit.(ExpressionHandler)
	}
	if offset != nil {
		if n, ok := offset.(NodeHandler); ok {
			nn.Offset = offset.(ExpressionHandler)
			nn.AddChild(n)
		} else if w, ok := offset.(*Wrapped); ok {
			nn.Offset = w.Value.(ExpressionHandler)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Offset = offset.(ExpressionHandler)
		}
	}
	return nn, nil
}



func NewFloatLiteral(
	) (*FloatLiteral, error) {
	fmt.Printf("NewFloatLiteral()\n",)
	nn := &FloatLiteral{}
	nn.SetKind(FloatLiteralKind)
	
	return nn, nil
}



func NewNullLiteral(
	) (*NullLiteral, error) {
	fmt.Printf("NewNullLiteral()\n",)
	nn := &NullLiteral{}
	nn.SetKind(NullLiteralKind)
	
	return nn, nil
}



func NewOnClause(
	expression interface{},
	) (*OnClause, error) {
	fmt.Printf("NewOnClause(%v)\n",expression,
		)
	nn := &OnClause{}
	nn.SetKind(OnClauseKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	return nn, nil
}



func NewWithClauseEntry(
	alias interface{},
	query interface{},
	) (*WithClauseEntry, error) {
	fmt.Printf("NewWithClauseEntry(%v, %v)\n",alias,
		query,
		)
	nn := &WithClauseEntry{}
	nn.SetKind(WithClauseEntryKind)
	
	if alias == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := alias.(NodeHandler); ok {
		nn.Alias = alias.(*Identifier)
		nn.AddChild(n)
	} else if w, ok := alias.(*Wrapped); ok {
		nn.Alias = w.Value.(*Identifier)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Alias = alias.(*Identifier)
	}
	if query == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := query.(NodeHandler); ok {
		nn.Query = query.(*Query)
		nn.AddChild(n)
	} else if w, ok := query.(*Wrapped); ok {
		nn.Query = w.Value.(*Query)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Query = query.(*Query)
	}
	return nn, nil
}



func NewJoin(
	lhs interface{},
	rhs interface{},
	clauselist interface{},
	jointype interface{},
	) (*Join, error) {
	fmt.Printf("NewJoin(%v, %v, %v, %v)\n",lhs,
		rhs,
		clauselist,
		jointype,
		)
	nn := &Join{}
	nn.SetKind(JoinKind)
	
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(TableExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(TableExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(TableExpressionHandler)
	}
	if rhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := rhs.(NodeHandler); ok {
		nn.RHS = rhs.(TableExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := rhs.(*Wrapped); ok {
		nn.RHS = w.Value.(TableExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.RHS = rhs.(TableExpressionHandler)
	}
	if clauselist != nil {
		if n, ok := clauselist.(NodeHandler); ok {
			nn.ClauseList = clauselist.(*OnOrUsingClauseList)
			nn.AddChild(n)
		} else if w, ok := clauselist.(*Wrapped); ok {
			nn.ClauseList = w.Value.(*OnOrUsingClauseList)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.ClauseList = clauselist.(*OnOrUsingClauseList)
		}
	}
	if jointype != nil {
		if n, ok := jointype.(NodeHandler); ok {
			nn.JoinType = jointype.(JoinType)
			nn.AddChild(n)
		} else if w, ok := jointype.(*Wrapped); ok {
			nn.JoinType = w.Value.(JoinType)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.JoinType = jointype.(JoinType)
		}
	}
	return nn, nil
}



func NewWithClause(
	with interface{},
	) (*WithClause, error) {
	fmt.Printf("NewWithClause(%v)\n",with,
		)
	nn := &WithClause{}
	nn.SetKind(WithClauseKind)
	
	nn.With = make([]*WithClauseEntry, 0, defaultCapacity)
	if n, ok := with.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := with.(*Wrapped); ok {
		newElem := w.Value.(*WithClauseEntry)
		nn.With = append(nn.With, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := with.(*WithClauseEntry)
		nn.With = append(nn.With, newElem)
	}
	return nn, nil
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
	) (*Having, error) {
	fmt.Printf("NewHaving()\n",)
	nn := &Having{}
	nn.SetKind(HavingKind)
	
	return nn, nil
}



func NewNamedType(
	typename interface{},
	typeparameters interface{},
	) (*NamedType, error) {
	fmt.Printf("NewNamedType(%v, %v)\n",typename,
		typeparameters,
		)
	nn := &NamedType{}
	nn.SetKind(NamedTypeKind)
	
	if typename == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := typename.(NodeHandler); ok {
		nn.TypeName = typename.(*PathExpression)
		nn.AddChild(n)
	} else if w, ok := typename.(*Wrapped); ok {
		nn.TypeName = w.Value.(*PathExpression)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.TypeName = typename.(*PathExpression)
	}
	if typeparameters != nil {
		if n, ok := typeparameters.(NodeHandler); ok {
			nn.TypeParameters = typeparameters.(*TypeParameterList)
			nn.AddChild(n)
		} else if w, ok := typeparameters.(*Wrapped); ok {
			nn.TypeParameters = w.Value.(*TypeParameterList)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TypeParameters = typeparameters.(*TypeParameterList)
		}
	}
	return nn, nil
}



func NewArrayType(
	elementtype interface{},
	typeparameters interface{},
	) (*ArrayType, error) {
	fmt.Printf("NewArrayType(%v, %v)\n",elementtype,
		typeparameters,
		)
	nn := &ArrayType{}
	nn.SetKind(ArrayTypeKind)
	
	if elementtype == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := elementtype.(NodeHandler); ok {
		nn.ElementType = elementtype.(TypeHandler)
		nn.AddChild(n)
	} else if w, ok := elementtype.(*Wrapped); ok {
		nn.ElementType = w.Value.(TypeHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.ElementType = elementtype.(TypeHandler)
	}
	if typeparameters != nil {
		if n, ok := typeparameters.(NodeHandler); ok {
			nn.TypeParameters = typeparameters.(*TypeParameterList)
			nn.AddChild(n)
		} else if w, ok := typeparameters.(*Wrapped); ok {
			nn.TypeParameters = w.Value.(*TypeParameterList)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TypeParameters = typeparameters.(*TypeParameterList)
		}
	}
	return nn, nil
}



func NewStructField(
	name interface{},
	) (*StructField, error) {
	fmt.Printf("NewStructField(%v)\n",name,
		)
	nn := &StructField{}
	nn.SetKind(StructFieldKind)
	
	if name != nil {
		if n, ok := name.(NodeHandler); ok {
			nn.Name = name.(*Identifier)
			nn.AddChild(n)
		} else if w, ok := name.(*Wrapped); ok {
			nn.Name = w.Value.(*Identifier)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Name = name.(*Identifier)
		}
	}
	return nn, nil
}



func NewStructType(
	structfields interface{},
	typeparameterlist interface{},
	) (*StructType, error) {
	fmt.Printf("NewStructType(%v, %v)\n",structfields,
		typeparameterlist,
		)
	nn := &StructType{}
	nn.SetKind(StructTypeKind)
	
	if structfields == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := structfields.(NodeHandler); ok {
		nn.StructFields = structfields.(*StructField)
		nn.AddChild(n)
	} else if w, ok := structfields.(*Wrapped); ok {
		nn.StructFields = w.Value.(*StructField)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.StructFields = structfields.(*StructField)
	}
	if typeparameterlist != nil {
		if n, ok := typeparameterlist.(NodeHandler); ok {
			nn.TypeParameterList = typeparameterlist.(*TypeParameterList)
			nn.AddChild(n)
		} else if w, ok := typeparameterlist.(*Wrapped); ok {
			nn.TypeParameterList = w.Value.(*TypeParameterList)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TypeParameterList = typeparameterlist.(*TypeParameterList)
		}
	}
	return nn, nil
}



func NewCastExpression(
	expr interface{},
	typ interface{},
	format interface{},
	issafecast interface{},
	) (*CastExpression, error) {
	fmt.Printf("NewCastExpression(%v, %v, %v, %v)\n",expr,
		typ,
		format,
		issafecast,
		)
	nn := &CastExpression{}
	nn.SetKind(CastExpressionKind)
	
	if expr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expr.(NodeHandler); ok {
		nn.Expr = expr.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expr.(*Wrapped); ok {
		nn.Expr = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expr = expr.(ExpressionHandler)
	}
	if typ == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := typ.(NodeHandler); ok {
		nn.Type = typ.(TypeHandler)
		nn.AddChild(n)
	} else if w, ok := typ.(*Wrapped); ok {
		nn.Type = w.Value.(TypeHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Type = typ.(TypeHandler)
	}
	if format != nil {
		if n, ok := format.(NodeHandler); ok {
			nn.Format = format.(*FormatClause)
			nn.AddChild(n)
		} else if w, ok := format.(*Wrapped); ok {
			nn.Format = w.Value.(*FormatClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Format = format.(*FormatClause)
		}
	}
	if issafecast != nil {
		if n, ok := issafecast.(NodeHandler); ok {
			nn.IsSafeCast = issafecast.(bool)
			nn.AddChild(n)
		} else if w, ok := issafecast.(*Wrapped); ok {
			nn.IsSafeCast = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IsSafeCast = issafecast.(bool)
		}
	}
	return nn, nil
}



func NewSelectAs(
	typename interface{},
	asmode interface{},
	) (*SelectAs, error) {
	fmt.Printf("NewSelectAs(%v, %v)\n",typename,
		asmode,
		)
	nn := &SelectAs{}
	nn.SetKind(SelectAsKind)
	
	if typename != nil {
		if n, ok := typename.(NodeHandler); ok {
			nn.TypeName = typename.(*PathExpression)
			nn.AddChild(n)
		} else if w, ok := typename.(*Wrapped); ok {
			nn.TypeName = w.Value.(*PathExpression)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TypeName = typename.(*PathExpression)
		}
	}
	if asmode != nil {
		if n, ok := asmode.(NodeHandler); ok {
			nn.AsMode = asmode.(AsMode)
			nn.AddChild(n)
		} else if w, ok := asmode.(*Wrapped); ok {
			nn.AsMode = w.Value.(AsMode)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.AsMode = asmode.(AsMode)
		}
	}
	return nn, nil
}



func NewRollup(
	expressions interface{},
	) (*Rollup, error) {
	fmt.Printf("NewRollup(%v)\n",expressions,
		)
	nn := &Rollup{}
	nn.SetKind(RollupKind)
	
	nn.Expressions = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := expressions.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := expressions.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Expressions = append(nn.Expressions, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := expressions.(ExpressionHandler)
		nn.Expressions = append(nn.Expressions, newElem)
	}
	return nn, nil
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
	arguments interface{},
	orderby interface{},
	limitoffset interface{},
	nullhandlingmodifier interface{},
	distinct interface{},
	) (*FunctionCall, error) {
	fmt.Printf("NewFunctionCall(%v, %v, %v, %v, %v, %v)\n",function,
		arguments,
		orderby,
		limitoffset,
		nullhandlingmodifier,
		distinct,
		)
	nn := &FunctionCall{}
	nn.SetKind(FunctionCallKind)
	
	if function == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := function.(NodeHandler); ok {
		nn.Function = function.(*PathExpression)
		nn.AddChild(n)
	} else if w, ok := function.(*Wrapped); ok {
		nn.Function = w.Value.(*PathExpression)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Function = function.(*PathExpression)
	}
	nn.Arguments = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := arguments.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := arguments.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Arguments = append(nn.Arguments, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := arguments.(ExpressionHandler)
		nn.Arguments = append(nn.Arguments, newElem)
	}
	if orderby != nil {
		if n, ok := orderby.(NodeHandler); ok {
			nn.OrderBy = orderby.(*OrderBy)
			nn.AddChild(n)
		} else if w, ok := orderby.(*Wrapped); ok {
			nn.OrderBy = w.Value.(*OrderBy)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.OrderBy = orderby.(*OrderBy)
		}
	}
	if limitoffset != nil {
		if n, ok := limitoffset.(NodeHandler); ok {
			nn.LimitOffset = limitoffset.(*LimitOffset)
			nn.AddChild(n)
		} else if w, ok := limitoffset.(*Wrapped); ok {
			nn.LimitOffset = w.Value.(*LimitOffset)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.LimitOffset = limitoffset.(*LimitOffset)
		}
	}
	if nullhandlingmodifier != nil {
		if n, ok := nullhandlingmodifier.(NodeHandler); ok {
			nn.NullHandlingModifier = nullhandlingmodifier.(NullHandlingModifier)
			nn.AddChild(n)
		} else if w, ok := nullhandlingmodifier.(*Wrapped); ok {
			nn.NullHandlingModifier = w.Value.(NullHandlingModifier)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.NullHandlingModifier = nullhandlingmodifier.(NullHandlingModifier)
		}
	}
	if distinct != nil {
		if n, ok := distinct.(NodeHandler); ok {
			nn.Distinct = distinct.(bool)
			nn.AddChild(n)
		} else if w, ok := distinct.(*Wrapped); ok {
			nn.Distinct = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Distinct = distinct.(bool)
		}
	}
	return nn, nil
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
	elements interface{},
	) (*ArrayConstructor, error) {
	fmt.Printf("NewArrayConstructor(%v, %v)\n",typ,
		elements,
		)
	nn := &ArrayConstructor{}
	nn.SetKind(ArrayConstructorKind)
	
	if typ != nil {
		if n, ok := typ.(NodeHandler); ok {
			nn.Type = typ.(*ArrayType)
			nn.AddChild(n)
		} else if w, ok := typ.(*Wrapped); ok {
			nn.Type = w.Value.(*ArrayType)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Type = typ.(*ArrayType)
		}
	}
	nn.Elements = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := elements.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := elements.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Elements = append(nn.Elements, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := elements.(ExpressionHandler)
		nn.Elements = append(nn.Elements, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewStructConstructorArg(%v, %v)\n",expression,
		alias,
		)
	nn := &StructConstructorArg{}
	nn.SetKind(StructConstructorArgKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	if alias == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := alias.(NodeHandler); ok {
		nn.Alias = alias.(*Alias)
		nn.AddChild(n)
	} else if w, ok := alias.(*Wrapped); ok {
		nn.Alias = w.Value.(*Alias)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Alias = alias.(*Alias)
	}
	return nn, nil
}



func NewStructConstructorWithParens(
	fieldexpressions interface{},
	) (*StructConstructorWithParens, error) {
	fmt.Printf("NewStructConstructorWithParens(%v)\n",fieldexpressions,
		)
	nn := &StructConstructorWithParens{}
	nn.SetKind(StructConstructorWithParensKind)
	
	nn.FieldExpressions = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := fieldexpressions.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := fieldexpressions.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.FieldExpressions = append(nn.FieldExpressions, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := fieldexpressions.(ExpressionHandler)
		nn.FieldExpressions = append(nn.FieldExpressions, newElem)
	}
	return nn, nil
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
	fields interface{},
	) (*StructConstructorWithKeyword, error) {
	fmt.Printf("NewStructConstructorWithKeyword(%v, %v)\n",structtype,
		fields,
		)
	nn := &StructConstructorWithKeyword{}
	nn.SetKind(StructConstructorWithKeywordKind)
	
	if structtype != nil {
		if n, ok := structtype.(NodeHandler); ok {
			nn.StructType = structtype.(*StructType)
			nn.AddChild(n)
		} else if w, ok := structtype.(*Wrapped); ok {
			nn.StructType = w.Value.(*StructType)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.StructType = structtype.(*StructType)
		}
	}
	nn.Fields = make([]*StructConstructorArg, 0, defaultCapacity)
	if n, ok := fields.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := fields.(*Wrapped); ok {
		newElem := w.Value.(*StructConstructorArg)
		nn.Fields = append(nn.Fields, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := fields.(*StructConstructorArg)
		nn.Fields = append(nn.Fields, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewInExpression(%v, %v, %v, %v, %v)\n",lhs,
		inlist,
		query,
		unnestexpr,
		isnot,
		)
	nn := &InExpression{}
	nn.SetKind(InExpressionKind)
	
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(ExpressionHandler)
	}
	if inlist != nil {
		if n, ok := inlist.(NodeHandler); ok {
			nn.InList = inlist.(*InList)
			nn.AddChild(n)
		} else if w, ok := inlist.(*Wrapped); ok {
			nn.InList = w.Value.(*InList)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.InList = inlist.(*InList)
		}
	}
	if query != nil {
		if n, ok := query.(NodeHandler); ok {
			nn.Query = query.(*Query)
			nn.AddChild(n)
		} else if w, ok := query.(*Wrapped); ok {
			nn.Query = w.Value.(*Query)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Query = query.(*Query)
		}
	}
	if unnestexpr != nil {
		if n, ok := unnestexpr.(NodeHandler); ok {
			nn.UnnestExpr = unnestexpr.(*UnnestExpression)
			nn.AddChild(n)
		} else if w, ok := unnestexpr.(*Wrapped); ok {
			nn.UnnestExpr = w.Value.(*UnnestExpression)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.UnnestExpr = unnestexpr.(*UnnestExpression)
		}
	}
	if isnot != nil {
		if n, ok := isnot.(NodeHandler); ok {
			nn.IsNot = isnot.(bool)
			nn.AddChild(n)
		} else if w, ok := isnot.(*Wrapped); ok {
			nn.IsNot = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IsNot = isnot.(bool)
		}
	}
	return nn, nil
}



func NewInList(
	list interface{},
	) (*InList, error) {
	fmt.Printf("NewInList(%v)\n",list,
		)
	nn := &InList{}
	nn.SetKind(InListKind)
	
	nn.List = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := list.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := list.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.List = append(nn.List, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := list.(ExpressionHandler)
		nn.List = append(nn.List, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewBetweenExpression(%v, %v, %v, %v)\n",lhs,
		low,
		high,
		isnot,
		)
	nn := &BetweenExpression{}
	nn.SetKind(BetweenExpressionKind)
	
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(ExpressionHandler)
	}
	if low == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := low.(NodeHandler); ok {
		nn.Low = low.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := low.(*Wrapped); ok {
		nn.Low = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Low = low.(ExpressionHandler)
	}
	if high == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := high.(NodeHandler); ok {
		nn.High = high.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := high.(*Wrapped); ok {
		nn.High = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.High = high.(ExpressionHandler)
	}
	if isnot != nil {
		if n, ok := isnot.(NodeHandler); ok {
			nn.IsNot = isnot.(bool)
			nn.AddChild(n)
		} else if w, ok := isnot.(*Wrapped); ok {
			nn.IsNot = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IsNot = isnot.(bool)
		}
	}
	return nn, nil
}



func NewNumericLiteral(
	) (*NumericLiteral, error) {
	fmt.Printf("NewNumericLiteral()\n",)
	nn := &NumericLiteral{}
	nn.SetKind(NumericLiteralKind)
	
	return nn, nil
}



func NewBigNumericLiteral(
	) (*BigNumericLiteral, error) {
	fmt.Printf("NewBigNumericLiteral()\n",)
	nn := &BigNumericLiteral{}
	nn.SetKind(BigNumericLiteralKind)
	
	return nn, nil
}



func NewBytesLiteral(
	) (*BytesLiteral, error) {
	fmt.Printf("NewBytesLiteral()\n",)
	nn := &BytesLiteral{}
	nn.SetKind(BytesLiteralKind)
	
	return nn, nil
}



func NewDateOrTimeLiteral(
	stringliteral interface{},
	typekind interface{},
	) (*DateOrTimeLiteral, error) {
	fmt.Printf("NewDateOrTimeLiteral(%v, %v)\n",stringliteral,
		typekind,
		)
	nn := &DateOrTimeLiteral{}
	nn.SetKind(DateOrTimeLiteralKind)
	
	if stringliteral == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := stringliteral.(NodeHandler); ok {
		nn.StringLiteral = stringliteral.(*StringLiteral)
		nn.AddChild(n)
	} else if w, ok := stringliteral.(*Wrapped); ok {
		nn.StringLiteral = w.Value.(*StringLiteral)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.StringLiteral = stringliteral.(*StringLiteral)
	}
	if typekind != nil {
		if n, ok := typekind.(NodeHandler); ok {
			nn.TypeKind = typekind.(TypeKind)
			nn.AddChild(n)
		} else if w, ok := typekind.(*Wrapped); ok {
			nn.TypeKind = w.Value.(TypeKind)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TypeKind = typekind.(TypeKind)
		}
	}
	return nn, nil
}



func NewCaseValueExpression(
	arguments interface{},
	) (*CaseValueExpression, error) {
	fmt.Printf("NewCaseValueExpression(%v)\n",arguments,
		)
	nn := &CaseValueExpression{}
	nn.SetKind(CaseValueExpressionKind)
	
	nn.Arguments = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := arguments.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := arguments.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Arguments = append(nn.Arguments, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := arguments.(ExpressionHandler)
		nn.Arguments = append(nn.Arguments, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewCaseNoValueExpression(%v)\n",arguments,
		)
	nn := &CaseNoValueExpression{}
	nn.SetKind(CaseNoValueExpressionKind)
	
	nn.Arguments = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := arguments.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := arguments.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.Arguments = append(nn.Arguments, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := arguments.(ExpressionHandler)
		nn.Arguments = append(nn.Arguments, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewArrayElement(%v, %v)\n",array,
		position,
		)
	nn := &ArrayElement{}
	nn.SetKind(ArrayElementKind)
	
	if array == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := array.(NodeHandler); ok {
		nn.Array = array.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := array.(*Wrapped); ok {
		nn.Array = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Array = array.(ExpressionHandler)
	}
	if position == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := position.(NodeHandler); ok {
		nn.Position = position.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := position.(*Wrapped); ok {
		nn.Position = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Position = position.(ExpressionHandler)
	}
	return nn, nil
}



func NewBitwiseShiftExpression(
	lhs interface{},
	rhs interface{},
	isleftshift interface{},
	) (*BitwiseShiftExpression, error) {
	fmt.Printf("NewBitwiseShiftExpression(%v, %v, %v)\n",lhs,
		rhs,
		isleftshift,
		)
	nn := &BitwiseShiftExpression{}
	nn.SetKind(BitwiseShiftExpressionKind)
	
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(ExpressionHandler)
	}
	if rhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := rhs.(NodeHandler); ok {
		nn.RHS = rhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := rhs.(*Wrapped); ok {
		nn.RHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.RHS = rhs.(ExpressionHandler)
	}
	if isleftshift != nil {
		if n, ok := isleftshift.(NodeHandler); ok {
			nn.IsLeftShift = isleftshift.(bool)
			nn.AddChild(n)
		} else if w, ok := isleftshift.(*Wrapped); ok {
			nn.IsLeftShift = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IsLeftShift = isleftshift.(bool)
		}
	}
	return nn, nil
}



func NewDotGeneralizedField(
	expr interface{},
	path interface{},
	) (*DotGeneralizedField, error) {
	fmt.Printf("NewDotGeneralizedField(%v, %v)\n",expr,
		path,
		)
	nn := &DotGeneralizedField{}
	nn.SetKind(DotGeneralizedFieldKind)
	
	if expr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expr.(NodeHandler); ok {
		nn.Expr = expr.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expr.(*Wrapped); ok {
		nn.Expr = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expr = expr.(ExpressionHandler)
	}
	if path == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := path.(NodeHandler); ok {
		nn.Path = path.(*PathExpression)
		nn.AddChild(n)
	} else if w, ok := path.(*Wrapped); ok {
		nn.Path = w.Value.(*PathExpression)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Path = path.(*PathExpression)
	}
	return nn, nil
}



func NewDotIdentifier(
	expr interface{},
	path interface{},
	) (*DotIdentifier, error) {
	fmt.Printf("NewDotIdentifier(%v, %v)\n",expr,
		path,
		)
	nn := &DotIdentifier{}
	nn.SetKind(DotIdentifierKind)
	
	if expr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expr.(NodeHandler); ok {
		nn.Expr = expr.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expr.(*Wrapped); ok {
		nn.Expr = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expr = expr.(ExpressionHandler)
	}
	if path == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := path.(NodeHandler); ok {
		nn.Path = path.(*PathExpression)
		nn.AddChild(n)
	} else if w, ok := path.(*Wrapped); ok {
		nn.Path = w.Value.(*PathExpression)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Path = path.(*PathExpression)
	}
	return nn, nil
}



func NewDotStar(
	expr interface{},
	) (*DotStar, error) {
	fmt.Printf("NewDotStar(%v)\n",expr,
		)
	nn := &DotStar{}
	nn.SetKind(DotStarKind)
	
	if expr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expr.(NodeHandler); ok {
		nn.Expr = expr.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expr.(*Wrapped); ok {
		nn.Expr = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expr = expr.(ExpressionHandler)
	}
	return nn, nil
}



func NewDotStarWithModifiers(
	expr interface{},
	modifiers interface{},
	) (*DotStarWithModifiers, error) {
	fmt.Printf("NewDotStarWithModifiers(%v, %v)\n",expr,
		modifiers,
		)
	nn := &DotStarWithModifiers{}
	nn.SetKind(DotStarWithModifiersKind)
	
	if expr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expr.(NodeHandler); ok {
		nn.Expr = expr.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expr.(*Wrapped); ok {
		nn.Expr = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expr = expr.(ExpressionHandler)
	}
	if modifiers == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := modifiers.(NodeHandler); ok {
		nn.Modifiers = modifiers.(*StarModifiers)
		nn.AddChild(n)
	} else if w, ok := modifiers.(*Wrapped); ok {
		nn.Modifiers = w.Value.(*StarModifiers)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Modifiers = modifiers.(*StarModifiers)
	}
	return nn, nil
}



func NewExpressionSubquery(
	query interface{},
	modifier interface{},
	) (*ExpressionSubquery, error) {
	fmt.Printf("NewExpressionSubquery(%v, %v)\n",query,
		modifier,
		)
	nn := &ExpressionSubquery{}
	nn.SetKind(ExpressionSubqueryKind)
	
	if query == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := query.(NodeHandler); ok {
		nn.Query = query.(*Query)
		nn.AddChild(n)
	} else if w, ok := query.(*Wrapped); ok {
		nn.Query = w.Value.(*Query)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Query = query.(*Query)
	}
	if modifier == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := modifier.(NodeHandler); ok {
		nn.Modifier = modifier.(SubqueryModifier)
		nn.AddChild(n)
	} else if w, ok := modifier.(*Wrapped); ok {
		nn.Modifier = w.Value.(SubqueryModifier)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Modifier = modifier.(SubqueryModifier)
	}
	return nn, nil
}



func NewExtractExpression(
	lhs interface{},
	rhs interface{},
	timezone interface{},
	) (*ExtractExpression, error) {
	fmt.Printf("NewExtractExpression(%v, %v, %v)\n",lhs,
		rhs,
		timezone,
		)
	nn := &ExtractExpression{}
	nn.SetKind(ExtractExpressionKind)
	
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(ExpressionHandler)
	}
	if rhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := rhs.(NodeHandler); ok {
		nn.RHS = rhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := rhs.(*Wrapped); ok {
		nn.RHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.RHS = rhs.(ExpressionHandler)
	}
	if timezone != nil {
		if n, ok := timezone.(NodeHandler); ok {
			nn.TimeZone = timezone.(ExpressionHandler)
			nn.AddChild(n)
		} else if w, ok := timezone.(*Wrapped); ok {
			nn.TimeZone = w.Value.(ExpressionHandler)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TimeZone = timezone.(ExpressionHandler)
		}
	}
	return nn, nil
}



func NewIntervalExpression(
	intervalvalue interface{},
	datepartname interface{},
	datepartnameto interface{},
	) (*IntervalExpression, error) {
	fmt.Printf("NewIntervalExpression(%v, %v, %v)\n",intervalvalue,
		datepartname,
		datepartnameto,
		)
	nn := &IntervalExpression{}
	nn.SetKind(IntervalExpressionKind)
	
	if intervalvalue == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := intervalvalue.(NodeHandler); ok {
		nn.IntervalValue = intervalvalue.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := intervalvalue.(*Wrapped); ok {
		nn.IntervalValue = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.IntervalValue = intervalvalue.(ExpressionHandler)
	}
	if datepartname == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := datepartname.(NodeHandler); ok {
		nn.DatePartName = datepartname.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := datepartname.(*Wrapped); ok {
		nn.DatePartName = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.DatePartName = datepartname.(ExpressionHandler)
	}
	if datepartnameto != nil {
		if n, ok := datepartnameto.(NodeHandler); ok {
			nn.DatePartNameTo = datepartnameto.(ExpressionHandler)
			nn.AddChild(n)
		} else if w, ok := datepartnameto.(*Wrapped); ok {
			nn.DatePartNameTo = w.Value.(ExpressionHandler)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.DatePartNameTo = datepartnameto.(ExpressionHandler)
		}
	}
	return nn, nil
}



func NewNullOrder(
	nullsfirst interface{},
	) (*NullOrder, error) {
	fmt.Printf("NewNullOrder(%v)\n",nullsfirst,
		)
	nn := &NullOrder{}
	nn.SetKind(NullOrderKind)
	
	if nullsfirst != nil {
		if n, ok := nullsfirst.(NodeHandler); ok {
			nn.NullsFirst = nullsfirst.(bool)
			nn.AddChild(n)
		} else if w, ok := nullsfirst.(*Wrapped); ok {
			nn.NullsFirst = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.NullsFirst = nullsfirst.(bool)
		}
	}
	return nn, nil
}



func NewOnOrUsingClauseList(
	list interface{},
	) (*OnOrUsingClauseList, error) {
	fmt.Printf("NewOnOrUsingClauseList(%v)\n",list,
		)
	nn := &OnOrUsingClauseList{}
	nn.SetKind(OnOrUsingClauseListKind)
	
	nn.List = make([]NodeHandler, 0, defaultCapacity)
	if n, ok := list.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := list.(*Wrapped); ok {
		newElem := w.Value.(NodeHandler)
		nn.List = append(nn.List, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := list.(NodeHandler)
		nn.List = append(nn.List, newElem)
	}
	return nn, nil
}


func (n *OnOrUsingClauseList) AddChild(c NodeHandler) {
	n.List = append(n.List, c.(NodeHandler))
	n.Node.AddChild(c)
	c.SetParent(n)
}

func (n *OnOrUsingClauseList) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c == nil {
			continue
		}
		n.List = append(n.List, c.(NodeHandler))
		n.Node.AddChild(c)
		c.SetParent(n)
	}
}

func NewParenthesizedJoin(
	join interface{},
	sampleclause interface{},
	) (*ParenthesizedJoin, error) {
	fmt.Printf("NewParenthesizedJoin(%v, %v)\n",join,
		sampleclause,
		)
	nn := &ParenthesizedJoin{}
	nn.SetKind(ParenthesizedJoinKind)
	
	if join == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := join.(NodeHandler); ok {
		nn.Join = join.(*Join)
		nn.AddChild(n)
	} else if w, ok := join.(*Wrapped); ok {
		nn.Join = w.Value.(*Join)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Join = join.(*Join)
	}
	if sampleclause != nil {
		if n, ok := sampleclause.(NodeHandler); ok {
			nn.SampleClause = sampleclause.(SampleClause)
			nn.AddChild(n)
		} else if w, ok := sampleclause.(*Wrapped); ok {
			nn.SampleClause = w.Value.(SampleClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.SampleClause = sampleclause.(SampleClause)
		}
	}
	return nn, nil
}



func NewPartitionBy(
	partitioningexpressions interface{},
	) (*PartitionBy, error) {
	fmt.Printf("NewPartitionBy(%v)\n",partitioningexpressions,
		)
	nn := &PartitionBy{}
	nn.SetKind(PartitionByKind)
	
	nn.PartitioningExpressions = make([]ExpressionHandler, 0, defaultCapacity)
	if n, ok := partitioningexpressions.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := partitioningexpressions.(*Wrapped); ok {
		newElem := w.Value.(ExpressionHandler)
		nn.PartitioningExpressions = append(nn.PartitioningExpressions, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := partitioningexpressions.(ExpressionHandler)
		nn.PartitioningExpressions = append(nn.PartitioningExpressions, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewSetOperation(%v, %v, %v)\n",inputs,
		optype,
		distinct,
		)
	nn := &SetOperation{}
	nn.SetKind(SetOperationKind)
	
	nn.Inputs = make([]QueryExpressionHandler, 0, defaultCapacity)
	if n, ok := inputs.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := inputs.(*Wrapped); ok {
		newElem := w.Value.(QueryExpressionHandler)
		nn.Inputs = append(nn.Inputs, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := inputs.(QueryExpressionHandler)
		nn.Inputs = append(nn.Inputs, newElem)
	}
	if optype != nil {
		if n, ok := optype.(NodeHandler); ok {
			nn.OpType = optype.(SetOp)
			nn.AddChild(n)
		} else if w, ok := optype.(*Wrapped); ok {
			nn.OpType = w.Value.(SetOp)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.OpType = optype.(SetOp)
		}
	}
	if distinct != nil {
		if n, ok := distinct.(NodeHandler); ok {
			nn.Distinct = distinct.(bool)
			nn.AddChild(n)
		} else if w, ok := distinct.(*Wrapped); ok {
			nn.Distinct = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Distinct = distinct.(bool)
		}
	}
	return nn, nil
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
	fmt.Printf("NewStarExceptList(%v)\n",identifiers,
		)
	nn := &StarExceptList{}
	nn.SetKind(StarExceptListKind)
	
	nn.Identifiers = make([]*Identifier, 0, defaultCapacity)
	if n, ok := identifiers.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := identifiers.(*Wrapped); ok {
		newElem := w.Value.(*Identifier)
		nn.Identifiers = append(nn.Identifiers, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := identifiers.(*Identifier)
		nn.Identifiers = append(nn.Identifiers, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewStarModifiers(%v, %v)\n",exceptlist,
		replaceitems,
		)
	nn := &StarModifiers{}
	nn.SetKind(StarModifiersKind)
	
	if exceptlist != nil {
		if n, ok := exceptlist.(NodeHandler); ok {
			nn.ExceptList = exceptlist.(*StarExceptList)
			nn.AddChild(n)
		} else if w, ok := exceptlist.(*Wrapped); ok {
			nn.ExceptList = w.Value.(*StarExceptList)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.ExceptList = exceptlist.(*StarExceptList)
		}
	}
	nn.ReplaceItems = make([]*StarReplaceItem, 0, defaultCapacity)
	if n, ok := replaceitems.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := replaceitems.(*Wrapped); ok {
		newElem := w.Value.(*StarReplaceItem)
		nn.ReplaceItems = append(nn.ReplaceItems, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := replaceitems.(*StarReplaceItem)
		nn.ReplaceItems = append(nn.ReplaceItems, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewStarReplaceItem(%v, %v)\n",expression,
		alias,
		)
	nn := &StarReplaceItem{}
	nn.SetKind(StarReplaceItemKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	if alias == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := alias.(NodeHandler); ok {
		nn.Alias = alias.(*Identifier)
		nn.AddChild(n)
	} else if w, ok := alias.(*Wrapped); ok {
		nn.Alias = w.Value.(*Identifier)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Alias = alias.(*Identifier)
	}
	return nn, nil
}



func NewStarWithModifiers(
	modifiers interface{},
	) (*StarWithModifiers, error) {
	fmt.Printf("NewStarWithModifiers(%v)\n",modifiers,
		)
	nn := &StarWithModifiers{}
	nn.SetKind(StarWithModifiersKind)
	
	if modifiers == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := modifiers.(NodeHandler); ok {
		nn.Modifiers = modifiers.(*StarModifiers)
		nn.AddChild(n)
	} else if w, ok := modifiers.(*Wrapped); ok {
		nn.Modifiers = w.Value.(*StarModifiers)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Modifiers = modifiers.(*StarModifiers)
	}
	return nn, nil
}



func NewTableSubquery(
	subquery interface{},
	alias interface{},
	sampleclause interface{},
	) (*TableSubquery, error) {
	fmt.Printf("NewTableSubquery(%v, %v, %v)\n",subquery,
		alias,
		sampleclause,
		)
	nn := &TableSubquery{}
	nn.SetKind(TableSubqueryKind)
	
	if subquery == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := subquery.(NodeHandler); ok {
		nn.Subquery = subquery.(*Query)
		nn.AddChild(n)
	} else if w, ok := subquery.(*Wrapped); ok {
		nn.Subquery = w.Value.(*Query)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Subquery = subquery.(*Query)
	}
	if alias != nil {
		if n, ok := alias.(NodeHandler); ok {
			nn.Alias = alias.(*Alias)
			nn.AddChild(n)
		} else if w, ok := alias.(*Wrapped); ok {
			nn.Alias = w.Value.(*Alias)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Alias = alias.(*Alias)
		}
	}
	if sampleclause != nil {
		if n, ok := sampleclause.(NodeHandler); ok {
			nn.SampleClause = sampleclause.(*SampleClause)
			nn.AddChild(n)
		} else if w, ok := sampleclause.(*Wrapped); ok {
			nn.SampleClause = w.Value.(*SampleClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.SampleClause = sampleclause.(*SampleClause)
		}
	}
	return nn, nil
}



func NewUnaryExpression(
	operand interface{},
	op interface{},
	) (*UnaryExpression, error) {
	fmt.Printf("NewUnaryExpression(%v, %v)\n",operand,
		op,
		)
	nn := &UnaryExpression{}
	nn.SetKind(UnaryExpressionKind)
	
	if operand == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := operand.(NodeHandler); ok {
		nn.Operand = operand.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := operand.(*Wrapped); ok {
		nn.Operand = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Operand = operand.(ExpressionHandler)
	}
	if op != nil {
		if n, ok := op.(NodeHandler); ok {
			nn.Op = op.(UnaryOp)
			nn.AddChild(n)
		} else if w, ok := op.(*Wrapped); ok {
			nn.Op = w.Value.(UnaryOp)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Op = op.(UnaryOp)
		}
	}
	return nn, nil
}



func NewUnnestExpression(
	expression interface{},
	) (*UnnestExpression, error) {
	fmt.Printf("NewUnnestExpression(%v)\n",expression,
		)
	nn := &UnnestExpression{}
	nn.SetKind(UnnestExpressionKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	return nn, nil
}



func NewWindowClause(
	windows interface{},
	) (*WindowClause, error) {
	fmt.Printf("NewWindowClause(%v)\n",windows,
		)
	nn := &WindowClause{}
	nn.SetKind(WindowClauseKind)
	
	nn.Windows = make([]*WindowDefinition, 0, defaultCapacity)
	if n, ok := windows.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := windows.(*Wrapped); ok {
		newElem := w.Value.(*WindowDefinition)
		nn.Windows = append(nn.Windows, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := windows.(*WindowDefinition)
		nn.Windows = append(nn.Windows, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewWindowDefinition(%v, %v)\n",name,
		windowspec,
		)
	nn := &WindowDefinition{}
	nn.SetKind(WindowDefinitionKind)
	
	if name == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := name.(NodeHandler); ok {
		nn.Name = name.(*Identifier)
		nn.AddChild(n)
	} else if w, ok := name.(*Wrapped); ok {
		nn.Name = w.Value.(*Identifier)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Name = name.(*Identifier)
	}
	if windowspec == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := windowspec.(NodeHandler); ok {
		nn.WindowSpec = windowspec.(*WindowSpecification)
		nn.AddChild(n)
	} else if w, ok := windowspec.(*Wrapped); ok {
		nn.WindowSpec = w.Value.(*WindowSpecification)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.WindowSpec = windowspec.(*WindowSpecification)
	}
	return nn, nil
}



func NewWindowFrame(
	startexpr interface{},
	endexpr interface{},
	frameunit interface{},
	) (*WindowFrame, error) {
	fmt.Printf("NewWindowFrame(%v, %v, %v)\n",startexpr,
		endexpr,
		frameunit,
		)
	nn := &WindowFrame{}
	nn.SetKind(WindowFrameKind)
	
	if startexpr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := startexpr.(NodeHandler); ok {
		nn.StartExpr = startexpr.(*WindowFrameExpression)
		nn.AddChild(n)
	} else if w, ok := startexpr.(*Wrapped); ok {
		nn.StartExpr = w.Value.(*WindowFrameExpression)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.StartExpr = startexpr.(*WindowFrameExpression)
	}
	if endexpr == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := endexpr.(NodeHandler); ok {
		nn.EndExpr = endexpr.(*WindowFrameExpression)
		nn.AddChild(n)
	} else if w, ok := endexpr.(*Wrapped); ok {
		nn.EndExpr = w.Value.(*WindowFrameExpression)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.EndExpr = endexpr.(*WindowFrameExpression)
	}
	if frameunit == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := frameunit.(NodeHandler); ok {
		nn.FrameUnit = frameunit.(FrameUnit)
		nn.AddChild(n)
	} else if w, ok := frameunit.(*Wrapped); ok {
		nn.FrameUnit = w.Value.(FrameUnit)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.FrameUnit = frameunit.(FrameUnit)
	}
	return nn, nil
}



func NewWindowFrameExpression(
	expression interface{},
	boundarytype interface{},
	) (*WindowFrameExpression, error) {
	fmt.Printf("NewWindowFrameExpression(%v, %v)\n",expression,
		boundarytype,
		)
	nn := &WindowFrameExpression{}
	nn.SetKind(WindowFrameExpressionKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	if boundarytype != nil {
		if n, ok := boundarytype.(NodeHandler); ok {
			nn.BoundaryType = boundarytype.(BoundaryType)
			nn.AddChild(n)
		} else if w, ok := boundarytype.(*Wrapped); ok {
			nn.BoundaryType = w.Value.(BoundaryType)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.BoundaryType = boundarytype.(BoundaryType)
		}
	}
	return nn, nil
}



func NewLikeExpression(
	lhs interface{},
	inlist interface{},
	isnot interface{},
	) (*LikeExpression, error) {
	fmt.Printf("NewLikeExpression(%v, %v, %v)\n",lhs,
		inlist,
		isnot,
		)
	nn := &LikeExpression{}
	nn.SetKind(LikeExpressionKind)
	
	if lhs == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := lhs.(NodeHandler); ok {
		nn.LHS = lhs.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := lhs.(*Wrapped); ok {
		nn.LHS = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.LHS = lhs.(ExpressionHandler)
	}
	if inlist == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := inlist.(NodeHandler); ok {
		nn.InList = inlist.(*InList)
		nn.AddChild(n)
	} else if w, ok := inlist.(*Wrapped); ok {
		nn.InList = w.Value.(*InList)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.InList = inlist.(*InList)
	}
	if isnot != nil {
		if n, ok := isnot.(NodeHandler); ok {
			nn.IsNot = isnot.(bool)
			nn.AddChild(n)
		} else if w, ok := isnot.(*Wrapped); ok {
			nn.IsNot = w.Value.(bool)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.IsNot = isnot.(bool)
		}
	}
	return nn, nil
}



func NewWindowSpecification(
	basewindowname interface{},
	partitionby interface{},
	orderby interface{},
	windowframe interface{},
	) (*WindowSpecification, error) {
	fmt.Printf("NewWindowSpecification(%v, %v, %v, %v)\n",basewindowname,
		partitionby,
		orderby,
		windowframe,
		)
	nn := &WindowSpecification{}
	nn.SetKind(WindowSpecificationKind)
	
	if basewindowname != nil {
		if n, ok := basewindowname.(NodeHandler); ok {
			nn.BaseWindowName = basewindowname.(*Identifier)
			nn.AddChild(n)
		} else if w, ok := basewindowname.(*Wrapped); ok {
			nn.BaseWindowName = w.Value.(*Identifier)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.BaseWindowName = basewindowname.(*Identifier)
		}
	}
	if partitionby != nil {
		if n, ok := partitionby.(NodeHandler); ok {
			nn.PartitionBy = partitionby.(*PartitionBy)
			nn.AddChild(n)
		} else if w, ok := partitionby.(*Wrapped); ok {
			nn.PartitionBy = w.Value.(*PartitionBy)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.PartitionBy = partitionby.(*PartitionBy)
		}
	}
	if orderby != nil {
		if n, ok := orderby.(NodeHandler); ok {
			nn.OrderBy = orderby.(*OrderBy)
			nn.AddChild(n)
		} else if w, ok := orderby.(*Wrapped); ok {
			nn.OrderBy = w.Value.(*OrderBy)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.OrderBy = orderby.(*OrderBy)
		}
	}
	if windowframe != nil {
		if n, ok := windowframe.(NodeHandler); ok {
			nn.WindowFrame = windowframe.(*WindowFrame)
			nn.AddChild(n)
		} else if w, ok := windowframe.(*Wrapped); ok {
			nn.WindowFrame = w.Value.(*WindowFrame)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.WindowFrame = windowframe.(*WindowFrame)
		}
	}
	return nn, nil
}



func NewWithOffset(
	alias interface{},
	) (*WithOffset, error) {
	fmt.Printf("NewWithOffset(%v)\n",alias,
		)
	nn := &WithOffset{}
	nn.SetKind(WithOffsetKind)
	
	if alias != nil {
		if n, ok := alias.(NodeHandler); ok {
			nn.Alias = alias.(*Alias)
			nn.AddChild(n)
		} else if w, ok := alias.(*Wrapped); ok {
			nn.Alias = w.Value.(*Alias)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Alias = alias.(*Alias)
		}
	}
	return nn, nil
}



func NewTypeParameterList(
	parameters interface{},
	) (*TypeParameterList, error) {
	fmt.Printf("NewTypeParameterList(%v)\n",parameters,
		)
	nn := &TypeParameterList{}
	nn.SetKind(TypeParameterListKind)
	
	nn.Parameters = make([]LeafHandler, 0, defaultCapacity)
	if n, ok := parameters.(NodeHandler); ok {
		nn.AddChild(n)
	} else if w, ok := parameters.(*Wrapped); ok {
		newElem := w.Value.(LeafHandler)
		nn.Parameters = append(nn.Parameters, newElem)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		newElem := parameters.(LeafHandler)
		nn.Parameters = append(nn.Parameters, newElem)
	}
	return nn, nil
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
	fmt.Printf("NewSampleClause(%v, %v, %v)\n",samplemethod,
		samplesize,
		samplesuffix,
		)
	nn := &SampleClause{}
	nn.SetKind(SampleClauseKind)
	
	if samplemethod == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := samplemethod.(NodeHandler); ok {
		nn.SampleMethod = samplemethod.(*Identifier)
		nn.AddChild(n)
	} else if w, ok := samplemethod.(*Wrapped); ok {
		nn.SampleMethod = w.Value.(*Identifier)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.SampleMethod = samplemethod.(*Identifier)
	}
	if samplesize == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := samplesize.(NodeHandler); ok {
		nn.SampleSize = samplesize.(*SampleSize)
		nn.AddChild(n)
	} else if w, ok := samplesize.(*Wrapped); ok {
		nn.SampleSize = w.Value.(*SampleSize)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.SampleSize = samplesize.(*SampleSize)
	}
	if samplesuffix != nil {
		if n, ok := samplesuffix.(NodeHandler); ok {
			nn.SampleSuffix = samplesuffix.(*SampleSuffix)
			nn.AddChild(n)
		} else if w, ok := samplesuffix.(*Wrapped); ok {
			nn.SampleSuffix = w.Value.(*SampleSuffix)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.SampleSuffix = samplesuffix.(*SampleSuffix)
		}
	}
	return nn, nil
}



func NewSampleSize(
	size interface{},
	partitionby interface{},
	unit interface{},
	) (*SampleSize, error) {
	fmt.Printf("NewSampleSize(%v, %v, %v)\n",size,
		partitionby,
		unit,
		)
	nn := &SampleSize{}
	nn.SetKind(SampleSizeKind)
	
	if size == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := size.(NodeHandler); ok {
		nn.Size = size.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := size.(*Wrapped); ok {
		nn.Size = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Size = size.(ExpressionHandler)
	}
	if partitionby != nil {
		if n, ok := partitionby.(NodeHandler); ok {
			nn.PartitionBy = partitionby.(*PartitionBy)
			nn.AddChild(n)
		} else if w, ok := partitionby.(*Wrapped); ok {
			nn.PartitionBy = w.Value.(*PartitionBy)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.PartitionBy = partitionby.(*PartitionBy)
		}
	}
	if unit != nil {
		if n, ok := unit.(NodeHandler); ok {
			nn.Unit = unit.(SampleSizeUnit)
			nn.AddChild(n)
		} else if w, ok := unit.(*Wrapped); ok {
			nn.Unit = w.Value.(SampleSizeUnit)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Unit = unit.(SampleSizeUnit)
		}
	}
	return nn, nil
}



func NewSampleSuffix(
	weight interface{},
	repeat interface{},
	) (*SampleSuffix, error) {
	fmt.Printf("NewSampleSuffix(%v, %v)\n",weight,
		repeat,
		)
	nn := &SampleSuffix{}
	nn.SetKind(SampleSuffixKind)
	
	if weight != nil {
		if n, ok := weight.(NodeHandler); ok {
			nn.Weight = weight.(*WithWeight)
			nn.AddChild(n)
		} else if w, ok := weight.(*Wrapped); ok {
			nn.Weight = w.Value.(*WithWeight)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Weight = weight.(*WithWeight)
		}
	}
	if repeat != nil {
		if n, ok := repeat.(NodeHandler); ok {
			nn.Repeat = repeat.(*RepeatableClause)
			nn.AddChild(n)
		} else if w, ok := repeat.(*Wrapped); ok {
			nn.Repeat = w.Value.(*RepeatableClause)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Repeat = repeat.(*RepeatableClause)
		}
	}
	return nn, nil
}



func NewWithWeight(
	alias interface{},
	) (*WithWeight, error) {
	fmt.Printf("NewWithWeight(%v)\n",alias,
		)
	nn := &WithWeight{}
	nn.SetKind(WithWeightKind)
	
	if alias != nil {
		if n, ok := alias.(NodeHandler); ok {
			nn.Alias = alias.(*Alias)
			nn.AddChild(n)
		} else if w, ok := alias.(*Wrapped); ok {
			nn.Alias = w.Value.(*Alias)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.Alias = alias.(*Alias)
		}
	}
	return nn, nil
}



func NewRepeatableClause(
	argument interface{},
	) (*RepeatableClause, error) {
	fmt.Printf("NewRepeatableClause(%v)\n",argument,
		)
	nn := &RepeatableClause{}
	nn.SetKind(RepeatableClauseKind)
	
	if argument == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := argument.(NodeHandler); ok {
		nn.Argument = argument.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := argument.(*Wrapped); ok {
		nn.Argument = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Argument = argument.(ExpressionHandler)
	}
	return nn, nil
}



func NewQualify(
	expression interface{},
	) (*Qualify, error) {
	fmt.Printf("NewQualify(%v)\n",expression,
		)
	nn := &Qualify{}
	nn.SetKind(QualifyKind)
	
	if expression == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := expression.(NodeHandler); ok {
		nn.Expression = expression.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := expression.(*Wrapped); ok {
		nn.Expression = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Expression = expression.(ExpressionHandler)
	}
	return nn, nil
}



func NewFormatClause(
	format interface{},
	timezoneexpr interface{},
	) (*FormatClause, error) {
	fmt.Printf("NewFormatClause(%v, %v)\n",format,
		timezoneexpr,
		)
	nn := &FormatClause{}
	nn.SetKind(FormatClauseKind)
	
	if format == nil {
		return nil, ErrMissingRequiredField
	}
	if n, ok := format.(NodeHandler); ok {
		nn.Format = format.(ExpressionHandler)
		nn.AddChild(n)
	} else if w, ok := format.(*Wrapped); ok {
		nn.Format = w.Value.(ExpressionHandler)
		nn.ExpandLoc(w.Loc.Start, w.Loc.End)
	} else {
		nn.Format = format.(ExpressionHandler)
	}
	if timezoneexpr != nil {
		if n, ok := timezoneexpr.(NodeHandler); ok {
			nn.TimeZoneExpr = timezoneexpr.(ExpressionHandler)
			nn.AddChild(n)
		} else if w, ok := timezoneexpr.(*Wrapped); ok {
			nn.TimeZoneExpr = w.Value.(ExpressionHandler)
			nn.ExpandLoc(w.Loc.Start, w.Loc.End)
		} else {
			nn.TimeZoneExpr = timezoneexpr.(ExpressionHandler)
		}
	}
	return nn, nil
}



type Visitor interface {
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
	VisitOnClause(*OnClause, interface{})
	VisitWithClauseEntry(*WithClauseEntry, interface{})
	VisitJoin(*Join, interface{})
	VisitWithClause(*WithClause, interface{})
	VisitHaving(*Having, interface{})
	VisitNamedType(*NamedType, interface{})
	VisitArrayType(*ArrayType, interface{})
	VisitStructField(*StructField, interface{})
	VisitStructType(*StructType, interface{})
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
	VisitIntervalExpression(*IntervalExpression, interface{})
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
	VisitWindowFrameExpression(*WindowFrameExpression, interface{})
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
	
}

func (n *QueryStatement) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitQueryStatement %s\n",
		n.SingleNodeDebugString())

	v.VisitQueryStatement(n, d)
}

func (n *Query) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitQuery %s\n",
		n.SingleNodeDebugString())

	v.VisitQuery(n, d)
}

func (n *Select) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSelect %s\n",
		n.SingleNodeDebugString())

	v.VisitSelect(n, d)
}

func (n *SelectList) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSelectList %s\n",
		n.SingleNodeDebugString())

	v.VisitSelectList(n, d)
}

func (n *SelectColumn) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSelectColumn %s\n",
		n.SingleNodeDebugString())

	v.VisitSelectColumn(n, d)
}

func (n *IntLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitIntLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitIntLiteral(n, d)
}

func (n *Identifier) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitIdentifier %s\n",
		n.SingleNodeDebugString())

	v.VisitIdentifier(n, d)
}

func (n *Alias) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitAlias %s\n",
		n.SingleNodeDebugString())

	v.VisitAlias(n, d)
}

func (n *PathExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitPathExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitPathExpression(n, d)
}

func (n *TablePathExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitTablePathExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitTablePathExpression(n, d)
}

func (n *FromClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitFromClause %s\n",
		n.SingleNodeDebugString())

	v.VisitFromClause(n, d)
}

func (n *WhereClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWhereClause %s\n",
		n.SingleNodeDebugString())

	v.VisitWhereClause(n, d)
}

func (n *BooleanLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitBooleanLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitBooleanLiteral(n, d)
}

func (n *AndExpr) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitAndExpr %s\n",
		n.SingleNodeDebugString())

	v.VisitAndExpr(n, d)
}

func (n *BinaryExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitBinaryExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitBinaryExpression(n, d)
}

func (n *StringLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStringLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitStringLiteral(n, d)
}

func (n *Star) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStar %s\n",
		n.SingleNodeDebugString())

	v.VisitStar(n, d)
}

func (n *OrExpr) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitOrExpr %s\n",
		n.SingleNodeDebugString())

	v.VisitOrExpr(n, d)
}

func (n *GroupingItem) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitGroupingItem %s\n",
		n.SingleNodeDebugString())

	v.VisitGroupingItem(n, d)
}

func (n *GroupBy) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitGroupBy %s\n",
		n.SingleNodeDebugString())

	v.VisitGroupBy(n, d)
}

func (n *OrderingExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitOrderingExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitOrderingExpression(n, d)
}

func (n *OrderBy) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitOrderBy %s\n",
		n.SingleNodeDebugString())

	v.VisitOrderBy(n, d)
}

func (n *LimitOffset) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitLimitOffset %s\n",
		n.SingleNodeDebugString())

	v.VisitLimitOffset(n, d)
}

func (n *FloatLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitFloatLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitFloatLiteral(n, d)
}

func (n *NullLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitNullLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitNullLiteral(n, d)
}

func (n *OnClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitOnClause %s\n",
		n.SingleNodeDebugString())

	v.VisitOnClause(n, d)
}

func (n *WithClauseEntry) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWithClauseEntry %s\n",
		n.SingleNodeDebugString())

	v.VisitWithClauseEntry(n, d)
}

func (n *Join) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitJoin %s\n",
		n.SingleNodeDebugString())

	v.VisitJoin(n, d)
}

func (n *WithClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWithClause %s\n",
		n.SingleNodeDebugString())

	v.VisitWithClause(n, d)
}

func (n *Having) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitHaving %s\n",
		n.SingleNodeDebugString())

	v.VisitHaving(n, d)
}

func (n *NamedType) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitNamedType %s\n",
		n.SingleNodeDebugString())

	v.VisitNamedType(n, d)
}

func (n *ArrayType) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitArrayType %s\n",
		n.SingleNodeDebugString())

	v.VisitArrayType(n, d)
}

func (n *StructField) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStructField %s\n",
		n.SingleNodeDebugString())

	v.VisitStructField(n, d)
}

func (n *StructType) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStructType %s\n",
		n.SingleNodeDebugString())

	v.VisitStructType(n, d)
}

func (n *CastExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitCastExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitCastExpression(n, d)
}

func (n *SelectAs) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSelectAs %s\n",
		n.SingleNodeDebugString())

	v.VisitSelectAs(n, d)
}

func (n *Rollup) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitRollup %s\n",
		n.SingleNodeDebugString())

	v.VisitRollup(n, d)
}

func (n *FunctionCall) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitFunctionCall %s\n",
		n.SingleNodeDebugString())

	v.VisitFunctionCall(n, d)
}

func (n *ArrayConstructor) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitArrayConstructor %s\n",
		n.SingleNodeDebugString())

	v.VisitArrayConstructor(n, d)
}

func (n *StructConstructorArg) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStructConstructorArg %s\n",
		n.SingleNodeDebugString())

	v.VisitStructConstructorArg(n, d)
}

func (n *StructConstructorWithParens) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStructConstructorWithParens %s\n",
		n.SingleNodeDebugString())

	v.VisitStructConstructorWithParens(n, d)
}

func (n *StructConstructorWithKeyword) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStructConstructorWithKeyword %s\n",
		n.SingleNodeDebugString())

	v.VisitStructConstructorWithKeyword(n, d)
}

func (n *InExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitInExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitInExpression(n, d)
}

func (n *InList) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitInList %s\n",
		n.SingleNodeDebugString())

	v.VisitInList(n, d)
}

func (n *BetweenExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitBetweenExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitBetweenExpression(n, d)
}

func (n *NumericLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitNumericLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitNumericLiteral(n, d)
}

func (n *BigNumericLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitBigNumericLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitBigNumericLiteral(n, d)
}

func (n *BytesLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitBytesLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitBytesLiteral(n, d)
}

func (n *DateOrTimeLiteral) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitDateOrTimeLiteral %s\n",
		n.SingleNodeDebugString())

	v.VisitDateOrTimeLiteral(n, d)
}

func (n *CaseValueExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitCaseValueExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitCaseValueExpression(n, d)
}

func (n *CaseNoValueExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitCaseNoValueExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitCaseNoValueExpression(n, d)
}

func (n *ArrayElement) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitArrayElement %s\n",
		n.SingleNodeDebugString())

	v.VisitArrayElement(n, d)
}

func (n *BitwiseShiftExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitBitwiseShiftExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitBitwiseShiftExpression(n, d)
}

func (n *DotGeneralizedField) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitDotGeneralizedField %s\n",
		n.SingleNodeDebugString())

	v.VisitDotGeneralizedField(n, d)
}

func (n *DotIdentifier) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitDotIdentifier %s\n",
		n.SingleNodeDebugString())

	v.VisitDotIdentifier(n, d)
}

func (n *DotStar) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitDotStar %s\n",
		n.SingleNodeDebugString())

	v.VisitDotStar(n, d)
}

func (n *DotStarWithModifiers) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitDotStarWithModifiers %s\n",
		n.SingleNodeDebugString())

	v.VisitDotStarWithModifiers(n, d)
}

func (n *ExpressionSubquery) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitExpressionSubquery %s\n",
		n.SingleNodeDebugString())

	v.VisitExpressionSubquery(n, d)
}

func (n *ExtractExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitExtractExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitExtractExpression(n, d)
}

func (n *IntervalExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitIntervalExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitIntervalExpression(n, d)
}

func (n *NullOrder) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitNullOrder %s\n",
		n.SingleNodeDebugString())

	v.VisitNullOrder(n, d)
}

func (n *OnOrUsingClauseList) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitOnOrUsingClauseList %s\n",
		n.SingleNodeDebugString())

	v.VisitOnOrUsingClauseList(n, d)
}

func (n *ParenthesizedJoin) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitParenthesizedJoin %s\n",
		n.SingleNodeDebugString())

	v.VisitParenthesizedJoin(n, d)
}

func (n *PartitionBy) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitPartitionBy %s\n",
		n.SingleNodeDebugString())

	v.VisitPartitionBy(n, d)
}

func (n *SetOperation) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSetOperation %s\n",
		n.SingleNodeDebugString())

	v.VisitSetOperation(n, d)
}

func (n *StarExceptList) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStarExceptList %s\n",
		n.SingleNodeDebugString())

	v.VisitStarExceptList(n, d)
}

func (n *StarModifiers) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStarModifiers %s\n",
		n.SingleNodeDebugString())

	v.VisitStarModifiers(n, d)
}

func (n *StarReplaceItem) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStarReplaceItem %s\n",
		n.SingleNodeDebugString())

	v.VisitStarReplaceItem(n, d)
}

func (n *StarWithModifiers) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitStarWithModifiers %s\n",
		n.SingleNodeDebugString())

	v.VisitStarWithModifiers(n, d)
}

func (n *TableSubquery) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitTableSubquery %s\n",
		n.SingleNodeDebugString())

	v.VisitTableSubquery(n, d)
}

func (n *UnaryExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitUnaryExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitUnaryExpression(n, d)
}

func (n *UnnestExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitUnnestExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitUnnestExpression(n, d)
}

func (n *WindowClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWindowClause %s\n",
		n.SingleNodeDebugString())

	v.VisitWindowClause(n, d)
}

func (n *WindowDefinition) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWindowDefinition %s\n",
		n.SingleNodeDebugString())

	v.VisitWindowDefinition(n, d)
}

func (n *WindowFrame) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWindowFrame %s\n",
		n.SingleNodeDebugString())

	v.VisitWindowFrame(n, d)
}

func (n *WindowFrameExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWindowFrameExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitWindowFrameExpression(n, d)
}

func (n *LikeExpression) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitLikeExpression %s\n",
		n.SingleNodeDebugString())

	v.VisitLikeExpression(n, d)
}

func (n *WindowSpecification) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWindowSpecification %s\n",
		n.SingleNodeDebugString())

	v.VisitWindowSpecification(n, d)
}

func (n *WithOffset) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWithOffset %s\n",
		n.SingleNodeDebugString())

	v.VisitWithOffset(n, d)
}

func (n *TypeParameterList) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitTypeParameterList %s\n",
		n.SingleNodeDebugString())

	v.VisitTypeParameterList(n, d)
}

func (n *SampleClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSampleClause %s\n",
		n.SingleNodeDebugString())

	v.VisitSampleClause(n, d)
}

func (n *SampleSize) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSampleSize %s\n",
		n.SingleNodeDebugString())

	v.VisitSampleSize(n, d)
}

func (n *SampleSuffix) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitSampleSuffix %s\n",
		n.SingleNodeDebugString())

	v.VisitSampleSuffix(n, d)
}

func (n *WithWeight) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitWithWeight %s\n",
		n.SingleNodeDebugString())

	v.VisitWithWeight(n, d)
}

func (n *RepeatableClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitRepeatableClause %s\n",
		n.SingleNodeDebugString())

	v.VisitRepeatableClause(n, d)
}

func (n *Qualify) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitQualify %s\n",
		n.SingleNodeDebugString())

	v.VisitQualify(n, d)
}

func (n *FormatClause) Accept(v Visitor, d interface{}) {
	fmt.Printf("VisitFormatClause %s\n",
		n.SingleNodeDebugString())

	v.VisitFormatClause(n, d)
}



// Operation is the base for new operations using visitors.
type Operation struct {
	// visitor is the visitor to passed when fallbacking to walk.
	visitor Visitor
}

func (o *Operation) VisitQueryStatement(n *QueryStatement, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitQuery(n *Query, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSelect(n *Select, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSelectList(n *SelectList, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSelectColumn(n *SelectColumn, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitIntLiteral(n *IntLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitIdentifier(n *Identifier, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitAlias(n *Alias, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitPathExpression(n *PathExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitTablePathExpression(n *TablePathExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitFromClause(n *FromClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWhereClause(n *WhereClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitBooleanLiteral(n *BooleanLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitAndExpr(n *AndExpr, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitBinaryExpression(n *BinaryExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStringLiteral(n *StringLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStar(n *Star, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitOrExpr(n *OrExpr, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitGroupingItem(n *GroupingItem, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitGroupBy(n *GroupBy, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitOrderingExpression(n *OrderingExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitOrderBy(n *OrderBy, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitLimitOffset(n *LimitOffset, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitFloatLiteral(n *FloatLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitNullLiteral(n *NullLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitOnClause(n *OnClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWithClauseEntry(n *WithClauseEntry, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitJoin(n *Join, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWithClause(n *WithClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitHaving(n *Having, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitNamedType(n *NamedType, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitArrayType(n *ArrayType, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStructField(n *StructField, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStructType(n *StructType, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitCastExpression(n *CastExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSelectAs(n *SelectAs, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitRollup(n *Rollup, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitFunctionCall(n *FunctionCall, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitArrayConstructor(n *ArrayConstructor, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStructConstructorArg(n *StructConstructorArg, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStructConstructorWithParens(n *StructConstructorWithParens, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStructConstructorWithKeyword(n *StructConstructorWithKeyword, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitInExpression(n *InExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitInList(n *InList, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitBetweenExpression(n *BetweenExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitNumericLiteral(n *NumericLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitBigNumericLiteral(n *BigNumericLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitBytesLiteral(n *BytesLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitDateOrTimeLiteral(n *DateOrTimeLiteral, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitCaseValueExpression(n *CaseValueExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitCaseNoValueExpression(n *CaseNoValueExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitArrayElement(n *ArrayElement, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitBitwiseShiftExpression(n *BitwiseShiftExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitDotGeneralizedField(n *DotGeneralizedField, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitDotIdentifier(n *DotIdentifier, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitDotStar(n *DotStar, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitDotStarWithModifiers(n *DotStarWithModifiers, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitExpressionSubquery(n *ExpressionSubquery, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitExtractExpression(n *ExtractExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitIntervalExpression(n *IntervalExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitNullOrder(n *NullOrder, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitOnOrUsingClauseList(n *OnOrUsingClauseList, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitParenthesizedJoin(n *ParenthesizedJoin, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitPartitionBy(n *PartitionBy, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSetOperation(n *SetOperation, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStarExceptList(n *StarExceptList, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStarModifiers(n *StarModifiers, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStarReplaceItem(n *StarReplaceItem, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitStarWithModifiers(n *StarWithModifiers, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitTableSubquery(n *TableSubquery, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitUnaryExpression(n *UnaryExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitUnnestExpression(n *UnnestExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWindowClause(n *WindowClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWindowDefinition(n *WindowDefinition, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWindowFrame(n *WindowFrame, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWindowFrameExpression(n *WindowFrameExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitLikeExpression(n *LikeExpression, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWindowSpecification(n *WindowSpecification, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWithOffset(n *WithOffset, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitTypeParameterList(n *TypeParameterList, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSampleClause(n *SampleClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSampleSize(n *SampleSize, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitSampleSuffix(n *SampleSuffix, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitWithWeight(n *WithWeight, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitRepeatableClause(n *RepeatableClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitQualify(n *Qualify, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
}
func (o *Operation) VisitFormatClause(n *FormatClause, d interface{}) {
	for _, c := range n.Children() {
		fmt.Printf("%s -> %s\n",
			n.SingleNodeDebugString(),
			c.(NodeStringer).SingleNodeDebugString())
		c.Accept(o.visitor, d)
	}
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
	OnClauseKind
	WithClauseEntryKind
	JoinKind
	WithClauseKind
	HavingKind
	NamedTypeKind
	ArrayTypeKind
	StructFieldKind
	StructTypeKind
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
	IntervalExpressionKind
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
	WindowFrameExpressionKind
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
	FormatClauseKind)

func (k NodeKind) String() string {
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
	case OnClauseKind:
		return "OnClause"
	case WithClauseEntryKind:
		return "WithClauseEntry"
	case JoinKind:
		return "Join"
	case WithClauseKind:
		return "WithClause"
	case HavingKind:
		return "Having"
	case NamedTypeKind:
		return "NamedType"
	case ArrayTypeKind:
		return "ArrayType"
	case StructFieldKind:
		return "StructField"
	case StructTypeKind:
		return "StructType"
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
	case IntervalExpressionKind:
		return "IntervalExpression"
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
	case WindowFrameExpressionKind:
		return "WindowFrameExpression"
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
 	}
	panic("unexpected kind")
 }
