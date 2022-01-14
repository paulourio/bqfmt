package zetasql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/paulourio/bqfmt/zetasql/lexer"
	"github.com/paulourio/bqfmt/zetasql/token"
)

type lexerTestCase struct {
	input  string
	tokens []string
}

var lexerTestCases = []lexerTestCase{
	{
		input:  "abcdef_",
		tokens: []string{"identifier(abcdef_)"},
	},
	{
		input:  "`5Customers`",
		tokens: []string{"identifier(`5Customers`)"},
	},
	{
		input:  "dataField",
		tokens: []string{"identifier(dataField)"},
	},
	{
		input:  "_dataField1",
		tokens: []string{"identifier(_dataField1)"},
	},
	{
		input:  "ADGROUP",
		tokens: []string{"identifier(ADGROUP)"},
	},
	{
		input:  "`tableName~`",
		tokens: []string{"identifier(`tableName~`)"},
	},
	{
		input:  "`GROUP`",
		tokens: []string{"identifier(`GROUP`)"},
	},
	{
		input:  "`project-name:dataset_name.table_name`",
		tokens: []string{"identifier(`project-name:dataset_name.table_name`)"},
	},
	{
		input: `"abc" "(abc)" "it's" 'it\'s' 'Title: "Boy"'`,
		tokens: []string{
			`string_literal("abc")`,
			`string_literal("(abc)")`,
			`string_literal("it's")`,
			`string_literal('it\'s')`,
			`string_literal('Title: "Boy"')`,
		},
	},
	{
		input: "\"\\``\" '\"' ' ` \" \\' \" '",
		tokens: []string{
			"string_literal(\"\\``\")",
			"string_literal('\"')",
			"string_literal(' ` \" \\' \" ')",
		},
	},
	{
		input: "'a\nb'",
		tokens: []string{
			"unterminated_string_literal('a)",
			"unterminated_bytes_literal(b')",
		},
	},
	{
		input: "\"a\nb\"",
		tokens: []string{
			"unterminated_string_literal(\"a)",
			"unterminated_bytes_literal(b\")",
		},
	},
	{
		input:  "'''a\nb'''",
		tokens: []string{"string_literal('''a\nb''')"},
	},
	{
		input:  "'''two\nlines'''",
		tokens: []string{"string_literal('''two\nlines''')"},
	},
	{
		input:  `"""abc"""`,
		tokens: []string{`string_literal("""abc""")`},
	},
	{
		input: `"""abc"\\""def""", """abc""\\"def"""`,
		tokens: []string{
			`string_literal("""abc"\\""def""")`,
			`,(,)`,
			`string_literal("""abc""\\"def""")`},
	},
	{
		input: `'''abc'\\''def''', '''abc''\\'def'''`,
		tokens: []string{
			`string_literal('''abc'\\''def''')`,
			`,(,)`,
			`string_literal('''abc''\\'def''')`},
	},
	{
		input:  `'''it's'''`,
		tokens: []string{`string_literal('''it's''')`},
	},
	{
		input:  `'''Title:"Boy"'''`,
		tokens: []string{`string_literal('''Title:"Boy"''')`},
	},
	{
		input:  `'''why\?'''`,
		tokens: []string{`string_literal('''why\?''')`},
	},
	{
		input:  `R"abc+"`,
		tokens: []string{`string_literal(R"abc+")`},
	},
	{
		input:  `r'''abc+'''`,
		tokens: []string{`string_literal(r'''abc+''')`},
	},
	{
		input:  `R"""abc+"""`,
		tokens: []string{`string_literal(R"""abc+""")`},
	},
	{
		input:  `'f\\(abc,(.*),def\\)'`,
		tokens: []string{`string_literal('f\\(abc,(.*),def\\)')`},
	},
	{
		input:  `r'f\(abc,(.*),def\)'`,
		tokens: []string{`string_literal(r'f\(abc,(.*),def\)')`},
	},
	{
		input:  `B"abc"`,
		tokens: []string{`bytes_literal(B"abc")`},
	},
	{
		input:  `B'''abc'''`,
		tokens: []string{`bytes_literal(B'''abc''')`},
	},
	{
		input:  `b"""abc"""`,
		tokens: []string{`bytes_literal(b"""abc""")`},
	},
	{
		input:  `br'abc+'`,
		tokens: []string{`bytes_literal(br'abc+')`},
	},
	{
		input:  `RB"abc+"`,
		tokens: []string{`bytes_literal(RB"abc+")`},
	},
	{
		input:  `RB'''abc'''`,
		tokens: []string{`bytes_literal(RB'''abc''')`},
	},
	{
		input:  `123`,
		tokens: []string{`integer_literal(123)`},
	},
	{
		input:  `0xABC`,
		tokens: []string{`integer_literal(0xABC)`},
	},
	{
		input:  `-123`,
		tokens: []string{`-(-)`, `integer_literal(123)`},
	},
	{
		input:  `+123`,
		tokens: []string{`+(+)`, `integer_literal(123)`},
	},
	{
		input:  `-0x12`,
		tokens: []string{`-(-)`, `integer_literal(0x12)`},
	},
	{
		input:  `066`,
		tokens: []string{`integer_literal(066)`},
	},
	{
		input:  `123.456e-67`,
		tokens: []string{`floating_point_literal(123.456e-67)`},
	},
	{
		input:  `.1E4`,
		tokens: []string{`floating_point_literal(.1E4)`},
	},
	{
		input:  `58.`,
		tokens: []string{`floating_point_literal(58.)`},
	},
	{
		input:  `4e2`,
		tokens: []string{`floating_point_literal(4e2)`},
	},
	{
		input:  `TRUE`,
		tokens: []string{`boolean_literal(TRUE)`},
	},
	{
		input:  `False`,
		tokens: []string{`boolean_literal(False)`},
	},
	{
		input:  "# comment\n",
		tokens: []string{"comment(# comment\n)"},
	},
	{
		input:  `# comment`,
		tokens: []string{`comment(# comment)`},
	},
	{
		input:  "#\u00a0comment",
		tokens: []string{"comment(#\u00a0comment)"},
	},
	{
		input: `""""a""""`,
		tokens: []string{
			"string_literal(\"\"\"\"a\"\"\")",
			"unterminated_string_literal(\")",
		},
	},
	{
		input:  `'a\b'`,
		tokens: []string{"string_literal('a\\b')"},
	},
	{
		input:  `br"\U1234"`,
		tokens: []string{`bytes_literal(br"\U1234")`},
	},
	{
		input: `132abc`,
		tokens: []string{
			`missing_whitespace_int_and_alias(132a)`,
			`identifier(bc)`,
		},
	},
	{
		input: `select $d`,
		tokens: []string{
			`select(select)`,
			`illegal_character($)`,
			`identifier(d)`,
		},
	},
	{
		input: "b\"\"\"line1\nline2\\\nline3\"\"\"",
		tokens: []string{
			"bytes_literal(b\"\"\"line1\nline2\\\nline3\"\"\")",
		},
	},
}

func TestString(t *testing.T) {
	for _, test := range lexerTestCases {
		t.Run(fmt.Sprintf("%#v", test.input), func(t *testing.T) {
			l := lexer.NewLexer([]byte(test.input))
			var output []string
			for {
				tok := l.Scan()
				if tok.Type == token.EOF {
					break
				}
				output = append(output, tokenString(tok))
			}
			// fmt.Println(test.input)
			assert.Equal(t, test.tokens, output)
		})
	}
}

func tokenString(tok *token.Token) string {
	return fmt.Sprintf("%s(%s)", token.TokMap.Id(tok.Type), tok.Lit)
}
