package literal

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidBytesLiteral  = errors.New("invalid bytes literal")
	ErrInvalidStringLiteral = errors.New("invalid string literal")
)

type StringStyle int

const (
	PreferSingleQuote StringStyle = iota
	PreferDoubleQuote
	AlwaysSingleQuote
	AlwaysDoubleQuote
)

type UnescapeError struct {
	Msg    string
	Offset int
}

func (e *UnescapeError) Error() string {
	return fmt.Sprintf("Syntax error: %s", e.Msg)
}

func ParseString(s string) (string, error) {
	isStrLiteral := maybeStringLiteral(s)
	isRawLiteral := maybeRawStringLiteral(s)

	if !isStrLiteral && !isRawLiteral {
		return "", ErrInvalidStringLiteral
	}

	offset := 0 // Offset to control the error position.

	if isRawLiteral {
		// Strip off the prefix 'r' from the raw string content before
		// parsing.
		s = s[1:]
		offset = 1
	}

	quotesLen := 1
	isTripleQuoted := maybeTripleQuotedStringLiteral(s)

	if isTripleQuoted {
		quotesLen = 3
	}

	offset += quotesLen

	e := &unescaper{
		source:    s,
		quotesLen: quotesLen,
		offset:    offset,
		isRaw:     isRawLiteral,
		isBytes:   false,
	}

	return e.Unescape()
}

func ParseBytes(s string) ([]byte, error) {
	isBytesLiteral := maybeBytesLiteral(s)
	isRawBytesLiteral := maybeRawBytesLiteral(s)

	if !isBytesLiteral && !isRawBytesLiteral {
		return []byte{}, ErrInvalidBytesLiteral
	}

	offset := 0 // Offset to control the error position.

	if isRawBytesLiteral {
		// Strip off the prefix 'r' from the raw string content before
		// parsing.
		s = s[2:]
		offset = 2
	} else {
		s = s[1:]
		offset = 1
	}

	quotesLen := 1
	isTripleQuoted := maybeTripleQuotedStringLiteral(s)

	if isTripleQuoted {
		quotesLen = 3
	}

	offset += quotesLen

	e := &unescaper{
		source:    s,
		quotesLen: quotesLen,
		offset:    offset,
		isRaw:     isRawBytesLiteral,
		isBytes:   isBytesLiteral,
	}

	r, err := e.Unescape()
	if err != nil {
		return []byte{}, err
	}

	return []byte(r), err
}

type unescaper struct {
	// Input
	source    string
	quotesLen int
	offset    int
	isRaw     bool
	isBytes   bool

	// Working data
	data   []byte // The string as bytes of slices.
	quotes []byte // The quotes of the string.
	out    []rune // The output.
	pos    int
}

func (u *unescaper) Unescape() (string, error) {
	str := []byte(u.source)
	u.quotes, u.data = str[:u.quotesLen], str[u.quotesLen:]

	if err := u.checkForClosingString(); err != nil {
		return "", err
	}

	// Strip off the closing quotes before unescaping.
	u.data = u.data[:len(u.data)-len(u.quotes)]

	if !u.isBytes {
		if !utf8.ValidString(u.source) {
			return "", &UnescapeError{
				Msg:    "Structurally invalid UTF8 string",
				Offset: u.offset,
			}
		}
	}

	u.out = make([]rune, 0, len(u.data))
	u.pos = 0

	for u.pos < len(u.data) {
		if u.data[u.pos] != '\\' {
			u.consumeRegularChar()
		} else if err := u.consumeEscapedSequence(); err != nil {
			return "", err
		}
	}

	return string(u.out), nil
}

// checkForClosingString returns nil when the following conditions are
// met:
// 	- <quotes> is a suffix of <str>
//  - No other unescaped occurrence of <quotes> inside <str>.
func (u *unescaper) checkForClosingString() error {
	// We need to walk byte for byte skipping unescaped characters.
	isClosed := false

	for i := 0; i+u.quotesLen <= len(u.data); i++ {
		if u.data[i] == '\\' {
			// Read past the escaped character.
			i++
			continue
		}

		isClosing := bytes.HasPrefix(u.data[i:], u.quotes)
		if isClosing && i+u.quotesLen < len(u.data) {
			return &UnescapeError{
				Msg:    "String cannot contain unescaped " + string(u.quotes),
				Offset: u.offset + i,
			}
		}

		isClosed = isClosing && i+u.quotesLen == len(u.data)
	}

	if !isClosed {
		return &UnescapeError{
			Msg:    "String must end with " + string(u.quotes),
			Offset: u.offset + len(u.data),
		}
	}

	return nil
}

