package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/singurty/lox/parser"
	"github.com/singurty/lox/scanner"
)

var hadError bool

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %v [file]\n", os.Args[0])
		os.Exit(1)
	} else if len(os.Args) == 2 {
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
		if err == io.EOF {
			fmt.Println("\nExiting..")
			return
		} else if err != nil {
			panic(err)
		}
		run(text)
		hadError = false
	}
}

func runFile(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	run(string(content))
}

func run(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	//for _, token := range tokens{
	//	fmt.Println(token.String())
	//}
	parser := parser.New(tokens)
	expression := parser.Parse()
	fmt.Println(expression.String())
}
