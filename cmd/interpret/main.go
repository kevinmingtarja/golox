package main

import (
	"fmt"

	"github.com/kevinmingtarja/golox/interpreter"
	"github.com/kevinmingtarja/golox/parser"
)

func main() {
	// src := []byte("(1 + 2) * -4")
	src := []byte(`"hello" + " " + "world!"`)
	expr := parser.Parse(src)
	if expr == nil {
		return
	}
	fmt.Printf("%v\n", interpreter.Evaluate(expr))
}