func (u *unescaper) consumeRegularChar() {
	if u.data[u.pos] == '\r' {
		// Replace all types of newlines in different platforms,
		// i.e. '\r', '\n', '\r\n' are replaced with '\n'.
		u.out = append(u.out, '\n')
		u.pos++

		if u.pos < len(u.data) && u.data[u.pos] == '\n' {
			u.pos++
		}
	} else {
		// Regular character.
		u.out = append(u.out, rune(u.data[u.pos]))
		u.pos++
	}
}

func (u *unescaper) consumeEscapedSequence() error {
	// Current character is a backslash.
	if u.pos+1 >= len(u.data) {
		if u.isRaw {
			return &UnescapeError{
				Msg:    "Raw literals cannot end with odd number of \\",
				Offset: len(u.data),
			}
		}

		if u.isBytes {
			return &UnescapeError{
				Msg:    "Bytes literals cannot end \\",
				Offset: len(u.data),
			}
		}

		return &UnescapeError{
			Msg:    "String literals cannot end with \\",
			Offset: len(u.data),
		}
	}

	if u.isRaw {
		// For raw literals, all escapes are valid and those
		// characters ('\\' and the escaped character) come through
		// literally in the string.
		u.out = append(u.out, rune(u.data[u.pos]), rune(u.data[u.pos+1]))
		u.pos += 2

		return nil
	}

	// Any error that occurs in the escape is accounted to the start
	// of the escape.
	u.offset += u.pos

	// Read past the escape character.
	u.pos++

	switch u.data[u.pos] {
	case 'a':
		u.out = append(u.out, '\a')
	case 'b':
		u.out = append(u.out, '\b')
	case 'f':
		u.out = append(u.out, '\f')
	case 'n':
		u.out = append(u.out, '\n')
	case 'r':
		u.out = append(u.out, '\r')
	case 't':
		u.out = append(u.out, '\t')
	case 'v':
		u.out = append(u.out, '\v')
	case '\\':
		u.out = append(u.out, '\\')
	case '?':
		u.out = append(u.out, '?')
	case '\'':
		u.out = append(u.out, '\'')
	case '"':
		u.out = append(u.out, '"')
	case '`':
		u.out = append(u.out, '`')
	case '0', '1', '2', '3':
		return u.consumeOctal()
	case 'x', 'X':
		return u.consumeHex()
	case 'u':
		return u.consumeUnicodeShortForm()
	case 'U':
		return u.consumeUnicodeLongForm()
	case '\r', '\n':
		return &UnescapeError{
			Msg:    "Illegal escaped newline",
			Offset: u.offset,
		}
	default:
		return &UnescapeError{
			Msg:    "Illegal escape sequence: \\" + string(u.data[u.pos]),
			Offset: u.offset,
		}
	}

	u.pos++

	return nil
}

func (u *unescaper) consumeOctal() error {
	// Octal escape '\ddd': requires exactly 3 octal digits.
	// Note that the highest valid escape secquence is '\377'.
	// For string literals, octal and hex escape sequences are
	// interpreted as unicode code points, and the related
	// UTF8-escaped character is added to the destination. For
	// bytes literals, octal, and hex escape sequences are
	// interpreted as a single byte value.
	octalStart := u.pos

	if u.pos+2 >= len(u.data) {
		return &UnescapeError{
			Msg: "Illegal escape sequence: Octal escape must be " +
				"followed by 3 octal digits but saw: \\" +
				string(u.data[octalStart:]),
			Offset: u.offset,
		}
	}

	value := byte(0)
	octalEnd := u.pos + 2

	for ; u.pos <= octalEnd; u.pos++ {
		if isOctalDigit(u.data[u.pos]) {
			value = value*8 + (u.data[u.pos] - '0')
		} else {
			return &UnescapeError{
				Msg: "Illegal escape sequence: Octal escape must be " +
					"followed by 3 octal digits but saw: \\" +
					string(u.data[octalStart:octalEnd+1]),
				Offset: u.offset,
			}
		}
	}

	u.pos = octalEnd + 1
	u.out = append(u.out, rune(value))

	return nil
}

