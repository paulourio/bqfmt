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
		return "-"
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
	OrderAscending
	OrderDescending
)

type JoinType int

const (
	DefaultJoin JoinType = iota
	CrossJoin
	OuterJoin
	InnerJoin
	LeftJoin
	RightJoin
)

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
	ExistsSubqueryMedifier                         // EXISTS(SELECT ...)
)

type SetOp int

const (
	Union     SetOp = iota // UNION {ALL|DISTINCT}
	Except                 // EXCEPT {ALL|DISTINCT}
	Intersect              // INTERSECT {ALL|DISTINCt}
)

type BoundaryType int

const (
	UnboundedPreceding BoundaryType = iota
	OffsetPreceding
	CurrentRow
	OffsetFollowing
	UnboundedFollowing
)

type SampleSizeUnit int

const (
	RowsSampling SampleSizeUnit = iota
	PercentSampling
)

type TypeKind int

const (
	BoolKind TypeKind = iota
	IntegerKind
	FloatingPointKind
	OtherKind
)

type FrameUnit int

const (
	Rows FrameUnit = iota
	Range
)
