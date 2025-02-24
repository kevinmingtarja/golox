package parser

import (
	"fmt"

	"github.com/kevinmingtarja/golox/ast"
	"github.com/kevinmingtarja/golox/scanner"
	"github.com/kevinmingtarja/golox/token"
)

type parser struct {
	tokens  []token.Token
	current int
}

func (p *parser) init(src []byte) {
	sc := scanner.New(src, nil /* TODO */)
	p.tokens = sc.ScanTokens()
}

func Parse(src []byte) ast.Expr {
	var p parser
	defer func() {
		if e := recover(); e != nil {
			bail, ok := e.(bailout)
			if !ok {
				panic(e)
			} else if bail.msg != "" {
				fmt.Println("parse error:", bail.msg)
			}
		}
	}()

	p.init(src)
	return p.consumeExpression()
}

func (p *parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(typ token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

func (p *parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *parser) consumeExpression() ast.Expr {
	return p.consumeExpressionList()
}

func (p *parser) consumeExpressionList() ast.Expr {
	expr := p.consumeTernary()
	for p.match(token.COMMA) {
		op := p.previous()
		if op.Type != token.COMMA {
			panic("wrong type")
		}
		right := p.consumeTernary()
		expr = &ast.BinaryExpr{X: expr, Op: op, Y: right}
	}
	return expr
}

func (p *parser) consumeTernary() ast.Expr {
	expr := p.consumeEquality()
	if p.match(token.QUESTION) {
		op := p.previous()
		if op.Type != token.QUESTION {
			panic("wrong type")
		}

		left := p.consumeExpression()
		p.consume(token.COLON, "Expect ':' after expression.")
		right := p.consumeTernary()

		expr = &ast.ConditionalOperator{Cond: expr, X: left, Y: right}
	}
	return expr
}

func (p *parser) consumeEquality() ast.Expr {
	expr := p.consumeComparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		op := p.previous() // grab the matched op token
		right := p.consumeComparison()
		expr = &ast.BinaryExpr{X: expr, Op: op, Y: right}
	}
	return expr
}

func (p *parser) consumeComparison() ast.Expr {
	expr := p.consumeTerm()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		op := p.previous()
		right := p.consumeTerm()
		expr = &ast.BinaryExpr{X: expr, Op: op, Y: right}
	}
	return expr
}

func (p *parser) consumeTerm() ast.Expr {
	expr := p.consumeFactor()
	for p.match(token.MINUS, token.PLUS) {
		op := p.previous()
		right := p.consumeFactor()
		expr = &ast.BinaryExpr{X: expr, Op: op, Y: right}
	}
	return expr
}

func (p *parser) consumeFactor() ast.Expr {
	expr := p.consumeUnary()
	for p.match(token.SLASH, token.STAR) {
		op := p.previous()
		right := p.consumeUnary()
		expr = &ast.BinaryExpr{X: expr, Op: op, Y: right}
	}
	return expr
}

func (p *parser) consumeUnary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		op := p.previous()
		right := p.consumeUnary()
		return &ast.UnaryExpr{Op: op, X: right}
	}
	return p.consumePrimary()
}

func (p *parser) consumePrimary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.LiteralExpr{Value: false}
	}
	if p.match(token.TRUE) {
		return &ast.LiteralExpr{Value: true}
	}
	if p.match(token.NIL) {
		return &ast.LiteralExpr{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.LiteralExpr{Value: p.previous().Literal}
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.consumeExpression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.GroupingExpr{Expr: expr}
	}
	panic(bailout{token: p.peek(), msg: "Expect expression."})
}

type bailout struct {
	token token.Token
	msg   string
}

func (p *parser) consume(typ token.TokenType, msg string) token.Token {
	if p.check(typ) {
		return p.advance()
	}
	panic(bailout{token: p.peek(), msg: msg})
}

func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
