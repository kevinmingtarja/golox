package interpreter

import (
	"github.com/kevinmingtarja/golox/ast"
	"github.com/kevinmingtarja/golox/token"
)

func Evaluate(expr ast.Expr) any {
	switch e := expr.(type) {
	case *ast.LiteralExpr:
		return e.Value
	case *ast.GroupingExpr:
		return Evaluate(e.Expr)
	case *ast.UnaryExpr:
		x := Evaluate(e.X)

		switch e.Op.Type {
		case token.MINUS:
			return -x.(float64)
		case token.BANG:
			return !isTruthy(x)
		}

		// unreachable
		return nil
	case *ast.BinaryExpr:
		x := Evaluate(e.X)
		y := Evaluate(e.Y)

		switch e.Op.Type {
		case token.MINUS:
			return x.(float64) - y.(float64)
		case token.SLASH:
			return x.(float64) / y.(float64)
		case token.STAR:
			return x.(float64) * y.(float64)
		case token.PLUS:
			if _, ok := x.(float64); ok {
				if _, ok := y.(float64); ok {
					return x.(float64) + y.(float64)
				}
				return nil
			}
			if _, ok := x.(string); ok {
				if _, ok := y.(string); ok {
					return x.(string) + y.(string)
				}
				return nil
			}
		case token.GREATER:
			return x.(float64) > y.(float64)
		case token.GREATER_EQUAL:
			return x.(float64) >= y.(float64)
		case token.LESS:
			return x.(float64) < y.(float64)
		case token.LESS_EQUAL:
			return x.(float64) <= y.(float64)
		case token.BANG_EQUAL:
			return !isEqual(x, y)
		case token.EQUAL_EQUAL:
			return isEqual(x, y)
		}

		// unreachable
		return nil
	}
	return nil
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

func isEqual(x, y any) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil {
		return false
	}
	return x == y
}
