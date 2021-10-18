package interpreter

import (
	"strings"
	"testing"

	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/environment"
	"github.com/singurty/lox/parser"
	"github.com/singurty/lox/resolver"
	"github.com/singurty/lox/scanner"
	//	"github.com/augustoroman/hexdump" // to debug minor differences in text comparison
)

type testInputs []struct {
	input string
	expected string
}

func runTest(source string, t *testing.T) {
	// reset environment
	env = environment.Global()
	global = env
	locals = make(map[ast.Expr]int)
	scan := scanner.New(source)
	tokens := scan.ScanTokens()
	if scan.HadError {
		t.Fatal("scanner error")
	}
	parse := parser.New(tokens)
	if parse.HadError {
		t.Fatal("parser error")
	}
	statements := parse.Parse()
	resolver := resolver.NewResolver()
	err := resolver.Resolve(statements)
	if err != nil {
		t.Fatal(err)
	}
	err = Interpret(statements, resolver)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVariable(t *testing.T) {
	input := `
		var a = 2;
		var b = 3;
		a = b = a * b;
		var c = "hello";
		var d = c + " world";
	`
	runTest(input, t)
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

func TestLogicalOperators(t *testing.T) {
	input := `
	print "hi" or 2;
	print null or "yes";
	print true or 2;
	print false or "yes";
	print false or false;
	print null or null;
	print "hi" and 2;
	print null and "yes";
	print true and 2;
	print false and "yes";
	print false and false;
	print null and null;
	`
	expected := `
hi
yes
true
yes
false
null
2
null
2
false
false
null
`
	testInterpreterOutput(input, expected, t)
}

func TestWhileLoop(t *testing.T) {
	input := `
	var a = 0;
	while (a < 10) {
		print "loo";
		a = a + 1;
	}
	`
	expected := `
loo
loo
loo
loo
loo
loo
loo
loo
loo
loo
`
	testInterpreterOutput(input, expected, t)
}

func TestForLoop(t *testing.T) {
	input := `
		var a = 0;
		var temp;
		
		for (var b = 1; a < 10000; b = temp + b) {
		  print a;
		  temp = a;
		  a = b;
		}
	`
	expected := `
0
1
1
2
3
5
8
13
21
34
55
89
144
233
377
610
987
1597
2584
4181
6765
`
	testInterpreterOutput(input, expected, t)
}

func TestBreak(t *testing.T) {
	input := `
		for (var a = 0; ; a = a + 1) {
			var b = 0;
			while (b < 5) {
				if (a > 2) {
					break;
				}
				print "while";
				b = b + 1;
			}
			if (a > 5) {
				break;
				print "after break";
			}
			print a;
		}
	`
	expected := `
while
while
while
while
while
0
while
while
while
while
while
1
while
while
while
while
while
2
3
4
5
`
	testInterpreterOutput(input, expected, t)
}

func TestContinue(t *testing.T) {
	tests := testInputs{
		{`
			var a = 1;
			while (a < 10) {
				a = a + 1;
				if (a < 9) {
					continue;
				}
				print a;
				break;
			}
			`, "9"},
		{`
			for (var a = 1; a < 10; a = a + 1) {
				if (a < 9) {
					continue;
				}
				print a;
			}
			`, "9"},
	}
	testInterpreterOutputs(tests, t)
}

func TestFunction(t *testing.T) {
	tests := testInputs{
		{
`
fun sayHi(first, last) {
	print "Hi, " + first + " " + last + "!";
}
sayHi("Dear", "Reader");
`,
`
Hi, Dear Reader!
`,
		},
		{
`
fun count(n) {
	if (n > 1) count(n - 1);
	print n;
}
count(3);
`,
`
1
2
3
`,
		},
		{
`
fun add(a, b, c) {
	print a + b + c;
}
add(1, 2, 3);

`,
`
6
`,
		},
		{
// recursion
`
fun fib(n) {
	if (n <= 1) return n;
	return fib(n - 2) + fib(n - 1);
}

for (var i = 0; i < 20; i = i + 1) {
	print fib(i);
}
`,
`
0
1
1
2
3
5
8
13
21
34
55
89
144
233
377
610
987
1597
2584
4181
`,
		},
		{
// closure
`
fun makeCounter() {
  var i = 0;
  fun count() {
    i = i + 1;
    print i;
  }

  return count;
}

var counter = makeCounter();
counter();
counter();
counter();
counter();
counter();
`,
`
1
2
3
4
5
`,
		},
		{
// lambda
`
fun thrice(fn) {
  for (var i = 1; i <= 3; i = i + 1) {
    fn(i);
  }
}

thrice(fun (a) {
  print a;
});
`,
`
1
2
3
`,
		},
		{
// this is what our resolver should solve
`
var a = "global";
{
  fun showA() {
    print a;
  }

  showA();
  var a = "block";
  showA();
}
`,
`
global
global
`,
		},
	}
	testInterpreterOutputs(tests, t)
}

func TestClass(t *testing.T) {
	tests := testInputs{
		{
// test instance properties
`
class Bagel{}
var bagel = Bagel();
bagel.prop = "property";
print bagel.prop;
`,
`
property
`,
		},
		{
// test methods
`
class Bacon {
  eat() {
    print "Crunch crunch crunch!";
  }
}

Bacon().eat();
`,
`
Crunch crunch crunch!
`,
		},
		{
// test this
`
class Cake {
  taste() {
    var adjective = "delicious";
    print "The " + this.flavor + " cake is " + adjective + "!";
  }
}

var cake = Cake();
cake.flavor = "German chocolate";
cake.taste();
`,
`
The German chocolate cake is delicious!
`,
		},
		{
// test initializer
`
class Foo {
        init(first, second) {
                this.first = first;
                this.second = second;
        }
        display() {
                print this.first;
                print this.second;
        }
}

var bar = Foo("hello", "world");
bar.display();
`,
`
hello
world
`,
	},
	{
`
class Foo {
	init(first, second) {
		this.first = first;
		this.second = second;
	}
	display() {
		print this.first;
		print this.second;
	}
}

var bar = Foo("hello", "world");
bar.display();
bar.first = "bye";
bar.second = "earth";
bar.display();
bar.init("hello", "world").display();
bar.display();
`,
`
hello
world
bye
earth
hello
world
hello
world
`,
	},
	{
`
class Foo {
	init() {
		this.yo = "before";
		return;
		this.yo = "after";
	}
	display() {
		print this.yo;
	}
}

Foo().display();
`,
`
before
`,
	},
		}
	testInterpreterOutputs(tests, t)
}

func testInterpreterOutputs(tests testInputs, t *testing.T) {
	for _, test := range tests {
		testInterpreterOutput(test.input, test.expected, t)
	}
}

func testInterpreterOutput(input string, expected string, t *testing.T) {
	sb :=  &strings.Builder{}
	InterpreterOptions.PrintOutput = sb
	runTest(input, t)
	output := strings.Trim(sb.String(), "\n")
	expected = strings.Trim(expected, "\n")
	if output != expected {
		t.Errorf("Expected output to be : %v\nGot: %v\n",expected, output)
	}
}
