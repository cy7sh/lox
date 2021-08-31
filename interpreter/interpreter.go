package interpreter

import (
	"fmt"

	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/token"
)

type RuntimeError struct {
	line int
	where string
	message string
}

func Eval(node ast.Expr) (interface{}, error) {
	switch n := node.(type) {
		case *ast.Literal:
			return n.Value, nil
		case *ast.Grouping:
			return Eval(n.Expression)
		case *ast.Unary:
			right, err := Eval(n.Right)
			if err != nil {
				return nil, err
			}
			switch n.Operator.Type {
			case token.MINUS:
				err := checkNumberOperand(n.Operator, right)
				if err != nil {
					return nil, err
				}
				return -right.(float64), nil
			case token.BANG:
				return !isTrue(right), nil
			}
		case *ast.Binary:
			left, err := Eval(n.Left)
			if err != nil {
				return nil, err
			}
			right, err := Eval(n.Right)
			if err != nil {
				return nil, err
			}
			switch n.Operator.Type {
				case token.MINUS:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					return left.(float64) - right.(float64), nil
				case token.SLASH:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					if right.(float64) == 0 {
						return nil, &RuntimeError{line: n.Operator.Line, where: n.Operator.Lexeme, message: "Divide by zero"}
					}
					return left.(float64) / right.(float64), nil
				case token.STAR:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					return left.(float64) * right.(float64), nil
				case token.PLUS:
					switch l := left.(type) {
						case float64:
							switch r := right.(type) {
							case float64:
								return l + r, nil
							}
						case string:
							switch r := right.(type) {
							case string:
								return l + r, nil
							}
					}
					return nil, &RuntimeError{line: n.Operator.Line, where: n.Operator.Lexeme, message: "Operands must be eithier numbers or strings"}
				case token.GREATER:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					return left.(float64) > right.(float64), nil
				case token.GREATER_EQUAL:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					return left.(float64) >= right.(float64), nil
				case token.LESS:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					return left.(float64) < right.(float64), nil
				case token.LESS_EQUAL:
					err := checkNumberOperands(n.Operator, right, left)
					if err != nil {
						return nil, err
					}
					return left.(float64) <= right.(float64), nil
				case token.EQUAL_EQUAL:
					return isEqual(left, right), nil
				case token.BANG_EQUAL:
					return !isEqual(left, right), nil
			}
		case *ast.Ternary:
			conditon, err := Eval(n.Condition)
			if err != nil {
				return nil, err
			}
			if isTrue(conditon) {
				return Eval(n.Then)
			} else {
				return Eval(n.Else)
			}
	}
	return nil, &RuntimeError{message: "Error evaluating expression"}
}

func (err *RuntimeError) Error() string {
	if err.where == "" {
		return fmt.Sprintf("[Line %v] RuntimeError: %v", err.line, err.message)
	}
	return fmt.Sprintf("[Line %v] RuntimeError at \"%v\": %v", err.line, err.where, err.message)
}

func checkNumberOperand(operator token.Token, operand interface{}) error {
	_, ok := operand.(float64)
	if !ok {
		return &RuntimeError{line: operator.Line, where: operator.Lexeme, message: "Operand must be a number"}
	}
	return nil
}

func checkNumberOperands(operator token.Token, operand1, operand2 interface{}) error {
	_, ok := operand1.(float64)
	if !ok {
		return &RuntimeError{line: operator.Line, where: operator.Lexeme, message:"Operand must be a number"}
	}
	_, ok = operand2.(float64)
	if !ok {
		return &RuntimeError{line: operator.Line, where: operator.Lexeme, message:"Operand must be a number"}
	}
	return nil
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
