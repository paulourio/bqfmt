package parser_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paulourio/bqfmt/zetasql/ast"
	zerrors "github.com/paulourio/bqfmt/zetasql/errors"
	"github.com/paulourio/bqfmt/zetasql/lexer"
	"github.com/paulourio/bqfmt/zetasql/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Parallel()

	testPath := "./testdata/"

	st, err := os.Stat(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if !st.IsDir() {
		t.Fatalf("%s must be a directory of .test files", testPath)
	}

	err = filepath.Walk(
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
				fmt.Printf("found %s\n", path)

				d, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}

				runTest(t, path, string(d))
			}
			return nil
		})

	if err != nil {
		t.Fatal(err)
	}
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
			input = []byte(getInput(elements[0]))
			expectedError = elements[1]
		case 3:
			// Regular parse test case
			input = []byte(getInput(elements[0]))
			expectedDump = elements[1] + "\n"
			expectedUnparse = elements[2]
		default:
			t.Logf("Test case %s must have three sections, but found %d:",
				name, len(elements))
			continue
		}

		t.Run(name, func(t *testing.T) {
			l := lexer.NewLexer(input).WithoutComment()
			p := parser.NewParser()

			var (
				dump, errMsg, unparsed   string
				okErr, okDump, okUnparse bool
			)

			r, err := p.Parse(l)
			if err != nil {
				errMsg = zerrors.FormatError(err, string(input))
			}

			if expectedError == "" {
				okErr = assert.Nil(t, err)
			} else {
				okErr = assert.Equal(t, expectedError, errMsg)
				if !okErr {
					t.Log("Input:")
					t.Log(string(input))
					t.Log("---")
					t.Log("Raw")

					var perr *zerrors.Error

					if errors.As(err, &perr) {
						t.Log(perr.String())
					} else {
						t.Log(perr)
					}

					t.Log("original error:", err)
				}
			}

			if err == nil {
				dump = r.(ast.NodeStringer).DebugString("")
			}

			if err == nil {
				unparsed = ast.Unparse(r.(ast.NodeHandler))
			}

			okDump = assert.Equal(t, expectedDump, dump)
			okUnparse = assert.Equal(t, expectedUnparse, unparsed)

			if !okErr || !okDump || !okUnparse {
				t.Logf("Input: %s", string(input))
				if err != nil {
					t.Logf("Error: %s", err.Error())
				}
				if expectedDump != dump {
					t.Log("Input:")
					t.Log(string(input))
					t.Log("---")
					t.Log("Dump:")
					t.Log(dump)
					t.Log("---")
				}
				if unparsed != expectedUnparse {
					t.Log("Unparsed:")
					t.Log(unparsed)
					t.Log("---")
					t.Log("Expected:")
					t.Log(expectedUnparse)
					t.Log("---")
				}
			}

		})
	}
}

func getInput(section string) string {
	lines := strings.Split(section[1:], "\n")
	i := 0

	for ; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "# ") &&
			lines[i] != "" &&
			!strings.HasPrefix(lines[i], "[") {
			break
		}
	}

	return strings.Join(lines[i:], "\n")
}
