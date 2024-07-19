package environment

import (
	"fmt"

	"github.com/Drumstickz64/golox/assert"
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

func (e *Environment) GetAt(distance int, name string) any {
	value, ok := e.ancestor(distance).values[name]
	assert.That(ok, fmt.Sprintf("calls to environment.GetAt() always have existing variables, '%s' does not exist", name))
	return value
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
		return e.enclosing.Assign(name, value)
	}

	return errors.NewRuntimeError(name, fmt.Sprintf("undefined variable '%v'", name.Lexeme))
}

func (e *Environment) AssignAt(distance int, name token.Token, value any) {
	ancestor := e.ancestor(distance)
	_, ok := ancestor.values[name.Lexeme]
	assert.That(ok, fmt.Sprintf("calls to environment.AssignAt() always have existing variables, '%s' does not exist", name.Lexeme))
	ancestor.values[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	curr := e
	for i := 0; i < distance; i++ {
		assert.That(curr.enclosing != nil, fmt.Sprintf("distance passed to ancestor() is always correct, distance = %d", distance))
		curr = curr.enclosing
	}

	return curr
}
