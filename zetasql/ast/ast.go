package ast

import (
	"fmt"
	"strings"
)

type Node struct {
	// Loc holds the parser location.
	Loc

	parent   NodeHandler
	children []NodeHandler
	kind     NodeKind
}

// A NodeHandler handles all generic information for any node.
type NodeHandler interface {
	IsTableExpression() bool
	IsQueryExpression() bool
	IsExpression() bool
	IsType() bool
	IsLeaf() bool
	IsStatement() bool
	IsScriptStatement() bool
	IsSQLStatement() bool

	// Location range of the parsed SQL.
	GetLoc() Loc
	StartLoc() int
	EndLoc() int
	SetStartLoc(int)
	SetEndLoc(int)
	// ExpandLoc updates the parsing location expanding when the given
	// range is outside the current range.
	ExpandLoc(start int, end int)

	Kind() NodeKind
	SetKind(NodeKind)

	// Tree handlers.
	Parent() NodeHandler
	SetParent(NodeHandler)
	Children() []NodeHandler
	// AddChild adds child to the list of children.  The child element
	// must not be nil.
	AddChild(NodeHandler)
	// AddChildren all nodes to the list of children.  Elements that
	// are nil are ignored.
	AddChildren([]NodeHandler)

	// Accept accepts the visitor with an optional arbitrary data.
	Accept(Visitor, interface{})
}

type NodeStringer interface {
	String() string

	// DebugString returns a multi-line tree dump.  Parse locations
	// are represented as integer ranges.  When sql is passed, fragments
	// of the original sql are used instead of raw integer offsets.
	DebugString(sql string) string

	// SingleNodeDebugString returns an one-line description of the
	// node, including modifiers but without child nodes.
	SingleNodeDebugString() string
}

type Expression struct {
	Node
	parenthesized bool
}

type ExpressionHandler interface {
	NodeHandler

	IsExpression() bool
	IsParenthesized() bool
	SetParenthesized(bool)

	// IsAllowedInComparison returns true when the expression is allowed
	// to occur as a child of a comparison expression.  This is not
	// allowed for unparenthesized comparison expressions and operators
	// with a lower precedence level (AND, OR, and NOT).
	IsAllowedInComparison() bool
}

type Leaf struct {
	Expression
	image string
}

type LeafHandler interface {
	ExpressionHandler

	IsLeaf() bool
	Image() string
	SetImage(string)
}

// Loc stores half-open range [Start, End) of offsets.
type Loc struct {
	Start int
	End   int
}

type Statement struct {
	Node
}

// A StatementHandler handles information for all statements.
type StatementHandler interface {
	NodeHandler

	IsStatement() bool
	IsSQLStatement() bool
}

type QueryExpression struct {
	Node
	parenthesized bool
}

// A QueryExpressionHandler handles information of all query expression
// types. These are top-level syntactic constructs (outside individual
// SELECTs) making up a query.  These include Query itself, Select,
// UnionAll, etc.
type QueryExpressionHandler interface {
	NodeHandler

	IsQueryExpression() bool
	IsParenthesized() bool
	SetParenthesized(bool)
}

type TableExpression struct {
	Node
}

func (t *TableExpression) IsTableExpression() bool { return true }
func (t *TableExpression) GetAlias() *Alias        { return nil }

// A TableExpressionHandler handles information for all table
// expressions.  These are this that appear in the from clause and
// produce a stream of rows like a table.  This includes table scans,
// joins, and subqueries.
type TableExpressionHandler interface {
	NodeHandler

	IsTableExpression() bool

	// GetAlias returns the alias, if the particular table expression
	// has one.
	GetAlias() *Alias
}

type Type struct {
	Node
	parameters *TypeParameterList
}

type TypeHandler interface {
	NodeHandler

	IsType() bool
	SetTypeParameters(*TypeParameterList)
	TypeParameters() *TypeParameterList
}

func (n *Node) IsTableExpression() bool { return false }
func (n *Node) IsQueryExpression() bool { return false }
func (n *Node) IsExpression() bool      { return false }
func (n *Node) IsType() bool            { return false }
func (n *Node) IsLeaf() bool            { return false }
func (n *Node) IsStatement() bool       { return false }
func (n *Node) IsScriptStatement() bool { return false }
func (n *Node) IsSQLStatement() bool    { return false }
func (n *Node) Kind() NodeKind          { return n.kind }
func (n *Node) SetKind(k NodeKind)      { n.kind = k }
func (n *Node) Parent() NodeHandler     { return n.parent }
func (n *Node) SetParent(p NodeHandler) { n.parent = p }
func (n *Node) Children() []NodeHandler { return n.children }

func (n *Node) AddChild(c NodeHandler) {
	n.children = append(n.children, c)
	c.SetParent(n)
	n.ExpandLoc(c.StartLoc(), c.EndLoc())
}

