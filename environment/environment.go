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

func (e *Environment) Get(tok token.Token) (any, error) {
	value, ok := e.values[tok.Lexeme]
	if !ok {
		return nil, errors.NewRuntimeError(tok, fmt.Sprintf("undefined variable '%v'", tok.Lexeme))
	}

	return value, nil
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}
