package literal

import (
	"fmt"
	"strconv"
)

type Error struct {
	Msg    string
	Offset int
}

func ParseInt(s string) (int, error) {
	for i := 0; i < len(s); i++ {
		if !(s[i] >= '0' && s[i] <= 9) {
			return 0, &Error{
				Msg:    "Missing whitespace between literal and alias",
				Offset: i}
		}
	}

	return strconv.Atoi(s)
}

func ParseHex(s string) (int, error) {
	if len(s) <= 2 {
		return 0, &Error{"Hex int must have at least three digits", len(s) - 1}
	}

	if s[0] != '0' {
		return 0, &Error{"Hex int must begin with 0x", 0}
	}

	if s[1] != 'x' && s[1] != 'X' {
		return 0, &Error{"Hex int must begin with 0x", 1}
	}

	for i := 2; i < len(s); i++ {
		if !(s[i] >= '0' && s[i] <= '9') {
			return 0, &Error{
				Msg:    "Missing whitespace between literal and alias",
				Offset: i,
			}
		}
	}

	return strconv.Atoi(s)
}

func (e *Error) Error() string {
	return fmt.Sprintf("Syntax error: %s", e.Msg)
}
