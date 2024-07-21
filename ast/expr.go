package ast

import (
	"github.com/Drumstickz64/golox/token"
)

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitLogicalExpr(*LogicalExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitCallExpr(*CallExpr) (any, error)
	VisitGetExpr(*GetExpr) (any, error)
	VisitSetExpr(*SetExpr) (any, error)
	VisitSuperExpr(*SuperExpr) (any, error)
	VisitThisExpr(*ThisExpr) (any, error)
	VisitVariableExpr(*VariableExpr) (any, error)
	VisitAssignmentExpr(*AssignmentExpr) (any, error)
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

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l *LogicalExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(l)
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

type CallExpr struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (c *CallExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitCallExpr(c)
}

type GetExpr struct {
	Object Expr
	Name   token.Token
}

func (g *GetExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGetExpr(g)
}

type SetExpr struct {
	Object Expr
	Name   token.Token
	Value  Expr
}

func (s *SetExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitSetExpr(s)
}

type SuperExpr struct {
	Keyword token.Token
	Method  token.Token
}

func (s *SuperExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitSuperExpr(s)
}

type ThisExpr struct {
	Keyword token.Token
}

func (t *ThisExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitThisExpr(t)
}

type VariableExpr struct {
	Name token.Token
}

func (v *VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}

type AssignmentExpr struct {
	Name  token.Token
	Value Expr
}

func (a *AssignmentExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssignmentExpr(a)
}