func (u *unescaper) consumeHex() error {
	// Hex escape '\xhh': requires exactly 2 hex digits.
	// For strings literals, octal, and hex escape sequences are
	// interpreted as unicode code points, and the related
	// UTF8-escaped character is added to the destination.
	// For bytes literals, octal, and hex escape sequences are
	// interpreted as a single byte value.
	hexStart := u.pos + 1
	hexEnd := hexStart + 2

	if u.pos+2 >= len(u.data) {
		return &UnescapeError{
			Msg: "Illegal escape sequence: Hex escape must be " +
				"followed by 2 hex digits but saw: \\" +
				string(u.data[u.pos:]),
			Offset: u.offset,
		}
	}

	values := make([]byte, hex.DecodedLen(2))

	n, err := hex.Decode(values, u.data[hexStart:hexEnd])
	if err != nil {
		return &UnescapeError{
			Msg: "Illegal escape sequence: Hex escape must be " +
				"followed by 2 hex digits but saw: \\" +
				string(u.data[u.pos:hexEnd]),
			Offset: u.offset,
		}
	}

	for i := 0; i < n; i++ {
		u.out = append(u.out, rune(values[i]))
	}

	u.pos = hexEnd

	return nil
}

func (u *unescaper) consumeUnicodeShortForm() error {
	if u.isBytes {
		return &UnescapeError{
			Msg: "Illegal escape sequence: Unicode escape " +
				"sequence \\" + string(u.data[u.pos]) +
				" cannot be used in bytes literals",
			Offset: u.offset,
		}
	}

	// \uhhhh => Read 4 hex digits as a code point, then write
	// it as UTF-8 bytes.
	hexStart := u.pos + 1
	hexEnd := hexStart + 4

	if u.pos+4 >= len(u.data) {
		return &UnescapeError{
			Msg: "Illegal escape sequence: \\u must be followed " +
				"by 4 hex digits but saw: \\" +
				string(u.data[u.pos:]),
			Offset: u.offset,
		}
	}

	quoted := `"\u` + string(u.data[hexStart:hexEnd]) + `"`

	value, err := strconv.Unquote(quoted)
	if err != nil {
		return &UnescapeError{
			Msg: "Illegal escape sequence: \\u must be followed " +
				"by 4 hex digits but saw: \\" +
				string(u.data[u.pos:]),
			Offset: u.offset,
		}
	}

	if []rune(value)[0] == unicode.ReplacementChar {
		return &UnescapeError{
			Msg: `Illegal escape sequence: Unicode value \` +
				string(u.data[u.pos:hexEnd]) +
				` is invalid`,
			Offset: u.offset,
		}
	}

	u.out = append(u.out, []rune(value)...)
	u.pos = hexEnd + 1

	return nil
}

func (u *unescaper) consumeUnicodeLongForm() error {
	if u.isBytes {
		return &UnescapeError{
			Msg: "Illegal escape sequence: Unicode escape " +
				"sequence \\" + string(u.data[u.pos]) +
				" cannot be used in bytes literals",
			Offset: u.offset,
		}
	}

	// \uhhhh => Read 4 hex digits as a code point, then write
	// it as UTF-8 bytes.
	hexStart := u.pos + 1
	hexEnd := hexStart + 8

	if u.pos+8 >= len(u.data) {
		return &UnescapeError{
			Msg: "Illegal escape sequence: \\U must be followed " +
				"by 8 hex digits but saw: \\" +
				string(u.data[u.pos:]),
			Offset: u.offset,
		}
	}

	quoted := `"\U` + string(u.data[hexStart:hexEnd]) + `"`

	values := make([]byte, hex.DecodedLen(8))

	n, err := hex.Decode(values, u.data[hexStart:hexEnd])
	if err != nil {
		return &UnescapeError{
			Msg: "Illegal escape sequence: \\U must be followed " +
				"by 8 hex digits but saw: \\U" +
				string(u.data[hexStart:hexEnd]),
			Offset: u.offset,
		}
	}

	cp := uint32(0)
	for _, byt := range values[:n] {
		cp = (cp << 8) + uint32(byt)
	}

	if cp > 0x0010FFFF {
		return &UnescapeError{
			Msg: `Illegal escape sequence: Value of \U` +
				string(u.data[hexStart:hexEnd]) +
				` exceeds Unicode limit (0x0010FFFF)`,
			Offset: u.offset,
		}
	}

	value, err := strconv.Unquote(quoted)
	if err != nil {
		return &UnescapeError{
			Msg: "Illegal escape sequence: \\U must be followed " +
				"by 8 hex digits but saw: \\" +
				string(u.data[u.pos:]),
			Offset: u.offset,
		}
	}

	if []rune(value)[0] == unicode.ReplacementChar {
		return &UnescapeError{
			Msg: `Illegal escape sequence: Unicode value \` +
				string(u.data[u.pos:hexEnd]) +
				` is invalid`,
			Offset: u.offset,
		}
	}

	u.out = append(u.out, []rune(value)...)
	u.pos = hexEnd + 1

	return nil
}

func isOctalDigit(c byte) bool {
	return c >= '0' && c <= '7'
}

func maybeTripleQuotedStringLiteral(s string) bool {
	if len(s) >= 6 &&
		(strings.HasPrefix(s, "'''") && strings.HasSuffix(s, "'''") ||
			strings.HasPrefix(s, `"""`) && strings.HasSuffix(s, `"""`)) {
		return true
	}

	return false
}

func maybeStringLiteral(s string) bool {
	if (len(s) >= 2) &&
		(s[0] == s[len(s)-1]) &&
		(s[0] == '\'' || s[0] == '"') {
		return true
	}

	return false
}

func maybeRawStringLiteral(s string) bool {
	if (len(s) >= 3) &&
		(s[0] == 'r' || s[0] == 'R') &&
		(s[1] == s[len(s)-1]) &&
		(s[1] == '\'' || s[1] == '"') {
		return true
	}

	return false
}

func maybeBytesLiteral(s string) bool {
	if (len(s) >= 3) &&
		(s[0] == 'b' || s[0] == 'B') &&
		(s[1] == s[len(s)-1]) &&
		(s[1] == '\'' || s[1] == '"') {
		return true
	}

	return false
}

func maybeRawBytesLiteral(s string) bool {
	if len(s) >= 4 {
		low := strings.ToLower(s[:2])
		if (low == "rb" || low == "br") &&
			(s[2] == s[len(s)-1]) &&
			(s[2] == '\'' || s[2] == '"') {
			return true
		}
	}

	return false
}

func ToStringLiteral(src string) string {
	return Escape(src, PreferSingleQuote)
}

func Escape(src string, style StringStyle) string {
	var quote string

	switch style {
	case AlwaysSingleQuote:
		quote = `'`
	case AlwaysDoubleQuote:
		quote = `"`
	case PreferSingleQuote:
		hasSingle := strings.Contains(src, `'`)
		hasDouble := strings.Contains(src, `"`)

		if !hasSingle || hasDouble {
			quote = `'`
		} else {
			quote = `"`
		}
	case PreferDoubleQuote:
		hasSingle := strings.Contains(src, `'`)
		hasDouble := strings.Contains(src, `"`)

		if !hasDouble && hasSingle {
			quote = `"`
		} else {
			quote = `'`
		}
	}

	return quote + escape(src, rune(quote[0])) + quote
}

// escape returns src rewritten using C-style escape sequences.
// If escape non-zero, only escape the quote character that matches the
// escape rune.  This allows writing "ab'cd", or 'ab"cd' without extra
// escaping.
func escape(src string, escape rune) string {
	out := make([]rune, 0, len(src))

	for _, c := range src {
		switch c {
		case '\a':
			out = append(out, '\\', '\a')
		case '\b':
			out = append(out, '\\', '\b')
		case '\f':
			out = append(out, '\\', '\f')
		case '\n':
			out = append(out, '\\', '\n')
		case '\r':
			out = append(out, '\\', '\r')
		case '\t':
			out = append(out, '\\', '\t')
		case '\v':
			out = append(out, '\\', '\v')
		case '\\':
			out = append(out, '\\', '\\')
		case '\'', '"', '`':
			if escape == 0 || c == escape {
				out = append(out, '\\')
			}

			out = append(out, c)
		default:
			out = append(out, c)
		}
	}

	return string(out)
}
