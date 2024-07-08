package expr

import (
	"fmt"
)

type Printer struct{}

func NewPrinter() Printer {
	return Printer{}
}

func (p Printer) Print(exp Expr) string {
	res, _ := exp.Accept(p)
	return res.(string)
}

func (p Printer) VisitBinary(exp *Binary) (any, error) {
	return p.parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right), nil
}

func (p Printer) VisitGrouping(exp *Grouping) (any, error) {
	return p.parenthesize("group", exp.Expression), nil
}

func (p Printer) VisitLiteral(exp *Literal) (any, error) {
	if exp.Value == nil {
		return "nil", nil
	}

	return fmt.Sprint(exp.Value), nil
}

func (p Printer) VisitUnary(exp *Unary) (any, error) {
	return p.parenthesize(exp.Operator.Lexeme, exp.Right), nil
}

func (p Printer) parenthesize(name string, exps ...Expr) string {
	result := "(" + name

	for _, exp := range exps {
		res, _ := exp.Accept(p)
		result += " " + res.(string)
	}

	result += ")"

	return result
}
