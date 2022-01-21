package ast

import (
	"fmt"
)

type Wrapped struct {
	Value interface{}
	Loc   Loc
}

func WrapWithLoc(value interface{}, start, end int) (*Wrapped, error) {
	return &Wrapped{value, Loc{start, end}}, nil
}

func (w *Wrapped) String() string {
	return fmt.Sprintf("Wrapped(%v, %v)", w.Value, w.Loc)
}
