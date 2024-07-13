package ast

import (
	"fmt"
)

type Printer struct{}

func NewPrinter() Printer {
	return Printer{}
}

func (p Printer) Print(expr Expr) (string, error) {
	res, err := expr.Accept(p)
	if err != nil {
		return "", err
	}

	return res.(string), err
}

func (p Printer) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p Printer) VisitLogicalExpr(expr *LogicalExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p Printer) VisitGroupingExpr(expr *GroupingExpr) (any, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p Printer) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}

	return fmt.Sprint(expr.Value), nil
}

func (p Printer) VisitUnaryExpr(expr *UnaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p Printer) VisitVariableExpr(expr *VariableExpr) (any, error) {
	return p.parenthesize("var " + expr.Name.Lexeme), nil
}

func (p Printer) VisitAssignmentExpr(expr *AssignmentExpr) (any, error) {
	return p.parenthesize("=", &LiteralExpr{Value: expr.Name}, expr.Value), nil
}

func (p *Printer) parenthesize(name string, exps ...Expr) string {
	result := "(" + name

	for _, expr := range exps {
		res, _ := expr.Accept(p)
		result += " " + res.(string)
	}

	result += ")"

	return result
}
