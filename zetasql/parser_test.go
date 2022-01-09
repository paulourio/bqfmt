package zetasql_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paulourio/bqfmt/zetasql/ast"
	"github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/lexer"
	"github.com/paulourio/bqfmt/zetasql/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	testPath := "./testdata/"
	st, err := os.Stat(testPath)
	if err != nil {
		t.Fatal(err)
	}
	if !st.IsDir() {
		t.Fatalf("%s must be a directory of .test files", testPath)
	}
	filepath.Walk(
		testPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			if info.IsDir() {
				if path == testPath {
					return nil
				}
				return fs.SkipDir
			}
			if strings.HasSuffix(path, ".test") {
				fmt.Printf("found %s", path)
				d, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				runTest(t, path, string(d))
			}
			return nil
		})
}

func runTest(t *testing.T, path string, input string) {
	cases := strings.Split(input, "\n==\n")
	for i, testCase := range cases {
		name := fmt.Sprintf("%s:Case#%d", path, i+1)
		elements := strings.Split(testCase+"\n", "\n--\n")
		var (
			input           []byte
			expectedDump    string
			expectedUnparse string
			expectedError   string
		)
		switch len(elements) {
		case 2:
			// Syntax error test case
			input = []byte(elements[0][1:])
			expectedError = elements[1] + "\n"
		case 3:
			// Regular parse test case
			input = []byte(elements[0][1:])
			expectedDump = elements[1] + "\n"
			expectedUnparse = elements[2]
		default:
			t.Fatalf("Test case %s must have three sections, but found %d:",
				name, len(elements))
		}
		t.Run(name, func(t *testing.T) {
			l := lexer.NewLexer(input).WithoutComment()
			p := parser.NewParser()
			r, err := p.Parse(l)
			if err != nil {
				fmt.Println("Input:")
				fmt.Println(string(input))
				fmt.Println("---")
				fmt.Println("Error")
				fmt.Println(err.(*errors.Error).String())
				fmt.Println("Expected Error")
				fmt.Println(expectedError)
				fmt.Println("---")
				t.Fatalf("Error: %v", err)
			}

			dump := r.(ast.NodeStringer).DebugString("")
			if expectedDump != dump {
				fmt.Println("Input:")
				fmt.Println(string(input))
				fmt.Println("---")
				fmt.Println("Dump:")
				fmt.Println(dump)
				fmt.Println("---")
			}

			unparsed := ast.Unparse(r.(ast.NodeHandler))
			if unparsed != expectedUnparse {
				fmt.Println("Unparsed:")
				fmt.Println(unparsed)
				fmt.Println("---")
				fmt.Println("Expected:")
				fmt.Println(expectedUnparse)
				fmt.Println("---")
			}

			assert.Equal(t, expectedDump, dump)
			assert.Equal(t, expectedUnparse, unparsed)
		})
	}
}
