package parser

import (
	"fmt"

//	"github.com/davecgh/go-spew/spew" // to dump structs for debugging
	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/token"
)

/*
program        → block* EOF
declaration    → funDecl | varDecl | statement
classDecl      → "class" IDENTIFIER "(" function* "}"
funDecl        → "fun" function
function       → IDENTIFIER "(" parameters? ")" block
parameters     → IDENTIFIER ("," IDENTIFIER )*
varDecl        → "var" IDENTIFIER ("=" expression)? ";"
statement      → exprStmt | printStmt | block | forStmt | break | returnStmt
break          → "break" ";"
forStmt        → "for" "(" (varDecl | exprStmt | ";") expression? ";" expression? ")" statement
whileStmt      → "while" "(" expression ")" statement
ifStmt         → "if " "(" expression ")" statement ("else" statement)?
block          → "{" declaration* "}"
exprStmt       → expression ";"
printStmt      → "print" expression ";"
returnStmt     → "return" expression? ";"
expression     → assignment
assignment     → (call ".")? IDENTIFIER "=" assignment | logic_or
logic_or       → logic_and ("or" logic_and)*
logic_and      → ternary ("and" ternary)*
ternary        → equality "?" equality ":" equality
equality       → comparison ( ( "!=" | "==" ) comparison )*
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
term           → factor ( ( "-" | "+" ) factor )*
factor         → unary ( ( "/" | "*" ) unary )*
unary          → ( "!" | "-" ) unary | primary | call
lambda        → "fun" "(" parameters? ")" block
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )*
arguments      → expression ("," expression)*
primary        → NUMBER | STRING | IDENTIFIER | "true" | "false" | "nil" | "(" expression ")"
*/

type Parser struct {
	tokens []token.Token
	current int
	HadError bool
}

func New(tokens []token.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(token.VAR) {
		return p.variableDeclaration()
	}
	if p.match(token.FUN) {
		return p.functionDeclaration()
	}
	if p.match(token.CLASS) {
		return p.classDeclaration()
	}
	return p.statement()
}

func (p *Parser) variableDeclaration() *ast.Var {
	name := p.consume(token.IDENTIFIER, "Expected variable name")
	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}
	p.consume(token.SEMICOLON, "Expected \";\" after variable declaration")
	return &ast.Var{Name: name, Initializer: initializer}
}

func (p *Parser) functionDeclaration() *ast.Function {
	name := p.consume(token.IDENTIFIER, "Expected function name.")
	p.consume(token.LEFT_PAREN, "Expected \"(\" after function name.")
	parameters := make([]token.Token, 0)
	for !p.match(token.RIGHT_PAREN) && !p.isAtEnd() {
		param := p.consume(token.IDENTIFIER, "Expected parameter.")
		parameters = append(parameters, param)
		if p.match(token.RIGHT_PAREN) {
			break
		}
		p.consume(token.COMMA, "Expected \",\" after parameter.")
	}
	p.consume(token.LEFT_BRACE, "Expected \"{\" before function body.")
	body := p.block().Statements
	return &ast.Function{Name: name, Parameters: parameters, Body: body}
}

func (p *Parser) classDeclaration() *ast.Class {
	name := p.consume(token.IDENTIFIER, "Expected class name.")
	p.consume(token.LEFT_BRACE, "Expected \"{\" after before class body.")
	methods := make([]*ast.Function, 0)
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.functionDeclaration())
	}
	p.consume(token.RIGHT_BRACE, "Expected \"}\" after class bdoy.")
	return &ast.Class{Name: name, Methods: methods}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.PRINT) {
		expr := p.expression()
		p.consume(token.SEMICOLON, "Expected \";\" after expression")
		return &ast.PrintStmt{Expression: expr}
	}
	if p.match(token.IF) {
		p.consume(token.LEFT_PAREN, "Expected \"(\" after \"if\"")
		condition := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected \")\" after if condition")
		thenBranch := p.statement()
		var elseBranch ast.Stmt
		if p.match(token.ELSE) {
			elseBranch = p.statement()
		}
		return &ast.If{Condition: condition, ElseBranch: elseBranch, ThenBranch: thenBranch}
	}
	if p.match(token.LEFT_BRACE) {
		return p.block()
	}
	if p.match(token.BREAK) {
		p.consume(token.SEMICOLON, "Expected \";\" after \"break\"")
		return &ast.Break{}
	}
	if p.match(token.CONTINUE) {
		p.consume(token.SEMICOLON, "Expected \";\" after \"continue\"")
		return &ast.Continue{}
	}
	if p.match(token.WHILE) {
		p.consume(token.LEFT_PAREN, "Expected \"(\" after \"while\"")
		conditon := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected \")\" after condition")
		body := p.statement()
		return &ast.While{Condition: conditon, Body: body}
	}
	if p.match(token.FOR) {
		p.consume(token.LEFT_PAREN, "Expected \"(\" after \"for\"")
		var initializer ast.Stmt
		if p.match(token.SEMICOLON) {
			initializer = nil
		} else if p.check(token.VAR) {
			initializer = p.declaration()
		} else {
			initializer = p.expressionStatement()
		}
		var condition ast.Expr
		if !p.check(token.SEMICOLON) {
			condition = p.expression()
		}
		p.consume(token.SEMICOLON, "Expected \";\" after loop condition")
		var increment ast.Expr
		if !p.check(token.SEMICOLON) {
			increment = p.expression()
		}
		p.consume(token.RIGHT_PAREN, "Expected \")\" after increment expression")
		body := p.statement()
		if condition == nil {
			condition = &ast.Literal{Value: true}
		}
		loop := &ast.For{Body: body, Condition: condition, Initializer: initializer, Increment: increment}
		// wrap the loop in a block so that it gets its own scope
		statements := make([]ast.Stmt, 1)
		statements[0] = loop
		return &ast.Block{Statements: statements}
	}
	if p.match(token.RETURN) {
		keyword := p.previous()
		var value ast.Expr
		if !p.check(token.SEMICOLON) {
			value = p.expression()
		}
		p.consume(token.SEMICOLON, "Expected \";\" after return value")
		return &ast.Return{Keyword: keyword, Value: value}
	}
	return p.expressionStatement()
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	if p.isAtEnd() {
		p.current = 0
		expr := p.expression()
		return &ast.PrintStmt{Expression: expr}
	}
	p.consume(token.SEMICOLON, "Expected \";\" after expression")
	return &ast.ExprStmt{Expression: expr}
}

