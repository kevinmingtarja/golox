package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/kevinmingtarja/golox/scanner"
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
	return run(bytes)
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for fmt.Print("> "); scanner.Scan(); fmt.Print("> ") {
		line := scanner.Bytes()
		if bytes.Equal(line, []byte("exit")) {
			break
		}
		if err := run(line); err != nil {
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

func run(src []byte) error {
	scanner := scanner.New(src)
	scanner.ScanTokens()
	return nil
}
