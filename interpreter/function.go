package interpreter

import (
	"github.com/singurty/lox/ast"
	"github.com/singurty/lox/environment"
	"github.com/singurty/lox/token"
)

type callable interface {
	arity() int
	call([]interface{}) (interface{}, error)
	String() string
}

type nativeFunction struct {
	nativeCallable func([]interface{}) (interface{}, error)
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
	declaration *ast.Function
	closure *environment.Environment
	isInitializer bool
}

func (u *userFunction) arity() int {
	return len(u.declaration.Parameters)
}

func (u *userFunction) String() string {
	return "<fun " + u.declaration.Name.Lexeme + ">"
}

func (u *userFunction) call(arguments []interface{}) (interface{}, error) {
	if u.isInitializer {
		funCall(u.closure, u.declaration.Parameters, u.arity(), u.declaration.Body, arguments)
		return u.closure.GetAt(0, "this")
	}
	return funCall(u.closure, u.declaration.Parameters, u.arity(), u.declaration.Body, arguments)
}

func (u *userFunction) Bind(instance *Instance) (*userFunction, error) {
	env := environment.Local(u.closure)
	err := env.Define("this", instance)
	if err != nil {
		return nil, err
	}
	return &userFunction{declaration: u.declaration, closure: env, isInitializer: u.isInitializer}, nil
}

type lambda struct {
	declaration *ast.Lambda
	closure *environment.Environment
}

func (l *lambda) arity() int {
	return len(l.declaration.Parameters)
}

func (l *lambda) String() string {
	return "<lambda>"
}

func (l *lambda) call(arguments []interface{}) (interface{}, error) {
	return funCall(l.closure, l.declaration.Parameters, l.arity(), l.declaration.Body, arguments)
}

func funCall(closure *environment.Environment, parameters []token.Token, arity int, body []ast.Stmt, arguments []interface{}) (interface{}, error) {
	envFun := environment.Local(closure)
	for i := 0; i < arity; i++ {
		envFun.Define(parameters[i].Lexeme, arguments[i])
	}
	err := executeBlock(body, envFun)
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
