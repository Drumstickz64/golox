package interpreting

import (
	"fmt"
	"reflect"

	"github.com/Drumstickz64/golox/assert"
	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) Interpret(exp ast.Expr) error {
	value, err := exp.Accept(i)
	if err != nil {
		return err
	}

	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitLiteralExpr(exp *ast.LiteralExpr) (any, error) {
	return exp.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(exp *ast.GroupingExpr) (any, error) {
	return exp.Expression.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(exp *ast.UnaryExpr) (any, error) {
	right, err := exp.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch exp.Operator.Kind {
	case token.MINUS:
		if err := checkNumberOperandUnary(exp.Operator, right); err != nil {
			return nil, err
		}

		return -right.(float64), nil
	case token.BANG:
		return !isTruthy(right), nil
	}

	return assert.Unreachable(fmt.Sprintf("'%v' is a valid unary operator", exp.Operator.Kind)), nil
}

func (i *Interpreter) VisitBinaryExpr(exp *ast.BinaryExpr) (any, error) {
	left, err := exp.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	right, err := exp.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch exp.Operator.Kind {
	case token.PLUS:
		// is* functions are needed becuase calling reflect.TypeOf on a nil causes a panic
		if isNumber(left) && isNumber(right) {
			return left.(float64) + right.(float64), nil
		}

		if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		} else if isString(left) {
			return left.(string) + stringify(right), nil
		} else if isString(right) {
			return stringify(left) + right.(string), nil
		}

		return errors.NewRuntimeError(exp.Operator, "operands must be two numbers or one of them must be a string"), nil
	case token.MINUS:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.STAR:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.SLASH:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.GREATER:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.EQUAL_EQUAL:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left == right, nil
	case token.BANG_EQUAL:
		if err := checkNumberOperandBinary(exp.Operator, left, right); err != nil {
			return nil, err
		}
		return left != right, nil
	}

	return assert.Unreachable(fmt.Sprintf("'%v' is a valid binary operator", exp.Operator.Kind)), nil
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
