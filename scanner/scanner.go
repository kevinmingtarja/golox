// Package scanner implements a scanner for Lox source text.
package scanner

import (
	"github.com/kevinmingtarja/golox/token"
)

type errorHandler func(line int, message string)

type scanner struct {
	// immutable state
	src []byte
	err errorHandler

	// scanning state
	tokens   []token.Token
	start    int
	current  int
	line     int
	errCount int
}

func New(src []byte, err errorHandler) *scanner {
	return &scanner{
		src:     src,
		err:     err,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *scanner) error(line int, message string) {
	if s.err != nil {
		s.err(line, message)
	}
	s.errCount++
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
	c := s.advance()
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
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
		break
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL, nil)
		} else {
			s.addToken(token.EQUAL, nil)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL, nil)
		} else {
			s.addToken(token.GREATER, nil)
		}
		break
	case '/':
		if s.match('/') {
			// comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
		break
	default:
		s.error(s.line, "Unexpected character.")
		break
	}
}

// advance consumes the next character in the source
// and returns it.
func (s *scanner) advance() byte {
	s.current++
	return s.src[s.current-1]
}

// match consumes the next character in the source
// if it matches the expected character.
func (s *scanner) match(expected byte) bool {
	if s.current >= len(s.src) {
		return false
	}
	if s.src[s.current] != expected {
		return false
	}
	s.current++
	return true
}

// peek returns the byte at the current position without
// consuming it. If current is at the end of the source,
// it returns 0.
func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.src[s.current]
}

func (s *scanner) addToken(t token.TokenType, literal interface{}) {
	text := string(s.src[s.start:s.current])
	s.tokens = append(s.tokens, token.Token{
		Type: t, Lexeme: text, Literal: literal, Line: s.line,
	})
}
