package environment

import (
	"fmt"

	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func New() *Environment {
	return &Environment{
		enclosing: nil,
		values:    map[string]any{},
	}
}

func WithEnclosing(env *Environment) *Environment {
	return &Environment{
		values:    map[string]any{},
		enclosing: env,
	}
}

func (e *Environment) Get(name token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		value, err := e.enclosing.Get(name)
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	return nil, errors.NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'", name.Lexeme))
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name token.Token, value any) error {
	_, exists := e.values[name.Lexeme]
	if exists {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		_, exists := e.values[name.Lexeme]
		if exists {
			e.values[name.Lexeme] = value
			return nil
		}
	}

	return errors.NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'", name.Lexeme))
}
