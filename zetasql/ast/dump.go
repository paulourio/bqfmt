package ast

import (
	"fmt"
	"strings"
)

type dumper struct {
	node         NodeHandler
	sep          string
	maxDepth     int
	currentDepth int
	sql          string
	out          strings.Builder
}

func newDumper(n NodeHandler, sep string, maxDepth int, sql string) *dumper {
	return &dumper{n, sep, maxDepth, 0, sql, strings.Builder{}}
}

func (d *dumper) Dump() {
	if !d.dumpNode() {
		return
	}

	d.currentDepth++

	for _, c := range d.node.Children() {
		d.node = c
		d.Dump()
	}

	d.currentDepth--
}

func (d *dumper) String() string {
	return d.out.String()
}

func (d *dumper) dumpNode() bool {
	d.indent()
	d.out.WriteString(d.node.(NodeStringer).SingleNodeDebugString())
	d.out.WriteString(fmt.Sprintf(" [%s]", locationRangeAsString(d.node)))
	d.out.WriteString(d.sep)

	if d.currentDepth >= d.maxDepth {
		d.indent()
		d.out.WriteString(
			fmt.Sprintf("  Subtree skipped (reached max depth %d)",
				d.maxDepth))

		return false
	}

	return true
}

func (d *dumper) indent() {
	d.out.WriteString(strings.Repeat(" ", 2*d.currentDepth))
}

func locationRangeAsString(n NodeHandler) string {
	s := n.StartLoc()
	e := n.EndLoc()

	return fmt.Sprintf("%d-%d", s, e)
}
