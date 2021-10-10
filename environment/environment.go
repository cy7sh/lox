package environment

import (
	"errors"
)

type Environment struct {
	environment map[string]interface{}
	enclosing *Environment
}

func Global() *Environment {
	return  &Environment{environment: make(map[string]interface{}), enclosing: nil}
}

func Local(enclosing *Environment) *Environment {
	return &Environment{environment: make(map[string]interface{}), enclosing: enclosing}
}

func (e *Environment) Define(variable string, value interface{}) error {
	_, ok := e.environment[variable]
	if ok {
		return errors.New("Redeclaration of \"" + variable + "\"")
	}
	e.environment[variable] = value
	return nil
}

func (e *Environment) Assign(variable string, value interface{}) error {
	if _, ok := e.environment[variable]; ok {
		e.environment[variable] = value
		return nil
	} else {
		if e.enclosing != nil {
			return e.enclosing.Assign(variable, value)
		}
		return errors.New("Undefined variable \"" + variable + "\"")
	}
}

func (e *Environment) AssignAt(distance int, variable string, value interface{}) error {
	env := e.ancestor(distance).environment
	if _, ok := env[variable]; ok {
		env[variable] = value
		return nil
	} else {
		return errors.New("Undefined variable \"" + variable + "\"")
	}
}

func (e *Environment) Get(variable string) (interface{}, error) {
	value, ok := e.environment[variable]
	if ok {
		if value == nil {
			return nil, errors.New("Uninitialized variable \"" + variable + "\"")
		}
		return value, nil
	} else {
		if e.enclosing != nil {
			return e.enclosing.Get(variable)
		}
		return nil, errors.New("Undefined variable \"" + variable + "\"")
	}
}

func (e *Environment) GetAt(distance int, variable string) (interface{}, error) {
	if value, ok := e.ancestor(distance).environment[variable]; ok {
		if value == nil {
			return nil, errors.New("Uninitialized variable \"" + variable + "\"")
		}
		return value, nil
	} else {
		return nil, errors.New("Undefined variable \"" + variable + "\"")
	}
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}
