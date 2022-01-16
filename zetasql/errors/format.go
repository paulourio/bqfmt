package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/paulourio/bqfmt/zetasql/literal"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func FormatError(err error, sql string) string {
	lines := strings.Split(sql, "\n")
	desc, line, col := findDescription(err, lines)
	var data string
	if line <= 0 || line > len(lines)+1 {
		fmt.Println("FormatError ERROR! line=", line, " col=", col)
		fmt.Printf("FormatError ERROR! %s\n", err.Error())
		fmt.Println("Input:")
		fmt.Println(sql)
		fmt.Println("~~~")
		line = 1
		col = 1
	} else {
		if line == len(lines)+1 {
			data = ""
		} else {
			data = lines[line-1]
		}
	}
	ind := fmt.Sprintf("%s^", strings.Repeat(" ", col-1))

	return fmt.Sprintf(
		"ERROR: %s [at %d:%d]\n%s\n%s\n",
		desc, line, col, data, ind)
}

func findDescription(err error, lines []string) (desc string, line, col int) {
	switch e := err.(type) {
	case *Error:
		if e.Err == nil {
			desc, line, col = tokenError(lines, e)
		} else {
			desc, line, col = findDescription(e.Err, lines)
		}
	case *literal.UnescapeError:
		desc = e.Error()
		line, col = computeLineCol(lines, e.Offset)
	default:
		var serr *SyntaxError
		if errors.As(e, &serr) {
			desc = serr.Error()
			line, col = computeLineCol(lines, serr.Loc.Start)
		} else {
			desc = err.Error()
		}
	}

	fmt.Printf("findDescription(%#v) -> %s, %d:%d\n",
		err, desc, line, col)
	return
}

func tokenError(lines []string, err *Error) (desc string, line, col int) {
	if len(err.ExpectedTokens) > 0 {
		expected := "" // strings.Join(err.ExpectedTokens, ", ")
		if isExpected("JOIN", err.ExpectedTokens) {
			expected = "keyword JOIN"
		}
		tok := strings.ToUpper(string(err.ErrorToken.Lit))
		if token.IsReservedKeyword(tok) {
			desc = fmt.Sprintf("Unexpected keyword %s", tok)
		} else if expected != "" {
			desc = fmt.Sprintf("Expected %s but got %s", expected,
				DescribeToken(err.ErrorToken))
		} else {
			desc = fmt.Sprintf("Unexpected %s",
				DescribeToken(err.ErrorToken))
		}
	} else {
		desc = fmt.Sprintf("Unexpected %s", DescribeToken(err.ErrorToken))
	}
	line = err.ErrorToken.Pos.Line
	col = err.ErrorToken.Pos.Column

	switch err.ErrorToken.Type {
	case token.TokMap.Type("$"):
		line++
		col = 1
		desc = "Unexpected end of statement"
	case token.TokMap.Type("unterminated_escaped_identifier"):
		desc = "Unclosed identifier literal"
	case token.TokMap.Type("unterminated_raw_bytes_literal"),
		token.TokMap.Type("unterminated_bytes_literal"),
		token.TokMap.Type("unterminated_raw_string_literal"),
		token.TokMap.Type("unterminated_string_literal"),
		token.TokMap.Type("unterminated_triple_quoted_bytes_literal"),
		token.TokMap.Type("unterminated_triple_quoted_raw_bytes_literal"),
		token.TokMap.Type("unterminated_triple_quoted_raw_string_literal"),
		token.TokMap.Type("unterminated_triple_quoted_string_literal"):
		tname := token.TokMap.Id(err.ErrorToken.Type)
		tname = strings.ReplaceAll(tname, "_", " ")
		tname = strings.ReplaceAll(tname, "unterminated", "Unclosed")
		tname = strings.ReplaceAll(tname, "triple quoted", "triple-quoted")
		desc = tname
	case token.TokMap.Type("unterminated_comment"):
		desc = "Unclosed comment"
	case token.TokMap.Type("missing_whitespace_int_and_alias"),
		token.TokMap.Type("missing_whitespace_hex_and_alias"),
		token.TokMap.Type("missing_whitespace_float_and_alias"):
		desc = "Missing whitespace between literal and alias"
		line, col = computeMissingWhitespaceLocation(lines, err.ErrorToken)
	case token.TokMap.Type("illegal_character"):
		desc = fmt.Sprintf(`Illegal input character "%s"`,
			string(err.ErrorToken.Lit))
	default:
		if isExpected("$", err.ExpectedTokens) {
			value := string(err.ErrorToken.Lit)
			if !strings.HasPrefix(value, `'`) &&
				!strings.HasPrefix(value, `"`) {
				value = fmt.Sprintf("%q", value)
			}
			desc = fmt.Sprintf("Expected end of input but got %s %s",
				strings.ReplaceAll(
					token.TokMap.Id(err.ErrorToken.Type),
					"_", " "),
				value,
			)
		}
	}

	desc = "Syntax error: " + desc
	return
}

func isExpected(elem string, expected []string) bool {
	elem = strings.ToUpper(elem)
	for _, sym := range expected {
		if strings.ToUpper(sym) == elem {
			return true
		}
	}

	return false
}

func computeMissingWhitespaceLocation(
	lines []string, tok *token.Token) (line, col int) {
	// Tokenizer consumes only the first wrong character from input,
	// so the offset should always be at the end.
	return computeLineCol(lines, tok.Offset+len(tok.Lit)-1)
}

func computeLineCol(lines []string, offset int) (line, col int) {
	for line < len(lines) && offset > len(lines[line])-1 {
		offset -= len(lines[line]) + 1
		line++
	}

	col = offset
	col++
	line++
	fmt.Printf("computeLineCol(offset=%d) -> %d:%d\n", offset, line, col)

	return
}
