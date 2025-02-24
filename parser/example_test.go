package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kevinmingtarja/golox/ast"
)

// astPrinter is a visitor that simply prints the AST
type astPrinter struct {
	b *strings.Builder
}

// Visit implements the Visitor interface
func (ap astPrinter) Visit(expr ast.Expr) ast.Visitor {
	// depending on the type of the expression,
	// we will print it differently
	switch e := expr.(type) {
	case *ast.LiteralExpr:
		if e.Value == nil {
			ap.b.Write([]byte("nil"))
		}
		ap.b.Write(fmt.Appendf(nil, "%v", e.Value))
	case *ast.BinaryExpr:
		ap.parenthesize(e.Op.Lexeme, e.X, e.Y)
	case *ast.GroupingExpr:
		ap.parenthesize("group", e.Expr)
	case *ast.UnaryExpr:
		ap.parenthesize(e.Op.Lexeme, e.X)
	case *ast.ConditionalOperator:
		ap.parenthesize("?:", e.Cond, e.X, e.Y)
	default:
		return ap
	}

	return nil
}

func (ap *astPrinter) parenthesize(name string, exprs ...ast.Expr) {
	ap.b.Write([]byte("(" + name))
	for _, expr := range exprs {
		ap.b.Write([]byte(" "))
		// recursively walk the child ASTs
		ast.Walk(ap, expr)
	}
	ap.b.Write([]byte(")"))
}

var _ ast.Visitor = astPrinter{}

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		src  []byte
		want ast.Expr
	}{
		{
			src: []byte("1 + 2"),
		},
		{
			src: []byte("!true"),
		},
		{
			src: []byte("-123 + (45.67)"),
		},
		{
			src: []byte("(123 + 1, 345)"),
		},
		{
			src: []byte("3 > 2 ? 1 : 0"),
		},
		{
			src: []byte("(1 + 2"),
		},
	} {
		got := Parse(tc.src)
		fmt.Println(got)
		ap := astPrinter{
			b: &strings.Builder{},
		}
		// walk the ast from the root
		ast.Walk(ap, got)
		println("RESULT", ap.b.String())
		println()
	}
}
