package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/paulourio/bqfmt/zetasql/literal"
	"github.com/paulourio/bqfmt/zetasql/token"
)

func FormatError(err error, sql string) string {
	var data string

	lines := strings.Split(sql, "\n")
	desc, line, col := findDescription(err, lines)

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
	var (
		genErr *Error
		litErr *literal.UnescapeError
		synErr *SyntaxError
	)

	if errors.As(err, &genErr) {
		if genErr.Err == nil {
			desc, line, col = tokenError(lines, genErr)
		} else {
			desc, line, col = findDescription(genErr.Err, lines)
		}

		return
	}

	if errors.As(err, &litErr) {
		desc = litErr.Error()
		line, col = computeLineCol(lines, litErr.Offset)

		return
	}

	if errors.As(err, &synErr) {
		desc = synErr.Error()
		line, col = computeLineCol(lines, synErr.Loc.Start)
	} else {
		desc = err.Error()
	}

	return
}

func tokenError(lines []string, err *Error) (desc string, line, col int) {
	tokDesc := describeUnexpectedToken(err.ErrorToken)

	if len(err.ExpectedTokens) > 0 {
		expected := ""
		if isExpected("JOIN", err.ExpectedTokens) {
			expected = "keyword JOIN"
		}

		tok := strings.ToUpper(string(err.ErrorToken.Lit))
		if token.IsReservedKeyword(tok) { //nolint:gocritic
			desc = fmt.Sprintf("Unexpected keyword %s", tok)
		} else if expected != "" {
			desc = fmt.Sprintf("Expected %s but got %s", expected,
				tokDesc)
		} else {
			desc = fmt.Sprintf("Unexpected %s",
				tokDesc)
		}
	} else {
		desc = fmt.Sprintf("Unexpected %s", tokDesc)
	}

	line = err.ErrorToken.Pos.Line
	col = err.ErrorToken.Pos.Column

	switch err.ErrorToken.Type {
	case token.INVALID:
		desc = "Unexpected unknown/invalid token"
	case token.EOF, token.TokMap.Type("$"):
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
			desc = fmt.Sprintf("Expected end of input but got %s",
				describeUnexpectedToken(err.ErrorToken))
		}
	}

	desc = "Syntax error: " + desc

	return
}

func describeUnexpectedToken(tok *token.Token) string {
	name := strings.ReplaceAll(token.TokMap.Id(tok.Type), "_", " ")
	value := string(tok.Lit)
	nameEqualsValue := name == value

	if tok.IsReservedKeyword() {
		return fmt.Sprintf("keyword %s", strings.ToUpper(value))
	}

	if !strings.HasPrefix(value, `'`) &&
		!strings.HasPrefix(value, `"`) {
		value = fmt.Sprintf("%q", value)
	}

	if nameEqualsValue {
		return value
	}

	return fmt.Sprintf("%s %s", name, value)
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

	return
}
