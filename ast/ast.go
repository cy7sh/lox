package ast

import (
	"fmt"
	"strings"

	"github.com/singurty/lox/token"
)

type Expr interface {
	String() string
}

type Ternary struct {
	Condition Expr
	Then Expr
	Else Expr
}

func (t *Ternary) String() string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(t.Condition.String())
	sb.WriteString(") then (")
	sb.WriteString(t.Then.String())
	sb.WriteString(") else (")
	sb.WriteString(t.Else.String())
	sb.WriteString(")")
	return sb.String()
}

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

// pretty print for binary
func (b *Binary) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(b.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(b.Left.String())
	sb.WriteString(" ")
	sb.WriteString(b.Right.String())
	sb.WriteString(")")
	return sb.String()
}

// for parenthesized expressions
type Grouping struct {
	Expression Expr
}

func (g *Grouping) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(g.Expression.String())
	sb.WriteString(")")
	return sb.String()
}

type Literal struct {
	Value interface{}
}

func (l *Literal) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v", l.Value))
	return sb.String()
}

type Unary struct {
	Operator token.Token
	Right Expr
}

func (u *Unary) String() string {
	var sb strings.Builder
	sb.WriteString(u.Operator.Lexeme)
	sb.WriteString(u.Right.String())
	return sb.String()
}

type Assign struct {
	Name token.Token
	Value Expr
}

func (a *Assign) String() string {
	var sb strings.Builder
	sb.WriteString(a.Name.Lexeme)
	sb.WriteString(" = ")
	sb.WriteString(a.Value.String())
	return sb.String()
}

type Get struct {
	Name token.Token
	Object Expr
}

func (g *Get) String() string {
	return fmt.Sprintf("(get %v)", g.Name.Lexeme)
}

type Set struct {
	Name token.Token
	Object Expr
	Value Expr
}

func (s *Set) String() string {
	return fmt.Sprintf("(set %v %v)", s.Name.Lexeme, s.Value)
}

type Stmt interface {
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type Block struct {
	Statements []Stmt
}

type Var struct {
	Name token.Token
	Initializer Expr
}

type Variable struct {
	Name token.Token
}

func (v *Variable) String() string {
	return v.Name.Lexeme
}

type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type Logical struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func (l *Logical) String() string {
	var sb strings.Builder
	sb.WriteString(l.Left.String())
	sb.WriteString(" or ")
	sb.WriteString(l.Right.String())
	return sb.String()
}

type Lambda struct {
	Parameters []token.Token
	Body []Stmt

}

func (l *Lambda) String() string {
	var sb strings.Builder
	sb.WriteString("lambda parameters: ")
	return sb.String()
}

type This struct {
	Keyword token.Token
}

func (t *This) String() string {
	return "<this>"
}

type While struct {
	Condition Expr
	Body Stmt
}

// keep track of increment expression because it should be executed even when continuing
type For struct {
	Body Stmt
	Condition Expr
	Increment Expr
	Initializer Stmt
}

type Break struct {
}

type Continue struct {
}

type Call struct {
	Callee Expr
	Paren token.Token
	Arguments []Expr
}

func (c *Call) String() string {
	var sb strings.Builder
	sb.WriteString(c.Callee.String())
	sb.WriteString("(")
	for _, v := range c.Arguments {
		sb.WriteString(v.String())
		sb.WriteString(",")
	}
	return sb.String()
}

type Function struct {
	Name token.Token
	Parameters []token.Token
	Body []Stmt
}

type Return struct {
	Keyword token.Token
	Value Expr
}

type Class struct {
	Name token.Token
	Methods []*Function
}
