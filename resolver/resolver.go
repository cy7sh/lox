package resolver

import (
	"errors"

	"github.com/singurty/lox/token"
	"github.com/singurty/lox/ast"
)

type Resolver struct {
	scopes *stack
	locals *stack
}

type stack struct {
	data []map[interface{}]interface{}
}

func (r *Resolver) NewResolver() {
	r.scopes.data = make([]map[interface{}]interface{}, 0)
}

// add new scope
func (s *stack) push(entry map[interface{}]interface{}) {
	s.data = append(s.data, entry)
}

// remove current scope
func (s *stack) pop() map[interface{}]interface{} {
	n := len(s.data) -1
	entry := s.data[n]
	s.data = s.data[:n]
	return entry
}

// get the map at the top of the stack without removing it
func (s *stack) peek() map[interface{}]interface{} {
	return s.data[len(s.data) - 1]
}

func (s *stack) isEmpty() bool {
	return len(s.data) == 0
}

func (s *stack) len() int {
	return len(s.data)
}

func (s *stack) get(index int) map[interface{}]interface{} {
	return s.data[index]
}

func (r *Resolver) blockStmt(block *ast.Block) {
	r.beginScope()
	r.endScope()
}

func (r *Resolver) beginScope() {
	r.scopes.push(make(map[interface{}]interface{}))
}

func (r *Resolver) endScope() {

}

func (r *Resolver) resolve(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(statement ast.Stmt) {

}

func (r *Resolver) resolveExpr(statement ast.Expr) {

}

func (r *Resolver) varStmt(statement *ast.Var) {
	r.declare(statement.Name.Lexeme)
	if statement.Initializer != nil {
		r.resolveExpr(statement.Initializer)
	}
	r.define(statement.Name.Lexeme)
}

func (r *Resolver) declare(name string) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name] = false
}

func (r *Resolver) define(name string) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name] = true
}

func (r *Resolver) variableExpr(expr *ast.Variable) error {
	if !r.scopes.isEmpty() && !r.scopes.peek()[expr.Name.Lexeme].(bool) {
		return errors.New("Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := r.scopes.len() - 1; i >= 0; i-- {
		if _, ok := r.scopes.get(i)[name.Lexeme]; ok {
//			interpreter.resolve(expr, len(r.scopes) - 1 - i)
		}
	}
}

func (r *Resolver) assignExpr(expr *ast.Assign) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) functionStmt(stmt *ast.Function) {
	r.declare(stmt.Name.Lexeme)
	r.define(stmt.Name.Lexeme)
	r.resolveFunction(stmt)
}

func (r *Resolver) resolveFunction(function *ast.Function) {
	r.beginScope()
	for _, param := range function.Parameters {
		r.declare(param.Lexeme)
		r.define(param.Lexeme)
	}
	r.resolve(function.Body)
	r.endScope()
}

func (r *Resolver) expressionStmt(stmt *ast.ExprStmt) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) ifStmt(stmt *ast.If) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
}

func (r *Resolver) printStmt(stmt *ast.PrintStmt) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) returnStmt(stmt *ast.Return) {
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
}

func (r *Resolver) whileStmt(stmt *ast.While) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

func (r *Resolver) forStmt(stmt *ast.For) {
	r.resolveExpr(stmt.Condition)
	r.resolveExpr(stmt.Increment)
	r.resolveStmt(stmt.Initializer)
	r.resolveStmt(stmt.Body)
}

func (r *Resolver) binaryExpr(expr *ast.Binary) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
}

func (r *Resolver) callExpr(expr *ast.Call) {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
}

func (r *Resolver) groupgingExpr(expr *ast.Grouping) {
	r.resolveExpr(expr.Expression)
}

func (r *Resolver) literalExpr(expr *ast.Literal) {

}

func (r *Resolver) logicalExpr(expr *ast.Logical) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
}

func (r *Resolver) unaryExpr(expr *ast.Unary) {
	r.resolveExpr(expr.Right)
}
