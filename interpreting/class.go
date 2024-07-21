package interpreting

import (
	"fmt"

	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

type class struct {
	name       string
	superClass *class
	methods    map[string]*function
}

func (c *class) Call(interpreter *Interpreter, arguments []any) (any, error) {
	ins := NewInstance(c)

	initializer, ok := c.findMethod("init")
	if ok {
		initializer.bind(ins).Call(interpreter, arguments)
	}

	return ins, nil
}

func (c *class) Arity() int {
	initializer, ok := c.findMethod("init")
	if !ok {
		return 0
	}

	return initializer.Arity()
}

func (c *class) String() string {
	return fmt.Sprintf("<class %s>", c.name)
}

func (c *class) findMethod(name string) (*function, bool) {
	method, ok := c.methods[name]
	if ok {
		return method, true
	}

	if c.superClass != nil {
		return c.superClass.findMethod(name)
	}

	return nil, false
}

type Instance struct {
	class  *class
	fields map[string]any
}

func NewInstance(class *class) *Instance {
	return &Instance{
		class:  class,
		fields: map[string]any{},
	}
}

func (i *Instance) Get(name token.Token) (any, error) {
	value, ok := i.fields[name.Lexeme]
	if ok {
		return value, nil
	}

	method, ok := i.class.findMethod(name.Lexeme)
	if ok {
		return method.bind(i), nil
	}

	return nil, errors.NewRuntimeError(name, fmt.Sprintf("undefined property '%s'", name.Lexeme))
}

func (i *Instance) Set(name token.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *Instance) String() string {
	return fmt.Sprintf("<instance of class %s>", i.class.name)
}
