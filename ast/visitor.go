package ast

// A Visitor is a way to inject behavior into the AST traversal done by [Walk].
// The Visit method is called for each node in the AST.
type Visitor interface {
	Visit(expr Expr) Visitor
}

// Walk traverses an AST in depth-first order.
func Walk(v Visitor, expr Expr) {
	if v = v.Visit(expr); v == nil {
		return
	}

	// walk children
	switch e := expr.(type) {
	case *LiteralExpr:
		// nothing to do
	case *BinaryExpr:
		Walk(v, e.X)
		Walk(v, e.Y)
	case *GroupingExpr:
		Walk(v, e.Expr)
	case *UnaryExpr:
		Walk(v, e.X)
	default:
		panic("unreachable")
	}

	v.Visit(nil)
}
