package main

import (
	"os"
	"bufio"
	"fmt"
	"io"
	"github.com/singurty/lox/lox"
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
		lox.Run(text)
		hadError = false
	}
}

func runFile(file string) {
	content, err := os.ReadFile(file)
	if (err != nil) {
		panic(err)
	}
	lox.Run(string(content))
}

func reportError(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %v] Error%v : %v", line, where, message)
}
