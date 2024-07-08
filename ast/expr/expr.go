package expr

import "github.com/Drumstickz64/golox/token"

type ExprVisitor interface {
	VisitBinary(exp *Binary) (any, error)
	VisitGrouping(exp *Grouping) (any, error)
	VisitLiteral(exp *Literal) (any, error)
	VisitUnary(exp *Unary) (any, error)
}

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(u)
}
