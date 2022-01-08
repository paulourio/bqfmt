package ast

type Wrapped struct {
	Value interface{}
	Loc   Loc
}

func WrapWithLoc(value interface{}, start, end int) (*Wrapped, error) {
	return &Wrapped{value, Loc{start, end}}, nil
}
