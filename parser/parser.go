package parser

import (
	"fmt"

	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/token"
)

/*
program        → declaration* EOF
block          → "{" declaration* "}"
declaration    → varDecl | statement
varDecl        → "var" IDENTIFIER ("=" expression)? ";"
statement      → exprStmt | printStmt
exprStmt       → expression ";"
printStmt      → "print" expression ";"
expression     → ternary
assignment     → IDENTIFIER "=" assignment | ternary
ternary        → equality ? equality : equality
equality       → comparison ( ( "!=" | "==" ) comparison )*
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
term           → factor ( ( "-" | "+" ) factor )*
factor         → unary ( ( "/" | "*" ) unary )*
unary          → ( "!" | "-" ) unary | primary
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
		name := p.consume(token.IDENTIFIER, "Expected variable name")
		var initializer ast.Expr
		if p.match(token.EQUAL) {
			initializer = p.expression()
		}
		p.consume(token.SEMICOLON, "Expected \";\" after variable declaration")
		return &ast.Var{Name: name, Initializer: initializer}
	}
	return p.statement()
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.PRINT) {
		expr := p.expression()
		p.consume(token.SEMICOLON, "Expected \";\" after expression")
		return &ast.PrintStmt{Expression: expr}
	}
	if p.match(token.LEFT_BRACE) {
		var statements []ast.Stmt
		for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
			statements = append(statements, p.declaration())
		}
		p.consume(token.RIGHT_BRACE, "Expect \"}\" after block")
		return &ast.Block{Statements: statements}
	}
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expected \";\" after expression")
	return &ast.ExprStmt{Expression: expr}
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.ternary()
	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		switch e := expr.(type) {
		case *ast.Variable:
			name := e.Name
			return &ast.Assign{Name: name, Value: value}
		}
		p.handleError(equals, "Invalid assignment target")
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
	return p.primary()
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
	p.HadError = true
	if tk.Type == token.EOF {
		fmt.Printf("[Line %v] Error at end: %v\n", tk.Line, message)
	} else {
		fmt.Printf("[Line %v] Error at %v: %v\n", tk.Line, tk.Lexeme, message)
	}
	p.synchronize()
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
	if (p.isAtEnd()) {
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
