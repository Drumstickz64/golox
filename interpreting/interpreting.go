package interpreting

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Drumstickz64/golox/assert"
	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/environment"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

// used for storing AST nodes in maps because maps compare keys by value, but we need
// identical objects to be unique
type exprId string

type Interpreter struct {
	globals     *environment.Environment
	env         *environment.Environment
	locals      map[exprId]int
	isReturning bool
	returnValue any
}

func NewInterpreter() *Interpreter {
	globals := environment.New()

	globals.Define("clock", &nativeFunction{
		arity: 0,
		call: func(interpreter *Interpreter, arguments []any) (any, error) {
			return float64(time.Now().Unix()), nil
		},
	})

	globals.Define("str", &nativeFunction{
		arity: 1,
		call: func(interpreter *Interpreter, arguments []any) (any, error) {
			return stringify(arguments[0]), nil
		},
	})

	return &Interpreter{
		globals: globals,
		env:     globals,
		locals:  map[exprId]int{},
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) error {
	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}
	for _, argument := range expr.Arguments {
		value, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, value)
	}

	callable, ok := callee.(Callable)
	if !ok {
		return nil, errors.NewRuntimeError(expr.Paren, "can only call functions and classes")
	}

	if len(arguments) != callable.Arity() {
		return nil, errors.NewRuntimeError(expr.Paren, fmt.Sprintf("expected %d arguments but got %d instead", callable.Arity(), len(arguments)))
	}

	return callable.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*Instance)
	if !ok {
		return nil, errors.NewRuntimeError(expr.Name, "only instances can have properties")
	}

	return instance.Get(expr.Name)
}

func (i *Interpreter) VisitSetExpr(expr *ast.SetExpr) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*Instance)
	if !ok {
		return nil, errors.NewRuntimeError(expr.Name, "only instances can have fields")
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	instance.Set(expr.Name, value)

	return value, nil
}

func (i *Interpreter) VisitSuperExpr(expr *ast.SuperExpr) (any, error) {
	distance := i.locals[makeExprId(expr)]
	superClass := i.env.GetAt(distance, "super").(*class)
	object := i.env.GetAt(distance-1, "this").(*Instance) // env for 'this' is always one hop away from the env for 'super'
	method, ok := superClass.findMethod(expr.Method.Lexeme)
	if !ok {
		return nil, errors.NewRuntimeError(expr.Method, fmt.Sprintf("undefined property '%s'", expr.Method.Lexeme))
	}

	return method.bind(object), nil
}

