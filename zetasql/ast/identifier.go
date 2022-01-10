package ast

func ToIdentifierLiteral(id string) string {
	if len(id) > 0 && id[0] == '`' {
		return id
	}

	if IsValidUnquotedIdentifier(id) {
		return id
	}

	// TODO: escape quoted identifier
	return "`" + id + "`"
}

func IsValidUnquotedIdentifier(id string) bool {
	if len(id) == 0 {
		return false
	}

	if !isLetter(id[0]) && !(id[0] == '_') {
		return false
	}

	for _, c := range id {
		if !isAlphanum(byte(c)) {
			return false
		}
	}

	return !IsReservedKeyword(id)
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isAlphanum(c byte) bool {
	return isLetter(c) || ('0' <= c && c <= '9') || (c == '_')
}
