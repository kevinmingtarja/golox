package main

import (
	"bufio"
	"bytes"
	"errors"
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
	sc := bufio.NewScanner(os.Stdin)
	for fmt.Print("> "); sc.Scan(); fmt.Print("> ") {
		line := sc.Bytes()
		if bytes.Equal(line, []byte("exit")) {
			break
		}
		err := run(line)
		if err != nil {
			if !errors.As(err, &scanner.Error{}) {
				return err
			}
			// don't kill the entire session if the user makes a mistake
			fmt.Println(err.Error())
		}
	}
	if err := sc.Err(); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func run(src []byte) error {
	scanner := scanner.New(src, nil)
	scanner.ScanTokens()
	return nil
}