func (i *Interpreter) VisitThisExpr(expr *ast.ThisExpr) (any, error) {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Kind {
	case token.MINUS:
		if err := checkNumberOperandUnary(expr.Operator, right); err != nil {
			return nil, err
		}

		return -right.(float64), nil
	case token.BANG:
		return !isTruthy(right), nil
	}

	assert.Unreachable(fmt.Sprintf("'%v' is a valid unary operator", expr.Operator.Kind))
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Kind {
	case token.PLUS:
		// is* functions are needed becuase calling reflect.TypeOf on a nil causes a panic
		if isNumber(left) && isNumber(right) {
			return left.(float64) + right.(float64), nil
		}

		if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		}

		return errors.NewRuntimeError(expr.Operator, "operands must be two numbers or two strings"), nil
	case token.MINUS:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.STAR:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.SLASH:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		if right.(float64) == 0.0 {
			return nil, errors.NewRuntimeError(expr.Operator, "attempted to divide by zero")
		}
		return left.(float64) / right.(float64), nil
	case token.GREATER:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.EQUAL_EQUAL:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left == right, nil
	case token.BANG_EQUAL:
		if err := checkNumberOperandBinary(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left != right, nil
	}

	assert.Unreachable(fmt.Sprintf("'%v' is a valid binary operator", expr.Operator.Kind))
	return nil, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	assert.That(expr.Operator.Kind == token.OR || expr.Operator.Kind == token.AND, "logical operator is either 'and' or 'or'")

	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Kind == token.OR && isTruthy(left) ||
		expr.Operator.Kind == token.AND && !isTruthy(left) {
		return left, nil
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) (any, error) {
	return i.lookupVariable(expr.Name, expr)
}

func (i Interpreter) lookupVariable(name token.Token, expr ast.Expr) (any, error) {
	distance, ok := i.locals[makeExprId(expr)]
	if !ok {
		return i.globals.Get(name)
	}

	return i.env.GetAt(distance, name.Lexeme), nil
}

func (i *Interpreter) VisitAssignmentExpr(expr *ast.AssignmentExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[makeExprId(expr)]
	if ok {
		i.env.AssignAt(distance, expr.Name, value)
	} else {
		if err := i.globals.Assign(expr.Name, value); err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(stringify(value))

	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	var value any
	if stmt.Value != nil {
		var err error
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	i.isReturning = true
	i.returnValue = value
	return value, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	for {
		condition, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}

		if !isTruthy(condition) {
			return nil, nil
		}

		if err := i.execute(stmt.Body); err != nil {
			return nil, err
		}

		if i.isReturning {
			return nil, nil
		}
	}
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		if err := i.execute(stmt.ThenBranch); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if stmt.ElseBranch == nil {
		return nil, nil
	}

	if err := i.execute(stmt.ElseBranch); err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	_, err := i.evaluate(stmt.Expression)
	return nil, err
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	if err := i.executeBlock(stmt.Statements, environment.WithEnclosing(i.env)); err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitClassStmt(stmt *ast.ClassStmt) (any, error) {
	var superClass *class = nil
	if stmt.SuperClass != nil {
		superClassValue, err := i.evaluate(stmt.SuperClass)
		if err != nil {
			return nil, err
		}

		superClassInstance, ok := superClassValue.(*class)
		if !ok {
			return nil, errors.NewRuntimeError(stmt.SuperClass.Name, "superclass must be a class")
		}

		superClass = superClassInstance
	}

	i.env.Define(stmt.Name.Lexeme, nil)

	if superClass != nil {
		i.env = environment.WithEnclosing(i.env)
		i.env.Define("super", superClass)
	}

	methods := map[string]*function{}
	for _, method := range stmt.Methods {
		fun := &function{
			declaration:   method,
			closure:       i.env,
			isInitializer: method.Name.Lexeme == "init",
		}
		methods[method.Name.Lexeme] = fun
	}

	class := &class{
		name:       stmt.Name.Lexeme,
		superClass: superClass,
		methods:    methods,
	}

	if superClass != nil {
		i.env = i.env.Enclosing()
	}

	i.env.Assign(stmt.Name, class)

	return nil, nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	var value any = nil

	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.env.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) (any, error) {
	fun := &function{
		declaration:   stmt,
		closure:       i.env,
		isInitializer: false,
	}
	i.env.Define(stmt.Name.Lexeme, fun)
	return nil, nil
}

func checkNumberOperandUnary(operator token.Token, operand any) error {
	_, ok := operand.(float64)
	if ok {
		return nil
	}

	return errors.NewRuntimeError(operator, "operand must be a number")
}

func checkNumberOperandBinary(operator token.Token, left, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return nil
	}

	return errors.NewRuntimeError(operator, "operand must be a number")
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, env *environment.Environment) error {
	previousEnv := i.env
	i.env = env
	defer func() { i.env = previousEnv }()
	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}

		if i.isReturning {
			return nil
		}
	}

	return nil
}

func (i *Interpreter) execute(stmt ast.Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	id := makeExprId(expr)
	i.locals[id] = depth
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func isTruthy(item any) bool {
	if item == nil {
		return false
	}

	if item == false {
		return false
	}

	return true
}

func stringify(item any) string {
	if item == nil {
		return "nil"
	}

	return fmt.Sprint(item)
}

func isNumber(value any) bool {
	if value == nil {
		return false
	}

	return reflect.TypeOf(value).Kind() == reflect.Float64
}

func isString(value any) bool {
	if value == nil {
		return false
	}

	return reflect.TypeOf(value).Kind() == reflect.String
}

func makeExprId(expr ast.Expr) exprId {
	return exprId(fmt.Sprintf("%p", expr))
}
