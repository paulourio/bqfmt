package errors

import (
	"fmt"
	"strings"

	"github.com/paulourio/bqfmt/zetasql/literal"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func FormatError(err error, sql string) string {
	lines := strings.Split(sql, "\n")
	desc, line, col := findDescription(err, lines)
	data := lines[line-1]
	ind := fmt.Sprintf("%s^", strings.Repeat(" ", col-1))

	return fmt.Sprintf(
		"ERROR: Syntax error: %s [at %d:%d]\n%s\n%s\n",
		desc, line, col, data, ind)
}

func findDescription(err error, lines []string) (desc string, line, col int) {
	switch e := err.(type) {
	case *Error:
		if e.Err == nil {
			desc, line, col = tokenError(e)
		} else {
			desc, line, col = findDescription(e.Err, lines)
		}
	case *literal.UnescapeError:
		desc = e.Msg
		line, col = computeLineCol(lines, e.Offset)
	default:
		desc = e.Error()
	}

	return
}

func computeLineCol(lines []string, offset int) (line, col int) {
	for line < len(lines) && offset > len(lines[line])-1 {
		offset -= len(lines[line]) + 1
		line++
	}

	col = offset
	line++
	col++

	return
}

func tokenError(err *Error) (desc string, line, col int) {
	desc = fmt.Sprintf("Unexpected %s", DescribeToken(err.ErrorToken))
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
		token.TokMap.Type("missing_whitespace_float_and_alias"):
		desc = "Missing whitespace between literal and alias"
	case token.TokMap.Type("illegal_character"):
		desc = fmt.Sprintf(`Illegal input character "%s"`,
			string(err.ErrorToken.Lit))
	default:
		if isEOFexpected(err.ExpectedTokens) {
			desc = fmt.Sprintf("Expected end of input but got %s %s",
				strings.ReplaceAll(
					token.TokMap.Id(err.ErrorToken.Type),
					"_", " "),
				string(err.ErrorToken.Lit))
		}
	}

	return
}

func isEOFexpected(expected []string) bool {
	for _, sym := range expected {
		if sym == "$" {
			return true
		}
	}

	return false
}
