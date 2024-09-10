package main

import (
	"strings"

	"github.com/kevinmingtarja/golox/ast"
	"github.com/kevinmingtarja/golox/token"
)

type astPrinter struct {
	b *strings.Builder
}

func (ap astPrinter) Visit(expr ast.Expr) ast.Visitor {
	switch e := expr.(type) {
	case *ast.LiteralExpr:
		if e.Value == nil {
			ap.b.Write([]byte("nil"))
		}
		ap.b.Write([]byte(e.Value.(string)))
	case *ast.BinaryExpr:
		ap.parenthesize(e.Op.Lexeme, e.X, e.Y)
	case *ast.GroupingExpr:
		ap.parenthesize("group", e.Expr)
	case *ast.UnaryExpr:
		ap.parenthesize(e.Op.Lexeme, e.X)
	default:
		return ap
	}

	return nil
}

func (ap *astPrinter) parenthesize(name string, exprs ...ast.Expr) {
	ap.b.Write([]byte("(" + name))
	for _, expr := range exprs {
		ap.b.Write([]byte(" "))
		ast.Walk(ap, expr)
	}
	ap.b.Write([]byte(")"))
	return
}

var _ ast.Visitor = astPrinter{}

func main() {
	expr := &ast.BinaryExpr{
		X: &ast.UnaryExpr{
			Op: token.Token{Type: token.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			X:  &ast.LiteralExpr{Value: "123"},
		},
		Op: token.Token{Type: token.PLUS, Lexeme: "+", Literal: nil, Line: 1},
		Y:  &ast.LiteralExpr{Value: "45.67"},
	}
	ap := astPrinter{
		b: &strings.Builder{},
	}
	ast.Walk(ap, expr)
	println(ap.b.String())
}
