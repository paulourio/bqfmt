package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulourio/bqfmt/zetasql/ast"
	"github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/lexer"
	"github.com/paulourio/bqfmt/zetasql/parser"
)

var (
	debug = flag.Bool("d", false, "dump the AST tree instead of formatting")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	path := flag.Arg(0)
	if _, err := os.Stat(path); err != nil {
		usage()

		return
	}

	input, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	s := lexer.NewLexer(input).WithoutComment()
	p := parser.NewParser()

	r, err := p.Parse(s)
	if err != nil {
		fmt.Println(errors.FormatError(err, string(input)))

		return
	}

	if *debug {
		fmt.Println(r.(ast.NodeStringer).DebugString(""))
	} else {
		fmt.Println(ast.Sprint(r.(ast.NodeHandler), nil))
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: bqfmt [flags] path\n")
	flag.PrintDefaults()
}
