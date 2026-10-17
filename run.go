package main

import (
	"bufio"
	"fmt"
	"os"
)

func (l *lox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		l.run(text)
		l.hadError = false
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error occured: ", err)
	}
}

func (l *lox) runFile(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error occured: ", err)
		return
	}
	bytes, err := os.ReadFile(path)
	l.run(string(bytes))

	if l.hadError {
		os.Exit(65)
	}
}

func (l *lox) run(code string) {
	fmt.Println(code)
}

func (l *lox) error(line string, message string) {
	l.report(line, "", message)
}

func (l *lox) report(line string, where string, message string) {
	fmt.Printf("[line %s] Error %s : %s", line, where, message)
	l.hadError = true
}
