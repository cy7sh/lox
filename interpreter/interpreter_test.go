package interpreter

import (
	"strings"
	"testing"

	"github.com/singurty/lox/environment"
	"github.com/singurty/lox/parser"
	"github.com/singurty/lox/scanner"
)

func runTest(source string) {
	// reset environment
	env = environment.Global()
	scan := scanner.New(source)
	tokens := scan.ScanTokens()
	parse := parser.New(tokens)
	statements := parse.Parse()
	Interpret(statements)
}

func TestVariable(t *testing.T) {
	input := `
		var a = 2;
		var b = 3;
		a = b = a * b;
		var c = "hello";
		var d = c + " world";
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
	expected := `
inner a
outer b
global c
outer a
outer b
global c
global a
global b
global c`
	testInterpreterOutput(input, expected, t)
}

func TestConditionalExecution(t *testing.T) {
	input := `
	var a = 2;
	var b = 3;
	var c = 2;
	if (a == c) {
		a = a * b;
		print a;
	}
	if (a == 2) {
		a = 12;
		print a;
	} else {
		a = 10;
		print a;
	}
	if (a == 6) {
		a = 5;
		print a;
	} else if (a == 10) {
		b = 24;
		print b;
	}
	if (a == 6) {
		a = 39;
		print a;
	} else if (a == 30) {
		b = 24;
		print b;
	} else {
		c = 25;
		print c;
	}
	print a;
	print b;
	print c;
	`
	expected := `
6
10
24
25
10
24
25
`
	testInterpreterOutput(input, expected, t)
}

func testInterpreterOutput(input string, expected string, t *testing.T) {
	sb :=  &strings.Builder{}
	InterpreterOptions.PrintOutput = sb
	runTest(input)
	output := strings.Trim(sb.String(), "\n")
	expected = strings.Trim(expected, "\n")
	if output != expected {
		t.Errorf("Expected output to be : %v\nGot: %v\n", expected, output)
	}
}
