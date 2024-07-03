package ast

import (
	"fmt"
)

type Printer struct{}

func (p Printer) Print(exp Expr) string {
	return exp.Accept(p).(string)
}

func (p Printer) VisitBinary(exp *Binary) any {
	return p.parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right)
}

func (p Printer) VisitGrouping(exp *Grouping) any {
	return p.parenthesize("group", exp.Expression)
}

func (p Printer) VisitLiteral(exp *Literal) any {
	if exp.Value == nil {
		return "nil"
	}

	return fmt.Sprint(exp.Value)
}

func (p Printer) VisitUnary(exp *Unary) any {
	return p.parenthesize(exp.Operator.Lexeme, exp.Right)
}

func (p Printer) parenthesize(name string, exps ...Expr) string {
	result := "(" + name

	for _, exp := range exps {
		result += " " + exp.Accept(p).(string)
	}

	result += ")"

	return result
}

type RPNPrinter struct{}

func (p RPNPrinter) Print(exp Expr) string {
	return exp.Accept(p).(string)
}

func (p RPNPrinter) VisitBinary(exp *Binary) any {
	return p.rpnesize(exp.Operator.Lexeme, exp.Left, exp.Right)
}

func (p RPNPrinter) VisitGrouping(exp *Grouping) any {
	return exp.Expression.Accept(p)
}

func (p RPNPrinter) VisitLiteral(exp *Literal) any {
	if exp.Value == nil {
		return "nil"
	}

	return fmt.Sprint(exp.Value)
}

func (p RPNPrinter) VisitUnary(exp *Unary) any {
	return p.rpnesize(exp.Operator.Lexeme, exp.Right)
}

func (p RPNPrinter) rpnesize(name string, exps ...Expr) string {
	result := ""

	for _, exp := range exps {
		result += exp.Accept(p).(string) + " "
	}

	result += name

	return result
}
