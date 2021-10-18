The best language in the universe (yep, aliens too)

## Usage
```
$ lox [filename]
```
Starts an interactive shell if `filename` is omitted.

## Documentation
#### Variables
- Delcaration
```
var a;
```
- Assignment
```
a = 10;
```
- Delcaration and assignment
```
var a = 10;
```
### Booleans
There are two boolean primitives `true` and `false`. `null` is falsey; anything else is truthy.
### Blocks
Block is a statement containing other statements. Statements inside a block have their own environment with variables. Statements inside the block can access and modify variables declared outside the block. Variables declared inside the block are only accessible inside the block.
```
{
    statements
}
```
### While loops
```
while (condition)
    statement
```

### For loops
```
for (initializer; condition; increment)
    statement
```
`initilizer` can be variable declaration or an expression. If a variable is declared, it's scope is limited to the loop. It is evaluated before the loop starts. `condition` must be an expression. It is evaluated *before* each iteration. Loop terminates if the result is falsey. `increment` must be an expression. It is evaluated *after* each iteration.
### continue and break statements
```
continue;
```
Stop current iteration and continue with another iteration of this loop. `condition` and `increment` are evaluated.
```
break;
```
Break out of current loop and continue executing statements after the loop.
### if else statements
```
if (condition)
    statement
else
    satement
```
`condition` must be an expression. `else` is optional.
### Logical operators
```
print "hi" or 2;
print true or 2;
print false or "yes";
print "hi" and 2;
print true and 2;
print false and "yes";
```
Output:
```
hi
true
yes
2
2
false
```
### Ternary operators
```
<condition> ? <if expression> : <else expression>
```

### Functions
```
fun fib(n) {
	if (n <= 1) return n;
	return fib(n - 2) + fib(n - 1);
}

for (var i = 0; i < 20; i = i + 1) {
	print fib(i);
}
```
Output:
```
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
```
#### Closures
```
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
```
Output:
```
1
2
3
4
5
```
### Lambdas (Anonymous functions)
```
fun thrice(fn) {
  for (var i = 1; i <= 3; i = i + 1) {
    fn(i);
  }
}

thrice(fun (a) {
  print a;
});
```
Output:
```
1
2
3
```
### Classes
#### Properties
```
class Bagel{}
var bagel = Bagel();
bagel.prop = "property";
print bagel.prop;
```
Output:
```
property
```
#### Methods
```
class Bacon {
  eat() {
    print "Crunch crunch crunch!";
  }
}

Bacon().eat();
```
Output:
```
Crunch crunch crunch!
```
#### this
```
class Cake {
  taste() {
    var adjective = "delicious";
    print "The " + this.flavor + " cake is " + adjective + "!";
  }
}

var cake = Cake();
cake.flavor = "German chocolate";
cake.taste();
```
Output:
```
The German chocolate cake is delicious!
```
#### Initializers
```
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
```
Output:
```
hello
world
```
```
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
```
Output:
```
hello
world
bye
earth
hello
world
hello
world
```
