package interpreting

import (
	"fmt"
	"os"

	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/environment"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
}

type function struct {
	declaration *ast.FunctionStmt
	closure     *environment.Environment
}

func (f *function) Arity() int {
	return len(f.declaration.Parameters)
}

func (f *function) Call(interpreter *Interpreter, arguments []any) any {
	env := environment.WithEnclosing(f.closure)
	for i, param := range f.declaration.Parameters {
		env.Define(param.Lexeme, arguments[i])
	}

	if err := interpreter.executeBlock(f.declaration.Body, env); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
	}

	interpreter.isReturning = false
	return interpreter.returnValue
}

func (f *function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

type nativeFunction struct {
	arity int
	call  func(interpreter *Interpreter, arguments []any) any
}

func (f *nativeFunction) Arity() int {
	return f.arity
}

func (f *nativeFunction) Call(interpreter *Interpreter, arguments []any) any {
	return f.call(interpreter, arguments)
}

func (f *nativeFunction) String() string {
	return "<native fn>"
}
