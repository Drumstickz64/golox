package ast

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

func (p Printer) VisitBinaryExpr(exp *BinaryExpr) (any, error) {
	return p.parenthesize(exp.Operator.Lexeme, exp.Left, exp.Right), nil
}

func (p Printer) VisitGroupingExpr(exp *GroupingExpr) (any, error) {
	return p.parenthesize("group", exp.Expression), nil
}

func (p Printer) VisitLiteralExpr(exp *LiteralExpr) (any, error) {
	if exp.Value == nil {
		return "nil", nil
	}

	return fmt.Sprint(exp.Value), nil
}

func (p Printer) VisitUnaryExpr(exp *UnaryExpr) (any, error) {
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
