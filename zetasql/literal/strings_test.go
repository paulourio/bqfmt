package literal_test

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"

	"github.com/paulourio/bqfmt/zetasql/literal"
)

type stringTestCase struct {
	Inputs   []string
	Expected string
	Error    error
}

type bytesTestCase struct {
	Inputs   []string
	Expected []byte
	Error    error
}

var octalTestCases = []stringTestCase{
	{
		Inputs: []string{`'\0'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Octal escape must be followed by 3 octal digits but saw: \0`,
			Offset: 1},
	},
	{
		Inputs: []string{`'''abc\0'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Octal escape must be followed by 3 octal digits but saw: \0`,
			Offset: 6},
	},
	{
		Inputs: []string{`'\00'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Octal escape must be followed by 3 octal digits but saw: \00`,
			Offset: 1},
	},
	{
		Inputs: []string{`'''ab\00'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Octal escape must be followed by 3 octal digits but saw: \00`,
			Offset: 5},
	},
	{
		Inputs: []string{`'a\008'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Octal escape must be followed by 3 octal digits but saw: \008`,
			Offset: 2},
	},
	{
		Inputs: []string{`'''\008'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Octal escape must be followed by 3 octal digits but saw: \008`,
			Offset: 3},
	},
	{
		Inputs: []string{`'\400'`},
		Error: &literal.UnescapeError{
			Msg:    `Illegal escape sequence: \4`,
			Offset: 1},
	},
	{
		Inputs: []string{`'''\400'''`},
		Error: &literal.UnescapeError{
			Msg:    `Illegal escape sequence: \4`,
			Offset: 3},
	},
	{
		Inputs: []string{`'\777'`},
		Error: &literal.UnescapeError{
			Msg:    `Illegal escape sequence: \7`,
			Offset: 1},
	},
	{
		Inputs: []string{`'''\777'''`},
		Error: &literal.UnescapeError{
			Msg:    `Illegal escape sequence: \7`,
			Offset: 3},
	},
}

var hexTestCases = []stringTestCase{
	{
		Inputs: []string{`'\x'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x`,
			Offset: 1},
	},
	{
		Inputs: []string{`'\x0'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x0`,
			Offset: 1},
	},
	{
		Inputs: []string{`'\x0G'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x0G`,
			Offset: 1},
	},
	{
		Inputs: []string{`'''\x'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x`,
			Offset: 3},
	},
	{
		Inputs: []string{`'''\x0'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x0`,
			Offset: 3},
	},
	{
		Inputs: []string{`'''\x0G'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x0G`,
			Offset: 3},
	},
}

