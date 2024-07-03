package ast

import (
	"fmt"
)

type Printer struct{}

func (ap Printer) Print(exp Expr) string {
	return exp.Accept(ap).(string)
}

func (ap Printer) VisitBinary(exp *Binary) any {
	return ap.parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right)
}

func (ap Printer) VisitGrouping(exp *Grouping) any {
	return ap.parenthesize("group", exp.Expression)
}

func (ap Printer) VisitLiteral(exp *Literal) any {
	if exp.Value == nil {
		return "nil"
	}

	return fmt.Sprint(exp.Value)
}

func (ap Printer) VisitUnary(exp *Unary) any {
	return ap.parenthesize(exp.Operator.Lexeme, exp.Right)
}

func (ap Printer) parenthesize(name string, exps ...Expr) string {
	result := "(" + name

	for _, exp := range exps {
		result += " " + exp.Accept(ap).(string)
	}

	result += ")"

	return result
}
