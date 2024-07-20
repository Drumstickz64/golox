package interpreting

import (
	"fmt"

	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

type class struct {
	Name    string
	Methods map[string]*function
}

func (c *class) Call(interpreter *Interpreter, arguments []any) (any, error) {
	ins := NewInstance(c)

	return ins, nil
}

func (c *class) Arity() int {
	return 0
}

func (c *class) String() string {
	return fmt.Sprintf("<class %s>", c.Name)
}

func (c *class) findMethod(name string) (*function, bool) {
	method, ok := c.Methods[name]
	if ok {
		return method, true
	}

	return nil, false
}

type Instance struct {
	Class  *class
	Fields map[string]any
}

func NewInstance(class *class) *Instance {
	return &Instance{
		Class:  class,
		Fields: map[string]any{},
	}
}

func (i *Instance) Get(name token.Token) (any, error) {
	value, ok := i.Fields[name.Lexeme]
	if ok {
		return value, nil
	}

	method, ok := i.Class.findMethod(name.Lexeme)
	if ok {
		return method.bind(i), nil
	}

	return nil, errors.NewRuntimeError(name, fmt.Sprintf("undefined property '%s'", name.Lexeme))
}

func (i *Instance) Set(name token.Token, value any) {
	i.Fields[name.Lexeme] = value
}

func (i *Instance) String() string {
	return fmt.Sprintf("<instance of class %s>", i.Class.Name)
}
