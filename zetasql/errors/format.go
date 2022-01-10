package errors

import (
	"fmt"
	"strings"
)

func FormatError(err error, sql string) string {
	e := err.(*Error)
	line := e.ErrorToken.Pos.Line
	col := e.ErrorToken.Pos.Column
	desc := fmt.Sprintf("Unexpected %s", DescribeToken(e.ErrorToken))
	lines := strings.Split(sql, "\n")
	data := lines[line-1]

	if e.ErrorToken.Type == 1 {
		line++
		col = 1
		desc = "Unexpected end of statement"
		data = ""
	}

	ind := fmt.Sprintf("%s^", strings.Repeat(" ", col-1))

	return fmt.Sprintf(
		"ERROR: Syntax error: %s [at %d:%d]\n%s\n%s\n",
		desc, line, col, data, ind)
}