func (p *Parser) block() *ast.Block {
	var statements []ast.Stmt
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(token.RIGHT_BRACE, "Expected \"}\" after block")
	return &ast.Block{Statements: statements}
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()
	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		switch e := expr.(type) {
		case *ast.Variable:
			return &ast.Assign{Name: e.Name, Value: value}
		case *ast.Get:
			return &ast.Set{Name: e.Name, Object: e.Object, Value: value}
		}
		p.handleError(equals, "Invalid assignment target")
	}
	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()
	if p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		return &ast.Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.ternary()
	if p.match(token.AND) {
		operator := p.previous()
		right := p.ternary()
		return &ast.Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) ternary() ast.Expr {
	expr := p.equality()
	if p.match(token.QUESTION_MARK) {
		thenExpr := p.equality()
		if p.match(token.COLON) {
			elseExpr := p.equality()
			expr = &ast.Ternary{Condition: expr, Then: thenExpr, Else: elseExpr}
		} else {
			p.handleError(p.peek(), "Unterminated ternary operator")
		}
	}
	return expr
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		expr := &ast.Unary{Operator: operator, Right: right}
		return expr
	}
	return p.lambda()
}

func (p *Parser) lambda() ast.Expr {
	if p.match(token.FUN) {
		p.consume(token.LEFT_PAREN, "Expected \"(\" after \"fun\"")
		parameters := make([]token.Token, 0)
		for !p.match(token.RIGHT_PAREN) && !p.isAtEnd() {
			param := p.consume(token.IDENTIFIER, "Expected parameter")
			parameters = append(parameters, param)
			if p.match(token.RIGHT_PAREN) {
				break
			}
			p.consume(token.COMMA, "Expected \",\" after parameter")
		}
		p.consume(token.LEFT_BRACE, "Expected \"{\" before function body")
		body := p.block().Statements
		expr := &ast.Lambda{Parameters: parameters, Body: body}
		return expr
	}
	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()
	for {
		if p.match(token.LEFT_PAREN){
			expr = p.finishCall(expr)
		} else if p.match(token.DOT) {
			name := p.consume(token.IDENTIFIER, "Expected property name after \".\".")
			expr = &ast.Get{Object: expr, Name: name}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee ast.Expr) *ast.Call {
	arguments := make([]ast.Expr, 0)
	for !p.check(token.RIGHT_PAREN) && !p.isAtEnd() {
		for {
			arguments = append(arguments, p.expression())
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	if len(arguments) > 255 {
		p.reportError(p.peek(), "Can't have more than 255 arguments")
	}
	paren := p.consume(token.RIGHT_PAREN, "Expect \")\" after arguments")
	return &ast.Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return &ast.Literal{Value: true}
	}
	if p.match(token.NULL) {
		return &ast.Literal{Value: nil}
	}
	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}
	}
	if p.match(token.IDENTIFIER) {
		return &ast.Variable{Name:p.previous()}
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}
	if p.match(token.THIS) {
		return &ast.This{Keyword: p.previous()}
	}
	p.handleError(p.peek(), "Expected expression.")
	return nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		switch(p.peek().Type) {
		case token.CLASS:
		case token.FUN:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) handleError(tk token.Token, message string) {
	p.reportError(tk, message)
	p.synchronize()
}

func (p *Parser) reportError(tk token.Token, message string) {
	p.HadError = true
	if tk.Type == token.EOF {
		fmt.Printf("[Line %v] Error at end: %v\n", tk.Line, message)
	} else {
		fmt.Printf("[Line %v] Error at \"%v\": %v\n", tk.Line, tk.Lexeme, message)
	}
}

func (p *Parser) consume(tokenType token.Type, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.handleError(p.peek(), message)
	return token.Token{}
}

func (p *Parser) match(types ...token.Type) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
		return p.previous()
	}
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current - 1]
}
