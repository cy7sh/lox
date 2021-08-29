package main

import (
	"os"
	"bufio"
	"fmt"
	"io"
	"github.com/singurty/lox/scanner"
)

var hadError bool

func main() {
	if (len(os.Args) > 2) {
		fmt.Printf("Usage: %v [file]\n", os.Args[0])
		os.Exit(1)
	} else if (len(os.Args) == 2) {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">> ")
		text, err := reader.ReadString('\n')
		if (err == io.EOF) {
			fmt.Println("\nExiting..")
			return
		} else if (err != nil) {
			panic(err)
		}
		run(text)
		hadError = false
	}
}

func runFile(file string) {
	content, err := os.ReadFile(file)
	if (err != nil) {
		panic(err)
	}
	run(string(content))
}

func Error(line int, message string) {
	hadError = true
	fmt.Fprintf(os.Stderr, "[line %v] Error : %v", line, message)
}

func run(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens{
		fmt.Print(token.String())
	}
}
