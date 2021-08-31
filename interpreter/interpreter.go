package interpreter

import (
	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/token"
)

func Eval(node ast.Expr) interface{} {
	switch n := node.(type) {
		case *ast.Literal:
			return n.Value
		case *ast.Grouping:
			return Eval(n.Expression)
		case *ast.Unary:
			right := Eval(n.Right)
			switch n.Operator.Type {
			case token.MINUS:
				return -right.(float64)
			case token.BANG:
				return !isTrue(right)
			}
		case *ast.Binary:
			left := Eval(n.Left)
			right := Eval(n.Right)
			switch n.Operator.Type {
				case token.MINUS:
					return left.(float64) - right.(float64)
				case token.SLASH:
					return left.(float64) / right.(float64)
				case token.STAR:
					return left.(float64) * right.(float64)
				case token.PLUS:
					switch l := left.(type) {
						case float64:
							switch r := right.(type) {
							case float64:
								return l + r
							}
						case string:
							switch r := right.(type) {
							case string:
								return l + r
							}
					}
				case token.GREATER:
					return left.(float64) > right.(float64)
				case token.GREATER_EQUAL:
					return left.(float64) >= right.(float64)
				case token.LESS:
					return left.(float64) < right.(float64)
				case token.LESS_EQUAL:
					return left.(float64) <= right.(float64)
				case token.EQUAL:
					return isEqual(left, right)
				case token.BANG_EQUAL:
					return !isEqual(left, right)
			}
	}
	panic("Error evaluating expression")
}

func isTrue(value interface{}) bool {
	if value == nil {
		return false
	} else if b, ok := value.(bool); ok {
		return b
	}
	return false
}

func isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}
