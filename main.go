package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64) // EX_USAGE
	} else if len(args) == 2 {
		if err := runFile(args[1]); err != nil {
			fmt.Println(err)
			os.Exit(65)
		}
	} else {
		// REPL
		if err := runPrompt(); err != nil {
			fmt.Println(err)
			os.Exit(65)
		}
	}
}
