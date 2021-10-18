package interpreter

import (
	"github.com/singurty/lox/token"
)

type class struct {
	name string
	methods map[string]*userFunction
}

func newClass(name string, methods map[string]*userFunction) *class {
	return &class{name: name, methods: methods}
}

func (c *class) String() string {
	return "<class " + c.name + ">"
}

func (c *class) arity() int {
	initializer := c.findMethod("init")
	if initializer != nil {
		return initializer.arity()
	}
	return 0
}

func (c *class) call(arguments []interface{}) (interface{}, error) {
	instance := newInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer, err := initializer.Bind(instance)
		if err != nil {
			return nil, err
		}
		_, err = initializer.call(arguments)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (c *class) findMethod(name string) *userFunction {
	if value, ok := c.methods[name]; ok {
		return value
	} else {
		return nil
	}
}

type Instance struct {
	klass *class
	fields map[string]interface{}
}

func newInstance(klass *class) *Instance {
	return &Instance{klass: klass, fields: make(map[string]interface{})}
}

func (i *Instance) String() string {
	return "<instance " + i.klass.name + ">"
}

func (i *Instance) get(name token.Token) (interface{}, error) {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value, nil
	}
	method := i.klass.findMethod(name.Lexeme)
	if method != nil {
		method, err := method.Bind(i)
		if err != nil {
			return nil, err
		}
		return method, nil
	}
	return nil, &runtimeError{line: name.Line, message: "Undefined property \"" + name.Lexeme + "\"."}
}

func (i *Instance) set(name string, value interface{}) {
	i.fields[name] = value
}