func (n *Node) GetLoc() Loc {
	return n.Loc
}

func (n *Node) ExpandLoc(start int, end int) {
	s := n.StartLoc()
	e := n.EndLoc()

	if s == 0 && e == 0 {
		n.SetStartLoc(start)
		n.SetEndLoc(end)
	} else {
		if s > start {
			n.SetStartLoc(start)
		}
		if e < end {
			n.SetEndLoc(end)
		}
	}
}

func (n *Node) AddChildren(children []NodeHandler) {
	for _, c := range children {
		if c != nil {
			n.AddChild(c)
		}
	}
}

func (n *Node) Accept(v Visitor, d interface{}) {
	panic("Node.Accept() called")
}

func (n *Node) SingleNodeDebugString() string {
	return n.kind.String()
}

func (n *Node) DebugString(sql string) string {
	d := newDumper(n, "\n", 256, sql)
	d.Dump()

	return d.String()
}

func (n *Node) String() string {
	return n.SingleNodeDebugString()
}

func (l *Loc) StartLoc() int       { return l.Start }
func (l *Loc) SetStartLoc(pos int) { l.Start = pos }
func (l *Loc) EndLoc() int         { return l.End }
func (l *Loc) SetEndLoc(pos int)   { l.End = pos }

func (l *Leaf) IsLeaf() bool      { return true }
func (l *Leaf) Image() string     { return l.image }
func (l *Leaf) SetImage(d string) { l.image = d }

func (l *Leaf) SingleNodeDebugString() string {
	return fmt.Sprintf("%s(%s)", l.kind.String(), l.Image())
}

func (s *Statement) IsStatement() bool    { return true }
func (s *Statement) IsSQLStatement() bool { return true }

func (e *Expression) IsExpression() bool          { return true }
func (e *Expression) IsParenthesized() bool       { return e.parenthesized }
func (e *Expression) SetParenthesized(v bool)     { e.parenthesized = v }
func (e *Expression) IsAllowedInComparison() bool { return true }

func (e *QueryExpression) IsQueryExpression() bool { return true }
func (e *QueryExpression) IsParenthesized() bool   { return true }
func (e *QueryExpression) SetParenthesized(v bool) { e.parenthesized = v }

func (e *AndExpr) IsAllowedInComparison() bool { return e.IsParenthesized() }

func (e *OrExpr) IsAllowedInComparison() bool { return e.IsParenthesized() }

func (t *Type) IsType() bool { return true }

func (e *InExpression) IsAllowedInComparison() bool {
	return e.IsParenthesized()
}

func (e *BetweenExpression) IsAllowedInComparison() bool {
	return e.IsParenthesized()
}

func (e *UnaryExpression) IsAllowedInComparison() bool {
	return e.IsParenthesized() || e.Op != UnaryNot
}

func (e *LikeExpression) IsAllowedInComparison() bool {
	return e.IsParenthesized()
}

func (n *BinaryExpression) SingleNodeDebugString() string {
	var kind string = n.Op.String()

	if n.IsNot {
		switch n.Op {
		case BinaryIs:
			kind = "IS NOT"
		case BinaryLike:
			kind = "NOT LIKE"
		}
	}

	return fmt.Sprintf("%s(%s)", n.kind.String(), kind)
}

func (n *UnaryExpression) SingleNodeDebugString() string {
	return fmt.Sprintf("%s(%v)", n.kind.String(), n.Op)
}

func (n *BetweenExpression) SingleNodeDebugString() string {
	return fmt.Sprintf("%s(BETWEEN)", n.kind.String())
}

func (n *Identifier) SingleNodeDebugString() string {
	return fmt.Sprintf("%s(%s)", n.kind.String(),
		ToIdentifierLiteral(n.IDString))
}

func (n *Join) SingleNodeDebugString() string {
	return fmt.Sprintf("%s(%v)", n.kind.String(), n.JoinType)
}

func (n *FunctionCall) SingleNodeDebugString() string {
	var opts string
	if n.Distinct {
		opts = "(distinct=true)"
	}

	return fmt.Sprintf("%s%s", n.kind.String(), opts)
}

func (n *CastExpression) SingleNodeDebugString() string {
	var opts string
	if n.IsSafeCast {
		opts = "(return_null_on_error=true)"
	}

	return fmt.Sprintf("%s%s", n.kind.String(), opts)
}

func (n *TablePathExpression) GetAlias() *Alias { return n.Alias }

func (n *BooleanLiteral) SetImage(v string) {
	if strings.ToUpper(v) == "TRUE" {
		n.Value = true
	} else {
		n.Value = false
	}

	n.image = v
}

func (t *Type) SetTypeParameters(p *TypeParameterList) {
	t.parameters = p
}
func (t *Type) TypeParameters() *TypeParameterList {
	return t.parameters
}

var _ TypeHandler = (*Type)(nil)
var _ TypeHandler = (*NamedType)(nil)
