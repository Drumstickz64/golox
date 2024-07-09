package ast

import (
	"github.com/Drumstickz64/golox/token"
)

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitVariableExpr(*VariableExpr) (any, error)
}

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func (g *GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l *LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (u *UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

type VariableExpr struct {
	Name token.Token
}

func (v *VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}
