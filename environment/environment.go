package environment

import (
	"fmt"

	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

func New() Environment {
	return Environment{values: map[string]any{}}
}

type Environment struct {
	values map[string]any
}

func (e *Environment) Get(name token.Token) (any, error) {
	value, ok := e.values[name.Lexeme]
	if !ok {
		return nil, errors.NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'", name.Lexeme))
	}

	return value, nil
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name token.Token, value any) error {
	_, exists := e.values[name.Lexeme]
	if !exists {
		return errors.NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'", name.Lexeme))
	}

	e.values[name.Lexeme] = value

	return nil
}
