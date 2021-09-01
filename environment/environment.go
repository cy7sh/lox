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
	_, ok := e.environment[variable]
	if ok {
		e.environment[variable] = value
		return nil
	} else {
		return errors.New("Undefined variable \"" + variable + "\"")
	}
}

func (e *Environment) Get(variable string) (interface{}, error) {
	value, ok := e.environment[variable]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("Undefined variable \"" + variable + "\"")
	}
}
