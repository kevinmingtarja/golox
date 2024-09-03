// Package scanner implements a scanner for Lox source text.
package scanner

import (
	"fmt"
	"os"
	"strconv"

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
	for _, t := range s.tokens {
		fmt.Println(t)
	}
	fmt.Fprintln(os.Stderr, s.errCount)
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

	case ' ':
	case '\r':
	case '\t':
		// ignore whitespace.
		break

	case '\n':
		s.line++
		break

	case '"':
		s.string()
		break

	default:
		if isDigit(c) {
			s.number()
		} else if isAlhpa(c) {
			s.identifier()
		} else {
			s.error(s.line, "Unexpected character.")
		}
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

func (s *scanner) peekNext() byte {
	if s.current+1 >= len(s.src) {
		return 0
	}
	return s.src[s.current+1]
}

func (s *scanner) addToken(t token.TokenType, literal interface{}) {
	text := string(s.src[s.start:s.current])
	s.tokens = append(s.tokens, token.Token{
		Type: t, Lexeme: text, Literal: literal, Line: s.line,
	})
}

func (s *scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error(s.line, "Unterminated string.")
		return
	}

	// consume the closing quote
	s.advance()

	// trim the surrounding quotes
	value := string(s.src[s.start+1 : s.current-1])
	s.addToken(token.STRING, value)
}

func (s *scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// we need two characters lookahead to check if
	// there's a digit after the '.'
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	f, err := strconv.ParseFloat(string(s.src[s.start:s.current]), 64)
	if err != nil {
		s.error(s.line, "Invalid number.")
	}
	s.addToken(token.NUMBER, f)
}

func (s *scanner) identifier() {
	for isAlhpaNumeric(s.peek()) {
		s.advance()
	}

	ident := string(s.src[s.start:s.current])
	var tok token.TokenType
	// all our keywords are at least two characters long.
	if len(ident) > 1 {
		tok = token.Lookup(ident)
	} else {
		tok = token.IDENTIFIER
	}

	s.addToken(tok, nil)
}

func isAlhpa(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlhpaNumeric(c byte) bool {
	return isAlhpa(c) || isDigit(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
