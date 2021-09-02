package interpreter

import (
	"testing"

	"github.com/singurty/lox/scanner"
	"github.com/singurty/lox/parser"
)

func runTest(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()
	Interpret(statements)
}

func TestVariable(t *testing.T) {
	input := `
		var a = 2;
		var b = 3;
		print a;
		print b;
		a = b = a * b;
		print a;
		print b;
		var c = "hello";
		print c;
		var d = c + " world";
		print d;
	`
	runTest(input)
	// Output:
	// 2
	// 3
	// 6
	// 6
	// hello
	// hello world
	a, err := env.Get("a")
	if err != nil {
		t.Fatalf("Expected variable 'a' in env")
	}
	if a.(float64) != 6.0 {
		t.Errorf("Expected variable 'a' to be 6.0 got %v instead", a.(float64))
	}
	b, err := env.Get("b")
	if err != nil {
		t.Fatalf("Expected variable 'b' in env")
	}
	if a.(float64) != 6.0 {
		t.Errorf("Expected variable 'b' to be 6.0 got %v instead", b.(float64))
	}
	c, err := env.Get("c")
	if err != nil {
		t.Fatalf("Expected variable 'c' in env")
	}
	if c != "hello" {
		t.Errorf("Expected variable 'c' to be \"hello\" got \"%v\" instead", c)
	}
	d, err := env.Get("d")
	if err != nil {
		t.Fatalf("Expected variable 'd' in env")
	}
	if d != "hello world" {
		t.Errorf("Expected variable 'd' to be \"hello world\" got \"%v\" instead", d)
	}
}

func TestVariableScope(t *testing.T) {
	input := `
		var a = "global a";
		var b = "global b";
		var c = "global c";
		{
		  var a = "outer a";
		  var b = "outer b";
		  {
		    var a = "inner a";
		    print a;
		    print b;
		    print c;
		  }
		  print a;
		  print b;
		  print c;
		}
		print a;
		print b;
		print c;
	`
	runTest(input)
	// Output:
	// inner a
	// outer b
	// global c
	// outer a
	// outer b
	// global c
	// gloabl a
	// global b
	// global c
}
