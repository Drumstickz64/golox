package ast

import "github.com/Drumstickz64/golox/token"

type ExprVisitor interface {
	VisitBinary(exp *Binary) any
	VisitGrouping(exp *Grouping) any
	VisitLiteral(exp *Literal) any
	VisitUnary(exp *Unary) any
}

type Expr interface {
	Accept(ExprVisitor) any
}

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator token.Token
	Right Expr
}

func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}

