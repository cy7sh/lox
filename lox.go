package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/singurty/lox/interpreter"
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
		err = run(text)
		if err != nil {
			fmt.Printf(err.Error())
		}
		hadError = false
	}
}

func runFile(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = run(string(content))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func run(source string) error {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	if scanner.HadError {
		return nil
	}
	//for _, token := range tokens{
	//	fmt.Println(token.String())
	//}
	parser := parser.New(tokens)
	expression := parser.Parse()
	if parser.HadError {
		return nil
	}
	fmt.Println(expression.String())
	interpreted, err := interpreter.Eval(expression)
	if err != nil {
		return err
	}
	fmt.Println(interpreted)
	return nil
}
