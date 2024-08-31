// Package scanner implements a scanner for Lox source text.
package scanner

import "github.com/kevinmingtarja/golox/token"

type scanner struct {
	// immutable state
	src []byte

	// scanning state
	tokens  []token.Token
	start   int
	current int
	line    int
}

func New(src []byte) *scanner {
	return &scanner{
		src:     src,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *scanner) ScanTokens() {
	for !s.isAtEnd() {
		// the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.Token{
		Type: token.EOF, Lexeme: "", Literal: nil, Line: s.line,
	})
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *scanner) scanToken() {
	c := advance(s.src, &s.current)
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
		break
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
		break
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
		break
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
		break
	case ',':
		s.addToken(token.COMMA, nil)
		break
	case '.':
		s.addToken(token.DOT, nil)
		break
	case '-':
		s.addToken(token.MINUS, nil)
		break
	case '+':
		s.addToken(token.PLUS, nil)
		break
	case ';':
		s.addToken(token.SEMICOLON, nil)
		break
	case '*':
		s.addToken(token.STAR, nil)
		break
	default:
		// TODO
	}
}

func advance(src []byte, current *int) byte {
	*current++
	return src[*current-1]
}

func (s *scanner) addToken(t token.TokenType, literal interface{}) {
	text := string(s.src[s.start:s.current])
	s.tokens = append(s.tokens, token.Token{
		Type: t, Lexeme: text, Literal: literal, Line: s.line,
	})
}
