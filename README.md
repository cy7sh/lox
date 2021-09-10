The best language in the universe (yep, aliens too)

## Documentation
#### Variable
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
### Boolean
There are two boolean primitives `true` and `false`. `null` is falsey; anything else is truthy.
### Block
Block is a statement containing other statements. Statements inside a block have their own environment with variables. Statements inside the block can access and modify variables declared outside the block. Variables declared inside the block are only accessible inside the block.
```
{
    statements
}
```
### While loop
```
while (condition)
    statement
```

### For loop
```
for (initializer; condition; increment)
    statement
```
`initilizer` can be variable declaration or an expression. If a variable is declared, it's scope is limited to the loop. It is evaluated before the loop starts. `condition` must be an expression. It is evaluated *before* each iteration. Loop terminates if the result is falsey. `increment` must be an expression. It is evaluated *after* each iteration.
### Continue and break
```
continue;
```
Stop current iteration and continue with another iteration of this loop. `condition` and `increment` are evaluated.
```
break;
```
Break out of current loop and continue executing statements after the loop.
### if else
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
### Ternary operator
```
<condition> ? <if expression> : <else expression>
```