var invalidStringTestCases = []stringTestCase{
	{
		Inputs: []string{`'\'`},
		Error:  &literal.UnescapeError{"String must end with '", 3},
	},
	{
		Inputs: []string{`'abc\'`},
		Error:  &literal.UnescapeError{"String must end with '", 6},
	},
	{
		Inputs: []string{`'''\'''`},
		Error:  &literal.UnescapeError{"String must end with '''", 7},
	},
	{
		Inputs: []string{`'''abc\'''`},
		Error:  &literal.UnescapeError{"String must end with '''", 10},
	},
	{
		Inputs: []string{`A`, `'`, `"`, `a'`, `a"`},
		Error:  literal.ErrInvalidStringLiteral,
	},
	{
		Inputs: []string{`'''`, `''''`, `'''''`, `'''abc'`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped '`,
			Offset: 1,
		},
	},
	{
		Inputs: []string{`"""`, `""""`, `"""""`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped "`,
			Offset: 1,
		},
	},
	{
		Inputs: []string{`'''''''`, `''''''''''`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped '''`,
			Offset: 3,
		},
	},
	{
		Inputs: []string{`"""""""`, `""""""""""`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped """`,
			Offset: 3,
		},
	},
	{
		Inputs: []string{`'abc'def'`, `'abc''def'`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped '`,
			Offset: 4,
		},
	},
	{
		Inputs: []string{`'''abc'''def'''`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped '''`,
			Offset: 6,
		},
	},
	{
		Inputs: []string{`"""abc"""def"""`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped """`,
			Offset: 6,
		},
	},
	{
		Inputs: []string{`'abc`, `"abc`, `'''abc`, `"""abc`, `abc'`, `abc"`,
			`abc'''`, `abc"""`, `"abc'`, `'abc"`, `'''abc"`,
			`'''abc"""`, `"""abc'''`, `"""abc'`, "`abc`",
			// Trailing escapes.
			`'\`, `"\`, `''''''\`, `""""""\`, `''\\`, `""\\`, `''''''\\`,
			`""""""\\`},
		Error: literal.ErrInvalidStringLiteral,
	},
	{
		Inputs: []string{string([]byte{0, 'a', 'b'})}, // Unescaped 0 byte.
		Error:  literal.ErrInvalidStringLiteral,
	},
	{
		Inputs: []string{
			// These are C-escapes to define invalid strings.
			"'\xc1'", "'\xca'", "'\xcc'", "'\xFA'",
			"'\xc1\xca\x1b\x62\x19o\xcc\x04'", "'\xc2\xc0'",
			// These are all valid prefixes for UTF8 characters, but
			// characters are not complete.
			"'\xc2'", "'\xc3'", // Should be a 2-byte UTF8 character.
			"'\xe0'", "'\xe0\xac'", // Should be a 3-byte UTF8 character.
			// Should be a 4-byte UTF8 character.
			"'\xf0'", "'\xf0\x90'", "'\xf0\x90\x80'",
		},
		Error: &literal.UnescapeError{
			Msg:    `Structurally invalid UTF8 string`,
			Offset: 1,
		},
	},
	{
		Inputs: []string{`r"\"`},
		Error:  &literal.UnescapeError{`String must end with "`, 4},
	},
	{
		Inputs: []string{`r"\\\"`},
		Error:  &literal.UnescapeError{`String must end with "`, 6},
	},
	{
		Inputs: []string{`r'''`},
		Error:  &literal.UnescapeError{`String cannot contain unescaped '`, 2},
	},
	{
		Inputs: []string{`r\""`, `r`, `rb""`, `b""`},
		Error:  literal.ErrInvalidStringLiteral,
	},
	{
		Inputs:   []string{`''''''`, `""""""`},
		Expected: ``,
	},
	{
		Inputs:   []string{`''''"'''`},
		Expected: `'"`,
	},
	{
		Inputs:   []string{`'''''\'''\''''`},
		Expected: `''''''`,
	},
	{
		Inputs:   []string{`'''\''''`},
		Expected: `'`,
	},
	{
		Inputs:   []string{`'''\'\''''`},
		Expected: `''`,
	},
	{
		Inputs:   []string{`''''a'''`},
		Expected: `'a`,
	},
	{
		Inputs:   []string{`""""a"""`},
		Expected: `"a`,
	},
	{
		Inputs:   []string{`'''''a'''`},
		Expected: `''a`,
	},
	{
		Inputs:   []string{`"""""a"""`},
		Expected: `""a`,
	},
	{
		Inputs:   []string{`r'\n'`, `r"\n"`, `r'''\n'''`, `r"""\n"""`},
		Expected: `\n`,
	},
	{
		Inputs:   []string{`'\n'`, `"\n"`, `'''\n'''`, `"""\n"""`},
		Expected: "\n",
	},
	{
		Inputs:   []string{`r'\e'`, `r"\e"`, `r'''\e'''`, `r"""\e"""`},
		Expected: `\e`,
	},
	{
		Inputs: []string{`'\e'`, `"\e"`},
		Error:  &literal.UnescapeError{`Illegal escape sequence: \e`, 1},
	},
	{
		Inputs: []string{`'''\e'''`, `"""\e"""`},
		Error:  &literal.UnescapeError{`Illegal escape sequence: \e`, 3},
	},
	{
		Inputs:   []string{`r'\x0'`, `r"\x0"`, `r'''\x0'''`, `r"""\x0"""`},
		Expected: `\x0`,
	},
	{
		Inputs: []string{`'\x0'`, `"\x0"`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x0`,
			Offset: 1},
	},
	{
		Inputs: []string{`'''\x0'''`, `"""\x0"""`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Hex escape must be followed by 2 hex digits but saw: \x0`,
			Offset: 3},
	},
	{
		Inputs:   []string{`r'\''`},
		Expected: `\'`,
	},
	{
		Inputs:   []string{`'\''`},
		Expected: `'`,
	},
	{
		Inputs:   []string{`r"\""`},
		Expected: `\"`,
	},
	{
		Inputs:   []string{`"\""`},
		Expected: `"`,
	},
	{
		Inputs:   []string{`r'''''\''''`},
		Expected: `''\'`,
	},
	{
		Inputs:   []string{`'''''\''''`},
		Expected: `'''`,
	},
	{
		Inputs:   []string{`r"""""\""""`},
		Expected: `""\"`,
	},
	{
		Inputs:   []string{`"""""\""""`},
		Expected: `"""`,
	},
}

var newlineTestCases = []stringTestCase{
	{
		Inputs:   []string{"'''a\rb'''", "'''a\nb'''", "'''a\r\nb'''"},
		Expected: "a\nb",
	},
	{
		Inputs:   []string{"'''a\n\rb'''", "'''a\r\n\r\nb'''"},
		Expected: "a\n\nb",
	},
	{
		Inputs:   []string{`'''a\nb'''`},
		Expected: "a\nb",
	},
	{
		Inputs:   []string{`'''a\rb'''`},
		Expected: "a\rb",
	},
	{
		Inputs:   []string{`'''a\r\nb'''`},
		Expected: "a\r\nb",
	},
}

var utf8UnescapeTestCases = []stringTestCase{
	{
		Inputs:   []string{`"\u0030"`},
		Expected: "0",
	},
	{
		Inputs:   []string{`"\u00A3"`},
		Expected: "\xC2\xA3",
	},
	{
		Inputs:   []string{`"\u22FD"`},
		Expected: "\xE2\x8B\xBD",
	},
	{
		Inputs:   []string{`"\ud7FF"`},
		Expected: "\xED\x9F\xBF",
	},
	{
		Inputs:   []string{`"\u22FD"`},
		Expected: "\xE2\x8B\xBD",
	},
	{
		Inputs:   []string{`"\U00010000"`},
		Expected: "\xF0\x90\x80\x80",
	},
	{
		Inputs:   []string{`"\U0000E000"`},
		Expected: "\xEE\x80\x80",
	},
	{
		Inputs:   []string{`"\U0001DFFF"`},
		Expected: "\xF0\x9D\xBF\xBF",
	},
	{
		Inputs:   []string{`"\U0010FFFD"`},
		Expected: "\xF4\x8F\xBF\xBD",
	},
	{
		Inputs:   []string{`"\xAbCD"`},
		Expected: "\xc2\xab" + "CD",
	},
	{
		Inputs:   []string{`"\253CD"`},
		Expected: "\xc2\xab" + "CD",
	},
	{
		Inputs:   []string{`"\x4141"`},
		Expected: "A41",
	},
	{
		Inputs: []string{`"\u1"`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`\u must be followed by 4 hex digits but saw: \u1`,
			Offset: 1},
	},
	{
		Inputs: []string{`"\U1"`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`\U must be followed by 8 hex digits but saw: \U1`,
			Offset: 1},
	},
	{
		Inputs: []string{`"\Uffffff"`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`\U must be followed by 8 hex digits but saw: \Uffffff`,
			Offset: 1},
	},
	{
		Inputs: []string{`"\UFFFFFFFF0"`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Value of \UFFFFFFFF exceeds Unicode limit (0x0010FFFF)`,
			Offset: 1},
	},
}

var invalidBytesCases = []bytesTestCase{
	{
		Inputs: []string{`A`, `b'A`, `'A'`, `"A"`, `'''A'''`, `"""A"""`},
		Error:  literal.ErrInvalidBytesLiteral,
	},
	{
		Inputs: []string{`b'k\u0030'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Unicode escape sequence \u cannot be used in bytes literals`,
			Offset: 3,
		},
	},
	{
		Inputs: []string{`b'''\u0030'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Unicode escape sequence \u cannot be used in bytes literals`,
			Offset: 4,
		},
	},
	{
		Inputs: []string{`b'\U00000030'`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Unicode escape sequence \U cannot be used in bytes literals`,
			Offset: 2,
		},
	},
	{
		Inputs: []string{`b'''qwerty\U00000030'''`},
		Error: &literal.UnescapeError{
			Msg: `Illegal escape sequence: ` +
				`Unicode escape sequence \U cannot be used in bytes literals`,
			Offset: 10,
		},
	},
	{
		Inputs: []string{"b\"\"\"line1\nline2\\\nline3\"\"\""},
		Error: &literal.UnescapeError{
			Msg:    `Illegal escaped newline`,
			Offset: 15,
		},
	},
}

var invalidRawBytesCases = []bytesTestCase{
	{
		Inputs: []string{`r''`, `r''''''`, `rrb''`, `brb''`, `rb'a\e"`, `rb`,
			`br`, `rb"`},
		Error: literal.ErrInvalidBytesLiteral,
	},
	{
		Inputs: []string{`rb"\"`},
		Error: &literal.UnescapeError{
			Msg:    `String must end with "`,
			Offset: 5,
		},
	},
	{
		Inputs: []string{`br"\\\"`},
		Error: &literal.UnescapeError{
			Msg:    `String must end with "`,
			Offset: 7,
		},
	},
	{
		Inputs: []string{`rb""""`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped "`,
			Offset: 3,
		},
	},
	{
		Inputs: []string{`rb"xyz"""`},
		Error: &literal.UnescapeError{
			Msg:    `String cannot contain unescaped "`,
			Offset: 6,
		},
	},
}

func TestParseingOfAllEscapeCharacters(t *testing.T) {
	validEscapes := []byte{
		'a', 'b', 'f', 'n', 'r', 't', 'v', '\\',
		'?', '"', '\'', '`', 'u', 'U', 'x', 'X',
	}

	for e := byte(0); e < 255; e++ {
		if contains(validEscapes, e) {
			if e == '\'' {
				testParseString(t, fmt.Sprintf(`"a\%c0010ffff"`, e))
			}
			testParseString(t, fmt.Sprintf(`'a\%c0010ffff'`, e))
			testParseString(t, fmt.Sprintf(`'''a\%c0010ffff'''`, e))
		} else if unicode.IsDigit(rune(e)) {
			if e <= '3' {
				testParseString(t, fmt.Sprintf(`'a\%c00b'`, e))
				testParseString(t, fmt.Sprintf(`'''a\%c00b'''`, e))
			} else {
				testInvalidString(t, fmt.Sprintf(`'a\%c00b'`, e),
					&literal.UnescapeError{
						Msg:    fmt.Sprintf(`Illegal escape sequence: \%c`, e),
						Offset: 2,
					})
				testInvalidString(t, fmt.Sprintf(`'''a\%c00b'''`, e),
					&literal.UnescapeError{
						Msg:    fmt.Sprintf(`Illegal escape sequence: \%c`, e),
						Offset: 4,
					})
			}
		} else {
			var msg string
			switch e {
			case '\n', '\r':
				msg = `Illegal escaped newline`
			default:
				msg = `Illegal escape sequence: `
			}
			testInvalidString(t, fmt.Sprintf(`'a\%cb'`, e), msg)
			testInvalidString(t, fmt.Sprintf(`'''a\%cb'''`, e), msg)
		}
	}
}
func TestParseHexEscapes(t *testing.T) {
	for i := 0; i < 256; i++ {
		lead := fmt.Sprintf("%X", i/16)[0]
		end := fmt.Sprintf("%x", i%16)[0]
		testParseString(t, fmt.Sprintf(`'\x%c%c'`, lead, end))
		testParseString(t, fmt.Sprintf(`'''\x%c%c'''`, lead, end))
		testParseString(t, fmt.Sprintf(`'\X%c%c'`, lead, end))
		testParseString(t, fmt.Sprintf(`'''\X%c%c'''`, lead, end))
	}
	runStringTestCases(t, hexTestCases)
}

func TestParseOctalEscapes(t *testing.T) {
	for i := 0; i < 256; i++ {
		lead := (i / 64) + '0'
		mid := (i/8)%8 + '0'
		end := i%8 + '0'
		testParseString(t, fmt.Sprintf(`'%c%c%c'`, lead, mid, end))
		testParseString(t, fmt.Sprintf(`"%c%c%c"`, lead, mid, end))
		testParseString(t, fmt.Sprintf(`'''%c%c%c'''`, lead, mid, end))
	}
	runStringTestCases(t, octalTestCases)
}

func TestInvalidBytes(t *testing.T) {
	runBytesTestCases(t, invalidBytesCases)
}

func TestInvalidRawBytes(t *testing.T) {
	runBytesTestCases(t, invalidRawBytesCases)
}

func TestInvalidString(t *testing.T) {
	runStringTestCases(t, invalidStringTestCases)
}

func TestNewlines(t *testing.T) {
	runStringTestCases(t, newlineTestCases)
}

func TestUTF8Unescape(t *testing.T) {
	runStringTestCases(t, utf8UnescapeTestCases)
}

func TestUnescapeString(t *testing.T) {
	runBytesTestCases(t, invalidBytesCases)
}

func TestUnescapeError(t *testing.T) {
	testInvalidString(t, `'\ude8c'`,
		`Illegal escape sequence: Unicode value \ude8c is invalid`)
}

func TestValidRawString(t *testing.T) {
	cases := []string{
		``,
		`1`,
		`\x53`,
		`\x123`,
		`\001`,
		`a\44'A`,
		`a\e`,
		`\ea`,
		`\U1234`,
		`\u`,
		`\xc2\\`,
		`f\(abc',(.*),def\ \?`,
		`a\\"b`,
	}
	for i, input := range cases {
		name := fmt.Sprintf("#%d", i+1)
		t.Run(name, func(t *testing.T) {
			testRawStringValue(t, input)
		})
	}
}

func testParseString(t *testing.T, input string) bool {
	_, err := literal.ParseString(input)
	return assert.Nil(t, err)
}

func testParseBytes(t *testing.T, input string) bool {
	_, err := literal.ParseBytes(input)
	return assert.Nil(t, err)
}

func testInvalidString(t *testing.T, input string, expected interface{}) bool {
	r, err := literal.ParseString(input)
	assert.Equal(t, "", r)
	switch e := expected.(type) {
	case error:
		return assert.Equal(t, e, err)
	case string:
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return assert.Contains(t, msg, e)
	default:
		panic("unexpected type")
	}
}

func testRawStringValue(t *testing.T, unquoted string) {
	quote := '"'
	if strings.Contains(unquoted, `"`) {
		quote = '\''
	}
	r, err := literal.ParseString(
		fmt.Sprintf(`r%c%s%c`, quote, unquoted, quote))
	assert.Nil(t, err)
	assert.Equal(t, unquoted, r)
	r, err = literal.ParseString(
		fmt.Sprintf(`r%c%s%c`, quote, unquoted, quote))
	assert.Nil(t, err)
	assert.Equal(t, unquoted, r)
}

func contains(list []byte, element byte) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == element {
			return true
		}
	}

	return false
}

func runStringTestCases(t *testing.T, tests []stringTestCase) {
	for i, test := range tests {
		name := fmt.Sprintf("#%d", i+1)
		t.Run(name, func(t *testing.T) {
			for _, input := range test.Inputs {
				result, err := literal.ParseString(input)
				if !assert.Equal(t, test.Expected, result) ||
					!assert.Equal(t, test.Error, err) {
					t.Log("Input:", input)
				}
			}
		})
	}
}

func runBytesTestCases(t *testing.T, tests []bytesTestCase) {
	for i, test := range tests {
		name := fmt.Sprintf("#%d", i+1)
		t.Run(name, func(t *testing.T) {
			for _, input := range test.Inputs {
				result, err := literal.ParseBytes(input)
				var resultOk bool
				if test.Expected == nil {
					resultOk = assert.Equal(t, []byte{}, result)
				} else {
					resultOk = assert.Equal(t, test.Expected, result)
				}
				if !assert.Equal(t, test.Error, err) || !resultOk {
					t.Log("Input:", input)
				}
			}
		})
	}
}
