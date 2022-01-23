package ast

type BinaryOp int

const (
	BinaryEq         BinaryOp = iota // "="
	BinaryNE                         // "!="
	BinaryNE2                        // "<>"
	BinaryLT                         // "<"
	BinaryLE                         // "<="
	BinaryGT                         // ">"
	BinaryGE                         // ">="
	BinaryPlus                       // "+"
	BinaryMinus                      // "-"
	BinaryMultiply                   // "*"
	BinaryDivide                     // "/"
	BinaryIs                         // "IS"
	BinaryBitwiseOr                  // "|"
	BinaryBitwiseXor                 // "^"
	BinaryBitwiseAnd                 // "&"
	BinaryConcat                     // "||"
	BinaryLike                       // "LIKE"
)

func (o BinaryOp) String() string {
	switch o {
	case BinaryEq:
		return "="
	case BinaryNE:
		return "!="
	case BinaryNE2:
		return "<>"
	case BinaryLT:
		return "<"
	case BinaryLE:
		return "<="
	case BinaryGT:
		return ">"
	case BinaryGE:
		return ">="
	case BinaryPlus:
		return "+"
	case BinaryMinus:
		return "-"
	case BinaryMultiply:
		return "*"
	case BinaryDivide:
		return "/"
	case BinaryIs:
		return "IS"
	case BinaryBitwiseOr:
		return "|"
	case BinaryBitwiseXor:
		return "^"
	case BinaryBitwiseAnd:
		return "&"
	case BinaryConcat:
		return "||"
	case BinaryLike:
		return "LIKE"
	}

	panic("unknown binary op")
}

type UnaryOp int

const (
	NoUnaryOp UnaryOp = iota
	UnaryNot
	UnaryBitwiseNot
	UnaryMinus
	UnaryPlus
)

func (o UnaryOp) String() string {
	switch o {
	case NoUnaryOp:
		return "NOT_SET"
	case UnaryNot:
		return "NOT"
	case UnaryBitwiseNot:
		return "~"
	case UnaryMinus:
		return "-"
	case UnaryPlus:
		return "+"
	}

	panic("unknown unary op")
}

type OrderingSpec int

const (
	NoOrderingSpec OrderingSpec = iota
	AscendingOrder
	DescendingOrder
)

func (o OrderingSpec) String() string {
	switch o {
	case NoOrderingSpec, AscendingOrder:
		return "ASC"
	case DescendingOrder:
		return "DESC"
	}

	panic("unknown ordering spec")
}

type JoinType int

const (
	DefaultJoin JoinType = iota
	CommaJoin
	CrossJoin
	FullJoin
	InnerJoin
	LeftJoin
	RightJoin
)

func (t JoinType) String() string {
	switch t {
	case DefaultJoin:
		return "JOIN"
	case CommaJoin:
		return "COMMA"
	case CrossJoin:
		return "CROSS"
	case FullJoin:
		return "FULL"
	case InnerJoin:
		return "INNER"
	case LeftJoin:
		return "LEFT"
	case RightJoin:
		return "RIGHT"
	}

	panic("unknown join type")
}

type AsMode int

const (
	NoAsMode AsMode = iota
	AsStruct
	AsValue
	AsTypeName
)

type NullHandlingModifier int

const (
	DefaultNullHandling NullHandlingModifier = iota
	IgnoreNulls
	RespectNulls
)

type SubqueryModifier int

const (
	NoSubqueryModifier     SubqueryModifier = iota // (SELECT ...)
	ArraySubqueryModifier                          // ARRAY(SELECT ...)
	ExistsSubqueryModifier                         // EXISTS(SELECT ...)
)

func (m SubqueryModifier) String() string {
	return m.ToSQL()
}

func (m SubqueryModifier) ToSQL() string {
	switch m {
	case NoSubqueryModifier:
		return ""
	case ArraySubqueryModifier:
		return "ARRAY"
	case ExistsSubqueryModifier:
		return "EXISTS"
	}

	panic("unknown subquery modifier")
}

type SetOp int

const (
	Union     SetOp = iota // UNION {ALL|DISTINCT}
	Except                 // EXCEPT {ALL|DISTINCT}
	Intersect              // INTERSECT {ALL|DISTINCt}
)

func (o SetOp) String() string {
	return o.ToSQL()
}

func (o SetOp) ToSQL() string {
	switch o {
	case Union:
		return "UNION"
	case Except:
		return "EXCEPT"
	case Intersect:
		return "INTERSECT"
	}

	panic("unknown set operation")
}

type BoundaryType int

const (
	UnboundedPreceding BoundaryType = iota
	OffsetPreceding
	CurrentRow
	OffsetFollowing
	UnboundedFollowing
)

func (t BoundaryType) ToSQL() string {
	switch t {
	case UnboundedPreceding:
		return "UNBOUNDED PRECEDING"
	case OffsetPreceding:
		return "PRECEDING"
	case CurrentRow:
		return "CURRENT ROW"
	case OffsetFollowing:
		return "FOLLOWING"
	case UnboundedFollowing:
		return "UNBOUNDED FOLLOWING"
	}

	panic("unknown boundary type")
}

type SampleSizeUnit int

const (
	RowsSampling SampleSizeUnit = iota
	PercentSampling
)

func (s SampleSizeUnit) String() string {
	switch s {
	case RowsSampling:
		return "ROWS"
	case PercentSampling:
		return "PERCENT"
	}

	panic("unknown sample size unit")
}

type TypeKind int

const (
	BoolKind TypeKind = iota
	IntegerKind
	FloatingPointKind
	DateKind
	DateTimeKind
	TimeKind
	TimestampKind
	OtherKind
)

func (t TypeKind) ToSQL() string {
	switch t {
	case BoolKind:
		return "BOOL"
	case IntegerKind:
		return "INT64"
	case FloatingPointKind:
		return "FLOAT64"
	case DateKind:
		return "DATE"
	case DateTimeKind:
		return "DATETIME"
	case TimeKind:
		return "TIME"
	case TimestampKind:
		return "TIMESTAMP"
	case OtherKind:
		panic("cannot generate OtherKind as SQL")
	}

	panic("unknown type kind")
}

type FrameUnit int

const (
	Rows FrameUnit = iota
	Range
)

func (u FrameUnit) String() string {
	switch u {
	case Rows:
		return "ROWS"
	case Range:
		return "RANGE"
	}

	panic("unknown frame unit")
}

type NotKeyword bool

const (
	NotKeywordAbsent  NotKeyword = false
	NotKeywordPresent NotKeyword = true
)

type ShiftOp int

const (
	LeftShift ShiftOp = iota
	RightShift
)
