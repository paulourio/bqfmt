package lexer

import "github.com/paulourio/bqfmt/zetasql/token"

type LexerFilter struct {
	lex *Lexer
	id  token.Type
}

// WithoutComment returns a Lexer that ignores comments.
func (l *Lexer) WithoutComment() *LexerFilter {
	id := token.TokMap.Type("comment")
	return &LexerFilter{l, id}
}

func (l *LexerFilter) Scan() (tok *token.Token) {
	for {
		tok = l.lex.Scan()
		if tok.Type != l.id {
			break
		}
	}
	return
}
