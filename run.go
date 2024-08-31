package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for fmt.Print("> "); scanner.Scan(); fmt.Print("> ") {
		line := scanner.Text()
		if line == "exit" {
			break
		}
		if err := run(string(line)); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func run(source string) error {
	tokens := strings.Split(source, " ")
	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}
