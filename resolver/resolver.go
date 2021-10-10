package resolver

import (
	"errors"

	"github.com/singurty/lox/token"
	"github.com/singurty/lox/ast"
)

type Resolver struct {
	scopes []map[string]bool
}

func (r *Resolver) NewResolver() {
	r.scopes = make([]map[string]bool, 0)
}

// add new scope
func (r *Resolver) push(entry map[string]bool) {
	r.scopes = append(r.scopes, entry)
}

// remove current scope
func (r *Resolver) pop() map[string]bool {
	n := len(r.scopes) -1
	entry := r.scopes[n]
	r.scopes = r.scopes[:n]
	return entry
}

// add or change an entry on current scope
func put(scope map[string]bool, name string, defined bool) {
	scope[name] = defined
}

// get the map at the top of the stack without removing it
func (r *Resolver) peek() map[string]bool {
	return r.scopes[len(r.scopes) - 1]
}

func (r *Resolver) blockStmt(block *ast.Block) {
	r.beginScope()
	r.endScope()
}

func (r *Resolver) beginScope() {
	r.push(make(map[string]bool))
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
	if len(r.scopes) == 0 {
		return
	}
	put(r.peek(), name, false)
}

func (r *Resolver) define(name string) {
	if len(r.scopes) == 0 {
		return
	}
	put(r.peek(), name, true)
}

func (r *Resolver) variableExpr(expr *ast.Variable) error {
	if len(r.scopes) > 0 && !r.peek()[expr.Name.Lexeme] {
		return errors.New("Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
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
