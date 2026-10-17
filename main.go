package main

import (
	"fmt"
	"os"
)

type lox struct {
	hadError bool
}

func main() {
	l := lox{}
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: golox [script]")
	} else if len(args) == 1 {
		l.runFile(args[0])
	} else {
		l.runPrompt()
	}
}
