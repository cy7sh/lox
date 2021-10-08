package interpreter

import (
	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/environment"
)

type loxCallable func([]interface{}) (interface{}, error)

type callable interface {
	arity() int
	call([]interface{}) (interface{}, error)
	String() string
}

type nativeFunction struct {
	nativeCallable loxCallable
	arityNum int
}

func (n *nativeFunction) arity() int {
	return n.arityNum
}

func (n *nativeFunction) call(arguments []interface{}) (interface{}, error) {
	return n.nativeCallable(arguments)
}

func (n *nativeFunction) String() string {
	return "<native fun>"
}

type userFunction struct {
	delcaration *ast.Function
}

func (u *userFunction) arity() int {
	return len(u.delcaration.Parameters)
}

func (u *userFunction) String() string {
	return "<fun " + u.delcaration.Name.Lexeme + ">"
}

func (u *userFunction) call(arguments []interface{}) (interface{}, error) {
	envFun := environment.Local(global)
	for i := 0; i < u.arity(); i++ {
		envFun.Define(u.delcaration.Parameters[i].Lexeme, arguments[i])
	}
	err := executeBlock(u.delcaration.Body, envFun)
	if err != nil {
		returnValue, ok := err.(*returnError)
		if ok {
			return returnValue.value, nil
		} else {
			return nil, err
		}
	}
	return nil, nil
}
