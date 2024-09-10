package ast

import "github.com/kevinmingtarja/golox/token"

type Expr interface {
	exprNode()
}

// Expressions
type (
	BinaryExpr struct {
		X  Expr
		Op token.Token
		Y  Expr
	}

	GroupingExpr struct {
		Expr Expr
	}

	LiteralExpr struct {
		Value interface{}
	}

	UnaryExpr struct {
		Op token.Token
		X  Expr
	}
)

// exprNode() ensures that only expression nodes can be
// assigned to an Expr.
// inspired by: https://github.com/golang/go/blob/807e01db4840e25e4d98911b28a8fa54244b8dfa/src/go/ast/ast.go#L548
func (*BinaryExpr) exprNode()   {}
func (*GroupingExpr) exprNode() {}
func (*LiteralExpr) exprNode()  {}
func (*UnaryExpr) exprNode()    {}
